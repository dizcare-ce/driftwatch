package metrics_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/metrics"
)

func TestRecordRun_Clean(t *testing.T) {
	c := metrics.New()
	c.RecordRun(false, nil)

	if got := c.Runs.Load(); got != 1 {
		t.Fatalf("Runs: want 1, got %d", got)
	}
	if got := c.Clean.Load(); got != 1 {
		t.Fatalf("Clean: want 1, got %d", got)
	}
	if got := c.Drifted.Load(); got != 0 {
		t.Fatalf("Drifted: want 0, got %d", got)
	}
}

func TestRecordRun_Drifted(t *testing.T) {
	c := metrics.New()
	c.RecordRun(true, nil)

	if got := c.Drifted.Load(); got != 1 {
		t.Fatalf("Drifted: want 1, got %d", got)
	}
	if got := c.Clean.Load(); got != 0 {
		t.Fatalf("Clean: want 0, got %d", got)
	}
}

func TestRecordRun_Error(t *testing.T) {
	c := metrics.New()
	c.RecordRun(false, errors.New("boom"))

	if got := c.Errors.Load(); got != 1 {
		t.Fatalf("Errors: want 1, got %d", got)
	}
	if got := c.Clean.Load(); got != 0 {
		t.Fatalf("Clean should not increment on error, got %d", got)
	}
}

func TestLastRunAt_BeforeAnyRun(t *testing.T) {
	c := metrics.New()
	_, ok := c.LastRunAt()
	if ok {
		t.Fatal("expected ok=false before any run")
	}
}

func TestLastRunAt_AfterRun(t *testing.T) {
	c := metrics.New()
	c.RecordRun(false, nil)
	_, ok := c.LastRunAt()
	if !ok {
		t.Fatal("expected ok=true after a run")
	}
}

func TestWrite_ContainsAllFields(t *testing.T) {
	c := metrics.New()
	c.RecordRun(false, nil)
	c.RecordRun(true, nil)
	c.RecordRun(false, errors.New("err"))

	var buf bytes.Buffer
	if err := c.Write(&buf); err != nil {
		t.Fatalf("Write: %v", err)
	}

	line := buf.String()
	for _, want := range []string{"runs=3", "drifted=1", "clean=1", "errors=1", "last_run="} {
		if !strings.Contains(line, want) {
			t.Errorf("output missing %q; got: %s", want, line)
		}
	}
}
