package scheduler

import (
	"context"
	"log"
	"time"
)

type StatsService interface {
	GenerateStats(ctx context.Context) error
}

type Scheduler struct {
	statsService StatsService
	interval     time.Duration
}

func New(svc StatsService, interval time.Duration) *Scheduler {
	return &Scheduler{
		statsService: svc,
		interval:     interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("scheduler started, interval: %s", s.interval)

	for {
		select {
		case <-ticker.C:
			log.Println("running scheduled stats aggregation")
			if err := s.statsService.GenerateStats(ctx); err != nil {
				log.Printf("generation failed: %v", err)
			} else {
				log.Println("stats aggregation completed successfully")
			}
		case <-ctx.Done():
			log.Println("scheduler stopped")
			return
		}
	}
}
