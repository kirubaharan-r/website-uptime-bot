package api

import (
	"fmt"
	"net/http"
	"uptime-monitor/pkg/api/handlers"
	"uptime-monitor/pkg/storage"
)

// Server is the API server.
type Server struct {
	addr  string
	store storage.Store
}

// NewServer creates a new API server.
func NewServer(addr string, store storage.Store) *Server {
	return &Server{
		addr:  addr,
		store: store,
	}
}

// Start starts the web server.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/websites/check", handlers.CheckHandler(s.store))
	mux.HandleFunc("/api/websites/status", handlers.StatusHandler(s.store))

	fmt.Printf("Starting web server on %s\n", s.addr)
	return http.ListenAndServe(s.addr, mux)
}
