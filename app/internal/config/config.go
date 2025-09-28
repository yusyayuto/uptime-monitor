package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config captures runtime configuration sourced from environment variables.
type Config struct {
	Environment   string
	HTTPPort      int
	AWSRegion     string
	SitesTable    string
	StatusTable   string
	PollInterval  int
	NotifierTopic string
}

// Load reads environment variables and applies sane defaults for local runs.
func Load() (Config, error) {
	cfg := Config{
		Environment:   getEnv("APP_ENV", "dev"),
		AWSRegion:     getEnv("AWS_REGION", "ap-northeast-1"),
		SitesTable:    getEnv("SITES_TABLE", "ssk-dev-sites"),
		StatusTable:   getEnv("STATUS_TABLE", "ssk-dev-status"),
		NotifierTopic: getEnv("NOTIFIER_TOPIC_ARN", ""),
	}

	port, err := strconv.Atoi(getEnv("API_PORT", "8080"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid API_PORT: %w", err)
	}
	cfg.HTTPPort = port

	interval, err := strconv.Atoi(getEnv("POLL_INTERVAL_SECONDS", "60"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid POLL_INTERVAL_SECONDS: %w", err)
	}
	cfg.PollInterval = interval

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
