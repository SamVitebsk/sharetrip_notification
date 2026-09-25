package api

import (
	"sharetrip_notification/internal/api/openapi"
	"sharetrip_notification/internal/service"
)

func toCreateNotificationRequest(request openapi.CreateNotificationRequest) service.CreateNotificationRequest {
	return service.CreateNotificationRequest{
		RecipientID: request.RecipientId,
		Type:        request.Type,
		Payload:     request.Payload,
	}
}

func toCreateNotificationResponse(response service.CreateNotificationResponse) openapi.NotificationResponse {
	return openapi.NotificationResponse{
		Id:          response.ID,
		RecipientId: response.RecipientID,
		Type:        response.Type,
		Status:      response.Status,
		Payload:     response.Payload,
		CreatedAt:   response.CreatedAt,
	}
}

func toGetNotificationResponse(response service.GetNotificationResponse) openapi.NotificationResponse {
	return openapi.NotificationResponse{
		Id:          response.ID,
		RecipientId: response.RecipientID,
		Type:        response.Type,
		Status:      response.Status,
		Payload:     response.Payload,
		CreatedAt:   response.CreatedAt,
	}
}
