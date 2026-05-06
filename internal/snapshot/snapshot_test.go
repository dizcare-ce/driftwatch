package snapshot_test

import (
	"errors"
	"os"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/snapshot"
)

func makeResult(service string, drifted bool) drift.Result {
	diffs := []drift.Diff{}
	if drifted {
		diffs = append(diffs, drift.Diff{
			Field:    "replicas",
			Expected: "3",
			Actual:   "1",
		})
	}
	return drift.Result{Service: service, Diffs: diffs}
}

func TestSave_And_Load_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result := makeResult("api", true)
	if err := store.Save("api", result); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entry, err := store.Load("api")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if entry.Service != "api" {
		t.Errorf("Service: got %q, want %q", entry.Service, "api")
	}
	if len(entry.Result.Diffs) != 1 {
		t.Errorf("Diffs len: got %d, want 1", len(entry.Result.Diffs))
	}
	if entry.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestLoad_MissingSnapshot(t *testing.T) {
	dir := t.TempDir()
	store, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got: %v", err)
	}
}

func TestSave_OverwritesPreviousSnapshot(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)

	_ = store.Save("svc", makeResult("svc", true))
	_ = store.Save("svc", makeResult("svc", false))

	entry, err := store.Load("svc")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entry.Result.Diffs) != 0 {
		t.Errorf("expected no diffs after overwrite, got %d", len(entry.Result.Diffs))
	}
}

func TestNew_CreatesDirectoryIfAbsent(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/snapshots"

	_, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
