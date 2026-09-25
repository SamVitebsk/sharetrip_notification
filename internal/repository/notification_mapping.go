package repository

import (
	"sharetrip_notification/internal/domain"
	"sharetrip_notification/internal/repository/entity"
)

func toDomainNotification(e entity.Notification) domain.Notification {
	return domain.Notification{
		ID:          e.ID,
		RecipientID: e.RecipientID,
		Type:        domain.NotificationType(e.Type),
		Status:      domain.NotificationStatus(e.Status),
		Payload:     e.Payload,
		CreatedAt:   e.CreatedAt,
	}
}

func toEntityNotification(d domain.Notification) entity.Notification {
	return entity.Notification{
		ID:          d.ID,
		RecipientID: d.RecipientID,
		Type:        string(d.Type),
		Status:      string(d.Status),
		Payload:     d.Payload,
		CreatedAt:   d.CreatedAt,
	}
}
