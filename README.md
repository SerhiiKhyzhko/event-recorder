# event-recorder

Go REST API for recording user activity events and aggregating stats every 4 hours. Minimal React client lives in `web/`.

**Repository:** https://github.com/SerhiiKhyzhko/event-recorder

## Requirements

- **Docker:** Docker + Docker Compose
- **Local:** Go 1.26+, [golang-migrate](https://github.com/golang-migrate/migrate)

## Run with Docker

```bash
make docker-up
```

or

```bash
docker compose up --build
```

API: http://localhost:8080  
Migrations run automatically on container startup.

Stop:

```bash
make docker-down
```

## Run locally

1. Start Postgres:

```bash
make postgres
make createdb
make migrateup
```

2. Environment variables (optional, defaults exist):

```bash
cp .env.example .env
```

3. Start the API:

```bash
go run ./cmd
```

## Make commands

| Command | Description |
|---|---|
| `make docker-up` | Start via Docker Compose |
| `make docker-down` | Stop containers |
| `make postgres` | Run Postgres in Docker (local dev) |
| `make createdb` | Create `event_recorder` database |
| `make dropdb` | Drop database |
| `make migrateup` | Apply migrations |
| `make migratedown` | Roll back migrations |
| `make sqlc_win` | Regenerate sqlc code (Windows) |
| `make sqlc_mac` | Regenerate sqlc code (macOS/Linux) |

## Sample requests

Create an event:

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{"user_id": 42, "action": "page_view", "metadata": {"page": "/home"}}'
```

List events (dates in RFC3339):

```bash
curl "http://localhost:8080/events?user_id=42&start_date=2026-08-30T00:00:00Z&end_date=2026-08-31T00:00:00Z"
```

Get stats:

```bash
curl "http://localhost:8080/stats?user_id=42&start_date=2026-08-30T00:00:00Z&end_date=2026-08-31T00:00:00Z"
```

## Background job

Every `STATS_INTERVAL` (default **4h**), the service counts events from the previous interval, groups by `user_id`, and upserts results into `user_stats`.

## Frontend

Open `web/index.html` in a browser (API must be running on `:8080`).
