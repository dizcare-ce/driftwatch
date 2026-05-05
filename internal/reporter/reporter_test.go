package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/reporter"
)

func makeResults(drifted bool) []drift.Result {
	diffs := []drift.Diff{}
	if drifted {
		diffs = []drift.Diff{
			{Field: "image", Expected: "nginx:1.24", Actual: "nginx:1.23"},
		}
	}
	return []drift.Result{
		{Service: "api", Drifted: drifted, Diffs: diffs},
	}
}

func TestWrite_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	if err := r.Write(makeResults(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "[OK] api") {
		t.Errorf("expected OK status in output, got:\n%s", output)
	}
	if strings.Contains(output, "DRIFTED") {
		t.Errorf("did not expect DRIFTED in output")
	}
}

func TestWrite_TextDrifted(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	if err := r.Write(makeResults(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "[DRIFTED] api") {
		t.Errorf("expected DRIFTED status in output, got:\n%s", output)
	}
	if !strings.Contains(output, "image") {
		t.Errorf("expected field name in output")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatJSON, &buf)
	if err := r.Write(makeResults(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if report.DriftCount != 1 {
		t.Errorf("expected drift_count=1, got %d", report.DriftCount)
	}
	if len(report.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(report.Results))
	}
}

func TestWrite_DefaultsToStdout(t *testing.T) {
	// Ensure New with nil writer does not panic.
	r := reporter.New(reporter.FormatText, nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
