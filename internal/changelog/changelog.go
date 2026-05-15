// Package changelog records a human-readable log of drift state transitions
// for each service, allowing operators to review when a service first drifted,
// when it was resolved, and how many times it has changed state.
package changelog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry describes a single state-transition event for a service.
type Entry struct {
	Service   string    `json:"service"`
	State     string    `json:"state"` // "drifted" | "clean"
	DiffCount int       `json:"diff_count"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Log persists changelog entries to a NDJSON file per service.
type Log struct {
	mu  sync.Mutex
	dir string
}

// New returns a Log that stores entries under dir, creating it if absent.
func New(dir string) (*Log, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("changelog: create dir: %w", err)
	}
	return &Log{dir: dir}, nil
}

// Record appends an entry for the named service.
func (l *Log) Record(e Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	e.RecordedAt = time.Now().UTC()

	f, err := os.OpenFile(l.path(e.Service), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("changelog: open: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(e); err != nil {
		return fmt.Errorf("changelog: encode: %w", err)
	}
	return nil
}

// ReadAll returns all recorded entries for the named service, oldest first.
// If no entries exist, a nil slice and no error are returned.
func (l *Log) ReadAll(service string) ([]Entry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := os.ReadFile(l.path(service))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("changelog: read: %w", err)
	}

	var entries []Entry
	dec := json.NewDecoder(bytesReader(data))
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("changelog: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (l *Log) path(service string) string {
	return filepath.Join(l.dir, service+".ndjson")
}
