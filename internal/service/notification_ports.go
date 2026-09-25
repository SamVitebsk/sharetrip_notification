package service

import (
	"context"
	"sharetrip_notification/internal/domain"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	Save(ctx context.Context, n domain.Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (domain.Notification, error)
}
