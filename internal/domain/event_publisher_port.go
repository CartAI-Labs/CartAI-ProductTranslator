package domain

import "context"

// EventPublisherPort defines the interface to emit events to the external message broker.
type EventPublisherPort interface {
	PublishTranslation(ctx context.Context, event TranslationCompletedEvent) error
}
