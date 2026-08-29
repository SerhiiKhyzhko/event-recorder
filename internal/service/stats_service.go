package service

import (
	"context"
	"log"
	"time"

	"github.com/SerhiiKhyzhko/event-recorder/internal/domain"
)

type StatsService interface {
	GenerateStats(ctx context.Context) error
	GetStats(ctx context.Context, userID int64, start, end time.Time) ([]domain.UserStat, error)
}

type statsService struct {
	repo     domain.Repository
	interval time.Duration
}

func NewStatsService(repo domain.Repository, interval time.Duration) StatsService {
	return &statsService{
		repo:     repo,
		interval: interval,
	}
}

func (s *statsService) GenerateStats(ctx context.Context) error {
	periodEnd := time.Now().UTC()
	periodStart := periodEnd.Add(-s.interval)

	count, err := s.repo.CreateUserStats(ctx, periodStart, periodEnd)
	if err != nil {
		return err
	}
	log.Printf("generated stats for %d users (period %s - %s)", count, periodStart, periodEnd)
	return nil
}
func (s *statsService) GetStats(ctx context.Context, userID int64, start, end time.Time) ([]domain.UserStat, error) {
	if start.After(end) {
		return nil, domain.ErrInvalidDateRange
	}
	return s.repo.ListUserStats(ctx, userID, start, end)
}
