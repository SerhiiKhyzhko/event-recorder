-- name: CreateEvent :one
INSERT INTO events (user_id, action, metadata)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListEvents :many
SELECT * FROM events
WHERE user_id = $1
  AND created_at >= $2
  AND created_at < $3
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;