package store

import (
	"sync"
	"uptime-monitor/internal/models"
)

// InMemoryStore holds the data for the application in memory.
type InMemoryStore struct {
	mu       sync.RWMutex
	websites []models.Website
	checks   map[int]models.Check
}

// NewInMemoryStore creates a new InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		websites: make([]models.Website, 0),
		checks:   make(map[int]models.Check),
	}
}

// AddWebsite adds a website to the store.
func (s *InMemoryStore) AddWebsite(website models.Website) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.websites = append(s.websites, website)
}

// GetWebsites returns all websites from the store.
func (s *InMemoryStore) GetWebsites() []models.Website {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.websites
}

// AddCheck adds a check result to the store.
func (s *InMemoryStore) AddCheck(check models.Check) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks[check.WebsiteID] = check
}

// GetLatestCheck returns the latest check for a website.
func (s *InMemoryStore) GetLatestCheck(websiteID int) (models.Check, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	check, ok := s.checks[websiteID]
	return check, ok
}
