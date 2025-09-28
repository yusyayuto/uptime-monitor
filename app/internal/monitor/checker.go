package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/example/uptime-monitor/internal/id"
	"github.com/example/uptime-monitor/internal/store"
)

// Checker periodically fetches registered sites and records their status.
type Checker struct {
	Repo     store.Repository
	Client   *http.Client
	Logger   interface{ Printf(string, ...any) }
	Interval time.Duration
}

// Run starts the monitoring loop. It is expected to run in its own goroutine.
func (c *Checker) Run(ctx context.Context) error {
	if c.Client == nil {
		c.Client = &http.Client{Timeout: 10 * time.Second}
	}
	if c.Interval == 0 {
		c.Interval = time.Minute
	}
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	for {
		if err := c.pollOnce(ctx); err != nil {
			c.Logger.Printf("poll error: %v", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c *Checker) pollOnce(ctx context.Context) error {
	sites, err := c.Repo.ListSites(ctx)
	if err != nil {
		return err
	}
	for _, site := range sites {
		status := store.Status{SiteID: site.ID, CheckedAt: time.Now().UTC()}
		start := time.Now()
		resp, err := c.Client.Get(site.URL)
		if err != nil {
			status.Healthy = false
			status.Message = err.Error()
		} else {
			status.LatencyMS = time.Since(start).Milliseconds()
			status.Healthy = resp.StatusCode < 500
			status.Message = resp.Status
			_ = resp.Body.Close()
		}
		if err := c.Repo.PutStatus(ctx, status); err != nil {
			c.Logger.Printf("store status for %s: %v", site.URL, err)
		}
	}
	return nil
}

// SeedDemoData registers a sample site for quick-start demos.
func SeedDemoData(ctx context.Context, repo store.Repository) {
	sites, err := repo.ListSites(ctx)
	if err == nil && len(sites) > 0 {
		return
	}
	_ = repo.PutSite(ctx, store.Site{
		ID:        id.New(),
		URL:       "https://example.com",
		CreatedAt: time.Now().UTC(),
	})
}
