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

// TranslateTexts calls the real Google Cloud Translation API for an array of texts.
func (g *GCPTranslator) TranslateTexts(ctx context.Context, texts []string, targetLang string) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	target, err := language.Parse(targetLang)
	if err != nil {
		return nil, fmt.Errorf("invalid target language: %w", err)
	}

	// Passing nil for options enables Google's auto-detect for source language
	translations, err := g.client.Translate(ctx, texts, target, nil)
	if err != nil {
		return nil, fmt.Errorf("google translation failed: %w", err)
	}

	if len(translations) != len(texts) {
		return nil, fmt.Errorf("expected %d translations, got %d", len(texts), len(translations))
	}

	result := make([]string, len(texts))
	for i, t := range translations {
		result[i] = t.Text
	}

	return result, nil
}

// Close closes the underlying GCP client.
func (g *GCPTranslator) Close() error {
	return g.client.Close()
}
