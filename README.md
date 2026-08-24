# Pipeline Processing

A small app for running data pipelines: point it at a CSV or JSON file (a local path, a `data:` URI, or a URL), it validates/transforms/aggregates the records and exports the result to S3, and you can watch it run from the UI. Built as a Go backend + a React frontend, backed by Postgres.

## Quick start (Docker)

```bash
cp .env.example .env   # defaults are fine for local dev, tweak if you want
docker compose up --build
```

Once everything's up:

- UI — http://localhost:5173
- API — http://localhost:8080 (everything needs an `X-API-Key` header, see `backend/README.md`)
- pgAdmin — http://localhost:5050, if you want to poke at the database directly
- Swagger — http://localhost:8080/swagger/index.html

## Running things natively

If you'd rather run the backend and UI directly instead of through Docker (faster feedback loop while developing), each has its own README with setup instructions:

- [`backend/README.md`](backend/README.md)
- [`ui/README.md`](ui/README.md)

You'll still need a Postgres instance either way — `docker compose up db` on its own works fine if you just want the database.

## Layout

```
backend/    Go API + the pipeline engine itself
ui/         React frontend
docker-compose.yml   the whole stack for local dev
```

## How a job actually runs

Creating a job kicks off a pipeline in the background, made of five stages wired together with channels: ingest → validate → transform → aggregate → export. Each stage runs its own pool of goroutines, so a slow ingest doesn't stall everything else behind it. The UI polls progress while a job is running, and you can cancel one mid-flight.

If the backend gets killed while jobs are running, they don't just resume when it comes back — on startup it marks anything still `pending`/`running` as `failed`, since there's no in-flight state left to resume from.
