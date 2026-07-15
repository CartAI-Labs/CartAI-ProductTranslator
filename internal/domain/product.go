package domain

// TranslationRequestedEvent represents the payload received from Kafka
// when the backend requests a product translation.
type TranslationRequestedEvent struct {
	ProductID       string            `json:"productId"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Attributes      map[string]string `json:"attributes,omitempty"`
	TargetLanguages []string          `json:"targetLanguages"`
}

// ProductTranslation represents the translated fields for a specific language.
type ProductTranslation struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// TranslationCompletedEvent is the payload published back to Kafka.
type TranslationCompletedEvent struct {
	ProductID    string                        `json:"productId"`
	Translations map[string]ProductTranslation `json:"translations"`
}
