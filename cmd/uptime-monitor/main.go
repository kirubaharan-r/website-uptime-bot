package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"uptime-monitor/pkg/api"
	"uptime-monitor/pkg/config"
	"uptime-monitor/pkg/monitoring"
	"uptime-monitor/pkg/storage"
)

func main() {
	fmt.Println("Uptime Monitor Starting...")

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	store, err := storage.NewSQLiteStore(db)
	if err != nil {
		log.Fatalf("Error creating store: %s", err)
	}

	for _, website := range cfg.Websites {
		if err := store.AddWebsite(website); err != nil {
			log.Printf("Error adding website %s: %s", website.URL, err)
		}
	}

	ticker := time.NewTicker(1 * time.Minute)
	go monitoring.StartMonitoring(store, ticker)

	server := api.NewServer(cfg.Addr, store)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
