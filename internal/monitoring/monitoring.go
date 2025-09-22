package monitoring

import (
	"fmt"
	"net/http"
	"time"
	"uptime-monitor/internal/models"
	"uptime-monitor/internal/store"
)

// StartMonitoring starts the monitoring loop.
func StartMonitoring(s *store.InMemoryStore, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			fmt.Println("--- New Check Cycle ---")
			websites := s.GetWebsites()
			for _, website := range websites {
				go CheckWebsite(s, website)
			}
		}
	}
}

// PerformCheck performs a single check of a website.
func PerformCheck(url string) models.Check {
	check := models.Check{
		Timestamp: time.Now(),
	}

	start := time.Now()
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	check.ResponseTime = time.Since(start)

	if err != nil {
		check.Status = "down"
		fmt.Printf("Error checking %s: %s\n", url, err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			check.Status = "up"
		} else {
			check.Status = "down"
		}
		fmt.Printf("Checked %s: %s (%d)\n", url, check.Status, resp.StatusCode)
	}
	return check
}

// CheckWebsite performs a single check of a website and stores the result.
func CheckWebsite(s *store.InMemoryStore, website models.Website) {
	check := PerformCheck(website.URL)
	check.WebsiteID = website.ID
	s.AddCheck(check)
}
