package anomaly_test

import (
	"testing"

	"driftwatch/internal/anomaly"
	"driftwatch/internal/drift"
)

func makeResult(service string, diffs int) drift.Result {
	d := make([]drift.Diff, diffs)
	for i := range d {
		d[i] = drift.Diff{Field: fmt.Sprintf("field%d", i)}
	}
	return drift.Result{Service: service, Diffs: d}
}

func makeHistory(counts ...int) []drift.Result {
	var run []drift.Result
	for i, c := range counts {
		run = append(run, makeResult(fmt.Sprintf("svc%d", i), c))
	}
	return run
}

import "fmt"

func historyFor(service string, driftCounts []int) [][]drift.Result {
	var history [][]drift.Result
	for _, c := range driftCounts {
		history = append(history, []drift.Result{makeResult(service, c)})
	}
	return history
}

func TestDetect_InsufficientHistory_ReturnsNil(t *testing.T) {
	det := anomaly.New()
	// Only 2 snapshots, MinSamples defaults to 3.
	hist := historyFor("api", []int{1, 2})
	current := []drift.Result{makeResult("api", 10)}
	got := det.Detect(hist, current)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestDetect_NoAnomaly_ReturnsNil(t *testing.T) {
	det := anomaly.New()
	// Stable history: all zeros.
	hist := historyFor("api", []int{0, 0, 0, 0, 0})
	// Current is also 0 — stddev is 0, so skipped.
	current := []drift.Result{makeResult("api", 0)}
	got := det.Detect(hist, current)
	if len(got) != 0 {
		t.Fatalf("expected no anomalies, got %v", got)
	}
}

func TestDetect_HighZScore_FlagsAnomaly(t *testing.T) {
	det := anomaly.New()
	// History: small stable drift.
	hist := historyFor("api", []int{1, 1, 1, 2, 1})
	// Current: sudden spike.
	current := []drift.Result{makeResult("api", 20)}
	got := det.Detect(hist, current)
	if len(got) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(got))
	}
	if got[0].Service != "api" {
		t.Errorf("expected service 'api', got %q", got[0].Service)
	}
	if got[0].ZScore <= 2.0 {
		t.Errorf("expected z-score > 2.0, got %.2f", got[0].ZScore)
	}
}

func TestDetect_SortedByZScoreDescending(t *testing.T) {
	det := anomaly.New()
	hist := append(
		historyFor("alpha", []int{1, 1, 1, 1}),
		historyFor("beta", []int{2, 2, 2, 2})...,
	)
	// Build combined history slices.
	combined := make([][]drift.Result, 4)
	for i := 0; i < 4; i++ {
		combined[i] = []drift.Result{
			makeResult("alpha", []int{1, 1, 1, 1}[i]),
			makeResult("beta", []int{2, 2, 2, 2}[i]),
		}
	}
	_ = hist
	current := []drift.Result{
		makeResult("alpha", 30), // very high spike
		makeResult("beta", 10),  // moderate spike
	}
	got := det.Detect(combined, current)
	if len(got) < 2 {
		t.Fatalf("expected at least 2 anomalies, got %d", len(got))
	}
	if got[0].Service != "alpha" {
		t.Errorf("expected alpha first (higher z), got %q", got[0].Service)
	}
}

func TestDetect_UnknownService_Skipped(t *testing.T) {
	det := anomaly.New()
	hist := historyFor("api", []int{1, 1, 1, 1})
	// 'new-svc' has no history.
	current := []drift.Result{makeResult("new-svc", 99)}
	got := det.Detect(hist, current)
	if len(got) != 0 {
		t.Fatalf("expected no anomalies for unknown service, got %v", got)
	}
}

func TestDetect_CustomThreshold(t *testing.T) {
	det := anomaly.New()
	det.Threshold = 10.0 // extremely high — nothing should fire
	hist := historyFor("api", []int{1, 1, 1, 2, 1})
	current := []drift.Result{makeResult("api", 5)}
	got := det.Detect(hist, current)
	if len(got) != 0 {
		t.Fatalf("expected no anomalies with high threshold, got %v", got)
	}
}
