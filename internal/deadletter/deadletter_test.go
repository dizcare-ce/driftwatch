package deadletter_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"driftwatch/internal/deadletter"
)

func newQueue(t *testing.T) *deadletter.Queue {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "dlq")
	q, err := deadletter.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return q
}

func TestPush_And_All_RoundTrip(t *testing.T) {
	q := newQueue(t)

	e := deadletter.Entry{
		Service:  "api",
		Reason:   "timeout",
		Payload:  `{"env":"prod"}`,
		Attempts: 1,
	}
	if err := q.Push(e); err != nil {
		t.Fatalf("Push: %v", err)
	}

	entries, err := q.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	got := entries[0]
	if got.Service != "api" || got.Reason != "timeout" || got.Attempts != 1 {
		t.Errorf("unexpected entry: %+v", got)
	}
	if got.FailedAt.IsZero() {
		t.Error("FailedAt should be set")
	}
}

func TestPush_MultipleEntries_AllReturned(t *testing.T) {
	q := newQueue(t)

	for i, svc := range []string{"alpha", "beta", "gamma"} {
		err := q.Push(deadletter.Entry{
			Service:  svc,
			Reason:   "error",
			Attempts: i + 1,
			FailedAt: time.Now().UTC().Add(time.Duration(i) * time.Millisecond),
		})
		if err != nil {
			t.Fatalf("Push %s: %v", svc, err)
		}
	}

	entries, err := q.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
}

func TestAll_EmptyQueue_ReturnsNil(t *testing.T) {
	q := newQueue(t)

	entries, err := q.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("want 0 entries, got %d", len(entries))
	}
}

func TestPurge_RemovesAllEntries(t *testing.T) {
	q := newQueue(t)

	for _, svc := range []string{"x", "y"} {
		if err := q.Push(deadletter.Entry{Service: svc, Reason: "err"}); err != nil {
			t.Fatalf("Push: %v", err)
		}
	}

	if err := q.Purge(); err != nil {
		t.Fatalf("Purge: %v", err)
	}

	entries, err := q.All()
	if err != nil {
		t.Fatalf("All after purge: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("want 0 entries after purge, got %d", len(entries))
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dlq")
	if _, err := deadletter.New(dir); err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}
