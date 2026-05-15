package ruleengine_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/ruleengine"
)

func makeResult(service string, diffs int) drift.Result {
	d := make([]drift.Diff, diffs)
	for i := range d {
		d[i] = drift.Diff{Field: fmt.Sprintf("field%d", i), Want: "a", Got: "b"}
	}
	return drift.Result{Service: service, Diffs: d}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := ruleengine.New([]ruleengine.Rule{
		{Name: "bad", ServicePattern: "[invalid"},
	})
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestEvaluate_NoRules_ReturnsNil(t *testing.T) {
	e, err := ruleengine.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	violations := e.Evaluate(makeResult("api", 3))
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestEvaluate_RequireClean_Passes(t *testing.T) {
	e, _ := ruleengine.New([]ruleengine.Rule{
		{Name: "no-drift", RequireClean: true},
	})
	violations := e.Evaluate(makeResult("api", 0))
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestEvaluate_RequireClean_Fails(t *testing.T) {
	e, _ := ruleengine.New([]ruleengine.Rule{
		{Name: "no-drift", RequireClean: true},
	})
	violations := e.Evaluate(makeResult("api", 2))
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "no-drift" {
		t.Errorf("wrong rule name: %s", violations[0].Rule)
	}
}

func TestEvaluate_MaxDiffs_Passes(t *testing.T) {
	e, _ := ruleengine.New([]ruleengine.Rule{
		{Name: "max-two", MaxDiffs: 2},
	})
	violations := e.Evaluate(makeResult("svc", 2))
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestEvaluate_MaxDiffs_Fails(t *testing.T) {
	e, _ := ruleengine.New([]ruleengine.Rule{
		{Name: "max-two", MaxDiffs: 2},
	})
	violations := e.Evaluate(makeResult("svc", 5))
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestEvaluate_PatternFiltersService(t *testing.T) {
	e, _ := ruleengine.New([]ruleengine.Rule{
		{Name: "api-only", ServicePattern: "^api$", RequireClean: true},
	})
	// "worker" does not match pattern — no violation expected
	violations := e.Evaluate(makeResult("worker", 4))
	if len(violations) != 0 {
		t.Fatalf("expected no violations for non-matching service, got %v", violations)
	}
	// "api" matches — violation expected
	violations = e.Evaluate(makeResult("api", 1))
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for matching service, got %d", len(violations))
	}
}

func TestEvaluate_MultipleRules_MultipleViolations(t *testing.T) {
	e, _ := ruleengine.New([]ruleengine.Rule{
		{Name: "clean", RequireClean: true},
		{Name: "max-one", MaxDiffs: 1},
	})
	violations := e.Evaluate(makeResult("svc", 3))
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}
