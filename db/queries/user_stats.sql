-- name: CreateUserStats :many
INSERT INTO user_stats (user_id, period_start, period_end, event_count)
SELECT
    user_id,
    sqlc.arg(period_start)::timestamptz,
    sqlc.arg(period_end)::timestamptz,
    COUNT(*)::bigint AS event_count
FROM events
WHERE created_at >= sqlc.arg(period_start) AND created_at < sqlc.arg(period_end)
GROUP BY user_id
ON CONFLICT (user_id, period_start, period_end)
DO UPDATE SET event_count = EXCLUDED.event_count
RETURNING *;

-- name: ListUserStats :many
SELECT * FROM user_stats
WHERE user_id = $1
  AND period_start >= $2
  AND period_end <= $3
ORDER BY period_start DESC;