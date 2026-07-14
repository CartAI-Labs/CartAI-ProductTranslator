package application_test

import (
	"context"
	"testing"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/application"
	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
)

// MockTranslationPort is a simple manual mock of the domain.TranslationPort interface
type MockTranslationPort struct {
	TranslateTextFunc func(ctx context.Context, text string, sourceLang string, targetLang string) (string, error)
}

func (m *MockTranslationPort) TranslateText(ctx context.Context, text string, sourceLang string, targetLang string) (string, error) {
	if m.TranslateTextFunc != nil {
		return m.TranslateTextFunc(ctx, text, sourceLang, targetLang)
	}
	return text, nil
}

func TestTranslateProduct(t *testing.T) {
	// Arrange
	event := domain.TranslationRequestedEvent{
		ProductID:      "PROD-123",
		Name:           "Portátil",
		Description:    "Muy rápido",
		SourceLanguage: "es",
		TargetLanguage: "en",
	}

	mockPort := &MockTranslationPort{
		TranslateTextFunc: func(ctx context.Context, text string, sourceLang string, targetLang string) (string, error) {
			if text == "Portátil" {
				return "Laptop", nil
			}
			if text == "Muy rápido" {
				return "Very fast", nil
			}
			return text, nil
		},
	}

	service := application.NewTranslationService(mockPort)

	// Act
	translatedEvent, err := service.ProcessTranslation(context.Background(), event)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if translatedEvent.Name != "Laptop" {
		t.Errorf("Expected Name 'Laptop', got '%s'", translatedEvent.Name)
	}
	if translatedEvent.Description != "Very fast" {
		t.Errorf("Expected Description 'Very fast', got '%s'", translatedEvent.Description)
	}
}
