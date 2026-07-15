package domain

import "context"

// TranslationPort defines the interface for external translation services
// like Google Cloud Translation or DeepL.
type TranslationPort interface {
	Translate(ctx context.Context, values []string, targetLang string) ([]string, error)
}
