package notifier_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/notifier"
)

func makeResult(service string, diffs []drift.Diff) drift.Result {
	return drift.Result{Service: service, Diffs: diffs}
}

func TestNotify_LevelNone_WritesNothing(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(notifier.LevelNone, &buf)
	results := []drift.Result{
		makeResult("api", []drift.Diff{{Field: "image", Want: "v1", Got: "v2"}}),
	}
	count, err := n.Notify(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 notifications, got %d", count)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestNotify_LevelDrift_SkipsClean(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(notifier.LevelDrift, &buf)
	results := []drift.Result{
		makeResult("api", nil),
		makeResult("worker", []drift.Diff{{Field: "replicas", Want: "3", Got: "1"}}),
	}
	count, err := n.Notify(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 notification, got %d", count)
	}
	if !strings.Contains(buf.String(), "worker") {
		t.Errorf("expected worker in output, got %q", buf.String())
	}
}

func TestNotify_LevelAll_IncludesClean(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(notifier.LevelAll, &buf)
	results := []drift.Result{
		makeResult("api", nil),
		makeResult("worker", nil),
	}
	count, err := n.Notify(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 notifications, got %d", count)
	}
}

func TestParseLevel_Valid(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  notifier.Level
	}{
		{"all", notifier.LevelAll},
		{"drift", notifier.LevelDrift},
		{"none", notifier.LevelNone},
		{"DRIFT", notifier.LevelDrift},
	} {
		got, err := notifier.ParseLevel(tc.input)
		if err != nil {
			t.Errorf("ParseLevel(%q) error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseLevel(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	_, err := notifier.ParseLevel("verbose")
	if err == nil {
		t.Error("expected error for unknown level, got nil")
	}
}
