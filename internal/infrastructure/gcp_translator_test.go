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

	translatedText, err := translator.TranslateText(context.Background(), "Hola mundo", "es", "en")
	if err != nil {
		t.Fatalf("Failed to translate text: %v", err)
	}

	if translatedText != "Hello world" {
		t.Errorf("Expected 'Hello world', got '%s'", translatedText)
	}
}
