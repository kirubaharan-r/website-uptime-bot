package models

import "time"

// Website represents a website to be monitored.
type Website struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// SSLInfo represents SSL certificate information.
type SSLInfo struct {
	Subject    string    `json:"subject"`
	Issuer     string    `json:"issuer"`
	NotBefore  time.Time `json:"not_before"`
	NotAfter   time.Time `json:"not_after"`
	IsValid    bool      `json:"is_valid"`
}

// Check represents a single check of a website.
type Check struct {
	ID            int
	WebsiteID     int
	Timestamp     time.Time
	Status        string // "up" or "down"
	ResponseTime  time.Duration
	Headers       map[string][]string
	SSLInfo       *SSLInfo
}
