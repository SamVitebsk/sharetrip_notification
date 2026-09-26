package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GetNotificationResponse struct {
	ID          uuid.UUID
	RecipientID string
	Type        string
	Status      string
	Payload     map[string]any
	CreatedAt   time.Time
}

func (s *Service) GetNotification(ctx context.Context, id uuid.UUID) (GetNotificationResponse, error) {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return GetNotificationResponse{}, fmt.Errorf("ошибка поиска уведомления: %w", err)
	}

	return toGetNotificationResponse(n), nil
}
