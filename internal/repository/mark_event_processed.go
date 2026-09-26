package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (w *NotificationTx) MarkEventProcessed(ctx context.Context, eventID uuid.UUID) (bool, error) {
	tag, err := w.tx.Exec(ctx, `
		INSERT INTO processed_events (event_id)
		VALUES ($1)
		ON CONFLICT (event_id) DO NOTHING
	`, eventID)

	if err != nil {
		return false, fmt.Errorf("ошибка записи обработанного события: %w", mapPostgresError(err))
	}

	inserted := tag.RowsAffected() > 0
	return inserted, nil
}
