package trendanalyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Store persists Snapshots to disk so trend analysis survives restarts.
type Store struct {
	dir string
}

// NewStore returns a Store that writes snapshot files under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("trendanalyzer: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

type record struct {
	At    time.Time `json:"at"`
	Drift int       `json:"drift"`
}

// Append writes a single Snapshot entry to the store.
func (s *Store) Append(snap Snapshot) error {
	name := fmt.Sprintf("%d.json", snap.At.UnixNano())
	path := filepath.Join(s.dir, name)
	b, err := json.Marshal(record{At: snap.At, Drift: snap.Drift})
	if err != nil {
		return fmt.Errorf("trendanalyzer: marshal: %w", err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("trendanalyzer: write: %w", err)
	}
	return nil
}

// Load returns all stored Snapshots ordered oldest-first.
// At most limit entries are returned; pass 0 for no limit.
func (s *Store) Load(limit int) ([]Snapshot, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("trendanalyzer: read dir: %w", err)
	}

	// entries are already sorted lexicographically by filename (unix nano)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	if limit > 0 && len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}

	var snaps []Snapshot
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("trendanalyzer: read file %s: %w", e.Name(), err)
		}
		var rec record
		if err := json.Unmarshal(b, &rec); err != nil {
			return nil, fmt.Errorf("trendanalyzer: unmarshal %s: %w", e.Name(), err)
		}
		snaps = append(snaps, Snapshot{At: rec.At, Drift: rec.Drift})
	}
	return snaps, nil
}
