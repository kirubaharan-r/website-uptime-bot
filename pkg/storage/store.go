package storage

import (
	"sync"
	"uptime-monitor/pkg/models"
)

// Store is the interface for the storage layer.
type Store interface {
	AddWebsite(website models.Website) error
	GetWebsites() ([]models.Website, error)
	AddCheck(check models.Check) error
	GetLatestCheck(websiteID int) (models.Check, bool, error)
}

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
func (s *InMemoryStore) AddWebsite(website models.Website) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.websites = append(s.websites, website)
	return nil
}

// GetWebsites returns all websites from the store.
func (s *InMemoryStore) GetWebsites() ([]models.Website, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.websites, nil
}

// AddCheck adds a check result to the store.
func (s *InMemoryStore) AddCheck(check models.Check) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks[check.WebsiteID] = check
	return nil
}

// GetLatestCheck returns the latest check for a website.
func (s *InMemoryStore) GetLatestCheck(websiteID int) (models.Check, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	check, ok := s.checks[websiteID]
	return check, ok, nil
}
