package repository

import (
	"errors"

	"sharetrip_notification/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	postgresUniqueViolation     = "23505"
	postgresForeignKeyViolation = "23503"
	postgresNotNullViolation    = "23502"
	postgresCheckViolation      = "23514"
)

var (
	ErrInvalidReference    = errors.New("invalid reference data")
	ErrConstraintViolation = errors.New("constraint violation")
)

func mapPostgresError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.Join(service.ErrNotFound, err)
	}

	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return err
	}

	switch postgresError.Code {
	case postgresUniqueViolation:
		return errors.Join(service.ErrConflict, err)
	case postgresForeignKeyViolation:
		return errors.Join(ErrInvalidReference, err)
	case postgresNotNullViolation, postgresCheckViolation:
		return errors.Join(ErrConstraintViolation, err)
	default:
		return err
	}
}
