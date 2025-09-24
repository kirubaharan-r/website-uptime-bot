package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"uptime-monitor/pkg/models"
)

type mockStore struct{}

type Store interface {
	AddWebsite(website models.Website) error
	GetWebsites() ([]models.Website, error)
	AddCheck(check models.Check) error
	GetLatestCheck(websiteID int) (models.Check, bool, error)
}

func (m *mockStore) AddWebsite(website models.Website) error {
	return nil
}

func (m *mockStore) GetWebsites() ([]models.Website, error) {
	return []models.Website{
		{ID: 1, Name: "Google", URL: "https://www.google.com"},
	}, nil
}

func (m *mockStore) AddCheck(check models.Check) error {
	return nil
}

func (m *mockStore) GetLatestCheck(websiteID int) (models.Check, bool, error) {
	return models.Check{
		ID:           1,
		WebsiteID:    1,
		Timestamp:    time.Now(),
		Status:       "up",
		ResponseTime: 100 * time.Millisecond,
	}, true, nil
}

func TestStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/websites/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	var store Store = &mockStore{}
	handler := http.HandlerFunc(StatusHandler(store))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"website":{"id":1,"name":"Google","url":"https://www.google.com"},"status":"up","response_time":"100.00ms","headers":null,"ssl_info":null}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCheckHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/websites/check", strings.NewReader(`["https://www.google.com"]`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	var store Store = &mockStore{}
	handler := http.HandlerFunc(CheckHandler(store))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
