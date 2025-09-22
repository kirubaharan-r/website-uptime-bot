package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uptime-monitor/internal/models"
	"uptime-monitor/internal/store"
)

type websiteStatus struct {
	Website      models.Website `json:"website"`
	Status       string         `json:"status"`
	ResponseTime string         `json:"response_time"`
}

// StartServer starts the web server.
func StartServer(s *store.InMemoryStore) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", apiStatusHandler(s))

	fmt.Println("Starting web server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
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
