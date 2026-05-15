package costestimator_test

import (
	"testing"

	"driftwatch/internal/costestimator"
	"driftwatch/internal/drift"
)

func makeResult(service string, fields ...string) drift.Result {
	r := drift.Result{Service: service}
	for _, f := range fields {
		r.Diffs = append(r.Diffs, drift.Diff{Field: f, Expected: "a", Actual: "b"})
	}
	return r
}

func TestCompute_NoDrift_ReturnsEmpty(t *testing.T) {
	e := costestimator.New(costestimator.FieldWeight{}, 1.0)
	results := []drift.Result{makeResult("svc-a")}
	got := e.Compute(results)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d entries", len(got))
	}
}

func TestCompute_UsesFieldWeight(t *testing.T) {
	weights := costestimator.FieldWeight{"replicas": 5.0, "image": 2.0}
	e := costestimator.New(weights, 1.0)
	results := []drift.Result{makeResult("svc-a", "replicas", "image")}
	got := e.Compute(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 estimate, got %d", len(got))
	}
	if got[0].TotalCost != 7.0 {
		t.Errorf("expected total cost 7.0, got %.2f", got[0].TotalCost)
	}
}

func TestCompute_DefaultCostAppliedForUnknownField(t *testing.T) {
	e := costestimator.New(costestimator.FieldWeight{}, 3.0)
	results := []drift.Result{makeResult("svc-b", "unknown-field")}
	got := e.Compute(results)
	if got[0].TotalCost != 3.0 {
		t.Errorf("expected 3.0, got %.2f", got[0].TotalCost)
	}
}

func TestCompute_SortedByTotalCostDescending(t *testing.T) {
	weights := costestimator.FieldWeight{"replicas": 10.0}
	e := costestimator.New(weights, 1.0)
	results := []drift.Result{
		makeResult("cheap", "env"),
		makeResult("expensive", "replicas"),
	}
	got := e.Compute(results)
	if got[0].Service != "expensive" {
		t.Errorf("expected expensive first, got %s", got[0].Service)
	}
}

func TestCompute_BreakdownContainsAllFields(t *testing.T) {
	e := costestimator.New(costestimator.FieldWeight{"a": 1.0, "b": 2.0}, 0.5)
	results := []drift.Result{makeResult("svc", "a", "b", "c")}
	got := e.Compute(results)
	if len(got[0].Breakdown) != 3 {
		t.Errorf("expected 3 breakdown entries, got %d", len(got[0].Breakdown))
	}
}

func TestCompute_CaseInsensitiveWeightLookup(t *testing.T) {
	weights := costestimator.FieldWeight{"Replicas": 4.0}
	e := costestimator.New(weights, 1.0)
	results := []drift.Result{makeResult("svc", "replicas")}
	got := e.Compute(results)
	if got[0].TotalCost != 4.0 {
		t.Errorf("expected 4.0, got %.2f", got[0].TotalCost)
	}
}

func TestFormat_NoDrift_PrintsMessage(t *testing.T) {
	out := costestimator.Format(nil)
	if out != "no drift cost estimated\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_ContainsServiceName(t *testing.T) {
	estimate := []costestimator.Estimate{{Service: "my-service", TotalCost: 5.0}}
	out := costestimator.Format(estimate)
	if !containsStr(out, "my-service") {
		t.Errorf("expected service name in output, got:\n%s", out)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
