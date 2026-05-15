package trendanalyzer_test

import (
	"os"
	"testing"
	"time"

	"driftwatch/internal/trendanalyzer"
)

func newStore(t *testing.T) *trendanalyzer.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := trendanalyzer.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return st
}

func TestStore_AppendAndLoad_RoundTrip(t *testing.T) {
	st := newStore(t)
	snap := trendanalyzer.Snapshot{At: epoch, Drift: 5}
	if err := st.Append(snap); err != nil {
		t.Fatalf("Append: %v", err)
	}
	loaded, err := st.Load(0)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(loaded))
	}
	if loaded[0].Drift != 5 {
		t.Fatalf("expected drift 5, got %d", loaded[0].Drift)
	}
}

func TestStore_Load_OrderedOldestFirst(t *testing.T) {
	st := newStore(t)
	for i := 0; i < 5; i++ {
		err := st.Append(trendanalyzer.Snapshot{
			At:    epoch.Add(time.Duration(i) * time.Second),
			Drift: i,
		})
		if err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	loaded, err := st.Load(0)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for i, s := range loaded {
		if s.Drift != i {
			t.Fatalf("index %d: expected drift %d, got %d", i, i, s.Drift)
		}
	}
}

func TestStore_Load_RespectsLimit(t *testing.T) {
	st := newStore(t)
	for i := 0; i < 6; i++ {
		_ = st.Append(trendanalyzer.Snapshot{
			At:    epoch.Add(time.Duration(i) * time.Second),
			Drift: i,
		})
	}
	loaded, err := st.Load(3)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 3 {
		t.Fatalf("expected 3, got %d", len(loaded))
	}
	// should be the last 3 (drift 3,4,5)
	if loaded[0].Drift != 3 {
		t.Fatalf("expected first drift 3, got %d", loaded[0].Drift)
	}
}

func TestStore_Load_MissingDir_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	st, _ := trendanalyzer.NewStore(dir)
	// remove directory to simulate missing
	_ = os.RemoveAll(dir)
	loaded, err := st.Load(0)
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got %v", err)
	}
	if loaded != nil {
		t.Fatalf("expected nil snapshots, got %v", loaded)
	}
}

func TestStore_NewStore_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/trend"
	_, err := trendanalyzer.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}
