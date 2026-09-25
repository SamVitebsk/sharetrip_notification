package api

import (
	"sharetrip_notification/internal/api/openapi"
	"sharetrip_notification/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) GetNotificationsId(c *fiber.Ctx, id openapi.NotificationId) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return service.ErrInvalidInput
	}

	response, err := s.service.GetNotification(c.Context(), parsedUUID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(toGetNotificationResponse(response))
}
