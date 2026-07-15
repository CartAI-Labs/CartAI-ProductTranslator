package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
	"github.com/segmentio/kafka-go"
)

// KafkaPublisher is the concrete implementation of domain.EventPublisherPort
type KafkaPublisher struct {
	writer *kafka.Writer
}

// NewKafkaPublisher initializes a new Kafka writer attached to a specific topic.
func NewKafkaPublisher(brokers []string, topic string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("at least one broker is required")
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		BatchTimeout:           10 * time.Millisecond,
		AllowAutoTopicCreation: true,
	}

	return &KafkaPublisher{
		writer: writer,
	}, nil
}

// PublishTranslationCompleted serializes the completed event to JSON and pushes it to Kafka.
func (p *KafkaPublisher) PublishTranslation(ctx context.Context, event domain.TranslationCompletedEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal TranslationCompletedEvent: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.ProductID), // Ensure messages for the same product go to the same partition
		Value: bytes,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

// Close gracefully closes the Kafka writer connection.
func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
