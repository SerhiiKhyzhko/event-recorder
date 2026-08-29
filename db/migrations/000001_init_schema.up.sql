CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "action" varchar(100) NOT NULL,
  "metadata" jsonb,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_stats" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "period_start" timestamptz NOT NULL,
  "period_end" timestamptz NOT NULL,
  "event_count" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "events" ("user_id", "created_at");

CREATE UNIQUE INDEX ON "user_stats" ("user_id", "period_start", "period_end");
