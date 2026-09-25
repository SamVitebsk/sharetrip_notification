package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"time"

	"sharetrip_notification/internal/api"
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

	port, err := env.Int("PGPORT", 5432)
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return errors.New("PGPORT должен быть в диапазоне от 1 до 65535")
	}

	cfg := repository.Config{
		Host:     env.String("PGHOST", "localhost"),
		Port:     port,
		User:     env.String("PGUSER", "postgres"),
		Password: env.String("PGPASSWORD", "postgres"),
		DBName:   env.String("PGDATABASE", "notification_db"),
		SSLMode:  env.String("PGSSLMODE", "disable"),
	}

	pool, err := repository.ConnectPostgres(ctx, cfg.DSN())
	if err != nil {
		return errors.New("не удалось подключиться к PostgreSQL: проверьте PG*-настройки и доступность БД")
	}
	defer pool.Close()

	pg, err := repository.NewPostgres(pool)
	if err != nil {
		return err
	}

	svc, err := service.New(pg)
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

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting notification service on :%d", serverPort)
		serverErrors <- app.Listen(fmt.Sprintf(":%d", serverPort))
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		log.Println("Shutting down gracefully...")
		return app.ShutdownWithTimeout(10 * time.Second)
	}
}
