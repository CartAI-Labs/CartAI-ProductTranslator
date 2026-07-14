package application

import (
	"context"
	"fmt"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
)

// TranslationService orchestrates the translation process.
type TranslationService struct {
	translationPort domain.TranslationPort
}

// NewTranslationService creates a new instance of TranslationService with the given port.
func NewTranslationService(translationPort domain.TranslationPort) *TranslationService {
	return &TranslationService{
		translationPort: translationPort,
	}
}

// ProcessTranslation translates the Name and Description of a product event using the configured port.
func (s *TranslationService) ProcessTranslation(ctx context.Context, event domain.TranslationRequestedEvent) (domain.TranslationRequestedEvent, error) {
	translatedName, err := s.translationPort.TranslateText(ctx, event.Name, event.SourceLanguage, event.TargetLanguage)
	if err != nil {
		return event, fmt.Errorf("failed to translate product name: %w", err)
	}

	translatedDesc, err := s.translationPort.TranslateText(ctx, event.Description, event.SourceLanguage, event.TargetLanguage)
	if err != nil {
		return event, fmt.Errorf("failed to translate product description: %w", err)
	}

	// Create a copy of the event with translated fields
	translatedEvent := event
	translatedEvent.Name = translatedName
	translatedEvent.Description = translatedDesc
	translatedEvent.SourceLanguage = event.TargetLanguage // Now the source is the target

	return translatedEvent, nil
}
