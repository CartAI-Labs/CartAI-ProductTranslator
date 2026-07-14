package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
)

func TestParseTranslationRequestedEvent(t *testing.T) {
	// Arrange
	jsonPayload := `{
		"productId": "PROD-123",
		"name": "Ordenador Portátil",
		"description": "Un ordenador muy potente",
		"sourceLanguage": "es",
		"targetLanguage": "en"
	}`

	// Act
	var event domain.TranslationRequestedEvent
	err := json.Unmarshal([]byte(jsonPayload), &event)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if event.ProductID != "PROD-123" {
		t.Errorf("Expected ProductID 'PROD-123', got '%s'", event.ProductID)
	}
	if event.Name != "Ordenador Portátil" {
		t.Errorf("Expected Name 'Ordenador Portátil', got '%s'", event.Name)
	}
}
