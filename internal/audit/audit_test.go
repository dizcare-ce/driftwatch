package audit_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/audit"
)

func TestRecord_And_ReadAll_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	log, err := audit.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Millisecond)
	e := audit.Entry{
		Timestamp:   now,
		ServicesRun: 5,
		DriftCount:  2,
		ErrorCount:  0,
		DurationMs:  142,
	}
	if err := log.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := log.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	got := entries[0]
	if got.ServicesRun != 5 || got.DriftCount != 2 || got.DurationMs != 142 {
		t.Errorf("entry mismatch: %+v", got)
	}
}

func TestRecord_AppendsMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	log, _ := audit.New(dir)

	for i := 0; i < 3; i++ {
		_ = log.Record(audit.Entry{ServicesRun: i + 1, DurationMs: int64(i * 10)})
	}

	entries, err := log.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
	if entries[2].ServicesRun != 3 {
		t.Errorf("last entry ServicesRun want 3, got %d", entries[2].ServicesRun)
	}
}

func TestReadAll_MissingFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	log, _ := audit.New(dir)

	// Remove the file that New may have created (it doesn't, but be safe)
	_ = os.Remove(dir + "/audit.jsonl")

	entries, err := log.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("want empty slice, got %d entries", len(entries))
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/audit"
	_, err := audit.New(dir)
	if err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("directory was not created: %s", dir)
	}
}
