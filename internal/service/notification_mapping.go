package service

import "sharetrip_notification/internal/domain"

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
