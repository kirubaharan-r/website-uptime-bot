package main

import (
	"fmt"
	"os"
	"time"
	"uptime-monitor/internal/models"
	"uptime-monitor/internal/monitoring"
	"uptime-monitor/internal/server"
	"uptime-monitor/internal/store"
)

func main() {
	fmt.Println("Uptime Monitor Starting...")

	addr := os.Getenv("UPTIME_MONITOR_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	s := store.NewInMemoryStore()

	s.AddWebsite(models.Website{ID: 1, Name: "Google", URL: "https://www.google.com"})
	s.AddWebsite(models.Website{ID: 2, Name: "GitHub", URL: "https://www.github.com"})
	s.AddWebsite(models.Website{ID: 3, Name: "Invalid URL", URL: "https://a-very-invalid-url.com"})

	ticker := time.NewTicker(1 * time.Minute)
	go monitoring.StartMonitoring(s, ticker)

	server.StartServer(s, addr)
}
