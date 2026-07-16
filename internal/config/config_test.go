package config_test

import (
	"testing"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/config"
)

func fakeLookup(env map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

func TestLoad_Success(t *testing.T) {
	env := map[string]string{
		"KAFKA_BROKERS":               "broker1:9092,broker2:9092",
		"TRANSLATION_REQUESTED_TOPIC": "translation.requested",
		"TRANSLATION_COMPLETED_TOPIC": "translation.completed",
		"CONSUMER_GROUP_ID":           "translation-group",
	}

	cfg, err := config.Load(fakeLookup(env))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	wantBrokers := []string{"broker1:9092", "broker2:9092"}
	if len(cfg.KafkaBrokers) != len(wantBrokers) {
		t.Fatalf("Expected %d brokers, got %d", len(wantBrokers), len(cfg.KafkaBrokers))
	}
	for i, b := range wantBrokers {
		if cfg.KafkaBrokers[i] != b {
			t.Errorf("Expected broker[%d] = %s, got %s", i, b, cfg.KafkaBrokers[i])
		}
	}
	if cfg.TranslationRequestedTopic != "translation.requested" {
		t.Errorf("Expected TranslationRequestedTopic 'translation.requested', got '%s'", cfg.TranslationRequestedTopic)
	}
	if cfg.TranslationCompletedTopic != "translation.completed" {
		t.Errorf("Expected TranslationCompletedTopic 'translation.completed', got '%s'", cfg.TranslationCompletedTopic)
	}
	if cfg.ConsumerGroupID != "translation-group" {
		t.Errorf("Expected ConsumerGroupID 'translation-group', got '%s'", cfg.ConsumerGroupID)
	}
}

func TestLoad_MissingRequiredVar(t *testing.T) {
	tests := []struct {
		name    string
		missing string
	}{
		{"missing brokers", "KAFKA_BROKERS"},
		{"missing requested topic", "TRANSLATION_REQUESTED_TOPIC"},
		{"missing completed topic", "TRANSLATION_COMPLETED_TOPIC"},
		{"missing group id", "CONSUMER_GROUP_ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := map[string]string{
				"KAFKA_BROKERS":               "broker1:9092",
				"TRANSLATION_REQUESTED_TOPIC": "translation.requested",
				"TRANSLATION_COMPLETED_TOPIC": "translation.completed",
				"CONSUMER_GROUP_ID":           "translation-group",
			}
			delete(env, tt.missing)

			_, err := config.Load(fakeLookup(env))
			if err == nil {
				t.Fatalf("Expected error for missing %s, got nil", tt.missing)
			}
		})
	}
}

func TestLoad_BlankBrokersAreIgnored(t *testing.T) {
	env := map[string]string{
		"KAFKA_BROKERS":               "broker1:9092, ,broker2:9092,",
		"TRANSLATION_REQUESTED_TOPIC": "translation.requested",
		"TRANSLATION_COMPLETED_TOPIC": "translation.completed",
		"CONSUMER_GROUP_ID":           "translation-group",
	}

	cfg, err := config.Load(fakeLookup(env))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	want := []string{"broker1:9092", "broker2:9092"}
	if len(cfg.KafkaBrokers) != len(want) {
		t.Fatalf("Expected %d brokers, got %d (%v)", len(want), len(cfg.KafkaBrokers), cfg.KafkaBrokers)
	}
}
