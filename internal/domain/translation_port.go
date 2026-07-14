package domain

import "context"

// TranslationPort defines the interface for external translation services
// like Google Cloud Translation or DeepL.
type TranslationPort interface {
	TranslateText(ctx context.Context, text string, sourceLang string, targetLang string) (string, error)
}
