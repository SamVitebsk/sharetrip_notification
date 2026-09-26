package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type TripPublished struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	TripID     string    `json:"trip_id"`
	DriverID   string    `json:"driver_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type Consumer struct {
	reader *kafkago.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) Listen(ctx context.Context, handler func(context.Context, TripPublished) error) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}

		var event TripPublished
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("неверный формат сообщения (poison pill), пропускаем офсет %d: %v", msg.Offset, err)
			if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
				return commitErr
			}
			continue
		}

		if event.EventType != "TripPublished" {
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				return err
			}
			continue
		}

		for {
			err := handler(ctx, event)
			if err == nil {
				break
			}

			log.Printf("системная ошибка обработки бизнес-события %s: %v, повтор через 5 секунд...", event.EventID, err)

			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
