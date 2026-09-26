package service

import (
	"sharetrip_notification/internal/clients/kafka"
	"sharetrip_notification/internal/domain"
	"time"
)

func toCreateNotificationResponse(n domain.Notification) CreateNotificationResponse {
	return CreateNotificationResponse{
		ID:          n.ID,
		RecipientID: n.RecipientID,
		Type:        string(n.Type),
		Status:      string(n.Status),
		Payload:     n.Payload,
		CreatedAt:   n.CreatedAt,
	}
}

func toGetNotificationResponse(n domain.Notification) GetNotificationResponse {
	return GetNotificationResponse{
		ID:          n.ID,
		RecipientID: n.RecipientID,
		Type:        string(n.Type),
		Status:      string(n.Status),
		Payload:     n.Payload,
		CreatedAt:   n.CreatedAt,
	}
}

func toCreateNotificationRequestFromEvent(event kafka.TripPublished, now time.Time) domain.CreateNotificationRequest {
	return domain.CreateNotificationRequest{
		RecipientID: event.DriverID,
		Type:        "email",
		Payload:     map[string]any{"trip_id": event.TripID},
		Now:         now,
	}
}
