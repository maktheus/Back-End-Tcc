# Back-End-Tcc

This repository provides a reference implementation of a benchmark orchestration backend written in Go. It demonstrates a microservice layout with shared packages, HTTP handlers, queue-based workers, and automated tests.

## Prerequisites

* Go 1.22+

## Project layout

```
├── cmd/              # Entrypoints for the API gateway and individual microservices
├── pkg/              # Shared packages (configuration, logging, queue, storage, models)
├── services/         # Domain modules broken into handlers, services, and repositories
├── tests/            # Unit and integration tests with fixtures
└── README.md
```

## Getting started

1. Install dependencies and generate the Go module files:

   ```bash
   go mod tidy
   ```

2. Run the tests to verify the installation:

   ```bash
   go test ./...
   ```

3. Start the API gateway (which wires together all in-memory components) via:

   ```bash
   go run ./cmd/api
   ```

   The gateway exposes endpoints for authentication, agent registration, benchmark definitions, submissions, scoring, traces, and the leaderboard. The default port can be configured through environment variables such as `HTTP_PORT`.

4. Each microservice can also be started individually. For example, to run the orchestrator-only service:

   ```bash
   go run ./cmd/orchestrator
   ```

5. Observability is enabled by default through structured logging and an in-memory metrics recorder that instruments queue publishing, service handlers, and HTTP entrypoints. Metrics are captured per service (e.g., `orchestrator_submit_total`, `queue_messages_total`) and can be inspected during tests or exported by wiring the `pkg/observability/metrics` package into an external backend.

## Configuration

Configuration values are loaded from environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `APP_ENV` | `development` | Environment indicator logged by the services |
| `HTTP_PORT` | `8080` | Port used by each HTTP service |
| `QUEUE_BUFFER_SIZE` | `100` | Capacity hint for the in-memory queue |
| `STORAGE_DSN` | `memory://default` | Placeholder storage connection string |
| `JWT_SIGNING_SECRET` | `dev-secret` | Secret used for signing authentication tokens |

## Testing

Unit tests cover individual services such as orchestrator submission handling and scoring aggregation. Integration tests (`tests/integration/e2e_benchmark_flow_test.go`) exercise the full submission-to-scoring flow using the in-memory queue.

Run all tests with:

```bash
go test ./...
```

## Sample data

Fixtures located in `tests/fixtures/` provide sample payloads that can be used with tools such as `curl` or Postman to simulate submissions when experimenting with the API gateway.

## Next steps

This skeleton is intended to be extended with persistent storage backends, authentication tokens, real scoring logic, production-ready queue implementations, and real telemetry exporters (Prometheus/OpenTelemetry) connected to the existing logging and metrics hooks. The modular structure and clean interfaces make it straightforward to add these capabilities.
