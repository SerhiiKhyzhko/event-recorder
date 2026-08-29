#!/bin/sh
set -e

DB_URL="postgresql://${POSTGRES_USER:-root}:${POSTGRES_PASSWORD:-secret}@${POSTGRES_HOST:-postgres}:${POSTGRES_PORT:-5432}/${POSTGRES_DB:-event_recorder}?sslmode=${POSTGRES_SSLMODE:-disable}"

migrate -path /migrations -database "$DB_URL" up
exec /app/event-recorder
