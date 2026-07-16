# Project Context: cart-ai-product-translator

## 1. Definition and Purpose
**`cart-ai-product-translator`** is a high-performance, asynchronous microservice written in Go. Its sole responsibility is to consume translation requests from Kafka, translate textual data using the Google Cloud Translation API, and publish the enriched payload back to Kafka.

## 2. Strict Technology Stack
* **Core:** Go 1.26+
* **Dependency Management:** Go Modules (`go.mod`)
* **Messaging:** Apache Kafka (`github.com/segmentio/kafka-go`)
* **Translation API:** Google Cloud Translation (`cloud.google.com/go/translate`)
* **Testing:** Standard library (`testing`)

## 3. Architecture: Hexagonal (Ports & Adapters)
The project strictly implements Hexagonal Architecture tailored for Go idioms.
* **Forbidden:** Global state, heavy frameworks, deep nesting.
* **Mandatory:** Use of Go `interface{}` to define "Ports" in the domain layer.

### Directory Structure
```text
CartAI-ProductTranslator/
├── go.mod
├── cmd/
│   └── translator/
│       └── main.go                 # Composition Root and Entrypoint
├── internal/
│   ├── domain/                     # Business Entities (Structs) and Ports (Interfaces)
│   ├── application/                # Use Cases (Business Logic)
│   └── infrastructure/             # Adapters (Kafka, GCP)
└── .agents/
```

## 4. TDD Workflow
* **STRICT RULE (TDD):** Development must follow Red-Green-Refactor.
* Always write the failing test (`*_test.go`) before implementing the logic.
* Tests should reside in the same directory as the code they test. Use `package <name>_test` for black-box testing where appropriate.

## 5. Code Quality Rules
* **Formatting:** `gofmt` must be run on all code.
* **Linting:** Standard `go vet`.
* **Error Handling:** Idiomatic Go error handling (`if err != nil`). No panics in business logic.
* **English Only:** All code (variables, docstrings, comments), application logs, and markdown documentation MUST be written exclusively in English.

## 6. Git Flow & Agent Behavior
Governed by the global `~/.claude/CLAUDE.md` (commit/push only on explicit instruction, `[CAR-XX]` Linear prefix, no rushing review).
