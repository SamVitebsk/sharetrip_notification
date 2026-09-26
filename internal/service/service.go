package service

import (
	"errors"
)

type Service struct {
	repo     NotificationRepository
	txRunner TxRunner
}

func New(repo NotificationRepository, txRunner TxRunner) (*Service, error) {
	if repo == nil {
		return nil, errors.New("notification repository is required")
	}
	if txRunner == nil {
		return nil, errors.New("txRunner is required")
	}
	return &Service{repo: repo, txRunner: txRunner}, nil
}
