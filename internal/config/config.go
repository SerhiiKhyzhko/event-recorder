package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	HTTPAddr      string
	DBURL         string
	StatsInterval time.Duration
}

func Load() *Config {
	dbUser := getEnv("POSTGRES_USER", "root")
	dbPass := getEnv("POSTGRES_PASSWORD", "secret")
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbName := getEnv("POSTGRES_DB", "event_recorder")
	dbSSLMode := getEnv("POSTGRES_SSLMODE", "disable")
	statsInterval := getEnvTime("STATS_INTERVAL", "4h")

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dbSSLMode)

	return &Config{
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		DBURL:         dbURL,
		StatsInterval: statsInterval,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvTime(key, defaultValue string) time.Duration {
	value, err := time.ParseDuration(getEnv(key, defaultValue))
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return value
}
