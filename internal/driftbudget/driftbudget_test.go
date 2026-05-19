package driftbudget_test

import (
	"testing"
	"time"

	"driftwatch/internal/drift"
	"driftwatch/internal/driftbudget"
)

func makeResult(service string, diffs int) drift.Result {
	r := drift.Result{Service: service}
	for i := 0; i < diffs; i++ {
		r.Diffs = append(r.Diffs, drift.Diff{
			Field:    "field",
			Expected: "a",
			Actual:   "b",
		})
	}
	return r
}

func TestNew_InvalidLimit_ReturnsError(t *testing.T) {
	_, err := driftbudget.New(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestNew_InvalidWindow_ReturnsError(t *testing.T) {
	_, err := driftbudget.New(10, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestRecord_UnderLimit_ReturnsNil(t *testing.T) {
	b, _ := driftbudget.New(10, time.Minute)
	results := []drift.Result{makeResult("svc-a", 3)}
	if err := b.Record(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecord_ExceedsLimit_ReturnsError(t *testing.T) {
	b, _ := driftbudget.New(5, time.Minute)
	results := []drift.Result{makeResult("svc-a", 6)}
	if err := b.Record(results); err == nil {
		t.Fatal("expected error when budget exceeded")
	}
}

func TestRecord_CumulativeExceedsLimit_ReturnsError(t *testing.T) {
	b, _ := driftbudget.New(5, time.Minute)
	_ = b.Record([]drift.Result{makeResult("svc-a", 3)})
	if err := b.Record([]drift.Result{makeResult("svc-b", 3)}); err == nil {
		t.Fatal("expected error on cumulative breach")
	}
}

func TestRemaining_ReflectsUsage(t *testing.T) {
	b, _ := driftbudget.New(10, time.Minute)
	_ = b.Record([]drift.Result{makeResult("svc-a", 4)})
	if got := b.Remaining(); got != 6 {
		t.Fatalf("expected 6 remaining, got %d", got)
	}
}

func TestReset_RestoresBudget(t *testing.T) {
	b, _ := driftbudget.New(5, time.Minute)
	_ = b.Record([]drift.Result{makeResult("svc-a", 5)})
	b.Reset()
	if got := b.Remaining(); got != 5 {
		t.Fatalf("expected full budget after reset, got %d", got)
	}
}

func TestRecord_CleanResults_DoNotConsumesBudget(t *testing.T) {
	b, _ := driftbudget.New(5, time.Minute)
	_ = b.Record([]drift.Result{makeResult("svc-a", 0)})
	if got := b.Remaining(); got != 5 {
		t.Fatalf("clean result should not consume budget, got %d remaining", got)
	}
}
