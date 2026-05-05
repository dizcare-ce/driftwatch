package watcher_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/watcher"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestRun_DetectsNewFile(t *testing.T) {
	dir := t.TempDir()
	var calls atomic.Int32
	fw := watcher.New(dir, 20*time.Millisecond, func(_ string) { calls.Add(1) })

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go fw.Run(ctx) //nolint:errcheck
	time.Sleep(40 * time.Millisecond)

	writeFile(t, dir, "svc.yaml", "name: svc")
	time.Sleep(80 * time.Millisecond)

	if calls.Load() == 0 {
		t.Error("expected onChange to be called for new file")
	}
}

func TestRun_DetectsModifiedFile(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "svc.yaml", "name: svc")

	var calls atomic.Int32
	fw := watcher.New(dir, 20*time.Millisecond, func(_ string) { calls.Add(1) })

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go fw.Run(ctx) //nolint:errcheck
	time.Sleep(50 * time.Millisecond)

	// Modify file after initial snapshot
	if err := os.WriteFile(p, []byte("name: svc\nenv: prod"), 0o644); err != nil {
		t.Fatalf("modify: %v", err)
	}
	time.Sleep(80 * time.Millisecond)

	if calls.Load() == 0 {
		t.Error("expected onChange to be called after file modification")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	dir := t.TempDir()
	fw := watcher.New(dir, 10*time.Millisecond, func(_ string) {})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- fw.Run(ctx) }()

	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("watcher did not stop after context cancel")
	}
}

func TestRun_NoCallbackForUnchangedFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "svc.yaml", "name: svc")

	var calls atomic.Int32
	fw := watcher.New(dir, 20*time.Millisecond, func(_ string) { calls.Add(1) })

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	fw.Run(ctx) //nolint:errcheck

	if calls.Load() != 0 {
		t.Errorf("expected no onChange calls for unchanged file, got %d", calls.Load())
	}
}
