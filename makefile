ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
endif

postgres:
	docker run --name event_pg18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18.4-alpine3.24
createdb:
	docker exec -it event_pg18 createdb --username=root --owner=root event_recorder
dropdb:
	docker exec -it event_pg18 dropdb event_recorder
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/event_recorder?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/event_recorder?sslmode=disable" -verbose down
sqlc_win:
	docker run --rm -v "$(CURDIR):/src" -w /src sqlc/sqlc generate
sqlc_mac:
	sqlc generate
docker-up:
	docker compose up --build
docker-down:
	docker compose down

.PHONY: postgres createdb dropdb migrateup migratedown sqlc_win sqlc_mac docker-up docker-down