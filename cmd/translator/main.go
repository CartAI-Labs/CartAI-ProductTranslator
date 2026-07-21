package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/application"
	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/config"
	"github.com/CartAI-Labs/CartAI-ProductTranslator/internal/infrastructure"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(os.LookupEnv)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	translator, err := infrastructure.NewGCPTranslator()
	if err != nil {
		return fmt.Errorf("failed to create GCP translator: %w", err)
	}
	defer translator.Close()

	service := application.NewTranslationService(translator)

	publisher, err := infrastructure.NewTranslationCompletedPublisher(cfg.KafkaBrokers, cfg.TranslationCompletedTopic)
	if err != nil {
		return fmt.Errorf("failed to create translation completed publisher: %w", err)
	}
	defer publisher.Close()

	consumer, err := infrastructure.NewTranslationRequestedConsumer(
		cfg.KafkaBrokers,
		cfg.TranslationRequestedTopic,
		cfg.ConsumerGroupID,
		service,
		publisher,
	)
	if err != nil {
		return fmt.Errorf("failed to create translation requested consumer: %w", err)
	}
	defer consumer.Close()

	consumer.Start(ctx)

	return nil
}
