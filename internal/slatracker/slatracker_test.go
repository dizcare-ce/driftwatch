package slatracker_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/slatracker"
)

func makeResult(service string, diffs []drift.Diff) drift.Result {
	return drift.Result{Service: service, Diffs: diffs}
}

func makeDiff(field string) drift.Diff {
	return drift.Diff{Field: field, Expected: "a", Actual: "b"}
}

func TestReport_NoData_ReturnsError(t *testing.T) {
	tr := slatracker.New(time.Hour, 99.0)
	_, err := tr.Report("api")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestReport_AllCompliant_MeetsSLA(t *testing.T) {
	tr := slatracker.New(time.Hour, 95.0)
	for i := 0; i < 10; i++ {
		tr.Record(makeResult("api", nil))
	}
	rep, err := tr.Report("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rep.MeetsSLA {
		t.Errorf("expected SLA met, pct=%.1f", rep.CompliancePct)
	}
	if rep.CompliancePct != 100.0 {
		t.Errorf("expected 100%%, got %.1f", rep.CompliancePct)
	}
}

func TestReport_SomeDrifted_BelowThreshold(t *testing.T) {
	tr := slatracker.New(time.Hour, 95.0)
	// 8 clean, 2 drifted → 80%
	for i := 0; i < 8; i++ {
		tr.Record(makeResult("api", nil))
	}
	for i := 0; i < 2; i++ {
		tr.Record(makeResult("api", []drift.Diff{makeDiff("replicas")}))
	}
	rep, err := tr.Report("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.MeetsSLA {
		t.Errorf("expected SLA not met, pct=%.1f", rep.CompliancePct)
	}
	if rep.TotalChecks != 10 {
		t.Errorf("expected 10 checks, got %d", rep.TotalChecks)
	}
}

func TestReport_IndependentServicesDoNotInterfere(t *testing.T) {
	tr := slatracker.New(time.Hour, 90.0)
	tr.Record(makeResult("api", nil))
	tr.Record(makeResult("worker", []drift.Diff{makeDiff("image")}))

	apiRep, _ := tr.Report("api")
	if !apiRep.MeetsSLA {
		t.Error("api should meet SLA")
	}
	wRep, _ := tr.Report("worker")
	if wRep.MeetsSLA {
		t.Error("worker should not meet SLA")
	}
}

func TestPrune_RemovesOldEntries(t *testing.T) {
	tr := slatracker.New(50*time.Millisecond, 90.0)
	tr.Record(makeResult("api", nil))
	time.Sleep(80 * time.Millisecond)
	tr.Prune()
	_, err := tr.Report("api")
	if err == nil {
		t.Error("expected error after pruning stale entries")
	}
}

func TestReport_ExactlyAtThreshold_MeetsSLA(t *testing.T) {
	tr := slatracker.New(time.Hour, 80.0)
	// 8 clean, 2 drifted → exactly 80%
	for i := 0; i < 8; i++ {
		tr.Record(makeResult("svc", nil))
	}
	for i := 0; i < 2; i++ {
		tr.Record(makeResult("svc", []drift.Diff{makeDiff("cpu")}))
	}
	rep, _ := tr.Report("svc")
	if !rep.MeetsSLA {
		t.Errorf("expected SLA met at exactly threshold, pct=%.1f", rep.CompliancePct)
	}
}
