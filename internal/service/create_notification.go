package service

import (
	"context"
	"fmt"
	"time"

	"sharetrip_notification/internal/domain"

	"github.com/google/uuid"
)

type CreateNotificationRequest struct {
	RecipientID string
	Type        string
	Payload     map[string]any
}

type CreateNotificationResponse struct {
	ID          uuid.UUID
	RecipientID string
	Type        string
	Status      string
	Payload     map[string]any
	CreatedAt   time.Time
}

func (s *Service) Create(ctx context.Context, request CreateNotificationRequest) (CreateNotificationResponse, error) {
	domainReq := domain.CreateNotificationRequest{
		RecipientID: request.RecipientID,
		Type:        request.Type,
		Payload:     request.Payload,
		Now:         time.Now().UTC().Truncate(time.Microsecond),
	}

	response, err := domain.CreateNotification(domainReq)
	if err != nil {
		return CreateNotificationResponse{}, mapCreateError(err)
	}

	err = s.txRunner(ctx, func(ctx context.Context, tx RepositoryTx) error {
		return tx.Save(ctx, response.Notification)
	})
	if err != nil {
		return CreateNotificationResponse{}, fmt.Errorf("ошибка сохранения уведомления: %w", err)
	}

	return toCreateNotificationResponse(response.Notification), nil
}
