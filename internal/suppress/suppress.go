// Package suppress provides a mechanism to temporarily silence drift
// notifications for a named service. A suppression is stored with an
// expiry time so that alerts automatically resume after the window closes.
package suppress

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ErrAlreadySuppressed is returned when a suppression already exists and has
// not yet expired.
var ErrAlreadySuppressed = errors.New("suppress: service is already suppressed")

type entry struct {
	ExpiresAt time.Time `json:"expires_at"`
}

// Suppressor persists per-service suppression windows to disk.
type Suppressor struct {
	mu  sync.Mutex
	dir string
}

// New returns a Suppressor that stores entries under dir, creating it when
// absent.
func New(dir string) (*Suppressor, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Suppressor{dir: dir}, nil
}

// Suppress silences alerts for service for the given duration. It returns
// ErrAlreadySuppressed when an active suppression already exists.
func (s *Suppressor) Suppress(service string, d time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if active, _ := s.isActive(service); active {
		return ErrAlreadySuppressed
	}
	e := entry{ExpiresAt: time.Now().Add(d)}
	return s.write(service, e)
}

// IsActive reports whether service currently has an unexpired suppression.
func (s *Suppressor) IsActive(service string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	active, _ := s.isActive(service)
	return active
}

// Lift removes any suppression for service, whether expired or not.
func (s *Suppressor) Lift(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(s.path(service))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// isActive must be called with s.mu held.
func (s *Suppressor) isActive(service string) (bool, error) {
	e, err := s.read(service)
	if err != nil {
		return false, err
	}
	return time.Now().Before(e.ExpiresAt), nil
}

func (s *Suppressor) path(service string) string {
	return filepath.Join(s.dir, service+".json")
}

func (s *Suppressor) write(service string, e entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(service), b, 0o644)
}

func (s *Suppressor) read(service string) (entry, error) {
	b, err := os.ReadFile(s.path(service))
	if err != nil {
		return entry{}, err
	}
	var e entry
	return e, json.Unmarshal(b, &e)
}
