package repository

import (
	"fmt"
	"os"
)

type Config struct {
	DSN string
}

func LoadConfig() (Config, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return Config{}, fmt.Errorf("DATABASE_DSN is required")
	}

	return Config{
		DSN: dsn,
	}, nil
}
