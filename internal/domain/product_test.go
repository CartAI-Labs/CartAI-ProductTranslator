package domain_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/domain"
)

func TestParseTranslationRequestedEvent(t *testing.T) {
	// Arrange
	jsonPayload := `{
		"productId": "PROD-123",
		"name": "Ordenador Portátil",
		"description": "Un ordenador muy potente",
		"attributes": {"color": "Rojo"},
		"targetLanguages": ["en_US", "fr_FR"]
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
	if event.Description != "Un ordenador muy potente" {
		t.Errorf("Expected Description 'Un ordenador muy potente', got '%s'", event.Description)
	}
	if event.Attributes["color"] != "Rojo" {
		t.Errorf("Expected Attributes[\"color\"] 'Rojo', got '%s'", event.Attributes["color"])
	}
	if len(event.TargetLanguages) != 2 {
		t.Fatalf("Expected 2 target languages, got %d", len(event.TargetLanguages))
	}
	if event.TargetLanguages[0] != "en_US" || event.TargetLanguages[1] != "fr_FR" {
		t.Errorf("Expected TargetLanguages [en_US fr_FR], got %v", event.TargetLanguages)
	}
}

func TestMarshalTranslationRequestedEvent_OmitsEmptyAttributes(t *testing.T) {
	// Arrange: no Attributes set, so it's nil.
	event := domain.TranslationRequestedEvent{
		ProductID:       "PROD-123",
		Name:            "Ordenador Portátil",
		Description:     "Un ordenador muy potente",
		TargetLanguages: []string{"en_US"},
	}

	// Act
	bytes, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Assert
	if strings.Contains(string(bytes), "attributes") {
		t.Errorf("Expected 'attributes' key to be omitted when empty, got %s", string(bytes))
	}
}

func TestMarshalProductTranslation_OmitsEmptyAttributes(t *testing.T) {
	// Arrange: no Attributes set, so it's nil.
	translation := domain.ProductTranslation{
		Name:        "Laptop",
		Description: "Very fast",
	}

	// Act
	bytes, err := json.Marshal(translation)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Assert
	if strings.Contains(string(bytes), "attributes") {
		t.Errorf("Expected 'attributes' key to be omitted when empty, got %s", string(bytes))
	}
}
