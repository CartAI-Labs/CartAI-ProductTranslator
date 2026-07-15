package infrastructure_test

import (
	"context"
	"os"
	"testing"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/infrastructure"
)

func TestGCPTranslator_TranslateText(t *testing.T) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("Skipping GCP integration test because GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	translator, err := infrastructure.NewGCPTranslator()
	if err != nil {
		t.Fatalf("Failed to initialize GCP translator: %v", err)
	}
	defer translator.Close()

	translatedTexts, err := translator.Translate(context.Background(), []string{"Hola mundo", "Rojo"}, "en_US")
	if err != nil {
		t.Fatalf("Failed to translate texts: %v", err)
	}

	if len(translatedTexts) != 2 {
		t.Fatalf("Expected 2 translations, got %d", len(translatedTexts))
	}

	if translatedTexts[0] != "Hello World" && translatedTexts[0] != "Hello world" {
		t.Errorf("Expected 'Hello world', got '%s'", translatedTexts[0])
	}
	if translatedTexts[1] != "Red" && translatedTexts[1] != "red" {
		t.Errorf("Expected 'Red', got '%s'", translatedTexts[1])
	}
}
