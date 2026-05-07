package suppress_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"driftwatch/internal/suppress"
)

func newSuppressor(t *testing.T) *suppress.Suppressor {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "suppress")
	s, err := suppress.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestSuppress_And_IsActive(t *testing.T) {
	s := newSuppressor(t)
	if err := s.Suppress("api", time.Hour); err != nil {
		t.Fatalf("Suppress: %v", err)
	}
	if !s.IsActive("api") {
		t.Fatal("expected suppression to be active")
	}
}

func TestIsActive_UnknownService_ReturnsFalse(t *testing.T) {
	s := newSuppressor(t)
	if s.IsActive("unknown") {
		t.Fatal("expected inactive for unknown service")
	}
}

func TestSuppress_AlreadySuppressed_ReturnsError(t *testing.T) {
	s := newSuppressor(t)
	if err := s.Suppress("api", time.Hour); err != nil {
		t.Fatalf("first Suppress: %v", err)
	}
	if err := s.Suppress("api", time.Hour); err != suppress.ErrAlreadySuppressed {
		t.Fatalf("expected ErrAlreadySuppressed, got %v", err)
	}
}

func TestLift_RemovesSuppression(t *testing.T) {
	s := newSuppressor(t)
	_ = s.Suppress("api", time.Hour)
	if err := s.Lift("api"); err != nil {
		t.Fatalf("Lift: %v", err)
	}
	if s.IsActive("api") {
		t.Fatal("expected suppression to be inactive after Lift")
	}
}

func TestLift_MissingEntry_NoError(t *testing.T) {
	s := newSuppressor(t)
	if err := s.Lift("nonexistent"); err != nil {
		t.Fatalf("Lift on missing entry: %v", err)
	}
}

func TestSuppress_ExpiredEntry_IsInactive(t *testing.T) {
	s := newSuppressor(t)
	if err := s.Suppress("api", -time.Second); err != nil {
		t.Fatalf("Suppress: %v", err)
	}
	if s.IsActive("api") {
		t.Fatal("expected expired suppression to be inactive")
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "deep", "suppress")
	if _, err := suppress.New(dir); err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
}
