package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/application"
	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
	"github.com/segmentio/kafka-go"
)

// TranslationRequestedConsumer listens to a Kafka topic and delegates processing to the TranslationService.
type TranslationRequestedConsumer struct {
	reader    *kafka.Reader
	service   *application.TranslationService
	publisher domain.EventPublisherPort
}

// NewTranslationRequestedConsumer initializes a new Kafka consumer attached to a specific topic and group.
func NewTranslationRequestedConsumer(brokers []string, topic string, groupID string, service *application.TranslationService, publisher domain.EventPublisherPort) (*TranslationRequestedConsumer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("at least one broker is required")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &TranslationRequestedConsumer{
		reader:    reader,
		service:   service,
		publisher: publisher,
	}, nil
}

// Start begins a blocking loop that reads messages from Kafka.P
func (k *TranslationRequestedConsumer) Start(ctx context.Context) {
	log.Printf("Starting Kafka consumer on topic %s...", k.reader.Config().Topic)

	for {
		m, err := k.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				// Context was cancelled, gracefully shut down
				log.Println("Context cancelled, shutting down Kafka consumer")
				break
			}
			log.Printf("Error while reading message: %v", err)
			continue
		}

		// Decode the JSON message into our Domain Event
		var event domain.TranslationRequestedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event (offset %d): %v", m.Offset, err)
			continue
		}
		log.Printf("Received translation request for product %s (target languages: %v)", event.ProductID, event.TargetLanguages)

		// 🛑 INYECCIÓN DE DEPENDENCIAS: Llamamos a la Capa de Aplicación
		if k.service != nil {
			translation, err := k.service.ProcessTranslation(ctx, event)
			if err != nil {
				log.Printf("Failed to process translation for event %s: %v", event.ProductID, err)
				// Aquí en el futuro lo enviaremos al Topic de letras muertas (DLQ)
				continue
			}
			log.Printf("Successfully translated product %s to %d languages", translation.ProductID, len(translation.Translations))

			if k.publisher != nil {
				if err := k.publisher.PublishTranslation(ctx, translation); err != nil {
					log.Printf("Failed to publish translation completed event: %v", err)
				} else {
					log.Printf("Successfully published translation.completed for product %s", translation.ProductID)
				}
			}
		}
	}
}

// Close safely terminates the Kafka reader connection.
func (k *TranslationRequestedConsumer) Close() error {
	return k.reader.Close()
}
