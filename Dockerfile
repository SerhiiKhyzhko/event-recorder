FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /event-recorder ./cmd

FROM migrate/migrate AS migrate

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY --from=migrate /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /event-recorder /app/event-recorder
COPY db/migrations /migrations
COPY docker-entrypoint.sh /entrypoint.sh
RUN sed -i 's/\r$//' /entrypoint.sh && chmod +x /entrypoint.sh

WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
