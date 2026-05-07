package baseline_test

import (
	"testing"

	"github.com/driftwatch/internal/baseline"
)

func TestSet_And_Get_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s, err := baseline.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	in := baseline.Entry{Service: "api", Checksum: "abc123", Note: "initial"}
	if err := s.Set(in); err != nil {
		t.Fatalf("Set: %v", err)
	}

	out, ok, err := s.Get("api")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if out.Service != in.Service || out.Checksum != in.Checksum || out.Note != in.Note {
		t.Errorf("round-trip mismatch: got %+v", out)
	}
	if out.ApprovedAt.IsZero() {
		t.Error("ApprovedAt should be populated")
	}
}

func TestGet_MissingEntry_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	s, _ := baseline.New(dir)

	_, ok, err := s.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected ok=false for missing service")
	}
}

func TestSet_OverwritesPreviousEntry(t *testing.T) {
	dir := t.TempDir()
	s, _ := baseline.New(dir)

	s.Set(baseline.Entry{Service: "svc", Checksum: "old"})
	s.Set(baseline.Entry{Service: "svc", Checksum: "new"})

	out, ok, err := s.Get("svc")
	if err != nil || !ok {
		t.Fatalf("Get: ok=%v err=%v", ok, err)
	}
	if out.Checksum != "new" {
		t.Errorf("expected checksum 'new', got %q", out.Checksum)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	dir := t.TempDir()
	s, _ := baseline.New(dir)

	s.Set(baseline.Entry{Service: "svc", Checksum: "abc"})
	if err := s.Delete("svc"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, ok, err := s.Get("svc")
	if err != nil || ok {
		t.Errorf("expected entry gone after delete: ok=%v err=%v", ok, err)
	}
}

func TestDelete_MissingEntry_NoError(t *testing.T) {
	dir := t.TempDir()
	s, _ := baseline.New(dir)

	if err := s.Delete("ghost"); err != nil {
		t.Errorf("expected no error for missing entry, got: %v", err)
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := t.TempDir() + "/nested/baseline"
	_, err := baseline.New(dir)
	if err != nil {
		t.Fatalf("New should create nested dir: %v", err)
	}
}
