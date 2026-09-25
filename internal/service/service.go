package service

import (
	"errors"
)

type Service struct {
	repo NotificationRepository
}

func New(repo NotificationRepository) (*Service, error) {
	if repo == nil {
		return nil, errors.New("notification repository is required")
	}
	return &Service{repo: repo}, nil
}
