// Package healthcheck exposes a simple HTTP endpoint that reports the
// current health of the driftwatch daemon, including the last run time
// and whether any drift was detected.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/driftwatch/driftwatch/internal/metrics"
)

// Status is the JSON body returned by the health endpoint.
type Status struct {
	OK        bool      `json:"ok"`
	LastRunAt time.Time `json:"last_run_at,omitempty"`
	Drifted   int       `json:"drifted_services"`
	Total     int       `json:"total_services"`
	Message   string    `json:"message"`
}

// Handler returns an http.HandlerFunc that serves the health status.
func Handler(m *metrics.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lastRun := m.LastRunAt()
		total, drifted := m.Counts()

		s := Status{
			OK:        drifted == 0,
			LastRunAt: lastRun,
			Drifted:   drifted,
			Total:     total,
		}

		switch {
		case lastRun.IsZero():
			s.Message = "no run completed yet"
		case drifted > 0:
			s.Message = "drift detected"
		default:
			s.Message = "all services in sync"
		}

		w.Header().Set("Content-Type", "application/json")
		if !s.OK && !lastRun.IsZero() {
			w.WriteHeader(http.StatusConflict)
		}
		_ = json.NewEncoder(w).Encode(s)
	}
}

// Server wraps an http.Server configured to serve the health endpoint.
type Server struct {
	server *http.Server
}

// New creates a Server listening on addr (e.g. ":8080").
func New(addr string, m *metrics.Metrics) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", Handler(m))
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Close shuts down the server immediately.
func (s *Server) Close() error {
	return s.server.Close()
}
