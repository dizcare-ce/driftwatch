// Package history tracks drift check results over time,
// allowing driftwatch to surface trends and repeated violations.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// Entry records a single drift check run.
type Entry struct {
	Timestamp time.Time     `json:"timestamp"`
	Results   []drift.Result `json:"results"`
}

// History persists and retrieves drift check entries.
type History struct {
	dir string
}

// New returns a History that stores entries under dir.
func New(dir string) (*History, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	return &History{dir: dir}, nil
}

// Record writes an entry for the current run to disk.
func (h *History) Record(results []drift.Result) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Results:   results,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	name := entry.Timestamp.Format("20060102T150405Z") + ".json"
	path := filepath.Join(h.dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write: %w", err)
	}
	return nil
}

// Last returns the most recent n entries, oldest first.
func (h *History) Last(n int) ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(h.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	sort.Strings(matches)
	if len(matches) > n {
		matches = matches[len(matches)-n:]
	}
	entries := make([]Entry, 0, len(matches))
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("history: read %s: %w", p, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: unmarshal %s: %w", p, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
