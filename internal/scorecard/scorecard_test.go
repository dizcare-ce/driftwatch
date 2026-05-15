package scorecard_test

import (
	"testing"

	"driftwatch/internal/drift"
	"driftwatch/internal/scorecard"
)

func makeResult(service string, drifted, total int) drift.Result {
	return drift.Result{Service: service, Drifted: drifted, Total: total}
}

func TestBuild_Empty_ReturnsNil(t *testing.T) {
	entries := scorecard.Build(nil)
	if len(entries) != 0 {
		t.Fatalf("expected empty, got %d entries", len(entries))
	}
}

func TestBuild_NoDrift_GradeA(t *testing.T) {
	results := []drift.Result{makeResult("api", 0, 10)}
	entries := scorecard.Build(results)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Grade != "A" {
		t.Errorf("expected grade A, got %s", entries[0].Grade)
	}
	if entries[0].Score != 0 {
		t.Errorf("expected score 0, got %d", entries[0].Score)
	}
}

func TestBuild_FullDrift_GradeF(t *testing.T) {
	results := []drift.Result{makeResult("worker", 10, 10)}
	entries := scorecard.Build(results)
	if entries[0].Grade != "F" {
		t.Errorf("expected grade F, got %s", entries[0].Grade)
	}
	if entries[0].Score != 100 {
		t.Errorf("expected score 100, got %d", entries[0].Score)
	}
}

func TestBuild_PartialDrift_GradeC(t *testing.T) {
	// 4/8 = 50 → grade C
	results := []drift.Result{makeResult("svc", 4, 8)}
	entries := scorecard.Build(results)
	if entries[0].Grade != "C" {
		t.Errorf("expected grade C, got %s", entries[0].Grade)
	}
}

func TestBuild_SortedWorstFirst(t *testing.T) {
	results := []drift.Result{
		makeResult("clean", 0, 10),
		makeResult("bad", 10, 10),
		makeResult("mid", 5, 10),
	}
	entries := scorecard.Build(results)
	if entries[0].Service != "bad" {
		t.Errorf("expected bad first, got %s", entries[0].Service)
	}
	if entries[len(entries)-1].Service != "clean" {
		t.Errorf("expected clean last, got %s", entries[len(entries)-1].Service)
	}
}

func TestBuild_SkipsZeroTotalResults(t *testing.T) {
	results := []drift.Result{makeResult("ghost", 0, 0)}
	entries := scorecard.Build(results)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestSummary_NoEntries(t *testing.T) {
	s := scorecard.Summary(nil)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_CountsFailing(t *testing.T) {
	results := []drift.Result{
		makeResult("a", 10, 10),
		makeResult("b", 0, 10),
	}
	entries := scorecard.Build(results)
	s := scorecard.Summary(entries)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
