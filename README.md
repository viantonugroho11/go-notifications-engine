# Notifications Engine (Go)

A Go service for managing **notifications**, **templates**, **per-user logs**, and **inbox** records, built around **clean architecture** patterns. It exposes a REST API (Echo), persists data with **GORM + PostgreSQL**, publishes and consumes events via **Kafka** ([go-lib/kafka](https://github.com/viantonugroho11/go-lib)), and initializes a **Redis** client for caching or future cross-cutting concerns. Configuration is loaded with **[go-config-library](https://github.com/viantonugroho11/go-config-library)** (optional Consul + local config files).

The Go module name is `github.com/viantonugroho11/go-notifications-engine` (see `go.mod`); this repository implements a **notifications engine** on top of that layout.

## Features

- **HTTP API** — CRUD for users, notifications, notification templates, notification logs, and notification inbox.
- **Domain model** — Notifications carry channel, category, scheduling, JSON `data`, and related `notification_logs` per user.
- **Kafka** — Producer for outbound notification events; separate **consumer** binary for processing pipeline stages (e.g. `notification`, `sent`).
- **PostgreSQL** — GORM `AutoMigrate` for `users`, `notifications`, `notification_templates`, `notification_logs`, `notification_inbox`.
- **Redis** — Client wired at bootstrap (ready for rate limiting, locks, or cache).
- **Integrations (code paths)** — Email (gomail), Firebase/FCM helpers, JSON schema validation, and outbound HTTP clients for downstream services (see `internal/client`).

## Requirements

- **Go** — `go 1.24.0` / toolchain as declared in `go.mod` (Go 1.24+ recommended).
- **PostgreSQL** — compatible with GORM Postgres driver.
- **Kafka** — broker reachable from the app and consumer processes; ZooKeeper if you use the provided Compose stack.
- **Redis** — for the initialized client (optional for minimal API-only runs depending on your code paths).

## Quick start

### 1. Clone and dependencies

```bash
git clone <repository-url>
cd go-notifications-engine
go mod download
```

### 2. Infrastructure

Use Docker Compose for Postgres, ZooKeeper, Kafka, and Redis:

```bash
docker compose up -d --build
```

Default Compose ports: app **8080**, Postgres **5432**, Kafka **9092**, Redis **6379**, ZooKeeper **2181**.

### 3. Configuration

The application loads configuration via `go-config-library`:

- Service key: **`github.com/viantonugroho11/go-notifications-engine`** (see `cmd/app/main.go` and `cmd/consumer/main.go`).
- Optional **Consul**: set `CONSUL_URL` when using remote config.
- Local files: **`WithConfigFileSearchPaths("./config")`** — provide a `./config` directory with files your environment expects, or align paths with your deployment. A baseline example lives under **`configs/config.yaml`** (you may copy or symlink into `./config` for local runs).

The `Configuration` struct in `internal/config` defines nested sections (`app`, `database`, `kafka`, `redis`, etc.). Keep your YAML/Consul keys consistent with how `go-config-library` maps into that struct in your environment.

### 4. Run the HTTP server

```bash
go run ./cmd/app
```

On startup the app connects to PostgreSQL, runs **AutoMigrate**, starts the Echo server, initializes **Redis**, and sets up the **Kafka producer** (see `internal/bootstrap/app.go`).

### 5. Run a Kafka consumer (separate process)

Consumers are started explicitly by key:

```bash
go run ./cmd/consumer -consumer notification
# or
go run ./cmd/consumer -consumer sent
```

Available keys are defined in `internal/transport/event/kafka` (`notification`, `sent`). The `sent` consumer requires `kafka.topic_sent` (or equivalent) to be set in configuration — see `internal/infrastructure/broker/kafka/registry.go`.

## Project layout

```
cmd/
  app/           # HTTP API entrypoint
  consumer/      # Kafka consumer entrypoint (-consumer <key>)
internal/
  bootstrap/     # Wiring: DB, producer, Redis, Echo, consumer app
  config/        # Configuration types and helpers (DSN, Kafka, email, FCM, …)
  entity/        # Domain entities
  client/        # Outbound integrations (email, firebase, notification gateway, person service, …)
  infrastructure/
    database/postgres/   # GORM connect + AutoMigrate
    broker/kafka/        # Producer, consumer runner, registry
    cache/redis/         # Redis client
  repository/    # Interfaces + GORM models + postgres implementations
  transport/
    apis/        # Echo routes, HTTP handlers, DTOs
    event/kafka/ # Consumer handlers and event routing
  usecase/       # Application services and state machines
  shared/        # Shared utilities (e.g. JSON schema)
configs/
  config.yaml    # Example baseline (adjust for local/Compose)
```

## HTTP API

Base URL: `http://localhost:<port>` (port from config, e.g. **8080**).

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness — returns `ok` |
| POST, GET, GET/:id, PUT/:id, DELETE/:id | `/users` | User CRUD |
| POST, GET, GET/:id, PUT/:id, DELETE/:id | `/notifications` | Notification CRUD |
| POST, GET, GET/:id, PUT/:id, DELETE/:id | `/notification-templates` | Template CRUD |
| POST, GET, GET/:id, PUT/:id, DELETE/:id | `/notification-logs` | Log CRUD |
| POST, GET, GET/:id, PUT/:id, DELETE/:id | `/notification-inbox` | Inbox CRUD |

Example: create a notification (shape from `internal/transport/apis/dto/notification_dto.go`):

```bash
curl -sS -X POST "http://localhost:8080/notifications" \
  -H "Content-Type: application/json" \
  -d '{
    "event_key": "order.created",
    "notification_template_id": "<template-uuid>",
    "data": { "orderId": "123" },
    "channel": "email",
    "category": "transactional",
    "user_ids": ["<user-uuid>"]
  }'
```

Field names and allowed `channel` / `category` values should match your domain enums in `internal/entity/notifications`.

## Kafka

- **Producer** — `internal/infrastructure/broker/kafka/notification_producer.go` (`NewNotificationProducer`, typed publish API).
- **Consumer registry** — `internal/infrastructure/broker/kafka/registry.go` maps consumer keys to topics and group/client IDs from `config.Configuration`.
- **Handlers** — `internal/transport/event/kafka/handler/` (e.g. notification lifecycle updates, send pipeline).
- **Message contract** — Event handlers use `internal/entity/notifications` event payloads as decoded by `go-lib/kafka`.

Ensure broker addresses and topics in config match your cluster. For `sent`, set the sent topic in configuration or the consumer will fail validation at startup.

## Database

- **Migrations** — At runtime, `postgres.Migrate` runs GORM `AutoMigrate` for all registered models (see `internal/infrastructure/database/postgres/connection.go`).
- **Legacy makefile** — The root `makefile` contains example `migrate` CLI targets pointing at another database name; treat as a template or update to match this project before use.

## Docker

Build the API image:

```bash
docker build -t go-notifications-engine:latest .
```

The **Dockerfile** builds `./cmd/app` and copies `configs/` into the image. Ensure the runtime **working directory and config search paths** match how you deploy `go-config-library` (the code searches `./config` by default).

## Development

- **Format / vet / test**

```bash
go fmt ./...
go vet ./...
go test ./...
```

- **Contributing** — See [CONTRIBUTING.md](./CONTRIBUTING.md) for branch workflow, code layout expectations, and review guidelines.

## Troubleshooting

| Symptom | Check |
|--------|--------|
| `config load` error | `CONSUL_URL`, local `./config` files, and service name `github.com/viantonugroho11/go-notifications-engine` |
| Postgres connection failed | DSN / host / port / credentials; Compose network vs localhost |
| Kafka producer/consumer errors | `KAFKA_BROKERS` (or mapped `kafka.brokers`), topic names, consumer `-consumer` flag |
| Redis errors | Address, password, DB index |
| Missing `topic_sent` | Required for `-consumer sent` — set in config |

## License

MIT — see license terms in the repository if a `LICENSE` file is present.
