package api

import (
	"sharetrip_notification/internal/api/openapi"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) PostNotifications(c *fiber.Ctx) error {
	var request openapi.CreateNotificationRequest
	if err := c.BodyParser(&request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := s.service.Create(c.Context(), toCreateNotificationRequest(request))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(toCreateNotificationResponse(response))
}
