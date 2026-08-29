package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SerhiiKhyzhko/event-recorder/internal/domain"
	"github.com/SerhiiKhyzhko/event-recorder/internal/repository/postgres/db"
)

type Repository struct {
	pool *pgxpool.Pool
	*db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		Queries: db.New(pool),
	}
}

func (r *Repository) CreateEvent(ctx context.Context, e domain.Event) (domain.Event, error) {
	arg := db.CreateEventParams{
		UserID:   e.UserID,
		Action:   e.Action,
		Metadata: e.Metadata,
	}

	dbEvent, err := r.Queries.CreateEvent(ctx, arg)
	if err != nil {
		return domain.Event{}, err
	}

	return domain.Event{
		ID:        dbEvent.ID,
		UserID:    dbEvent.UserID,
		Action:    dbEvent.Action,
		Metadata:  dbEvent.Metadata,
		CreatedAt: dbEvent.CreatedAt.Time,
	}, nil
}

func (r *Repository) ListEvents(ctx context.Context, userID int64, start, end time.Time, limit, offset int32) ([]domain.Event, error) {
	arg := db.ListEventsParams{
		UserID:      userID,
		CreatedAt:   pgtype.Timestamptz{Time: start, Valid: true},
		CreatedAt_2: pgtype.Timestamptz{Time: end, Valid: true},
		Limit:       limit,
		Offset:      offset,
	}

	dbEvents, err := r.Queries.ListEvents(ctx, arg)
	if err != nil {
		return nil, err
	}

	var events []domain.Event
	for _, e := range dbEvents {
		events = append(events, domain.Event{
			ID:        e.ID,
			UserID:    e.UserID,
			Action:    e.Action,
			Metadata:  e.Metadata,
			CreatedAt: e.CreatedAt.Time,
		})
	}
	return events, nil
}

func (r *Repository) CreateUserStats(ctx context.Context, start, end time.Time) (int, error) {
	arg := db.CreateUserStatsParams{
		PeriodStart: pgtype.Timestamptz{Time: start, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: end, Valid: true},
	}

	rows, err := r.Queries.CreateUserStats(ctx, arg)
	if err != nil {
		return 0, err
	}

	return len(rows), nil
}

func (r *Repository) ListUserStats(ctx context.Context, userID int64, start, end time.Time) ([]domain.UserStat, error) {
	arg := db.ListUserStatsParams{
		UserID:      userID,
		PeriodStart: pgtype.Timestamptz{Time: start, Valid: true},
		PeriodEnd:   pgtype.Timestamptz{Time: end, Valid: true},
	}

	dbStats, err := r.Queries.ListUserStats(ctx, arg)
	if err != nil {
		return nil, err
	}

	var stats []domain.UserStat
	for _, s := range dbStats {
		stats = append(stats, domain.UserStat{
			ID:          s.ID,
			UserID:      s.UserID,
			PeriodStart: s.PeriodStart.Time,
			PeriodEnd:   s.PeriodEnd.Time,
			EventCount:  s.EventCount,
			CreatedAt:   s.CreatedAt.Time,
		})
	}
	return stats, nil
}
