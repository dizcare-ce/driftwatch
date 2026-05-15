package signalmap_test

import (
	"testing"

	"driftwatch/internal/drift"
	"driftwatch/internal/priority"
	"driftwatch/internal/signalmap"
)

func makeResult(service string, diffs []drift.Diff) drift.Result {
	return drift.Result{Service: service, Diffs: diffs}
}

func makeDiff(field, want, got string) drift.Diff {
	return drift.Diff{Field: field, Expected: want, Actual: got}
}

func scorer() *priority.Scorer {
	return priority.New()
}

func TestMap_EmptyResults_ReturnsNil(t *testing.T) {
	m := signalmap.New("driftwatch", scorer())
	sigs := m.Map(nil, priority.LevelLow)
	if len(sigs) != 0 {
		t.Fatalf("expected no signals, got %d", len(sigs))
	}
}

func TestMap_CleanResult_BelowMinLevel_Excluded(t *testing.T) {
	m := signalmap.New("dw", scorer())
	results := []drift.Result{makeResult("svc-a", nil)}
	sigs := m.Map(results, priority.LevelMedium)
	if len(sigs) != 0 {
		t.Fatalf("expected no signals for clean result below threshold, got %d", len(sigs))
	}
}

func TestMap_DriftedResult_AboveMinLevel_Included(t *testing.T) {
	m := signalmap.New("dw", scorer())
	results := []drift.Result{
		makeResult("svc-b", []drift.Diff{
			makeDiff("replicas", "3", "1"),
		}),
	}
	sigs := m.Map(results, priority.LevelLow)
	if len(sigs) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(sigs))
	}
	if sigs[0].Service != "svc-b" {
		t.Errorf("expected service svc-b, got %s", sigs[0].Service)
	}
	if sigs[0].Name != "dw.svc-b" {
		t.Errorf("unexpected signal name: %s", sigs[0].Name)
	}
}

func TestMap_NamePrefix_AppliedToAllSignals(t *testing.T) {
	m := signalmap.New("prod", scorer())
	results := []drift.Result{
		makeResult("alpha", []drift.Diff{makeDiff("image", "v1", "v2")}),
		makeResult("beta", []drift.Diff{makeDiff("port", "8080", "9090")}),
	}
	sigs := m.Map(results, priority.LevelLow)
	for _, s := range sigs {
		if len(s.Name) < 5 || s.Name[:5] != "prod." {
			t.Errorf("signal name %q missing expected prefix", s.Name)
		}
	}
}

func TestMap_MessageIsNonEmpty(t *testing.T) {
	m := signalmap.New("dw", scorer())
	results := []drift.Result{
		makeResult("svc-c", []drift.Diff{makeDiff("env", "prod", "staging")}),
	}
	sigs := m.Map(results, priority.LevelLow)
	if len(sigs) == 0 {
		t.Fatal("expected at least one signal")
	}
	if sigs[0].Message == "" {
		t.Error("expected non-empty message on signal")
	}
}
