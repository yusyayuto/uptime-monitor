package store

import (
	"context"
	"sync"
)

type memoryRepo struct {
	sites  map[string]Site
	status map[string]Status
	mu     sync.RWMutex
}

// NewInMemory returns a Repository backed by in-memory maps.
func NewInMemory() Repository {
	return &memoryRepo{
		sites:  make(map[string]Site),
		status: make(map[string]Status),
	}
}

func (m *memoryRepo) PutSite(_ context.Context, site Site) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sites[site.ID] = site
	return nil
}

func (m *memoryRepo) ListSites(_ context.Context) ([]Site, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sites := make([]Site, 0, len(m.sites))
	for _, s := range m.sites {
		sites = append(sites, s)
	}
	return sites, nil
}

func (m *memoryRepo) PutStatus(_ context.Context, status Status) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status[status.SiteID] = status
	return nil
}

func (m *memoryRepo) LatestStatus(_ context.Context, siteID string) (Status, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	status, ok := m.status[siteID]
	if !ok {
		return Status{}, ErrNotFound
	}
	return status, nil
}
