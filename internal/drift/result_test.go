package drift_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func TestResult_Summary_NoDrift(t *testing.T) {
	r := drift.Result{Service: "auth", Drifted: false}
	got := r.Summary()
	if !strings.Contains(got, "no drift") {
		t.Errorf("expected 'no drift' in summary, got: %s", got)
	}
}

func TestResult_Summary_SingleDiff(t *testing.T) {
	r := drift.Result{
		Service: "auth",
		Drifted: true,
		Diffs:   []drift.Diff{{Field: "image", Expected: "v1", Actual: "v2"}},
	}
	got := r.Summary()
	if !strings.Contains(got, "1 field") {
		t.Errorf("expected '1 field' in summary, got: %s", got)
	}
}

func TestResult_Summary_MultipleDiffs(t *testing.T) {
	r := drift.Result{
		Service: "auth",
		Drifted: true,
		Diffs: []drift.Diff{
			{Field: "image", Expected: "v1", Actual: "v2"},
			{Field: "replicas", Expected: "3", Actual: "1"},
		},
	}
	got := r.Summary()
	if !strings.Contains(got, "2 fields") {
		t.Errorf("expected '2 fields' in summary, got: %s", got)
	}
}

func TestDiff_Fields(t *testing.T) {
	d := drift.Diff{Field: "tag", Expected: "stable", Actual: "latest"}
	if d.Field != "tag" || d.Expected != "stable" || d.Actual != "latest" {
		t.Errorf("unexpected diff values: %+v", d)
	}
}
