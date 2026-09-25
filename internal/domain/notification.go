package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationStatus string
type NotificationType string

const (
	StatusCreated NotificationStatus = "created"
)

type Notification struct {
	ID          uuid.UUID
	RecipientID string
	Type        NotificationType
	Status      NotificationStatus
	Payload     map[string]any
	CreatedAt   time.Time
}
