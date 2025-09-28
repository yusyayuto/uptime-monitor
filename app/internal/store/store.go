package store

import (
	"context"
	"errors"
	"time"
)

// Site represents a monitored endpoint.
type Site struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// Status represents the last observation for a site.
type Status struct {
	SiteID    string    `json:"site_id"`
	CheckedAt time.Time `json:"checked_at"`
	LatencyMS int64     `json:"latency_ms"`
	Healthy   bool      `json:"healthy"`
	Message   string    `json:"message"`
}

// Repository defines storage operations required by the API and worker.
type Repository interface {
	PutSite(ctx context.Context, site Site) error
	ListSites(ctx context.Context) ([]Site, error)
	PutStatus(ctx context.Context, status Status) error
	LatestStatus(ctx context.Context, siteID string) (Status, error)
}

var ErrNotFound = errors.New("site not found")
