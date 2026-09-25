package domain

import "errors"

var (
	ErrEmptyRecipient = errors.New("empty recipient id")
	ErrEmptyType      = errors.New("empty notification type")
	ErrInvalidTime    = errors.New("invalid time")
)
