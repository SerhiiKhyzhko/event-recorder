package domain

import (
	"context"
	"encoding/json"
	"time"
)

type Event struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Action    string          `json:"action"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
}

type Repository interface {
	CreateEvent(ctx context.Context, event Event) (Event, error)
	ListEvents(ctx context.Context, userID int64, start, end time.Time, limit, offset int32) ([]Event, error)
	CreateUserStats(ctx context.Context, start, end time.Time) (int, error)
	ListUserStats(ctx context.Context, userID int64, start, end time.Time) ([]UserStat, error)
}
