package domain

import "time"

type UserStat struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	EventCount  int64     `json:"event_count"`
	CreatedAt   time.Time `json:"created_at"`
}
