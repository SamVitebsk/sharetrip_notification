package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID          uuid.UUID
	RecipientID string
	Type        string
	Status      string
	Payload     map[string]any
	CreatedAt   time.Time
}
