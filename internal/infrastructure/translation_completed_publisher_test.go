package infrastructure_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/infrastructure"
	"github.com/segmentio/kafka-go"
)

func TestTranslationCompletedPublisher_PublishTranslation(t *testing.T) {
	if os.Getenv("RUN_KAFKA_INTEGRATION_TESTS") == "" {
		t.Skip("Skipping Kafka integration test. Set RUN_KAFKA_INTEGRATION_TESTS=1 to run.")
	}

	topic := "translation.completed"

	publisher, err := infrastructure.NewTranslationCompletedPublisher([]string{"localhost:9092"}, topic)
	if err != nil {
		t.Fatalf("Failed to initialize Kafka publisher: %v", err)
	}
	defer publisher.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		GroupID:  "translation-completed-publisher-test",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	event := domain.TranslationCompletedEvent{
		ProductID: "PROD-123",
		Translations: map[string]domain.ProductTranslation{
			"en_US": {
				Name:        "Laptop",
				Description: "Very fast",
				Attributes:  map[string]string{"color": "Red"},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := publisher.PublishTranslation(ctx, event); err != nil {
		t.Fatalf("Failed to publish translation completed event: %v", err)
	}

	m, err := reader.ReadMessage(ctx)
	if err != nil {
		t.Fatalf("Failed to read back published message: %v", err)
	}

	if string(m.Key) != event.ProductID {
		t.Errorf("Expected message key '%s', got '%s'", event.ProductID, string(m.Key))
	}

	var got domain.TranslationCompletedEvent
	if err := json.Unmarshal(m.Value, &got); err != nil {
		t.Fatalf("Failed to unmarshal published message: %v", err)
	}

	if got.ProductID != event.ProductID {
		t.Errorf("Expected ProductID '%s', got '%s'", event.ProductID, got.ProductID)
	}
	if got.Translations["en_US"].Name != "Laptop" {
		t.Errorf("Expected Name 'Laptop', got '%s'", got.Translations["en_US"].Name)
	}
}
