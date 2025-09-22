package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uptime-monitor/internal/models"
	"uptime-monitor/internal/monitoring"
	"uptime-monitor/internal/store"
)

type websiteStatus struct {
	Website      models.Website `json:"website"`
	Status       string         `json:"status"`
	ResponseTime string         `json:"response_time"`
}

// StartServer starts the web server.
func StartServer(s *store.InMemoryStore, addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", apiStatusHandler(s))
	mux.HandleFunc("/api/check", apiCheckHandler())

	fmt.Printf("Starting web server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func apiCheckHandler() http.HandlerFunc {
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
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func apiStatusHandler(s *store.InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		websites := s.GetWebsites()
		status := make([]websiteStatus, len(websites))
		for i, website := range websites {
			check, ok := s.GetLatestCheck(website.ID)
			if !ok {
				continue
			}
			status[i] = websiteStatus{
				Website:      website,
				Status:       check.Status,
				ResponseTime: fmt.Sprintf("%.2fms", float64(check.ResponseTime.Nanoseconds())/1e6),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
