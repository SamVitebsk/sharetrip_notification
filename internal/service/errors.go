package service

import (
	"errors"

	"sharetrip_notification/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("notification not found")
	ErrConflict     = errors.New("conflict operation")
)

func mapCreateError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmptyRecipient),
		errors.Is(err, domain.ErrEmptyType),
		errors.Is(err, domain.ErrInvalidTime):
		return errors.Join(ErrInvalidInput, err)
	default:
		return err
	}
}
