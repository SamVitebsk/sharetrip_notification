package service

import (
	"context"
	"sharetrip_notification/internal/domain"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (domain.Notification, error)
}

type RepositoryTx interface {
	Save(ctx context.Context, n domain.Notification) error
	MarkEventProcessed(ctx context.Context, eventID uuid.UUID) (bool, error)
}

type TxRunner func(ctx context.Context, operation func(context.Context, RepositoryTx) error) error
