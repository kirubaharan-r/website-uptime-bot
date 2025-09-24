package config

import (
	"encoding/json"
	"os"
	"uptime-monitor/pkg/models"
)

// Config holds the application configuration.
type Config struct {
	Addr     string          `json:"addr"`
	DBPath   string          `json:"db_path"`
	Websites []models.Website `json:"websites"`
}

// LoadConfig loads the configuration from a file.
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
