package changelog

import (
	"os"
	"path/filepath"
	"testing"
)

func newLog(t *testing.T) *Log {
	t.Helper()
	dir := t.TempDir()
	l, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return l
}

func TestRecord_And_ReadAll_RoundTrip(t *testing.T) {
	l := newLog(t)

	e := Entry{Service: "svc-a", State: "drifted", DiffCount: 2}
	if err := l.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := l.ReadAll("svc-a")
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].State != "drifted" || entries[0].DiffCount != 2 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
	if entries[0].RecordedAt.IsZero() {
		t.Error("RecordedAt should be set")
	}
}

func TestRecord_AppendsMultipleEntries(t *testing.T) {
	l := newLog(t)

	for _, state := range []string{"drifted", "clean", "drifted"} {
		if err := l.Record(Entry{Service: "svc-b", State: state}); err != nil {
			t.Fatalf("Record(%s): %v", state, err)
		}
	}

	entries, err := l.ReadAll("svc-b")
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
	if entries[0].State != "drifted" || entries[1].State != "clean" || entries[2].State != "drifted" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestReadAll_MissingService_ReturnsNil(t *testing.T) {
	l := newLog(t)

	entries, err := l.ReadAll("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("want nil, got %v", entries)
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "changelog")
	if _, err := New(dir); err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

func TestRecord_ServicesAreIsolated(t *testing.T) {
	l := newLog(t)

	_ = l.Record(Entry{Service: "alpha", State: "drifted", DiffCount: 1})
	_ = l.Record(Entry{Service: "beta", State: "clean", DiffCount: 0})

	alpha, _ := l.ReadAll("alpha")
	beta, _ := l.ReadAll("beta")

	if len(alpha) != 1 || alpha[0].Service != "alpha" {
		t.Errorf("alpha contaminated: %v", alpha)
	}
	if len(beta) != 1 || beta[0].Service != "beta" {
		t.Errorf("beta contaminated: %v", beta)
	}
}
