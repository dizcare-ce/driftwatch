// Package snapshot provides functionality for persisting and loading drift
// results to disk, enabling comparison between runs to detect new drift.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a persisted snapshot of a single service drift result.
type Entry struct {
	Service   string      `json:"service"`
	CapturedAt time.Time  `json:"captured_at"`
	Result    drift.Result `json:"result"`
}

// Store handles reading and writing snapshots to a directory.
type Store struct {
	dir string
}

// New creates a new Store backed by the given directory.
// The directory is created if it does not exist.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir %q: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Save writes the drift result for the named service to disk.
func (s *Store) Save(service string, result drift.Result) error {
	entry := Entry{
		Service:    service,
		CapturedAt: time.Now().UTC(),
		Result:     result,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal %q: %w", service, err)
	}
	p := s.path(service)
	if err := os.WriteFile(p, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %q: %w", p, err)
	}
	return nil
}

// Load reads the most recent persisted drift result for the named service.
// Returns os.ErrNotExist if no snapshot has been saved yet.
func (s *Store) Load(service string) (Entry, error) {
	p := s.path(service)
	data, err := os.ReadFile(p)
	if err != nil {
		return Entry{}, fmt.Errorf("snapshot: read %q: %w", p, err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, fmt.Errorf("snapshot: unmarshal %q: %w", p, err)
	}
	return entry, nil
}

func (s *Store) path(service string) string {
	return filepath.Join(s.dir, service+".json")
}
