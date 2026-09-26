package kafka

import (
	"sharetrip_notification/internal/env"
)

type Config struct {
	Brokers []string
	GroupID string
	Topic   string
}

func LoadConfig() Config {
	return Config{
		Brokers: env.StringSlice("KAFKA_BROKERS", []string{"localhost:29092"}),
		GroupID: env.String("KAFKA_GROUP_ID", "notification-service"),
		Topic:   env.String("KAFKA_TOPIC", "trip.events"),
	}
}
