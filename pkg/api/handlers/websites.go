package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uptime-monitor/pkg/models"
	"uptime-monitor/pkg/monitoring"
	"uptime-monitor/pkg/storage"
)

type websiteStatus struct {
	Website      models.Website        `json:"website"`
	Status       string                `json:"status"`
	ResponseTime string                `json:"response_time"`
	Headers      map[string][]string   `json:"headers"`
	SSLInfo      *models.SSLInfo       `json:"ssl_info"`
}

// CheckHandler handles the checking of websites.
func CheckHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var urls []string
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		results := make([]websiteStatus, len(urls))
		for i, url := range urls {
			check := monitoring.PerformCheck(url)
			results[i] = websiteStatus{
				Website: models.Website{
					URL: url,
				},
				Status:       check.Status,
				ResponseTime: fmt.Sprintf("%.2fms", float64(check.ResponseTime.Nanoseconds())/1e6),
				Headers:      check.Headers,
				SSLInfo:      check.SSLInfo,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

// StatusHandler handles the status of websites.
func StatusHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		websites, err := store.GetWebsites()
		if err != nil {
			http.Error(w, "Error getting websites", http.StatusInternalServerError)
			return
		}
		status := make([]websiteStatus, len(websites))
		for i, website := range websites {
			check, ok, err := store.GetLatestCheck(website.ID)
			if err != nil {
				http.Error(w, "Error getting latest check", http.StatusInternalServerError)
				return
			}
			if !ok {
				continue
			}
			status[i] = websiteStatus{
				Website:      website,
				Status:       check.Status,
				ResponseTime: fmt.Sprintf("%.2fms", float64(check.ResponseTime.Nanoseconds())/1e6),
				Headers:      check.Headers,
				SSLInfo:      check.SSLInfo,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
