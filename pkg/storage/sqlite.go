package storage

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"uptime-monitor/pkg/models"
)

// SQLiteStore holds the data for the application in a SQLite database.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLiteStore.
func NewSQLiteStore(db *sql.DB) (*SQLiteStore, error) {
	// Create tables if they don't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS websites (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS checks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			website_id INTEGER NOT NULL,
			timestamp DATETIME NOT NULL,
			status TEXT NOT NULL,
			response_time INTEGER NOT NULL,
			headers TEXT,
			ssl_info TEXT,
			FOREIGN KEY (website_id) REFERENCES websites (id)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

// AddWebsite adds a website to the store.
func (s *SQLiteStore) AddWebsite(website models.Website) error {
	_, err := s.db.Exec("INSERT INTO websites (name, url) VALUES (?, ?)", website.Name, website.URL)
	return err
}

// GetWebsites returns all websites from the store.
func (s *SQLiteStore) GetWebsites() ([]models.Website, error) {
	rows, err := s.db.Query("SELECT id, name, url FROM websites")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var websites []models.Website
	for rows.Next() {
		var website models.Website
		if err := rows.Scan(&website.ID, &website.Name, &website.URL); err != nil {
			return nil, err
		}
		websites = append(websites, website)
	}

	return websites, nil
}

// AddCheck adds a check result to the store.
func (s *SQLiteStore) AddCheck(check models.Check) error {
	headers, err := json.Marshal(check.Headers)
	if err != nil {
		return err
	}

	sslInfo, err := json.Marshal(check.SSLInfo)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO checks (website_id, timestamp, status, response_time, headers, ssl_info) VALUES (?, ?, ?, ?, ?, ?)",
		check.WebsiteID, check.Timestamp, check.Status, check.ResponseTime, headers, sslInfo)
	return err
}

// GetLatestCheck returns the latest check for a website.
func (s *SQLiteStore) GetLatestCheck(websiteID int) (models.Check, bool, error) {
	row := s.db.QueryRow("SELECT id, website_id, timestamp, status, response_time, headers, ssl_info FROM checks WHERE website_id = ? ORDER BY timestamp DESC LIMIT 1", websiteID)

	var check models.Check
	var headers, sslInfo []byte
	err := row.Scan(&check.ID, &check.WebsiteID, &check.Timestamp, &check.Status, &check.ResponseTime, &headers, &sslInfo)
	if err == sql.ErrNoRows {
		return models.Check{}, false, nil
	}
	if err != nil {
		return models.Check{}, false, err
	}

	if err := json.Unmarshal(headers, &check.Headers); err != nil {
		return models.Check{}, false, err
	}
	if err := json.Unmarshal(sslInfo, &check.SSLInfo); err != nil {
		return models.Check{}, false, err
	}

	return check, true, nil
}
