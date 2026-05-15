package staledetector_test

import (
	"testing"
	"time"

	"driftwatch/internal/drift"
	"driftwatch/internal/staledetector"
)

func makeResult(service string, checkedAt time.Time) drift.Result {
	return drift.Result{Service: service, CheckedAt: checkedAt}
}

func TestDetect_AllFresh_ReturnsNil(t *testing.T) {
	d := staledetector.New(5 * time.Minute)
	now := time.Now()
	results := []drift.Result{
		makeResult("api", now.Add(-1*time.Minute)),
		makeResult("worker", now.Add(-2*time.Minute)),
	}
	got := d.Detect(results)
	if len(got) != 0 {
		t.Fatalf("expected no stale entries, got %d", len(got))
	}
}

func TestDetect_StaleEntry_Returned(t *testing.T) {
	d := staledetector.New(5 * time.Minute)
	now := time.Now()
	results := []drift.Result{
		makeResult("api", now.Add(-10*time.Minute)),
		makeResult("worker", now.Add(-1*time.Minute)),
	}
	got := d.Detect(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 stale entry, got %d", len(got))
	}
	if got[0].Service != "api" {
		t.Errorf("expected stale service api, got %s", got[0].Service)
	}
}

func TestDetect_ZeroTimestamp_AlwaysStale(t *testing.T) {
	d := staledetector.New(5 * time.Minute)
	results := []drift.Result{
		makeResult("ghost", time.Time{}),
	}
	got := d.Detect(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 stale entry, got %d", len(got))
	}
	if !got[0].LastSeen.IsZero() {
		t.Errorf("expected zero LastSeen")
	}
}

func TestDetect_MultipleStale_AllReturned(t *testing.T) {
	d := staledetector.New(5 * time.Minute)
	now := time.Now()
	results := []drift.Result{
		makeResult("a", now.Add(-6*time.Minute)),
		makeResult("b", now.Add(-7*time.Minute)),
		makeResult("c", now.Add(-1*time.Minute)),
	}
	got := d.Detect(results)
	if len(got) != 2 {
		t.Fatalf("expected 2 stale entries, got %d", len(got))
	}
}

func TestIsStale_FreshResult_ReturnsFalse(t *testing.T) {
	d := staledetector.New(5 * time.Minute)
	r := makeResult("api", time.Now().Add(-1*time.Minute))
	if d.IsStale(r) {
		t.Error("expected fresh result to not be stale")
	}
}

func TestIsStale_OldResult_ReturnsTrue(t *testing.T) {
	d := staledetector.New(5 * time.Minute)
	r := makeResult("api", time.Now().Add(-10*time.Minute))
	if !d.IsStale(r) {
		t.Error("expected old result to be stale")
	}
}

func TestDetect_EmptyInput_ReturnsNil(t *testing.T) {
	d := staledetector.New(5 * time.Minute)
	got := d.Detect(nil)
	if got != nil {
		t.Errorf("expected nil for empty input, got %v", got)
	}
}
