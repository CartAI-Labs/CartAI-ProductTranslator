# Cart·AI Product Translator

<img src="docs/assets/logoCartAI.png" alt="Cart·AI Shopping Logo" width="120" />

**Cart·AI Product Translator** is a high-performance, asynchronous microservice written in **Go (Golang)**. It is a core component of the broader **Cart·AI Ecosystem**, designed to work in tandem with the primary Java Spring Boot backend and the AI-driven Python worker.

## 🔗 Ecosystem Integration

This service operates as a decoupled enrichment step in the Cart·AI event-driven architecture:
1. The **Cart·AI Backend** publishes a `translation.requested` event to Apache Kafka when a product requires internationalization.
2. This **Product Translator** service consumes the event asynchronously.
3. It performs the translation using the Google Cloud Translation API.
4. Finally, it publishes a `translation.completed` event back to Kafka, which the main backend consumes to update the catalog.

## 🏛️ System Architecture

This project strictly follows **Hexagonal Architecture (Ports & Adapters)** adapted for Go idioms, ensuring that the business logic remains fully agnostic to external frameworks and message brokers.

- `cmd/`: Application entry points and dependency injection (Composition Root).
- `internal/domain/`: Core business models (Structs) and interfaces (Ports).
- `internal/application/`: Business use cases orchestrating the translation flow.
- `internal/infrastructure/`: Secondary adapters (Kafka Consumers/Producers, GCP Translation client).

## 🛠️ Technology Stack

- **Core:** Go 1.26+
- **Messaging:** Apache Kafka (`github.com/segmentio/kafka-go`)
- **Cloud:** Google Cloud Translation API (`cloud.google.com/go/translate`)
- **Methodology:** Strict TDD (Test-Driven Development) using the native Go `testing` package.

## 🚀 Development

To run the unit tests (following the Red-Green-Refactor TDD loop):
```bash
go test ./... -v
```

To compile the standalone binary:
```bash
go build -o bin/translator cmd/translator/main.go
```

---

## 👨‍💻 Author

**Roberto Díaz**  
*Senior Software Engineer & Architect*

## 📄 License

Copyright (C) 2026 Roberto Díaz. All rights reserved.  
Licensed under the GNU General Public License v3.0. See LICENSE for details.
