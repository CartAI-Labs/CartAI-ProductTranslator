package application_test

import (
	"context"
	"testing"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/application"
	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
)

// MockTranslationPort is a simple manual mock of the domain.TranslationPort interface
type MockTranslationPort struct {
	TranslateFunc func(ctx context.Context, values []string, targetLang string) ([]string, error)
}

func (m *MockTranslationPort) Translate(ctx context.Context, values []string, targetLang string) ([]string, error) {
	if m.TranslateFunc != nil {
		return m.TranslateFunc(ctx, values, targetLang)
	}
	return values, nil
}

func TestTranslateProduct(t *testing.T) {
	// Arrange
	event := domain.TranslationRequestedEvent{
		ProductID:       "PROD-123",
		Name:            "Portátil",
		Description:     "Muy rápido",
		Attributes:      map[string]string{"color": "Rojo"},
		TargetLanguages: []string{"en_US", "fr_FR"},
	}

	mockPort := &MockTranslationPort{
		TranslateFunc: func(ctx context.Context, values []string, targetLang string) ([]string, error) {
			// Dummy logic for mock
			result := make([]string, len(values))
			for i, value := range values {
				if targetLang == "en_US" {
					if value == "Portátil" {
						result[i] = "Laptop"
					} else if value == "Muy rápido" {
						result[i] = "Very fast"
					} else if value == "Rojo" {
						result[i] = "Red"
					}
				} else if targetLang == "fr_FR" {
					if value == "Portátil" {
						result[i] = "Ordinateur portable"
					} else if value == "Muy rápido" {
						result[i] = "Très rapide"
					} else if value == "Rojo" {
						result[i] = "Rouge"
					}
				}
			}
			return result, nil
		},
	}

	service := application.NewTranslationService(mockPort)

	// Act
	translatedEvent, err := service.ProcessTranslation(context.Background(), event)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	enTranslation, exists := translatedEvent.Translations["en_US"]
	if !exists {
		t.Fatalf("Expected 'en_US' translation to exist")
	}
	if enTranslation.Name != "Laptop" {
		t.Errorf("Expected Name 'Laptop', got '%s'", enTranslation.Name)
	}
	if enTranslation.Attributes["color"] != "Red" {
		t.Errorf("Expected color 'Red', got '%s'", enTranslation.Attributes["color"])
	}

	frTranslation, exists := translatedEvent.Translations["fr_FR"]
	if !exists {
		t.Fatalf("Expected 'fr_FR' translation to exist")
	}
	if frTranslation.Description != "Très rapide" {
		t.Errorf("Expected Description 'Très rapide', got '%s'", frTranslation.Description)
	}
}
