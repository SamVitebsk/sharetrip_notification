package service

import (
	"context"
	"log"
	"time"

	"sharetrip_notification/internal/clients/kafka"
	"sharetrip_notification/internal/domain"

	"github.com/google/uuid"
)

func (s *Service) ProcessTripPublished(ctx context.Context, event kafka.TripPublished) error {
	eventID, err := uuid.Parse(event.EventID)
	if err != nil {
		log.Printf("некорректный event_id (poison pill): %v", err)
		return nil
	}

	eventTime := event.OccurredAt
	if eventTime.IsZero() {
		eventTime = time.Now().UTC()
	}
	eventTime = eventTime.Truncate(time.Microsecond)

	domainReq := toCreateNotificationRequestFromEvent(event, eventTime)

	response, err := domain.CreateNotification(domainReq)
	if err != nil {
		log.Printf("неверные данные для уведомления (poison pill): %v", err)
		return nil
	}

	return s.txRunner(ctx, func(ctx context.Context, tx RepositoryTx) error {
		inserted, err := tx.MarkEventProcessed(ctx, eventID)
		if err != nil {
			return err
		}

		if !inserted {
			return nil
		}

		return tx.Save(ctx, response.Notification)
	})
}
