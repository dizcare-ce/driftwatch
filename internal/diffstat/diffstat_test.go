package diffstat_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/diffstat"
	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeResult(name string, diffs []drift.Diff) drift.Result {
	return drift.Result{Service: name, Diffs: diffs}
}

func makeDiff(field, want, got string) drift.Diff {
	return drift.Diff{Field: field, Expected: want, Actual: got}
}

func TestCompute_Empty(t *testing.T) {
	s := diffstat.Compute(nil)
	if s.Total != 0 || s.Clean != 0 || s.Drifted != 0 || s.Diffs != 0 {
		t.Fatalf("expected zero stats, got %+v", s)
	}
}

func TestCompute_AllClean(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", nil),
		makeResult("svc-b", nil),
	}
	s := diffstat.Compute(results)
	if s.Total != 2 || s.Clean != 2 || s.Drifted != 0 || s.Diffs != 0 {
		t.Fatalf("unexpected stats: %+v", s)
	}
}

func TestCompute_SomeDrifted(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", nil),
		makeResult("svc-b", []drift.Diff{makeDiff("image", "v1", "v2")}),
		makeResult("svc-c", []drift.Diff{makeDiff("replicas", "3", "1"), makeDiff("port", "8080", "9090")}),
	}
	s := diffstat.Compute(results)
	if s.Total != 3 {
		t.Fatalf("expected Total=3, got %d", s.Total)
	}
	if s.Clean != 1 {
		t.Fatalf("expected Clean=1, got %d", s.Clean)
	}
	if s.Drifted != 2 {
		t.Fatalf("expected Drifted=2, got %d", s.Drifted)
	}
	if s.Diffs != 3 {
		t.Fatalf("expected Diffs=3, got %d", s.Diffs)
	}
	if s.AvgDiffs != 1.5 {
		t.Fatalf("expected AvgDiffs=1.5, got %f", s.AvgDiffs)
	}
}

func TestDriftRate_NoResults(t *testing.T) {
	rate := diffstat.DriftRate(diffstat.Stats{})
	if rate != 0 {
		t.Fatalf("expected 0, got %f", rate)
	}
}

func TestDriftRate_HalfDrifted(t *testing.T) {
	s := diffstat.Stats{Total: 4, Drifted: 2}
	rate := diffstat.DriftRate(s)
	if rate != 0.5 {
		t.Fatalf("expected 0.5, got %f", rate)
	}
}

func TestDriftRate_AllDrifted(t *testing.T) {
	s := diffstat.Stats{Total: 3, Drifted: 3}
	rate := diffstat.DriftRate(s)
	if rate != 1.0 {
		t.Fatalf("expected 1.0, got %f", rate)
	}
}
