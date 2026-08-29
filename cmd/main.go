package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SerhiiKhyzhko/event-recorder/internal/config"
	"github.com/SerhiiKhyzhko/event-recorder/internal/handler"
	"github.com/SerhiiKhyzhko/event-recorder/internal/repository/postgres"
	"github.com/SerhiiKhyzhko/event-recorder/internal/scheduler"
	"github.com/SerhiiKhyzhko/event-recorder/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	pool := initDB(cfg.DBURL)
	defer pool.Close()

	repo := postgres.NewRepository(pool)
	eventSvc := service.NewEventService(repo)
	statsSvc := service.NewStatsService(repo, cfg.StatsInterval)

	eventHandler := handler.NewEventHandler(eventSvc)
	statsHandler := handler.NewStatsHandler(statsSvc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched := scheduler.New(statsSvc, cfg.StatsInterval)
	go sched.Run(ctx)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: handler.NewRouter(eventHandler, statsHandler),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	//graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("shutdown signal received")

	cancel() //stop the scheduler

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func initDB(dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	//DB conn check
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return pool
}
