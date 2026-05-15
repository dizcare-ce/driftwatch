package scorecard_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"driftwatch/internal/drift"
	"driftwatch/internal/scorecard"
)

func buildEntries(t *testing.T) []scorecard.Entry {
	t.Helper()
	return scorecard.Build([]drift.Result{
		makeResult("api", 3, 10),
		makeResult("worker", 0, 5),
	})
}

func TestFormat_Text_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	if err := scorecard.Format(&buf, buildEntries(t), "text"); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"SERVICE", "SCORE", "GRADE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormat_Text_ContainsServiceName(t *testing.T) {
	var buf bytes.Buffer
	_ = scorecard.Format(&buf, buildEntries(t), "text")
	if !strings.Contains(buf.String(), "api") {
		t.Error("expected service name 'api' in text output")
	}
}

func TestFormat_JSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := scorecard.Format(&buf, buildEntries(t), "json"); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	var out []scorecard.Entry
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestFormat_Default_FallsBackToText(t *testing.T) {
	var buf bytes.Buffer
	if err := scorecard.Format(&buf, buildEntries(t), ""); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Error("expected text output for empty format string")
	}
}

func TestFormat_EmptyEntries_WritesHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	if err := scorecard.Format(&buf, nil, "text"); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Error("expected header even for empty entries")
	}
}
