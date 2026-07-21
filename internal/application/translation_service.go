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

// ProcessTranslation translates the Name, Description and Attributes of a product for all requested languages.
func (s *TranslationService) ProcessTranslation(ctx context.Context, event domain.TranslationRequestedEvent) (domain.TranslationCompletedEvent, error) {
	productTranslations := domain.TranslationCompletedEvent{
		ProductID:    event.ProductID,
		Translations: make(map[string]domain.ProductTranslation),
	}

	// For each target language, we make one batch API call
	for _, targetLang := range event.TargetLanguages {
		// Prepare the array of texts to translate:
		// [0] = Name, [1] = Description, [2..N] = Attributes
		values := []string{event.Name, event.Description}

		// We need to keep track of attribute keys to map them back
		attrKeys := make([]string, 0, len(event.Attributes))
		for k, v := range event.Attributes {
			attrKeys = append(attrKeys, k)
			values = append(values, v)
		}

		// Call the translation port
		translatedValues, err := s.translationPort.Translate(ctx, values, targetLang)
		if err != nil {
			return productTranslations, fmt.Errorf("failed to translate to %s: %w", targetLang, err)
		}
		if len(translatedValues) != len(values) {
			return productTranslations, fmt.Errorf("translation mismatch for %s: expected %d, got %d", targetLang, len(values), len(translatedValues))
		}

		// Map back the results
		translation := domain.ProductTranslation{
			Name:        translatedValues[0],
			Description: translatedValues[1],
			Attributes:  make(map[string]string),
		}

		for i, k := range attrKeys {
			translation.Attributes[k] = translatedValues[2+i]
		}

		productTranslations.Translations[targetLang] = translation
	}

	return productTranslations, nil
}
