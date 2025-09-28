package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/uptime-monitor/internal/config"
	"github.com/example/uptime-monitor/internal/logging"
	"github.com/example/uptime-monitor/internal/monitor"
	"github.com/example/uptime-monitor/internal/store"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	repo := store.NewDynamo(ctx, cfg.AWSRegion, cfg.SitesTable, cfg.StatusTable)
	logger := logging.New("worker")

	monitor.SeedDemoData(ctx, repo)

	checker := &monitor.Checker{
		Repo:     repo,
		Logger:   logger,
		Interval: time.Duration(cfg.PollInterval) * time.Second,
	}

	if err := checker.Run(ctx); err != nil {
		logger.Printf("monitor stopped: %v", err)
	}
}
