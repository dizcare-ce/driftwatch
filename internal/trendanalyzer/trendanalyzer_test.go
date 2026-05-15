package trendanalyzer_test

import (
	"testing"
	"time"

	"driftwatch/internal/trendanalyzer"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeSnapshots(drifts ...int) []trendanalyzer.Snapshot {
	ss := make([]trendanalyzer.Snapshot, len(drifts))
	for i, d := range drifts {
		ss[i] = trendanalyzer.Snapshot{At: epoch.Add(time.Duration(i) * time.Hour), Drift: d}
	}
	return ss
}

func TestAnalyse_Empty_ReturnsStable(t *testing.T) {
	a := trendanalyzer.New(10)
	r := a.Analyse(nil)
	if r.Direction != trendanalyzer.DirectionStable {
		t.Fatalf("expected stable, got %s", r.Direction)
	}
	if r.Window != 0 {
		t.Fatalf("expected window 0, got %d", r.Window)
	}
}

func TestAnalyse_SingleSnapshot(t *testing.T) {
	a := trendanalyzer.New(10)
	r := a.Analyse(makeSnapshots(3))
	if r.Window != 1 {
		t.Fatalf("expected window 1, got %d", r.Window)
	}
	if r.MinDrifts != 3 || r.MaxDrifts != 3 {
		t.Fatalf("unexpected min/max: %d/%d", r.MinDrifts, r.MaxDrifts)
	}
	if r.Direction != trendanalyzer.DirectionStable {
		t.Fatalf("expected stable, got %s", r.Direction)
	}
}

func TestAnalyse_Increasing(t *testing.T) {
	a := trendanalyzer.New(10)
	r := a.Analyse(makeSnapshots(1, 2, 3, 4, 5, 6))
	if r.Direction != trendanalyzer.DirectionIncreasing {
		t.Fatalf("expected increasing, got %s", r.Direction)
	}
}

func TestAnalyse_Decreasing(t *testing.T) {
	a := trendanalyzer.New(10)
	r := a.Analyse(makeSnapshots(6, 5, 4, 3, 2, 1))
	if r.Direction != trendanalyzer.DirectionDecreasing {
		t.Fatalf("expected decreasing, got %s", r.Direction)
	}
}

func TestAnalyse_Stable(t *testing.T) {
	a := trendanalyzer.New(10)
	r := a.Analyse(makeSnapshots(3, 3, 3, 3))
	if r.Direction != trendanalyzer.DirectionStable {
		t.Fatalf("expected stable, got %s", r.Direction)
	}
}

func TestAnalyse_RespectsMaxWindow(t *testing.T) {
	a := trendanalyzer.New(3)
	ss := makeSnapshots(10, 10, 10, 1, 1, 1) // last 3 are low
	r := a.Analyse(ss)
	if r.Window != 3 {
		t.Fatalf("expected window 3, got %d", r.Window)
	}
	if r.MaxDrifts != 1 {
		t.Fatalf("expected max 1 within window, got %d", r.MaxDrifts)
	}
}

func TestAnalyse_AvgDrifts(t *testing.T) {
	a := trendanalyzer.New(10)
	r := a.Analyse(makeSnapshots(2, 4, 6))
	if r.AvgDrifts != 4.0 {
		t.Fatalf("expected avg 4.0, got %f", r.AvgDrifts)
	}
}

func TestAnalyse_FirstAndLastTimestamps(t *testing.T) {
	a := trendanalyzer.New(10)
	ss := makeSnapshots(1, 2, 3)
	r := a.Analyse(ss)
	if !r.First.Equal(ss[0].At) {
		t.Fatalf("unexpected First: %v", r.First)
	}
	if !r.Last.Equal(ss[2].At) {
		t.Fatalf("unexpected Last: %v", r.Last)
	}
}
