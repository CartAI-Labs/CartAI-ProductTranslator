package domain

// TranslationRequestedEvent represents the payload received from Kafka
// when the backend requests a product translation.
type TranslationRequestedEvent struct {
	ProductID      string `json:"productId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
}
