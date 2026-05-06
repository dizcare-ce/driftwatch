package history_test

import (
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
	"github.com/driftwatch/driftwatch/internal/history"
)

func makeResults(drifted bool) []drift.Result {
	diffs := []drift.Diff{}
	if drifted {
		diffs = []drift.Diff{{Field: "replicas", Want: "3", Got: "1"}}
	}
	return []drift.Result{
		{ServiceName: "api", Diffs: diffs},
	}
}

func TestRecord_And_Last_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	h, err := history.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	results := makeResults(true)
	if err := h.Record(results); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := h.Last(10)
	if err != nil {
		t.Fatalf("Last: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if len(entries[0].Results) != 1 {
		t.Fatalf("want 1 result, got %d", len(entries[0].Results))
	}
	if entries[0].Results[0].ServiceName != "api" {
		t.Errorf("want service api, got %s", entries[0].Results[0].ServiceName)
	}
}

func TestLast_LimitsEntries(t *testing.T) {
	dir := t.TempDir()
	h, err := history.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Millisecond) // ensure distinct timestamps
		if err := h.Record(makeResults(false)); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	entries, err := h.Last(3)
	if err != nil {
		t.Fatalf("Last: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("want 3 entries, got %d", len(entries))
	}
}

func TestLast_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	h, err := history.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	entries, err := h.Last(5)
	if err != nil {
		t.Fatalf("Last: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("want 0 entries, got %d", len(entries))
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/history"
	_, err := history.New(dir)
	if err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
}
