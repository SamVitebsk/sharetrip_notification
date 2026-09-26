package api

import (
	"errors"
	"log"

	"sharetrip_notification/internal/api/openapi"
	"sharetrip_notification/internal/service"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "internal server error"

	var fiberErr *fiber.Error

	switch {
	case errors.Is(err, service.ErrInvalidInput):
		status, message = fiber.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrNotFound):
		status, message = fiber.StatusNotFound, service.ErrNotFound.Error()
	case errors.Is(err, service.ErrConflict):
		status, message = fiber.StatusConflict, service.ErrConflict.Error()
	case errors.As(err, &fiberErr):
		status, message = fiberErr.Code, fiberErr.Message
	}

	if status == fiber.StatusInternalServerError {
		log.Printf("Внутренняя ошибка (HTTP 500): метод=%s маршрут=%q ошибка=%v", c.Method(), c.Route().Path, err)
	}

	return c.Status(status).JSON(openapi.ErrorResponse{
		Error: message,
	})
}
