package infrastructure

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// GCPTranslator is the concrete implementation of domain.TranslationPort
type GCPTranslator struct {
	client *translate.Client
}

// NewGCPTranslator initializes a new GCP Translation client.
// It automatically picks up the GOOGLE_APPLICATION_CREDENTIALS environment variable.
func NewGCPTranslator() (*GCPTranslator, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create translate client: %w", err)
	}

	return &GCPTranslator{
		client: client,
	}, nil
}

// TranslateText calls the real Google Cloud Translation API.
func (g *GCPTranslator) TranslateText(ctx context.Context, text string, sourceLang string, targetLang string) (string, error) {
	target, err := language.Parse(targetLang)
	if err != nil {
		return "", fmt.Errorf("invalid target language: %w", err)
	}

	source, err := language.Parse(sourceLang)
	if err != nil {
		return "", fmt.Errorf("invalid source language: %w", err)
	}

	translations, err := g.client.Translate(ctx, []string{text}, target, &translate.Options{
		Source: source,
	})
	if err != nil {
		return "", fmt.Errorf("google translation failed: %w", err)
	}

	if len(translations) == 0 {
		return "", fmt.Errorf("no translations returned from google")
	}

	return translations[0].Text, nil
}

// Close closes the underlying GCP client.
func (g *GCPTranslator) Close() error {
	return g.client.Close()
}
