# Backend

Go API for the pipeline app. Plain `net/http` (no router library), Postgres via `lib/pq`, S3 via the AWS SDK (points at LocalStack for local dev).

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
- `internal/models` — the DB row shapes. Kept as a separate package from `internal/pipeline`'s own types specifically to avoid an import cycle (both `job` and `pipeline` depend on it).
- `internal/config`, `internal/db` — env-based config loading, connection pool setup.
- `migrations/` — plain SQL, run through `golang-migrate`.
- `docs/` — generated swagger output, don't hand-edit it.

## A few things worth knowing

- Job IDs are UUIDs (`gen_random_uuid()`, native since Postgres 13, no extension required).
- Ingesters/transformers/exporters are a small plugin registry (`internal/pipeline/interface.go`), each registered via `init()` in its own file. Right now that's csv/json ingest, upper/lowercase transforms, and S3 export — adding a new one is just implementing the relevant interface and calling `Register*` on it.
- File-path sources are sandboxed to `SandboxDir` (`internal/pipeline/ingest.go`) to block path traversal — only `data:` URIs, `http(s)://` URLs, and paths under that directory are allowed.
- `golangci-lint run ./...` currently reports a handful of pre-existing `errcheck`/`bodyclose` findings (unchecked `rows.Close()`, `json.Encode()`, etc.) that haven't been cleaned up yet.
