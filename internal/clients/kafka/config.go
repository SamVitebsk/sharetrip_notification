package kafka

import (
	"fmt"
	"os"
)

type Config struct {
	Brokers []string
	GroupID string
	Topic   string
}

func LoadConfig() (Config, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return Config{}, fmt.Errorf("KAFKA_BROKERS is required")
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		return Config{}, fmt.Errorf("KAFKA_GROUP_ID is required")
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "trip-events"
	}

	return Config{
		Brokers: []string{brokers},
		GroupID: groupID,
		Topic:   topic,
	}, nil
}
