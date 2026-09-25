package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreateNotificationRequest struct {
	RecipientID string
	Type        string
	Payload     map[string]any
	Now         time.Time
}

type CreateNotificationResponse struct {
	Notification Notification
}

func CreateNotification(req CreateNotificationRequest) (CreateNotificationResponse, error) {
	if req.RecipientID == "" {
		return CreateNotificationResponse{}, ErrEmptyRecipient
	}
	if req.Type == "" {
		return CreateNotificationResponse{}, ErrEmptyType
	}
	if req.Now.IsZero() {
		return CreateNotificationResponse{}, ErrInvalidTime
	}

	notif := Notification{
		ID:          uuid.New(),
		RecipientID: req.RecipientID,
		Type:        NotificationType(req.Type),
		Status:      StatusCreated,
		Payload:     deepCopyMap(req.Payload),
		CreatedAt:   req.Now,
	}

	return CreateNotificationResponse{Notification: notif}, nil
}

func deepCopyMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = deepCopyValue(v)
	}
	return dst
}

func deepCopyValue(val any) any {
	switch v := val.(type) {
	case map[string]any:
		return deepCopyMap(v)
	case []any:
		dst := make([]any, len(v))
		for i, item := range v {
			dst[i] = deepCopyValue(item)
		}
		return dst
	default:
		return v
	}
}
