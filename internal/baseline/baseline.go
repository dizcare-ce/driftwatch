// Package baseline stores and retrieves approved drift baselines.
// A baseline records a known-acceptable drift state so that repeated
// alerts for the same drift are suppressed until the definition changes.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry is a persisted baseline record for a single service.
type Entry struct {
	Service   string    `json:"service"`
	Checksum  string    `json:"checksum"`
	ApprovedAt time.Time `json:"approved_at"`
	Note      string    `json:"note,omitempty"`
}

// Store manages baseline entries on disk.
type Store struct {
	mu  sync.RWMutex
	dir string
}

// New returns a Store rooted at dir, creating it if necessary.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(service string) string {
	return filepath.Join(s.dir, service+".json")
}

// Set writes or replaces the baseline entry for service.
func (s *Store) Set(e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e.ApprovedAt = time.Now().UTC()
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(e.Service), b, 0o644)
}

// Get returns the baseline entry for service, or false if none exists.
func (s *Store) Get(service string) (Entry, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, err := os.ReadFile(s.path(service))
	if errors.Is(err, os.ErrNotExist) {
		return Entry{}, false, nil
	}
	if err != nil {
		return Entry{}, false, err
	}
	var e Entry
	if err := json.Unmarshal(b, &e); err != nil {
		return Entry{}, false, err
	}
	return e, true, nil
}

// Delete removes the baseline for service. A missing entry is not an error.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(s.path(service))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
