package repository

import (
	"context"
	"fmt"

	"sharetrip_notification/internal/domain"
)

func (p *Postgres) Save(ctx context.Context, n domain.Notification) error {
	entity := toEntityNotification(n)

	_, err := p.pool.Exec(ctx, `
		INSERT INTO notifications (
			id, recipient_id, type, status, payload, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		entity.ID,
		entity.RecipientID,
		entity.Type,
		entity.Status,
		entity.Payload,
		entity.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert notification: %w", mapPostgresError(err))
	}

	return nil
}
