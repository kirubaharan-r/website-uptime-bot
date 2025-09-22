package models

import "time"

// Website represents a website to be monitored.
type Website struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Check represents a single check of a website.
type Check struct {
	ID         int
	WebsiteID  int
	Timestamp  time.Time
	Status     string // "up" or "down"
	ResponseTime time.Duration
}
