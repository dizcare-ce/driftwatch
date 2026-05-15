package convergence

import (
	"strings"
	"testing"
	"time"

	"driftwatch/internal/drift"
)

var epoch = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func makeResult(svc string, checkedAt time.Time, diffs int) drift.Result {
	d := make([]drift.Diff, diffs)
	for i := range d {
		d[i] = drift.Diff{Field: "f", Expected: "a", Actual: "b"}
	}
	return drift.Result{Service: svc, CheckedAt: checkedAt, Diffs: d}
}

func TestAnalyse_Empty_ReturnsNil(t *testing.T) {
	got := Analyse(nil, epoch)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestAnalyse_NoDriftInHistory_ReturnsNil(t *testing.T) {
	snap := []drift.Result{makeResult("api", epoch.Add(-time.Hour), 0)}
	got := Analyse([][]drift.Result{snap}, epoch)
	if len(got) != 0 {
		t.Fatalf("expected no estimates, got %d", len(got))
	}
}

func TestAnalyse_SingleDriftedService_Stable(t *testing.T) {
	snaps := [][]drift.Result{
		{makeResult("svc", epoch.Add(-3*time.Hour), 2)},
		{makeResult("svc", epoch.Add(-2*time.Hour), 2)},
		{makeResult("svc", epoch.Add(-time.Hour), 2)},
	}
	got := Analyse(snaps, epoch)
	if len(got) != 1 {
		t.Fatalf("expected 1 estimate, got %d", len(got))
	}
	if got[0].Trend != "stable" {
		t.Errorf("expected stable, got %s", got[0].Trend)
	}
	if !got[0].ProjectedClean.IsZero() {
		t.Error("expected no projection for stable trend")
	}
}

func TestAnalyse_ImprovingTrend_HasProjection(t *testing.T) {
	snaps := [][]drift.Result{
		{makeResult("svc", epoch.Add(-3*time.Hour), 6)},
		{makeResult("svc", epoch.Add(-2*time.Hour), 4)},
		{makeResult("svc", epoch.Add(-time.Hour), 2)},
	}
	got := Analyse(snaps, epoch)
	if len(got) != 1 {
		t.Fatalf("expected 1 estimate, got %d", len(got))
	}
	if got[0].Trend != "improving" {
		t.Errorf("expected improving, got %s", got[0].Trend)
	}
	if got[0].ProjectedClean.IsZero() {
		t.Error("expected a projected clean time for improving trend")
	}
	if !got[0].ProjectedClean.After(epoch) {
		t.Error("projected clean time should be in the future")
	}
}

func TestAnalyse_WorseningTrend_NoProjection(t *testing.T) {
	snaps := [][]drift.Result{
		{makeResult("svc", epoch.Add(-3*time.Hour), 1)},
		{makeResult("svc", epoch.Add(-2*time.Hour), 3)},
		{makeResult("svc", epoch.Add(-time.Hour), 5)},
	}
	got := Analyse(snaps, epoch)
	if len(got) != 1 {
		t.Fatalf("expected 1 estimate, got %d", len(got))
	}
	if got[0].Trend != "worsening" {
		t.Errorf("expected worsening, got %s", got[0].Trend)
	}
	if !got[0].ProjectedClean.IsZero() {
		t.Error("expected no projection for worsening trend")
	}
}

func TestEstimate_String_ContainsServiceName(t *testing.T) {
	e := Estimate{
		Service:       "payments",
		DriftDuration: 90 * time.Minute,
		Trend:         "stable",
	}
	if !strings.Contains(e.String(), "payments") {
		t.Errorf("String() missing service name: %s", e.String())
	}
}

func TestEstimate_String_WithProjection(t *testing.T) {
	e := Estimate{
		Service:        "auth",
		DriftDuration:  2 * time.Hour,
		Trend:          "improving",
		ProjectedClean: epoch.Add(3 * time.Hour),
	}
	s := e.String()
	if !strings.Contains(s, "improving") {
		t.Errorf("String() missing trend: %s", s)
	}
	if !strings.Contains(s, "2024") {
		t.Errorf("String() missing projected time: %s", s)
	}
}
