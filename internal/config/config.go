package config

import (
	"fmt"
	"strings"
)

// Config holds the environment-driven settings needed to wire up the service.
type Config struct {
	KafkaBrokers              []string
	TranslationRequestedTopic string
	TranslationCompletedTopic string
	ConsumerGroupID           string
}

// Load reads the required settings using the given lookup function (typically os.LookupEnv).
// It returns an error naming the first missing required variable it finds.
func Load(lookup func(string) (string, bool)) (Config, error) {
	brokersRaw, err := require(lookup, "KAFKA_BROKERS")
	if err != nil {
		return Config{}, err
	}

	requestedTopic, err := require(lookup, "TRANSLATION_REQUESTED_TOPIC")
	if err != nil {
		return Config{}, err
	}

	completedTopic, err := require(lookup, "TRANSLATION_COMPLETED_TOPIC")
	if err != nil {
		return Config{}, err
	}

	groupID, err := require(lookup, "CONSUMER_GROUP_ID")
	if err != nil {
		return Config{}, err
	}

	return Config{
		KafkaBrokers:              splitBrokers(brokersRaw),
		TranslationRequestedTopic: requestedTopic,
		TranslationCompletedTopic: completedTopic,
		ConsumerGroupID:           groupID,
	}, nil
}

func require(lookup func(string) (string, bool), key string) (string, error) {
	v, ok := lookup(key)
	if !ok || v == "" {
		return "", fmt.Errorf("missing required environment variable %s", key)
	}
	return v, nil
}

// splitBrokers turns a comma-separated broker list into a clean slice,
// trimming whitespace and dropping empty entries.
func splitBrokers(raw string) []string {
	parts := strings.Split(raw, ",")
	brokers := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		brokers = append(brokers, p)
	}
	return brokers
}
