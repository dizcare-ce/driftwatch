package staledetector_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"driftwatch/internal/staledetector"
)

func TestReport_NoEntries_PrintsCurrentMessage(t *testing.T) {
	var buf bytes.Buffer
	if err := staledetector.Report(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "current") {
		t.Errorf("expected 'current' in output, got: %s", buf.String())
	}
}

func TestReport_WithEntries_ContainsServiceName(t *testing.T) {
	now := time.Now()
	entries := []staledetector.Entry{
		{Service: "payments", LastSeen: now.Add(-15 * time.Minute), Age: 15 * time.Minute},
	}
	var buf bytes.Buffer
	if err := staledetector.Report(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "payments") {
		t.Errorf("expected service name in output, got: %s", out)
	}
}

func TestReport_WithEntries_ContainsHeader(t *testing.T) {
	entries := []staledetector.Entry{
		{Service: "api", LastSeen: time.Time{}, Age: 30 * time.Minute},
	}
	var buf bytes.Buffer
	_ = staledetector.Report(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "STALE SERVICES") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestReport_ZeroLastSeen_PrintsNever(t *testing.T) {
	entries := []staledetector.Entry{
		{Service: "ghost", LastSeen: time.Time{}, Age: 1 * time.Hour},
	}
	var buf bytes.Buffer
	_ = staledetector.Report(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "never") {
		t.Errorf("expected 'never' in output for zero timestamp, got: %s", out)
	}
}

func TestReport_AgeFormatting_HoursMinutesSeconds(t *testing.T) {
	entries := []staledetector.Entry{
		{Service: "svc", LastSeen: time.Now().Add(-2*time.Hour - 5*time.Minute - 3*time.Second), Age: 2*time.Hour + 5*time.Minute + 3*time.Second},
	}
	var buf bytes.Buffer
	_ = staledetector.Report(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "2h") {
		t.Errorf("expected hours in age format, got: %s", out)
	}
}
