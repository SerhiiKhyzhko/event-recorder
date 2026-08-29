package service

import (
	"context"
	"time"

	"github.com/SerhiiKhyzhko/event-recorder/internal/domain"
)

type EventService interface {
	CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error)
	GetEvents(ctx context.Context, userID int64, start, end time.Time, limit, offset int32) ([]domain.Event, error)
}

type eventService struct {
	repo domain.Repository
}

func NewEventService(repo domain.Repository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	return s.repo.CreateEvent(ctx, event)
}

func (s *eventService) GetEvents(ctx context.Context, userID int64, start, end time.Time, limit, offset int32) ([]domain.Event, error) {
	if start.After(end) {
		return nil, domain.ErrInvalidDateRange
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListEvents(ctx, userID, start, end, limit, offset)
}