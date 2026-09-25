package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"sharetrip_notification/internal/api/openapi"
	"sharetrip_notification/internal/service"
)

type Server struct {
	service *service.Service
}

var _ openapi.ServerInterface = (*Server)(nil)

func NewServer(svc *service.Service) (*Server, error) {
	if svc == nil {
		return nil, errors.New("notification service is required")
	}
	return &Server{service: svc}, nil
}

func RegisterRoutes(router fiber.Router, server *Server) {
	openapi.RegisterHandlers(router, server)
}
