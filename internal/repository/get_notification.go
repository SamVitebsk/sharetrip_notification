package repository

import (
	"context"
	"fmt"

	"sharetrip_notification/internal/domain"
	"sharetrip_notification/internal/repository/entity"

	"github.com/google/uuid"
)

func (p *Postgres) FindByID(ctx context.Context, id uuid.UUID) (domain.Notification, error) {
	var ent entity.Notification

	err := p.pool.QueryRow(ctx, `
		SELECT id, recipient_id, type, status, payload, created_at
		FROM notifications
		WHERE id = $1
	`, id).Scan(
		&ent.ID,
		&ent.RecipientID,
		&ent.Type,
		&ent.Status,
		&ent.Payload,
		&ent.CreatedAt,
	)

	if err != nil {
		return domain.Notification{}, fmt.Errorf("select notification: %w", mapPostgresError(err))
	}

	return toDomainNotification(ent), nil
}
