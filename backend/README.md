# Backend

Go API for the pipeline app. Plain `net/http` (no router library) and Postgres via `lib/pq`.

## Functionality

This backend powers a highly concurrent, streaming data processing pipeline. Its core capabilities include:
- **Data Ingestion:** Reads data from various sources (CSV, JSON, APIs, data URIs) simultaneously.
- **Validation:** Ensures data integrity by filtering out empty or invalid records on the fly.
- **Transformation:** Mutates data (e.g., lowercase, uppercase) in real-time as it flows through the system.
- **Aggregation:** Groups and computes metrics (Sum, Average, Count) across millions of records efficiently.
- **Exporting:** Saves the final aggregated results and raw records to external storage.
- **Real-Time Tracking:** Provides live progress updates (processed counts, error counts) via polling.

The pipeline is built on a "streaming" architecture. Records are passed between stages one-by-one via Go channels, meaning massive datasets can be processed with a very small memory footprint.

## Running it

Needs a reachable Postgres. Easiest is `docker compose up db` from the repo root; otherwise point the `POSTGRES_*` env vars at whatever you've got.

```bash
cp ../.env.example .env   # or just export the vars yourself
make migrate_up           # applies all 5 migrations
make run                  # regenerates swagger docs, builds, runs
```

Listens on `:8080` by default (`PORT` env var). Every request needs an `X-API-Key` header — defaults to `secret-pipeline-key` if `API_KEY` isn't set, which is fine locally but should obviously be changed for anything real.

Swagger UI is served at `/swagger/index.html` once it's running.

## Make targets

- `make build` — regenerate swagger docs, then compile
- `make test` — `go test -race -cover ./...`
- `make lint` — golangci-lint
- `make fmt` — gofmt/goimports
- `make migrate_up` / `make migrate_down` — apply/roll back migrations
- `make create_migration name=whatever` — scaffold a new migration pair
- `make swagger` — just regenerate the docs (needed after changing any `@Param`/`@Success` annotations)

## Layout

- `cmd/server` — entrypoint. Wires config → DB → server together, handles graceful shutdown on SIGTERM/interrupt.
- `internal/job` — HTTP handlers, the job service (owns the in-memory bookkeeping for jobs currently running — cancel funcs, progress trackers), and the Postgres repositories.
- `internal/pipeline` — the actual pipeline engine: ingest/validate/transform/aggregate/export, each stage a pool of goroutines connected by channels.
- `internal/middleware` — API key auth, CORS, rate limiting, security headers, request logging (with correlation IDs).
- `internal/models` — the DB row shapes.
- `internal/config`, `internal/db` — env-based config loading, connection pool setup.
- `migrations/` — plain SQL, run through `golang-migrate`.
- `docs/` — generated swagger output, don't hand-edit it.

## Architecture Highlights

- **Plugin System:** Ingesters, transformers, and exporters are built on an interface registry (`internal/pipeline/interface.go`), making it incredibly easy to add new data sources or transformation rules.
- **Security:** Built-in Path Traversal protection ensures local file reads cannot escape the designated sandbox directory.
