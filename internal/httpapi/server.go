// Package httpapi provides a lightweight HTTP server exposing portwatch
// runtime state via a simple JSON API.
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Server is a minimal HTTP API server.
type Server struct {
	addr    string
	reg     *metrics.Registry
	httpSrv *http.Server
}

// NewServer creates a new Server listening on addr.
func NewServer(addr string, reg *metrics.Registry) *Server {
	s := &Server{addr: addr, reg: reg}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/metrics", s.handleMetrics)
	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return s
}

// Start begins serving in a background goroutine.
func (s *Server) Start() error {
	go func() { _ = s.httpSrv.ListenAndServe() }()
	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = s.reg.WriteTo(w)
}
