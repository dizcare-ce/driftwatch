// Package deadletter stores drift check failures that could not be
// processed or delivered so they can be inspected and replayed later.
package deadletter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry represents a single failed event stored in the dead-letter queue.
type Entry struct {
	Service   string    `json:"service"`
	Reason    string    `json:"reason"`
	Payload   string    `json:"payload"`
	FailedAt  time.Time `json:"failed_at"`
	Attempts  int       `json:"attempts"`
}

// Queue persists failed entries to a directory on disk.
type Queue struct {
	mu  sync.Mutex
	dir string
}

// New creates a Queue rooted at dir, creating the directory if absent.
func New(dir string) (*Queue, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("deadletter: create dir: %w", err)
	}
	return &Queue{dir: dir}, nil
}

// Push appends an entry to the queue.
func (q *Queue) Push(e Entry) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if e.FailedAt.IsZero() {
		e.FailedAt = time.Now().UTC()
	}

	name := fmt.Sprintf("%d_%s.json", e.FailedAt.UnixNano(), sanitise(e.Service))
	path := filepath.Join(q.dir, name)

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("deadletter: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// All returns every entry currently in the queue, ordered by filename.
func (q *Queue) All() ([]Entry, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	glob := filepath.Join(q.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("deadletter: glob: %w", err)
	}

	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("deadletter: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("deadletter: unmarshal %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Purge removes all entries from the queue.
func (q *Queue) Purge() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	glob := filepath.Join(q.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return fmt.Errorf("deadletter: glob: %w", err)
	}
	for _, m := range matches {
		if err := os.Remove(m); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("deadletter: remove %s: %w", m, err)
		}
	}
	return nil
}

// sanitise replaces characters that are unsafe in filenames.
func sanitise(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			out[i] = c
		} else {
			out[i] = '_'
		}
	}
	return string(out)
}
