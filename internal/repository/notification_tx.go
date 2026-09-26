package repository

import (
	"context"
	"fmt"
	"sharetrip_notification/internal/service"

	"github.com/jackc/pgx/v5"
)

type NotificationTx struct {
	tx pgx.Tx
}

func (p *Postgres) RunInTx(ctx context.Context, operation func(context.Context, service.RepositoryTx) error) error {
	if operation == nil {
		return fmt.Errorf("операция транзакции обязательна")
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ошибка старта транзакции: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	wrapper := &NotificationTx{tx: tx}

	if err := operation(ctx, wrapper); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %w", err)
	}
	return nil
}
