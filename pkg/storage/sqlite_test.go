package storage

import (
	"database/sql"
	"testing"
	"time"
	"uptime-monitor/pkg/models"
)

func TestSQLiteStore(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()

	store, err := NewSQLiteStore(db)
	if err != nil {
		t.Fatalf("Error creating store: %s", err)
	}

	website := models.Website{
		Name: "Google",
		URL:  "https://www.google.com",
	}

	if err := store.AddWebsite(website); err != nil {
		t.Fatalf("Error adding website: %s", err)
	}

	websites, err := store.GetWebsites()
	if err != nil {
		t.Fatalf("Error getting websites: %s", err)
	}

	if len(websites) != 1 {
		t.Fatalf("Expected 1 website, got %d", len(websites))
	}

	if websites[0].Name != website.Name {
		t.Fatalf("Expected website name %s, got %s", website.Name, websites[0].Name)
	}

	check := models.Check{
		WebsiteID:    websites[0].ID,
		Timestamp:    time.Now(),
		Status:       "up",
		ResponseTime: 100 * time.Millisecond,
	}

	if err := store.AddCheck(check); err != nil {
		t.Fatalf("Error adding check: %s", err)
	}

	latestCheck, ok, err := store.GetLatestCheck(websites[0].ID)
	if err != nil {
		t.Fatalf("Error getting latest check: %s", err)
	}

	if !ok {
		t.Fatalf("Expected to find a check, but didn't")
	}

	if latestCheck.Status != check.Status {
		t.Fatalf("Expected check status %s, got %s", check.Status, latestCheck.Status)
	}
}
