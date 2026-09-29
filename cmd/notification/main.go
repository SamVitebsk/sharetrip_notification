package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"sharetrip_notification/internal/api"
	"sharetrip_notification/internal/clients/kafka"
	"sharetrip_notification/internal/env"
	"sharetrip_notification/internal/repository"
	"sharetrip_notification/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverPort, err := env.Int("HTTP_PORT", 9290)
	if err != nil {
		return err
	}
	if serverPort < 1 || serverPort > 65535 {
		return errors.New("HTTP_PORT должен быть в диапазоне от 1 до 65535")
	}

	cfg, err := repository.LoadConfig()
	if err != nil {
		return err
	}

	pool, err := repository.ConnectPostgres(ctx, cfg.DSN)
	if err != nil {
		return errors.New("не удалось подключиться к PostgreSQL: проверьте PG*-настройки и доступность БД")
	}
	defer pool.Close()

	pg, err := repository.NewPostgres(pool)
	if err != nil {
		return err
	}

	svc, err := service.New(pg, pg.RunInTx)
	if err != nil {
		return err
	}

	server, err := api.NewServer(svc)
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})
	app.Use(logger.New())
	app.Use(func(request *fiber.Ctx) error {
		request.SetUserContext(ctx)
		return request.Next()
	})

	api.RegisterRoutes(app, server)

	kafkaCfg, err := kafka.LoadConfig()
	if err != nil {
		return err
	}
	consumer := kafka.NewConsumer(kafkaCfg.Brokers, kafkaCfg.Topic, kafkaCfg.GroupID)
	defer func() { _ = consumer.Close() }()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Printf("Запуск сервиса уведомлений (HTTP) на порту :%d", serverPort)
		if err := app.Listen(fmt.Sprintf(":%d", serverPort)); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		log.Printf("Запуск Kafka Consumer, слушаем топик: %s", kafkaCfg.Topic)
		if err := consumer.Listen(gCtx, svc.ProcessTripPublished); err != nil {
			log.Printf("Ошибка Kafka Consumer: %v", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		log.Println("Плавное завершение работы сервиса (graceful shutdown)...")
		return app.ShutdownWithTimeout(10 * time.Second)
	})

	return g.Wait()
}
