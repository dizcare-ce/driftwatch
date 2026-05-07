package baseline_test

import (
	"testing"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/cache"
	"github.com/driftwatch/internal/drift"
)

func makeResult(service string, drifted bool) drift.Result {
	r := drift.Result{Service: service}
	if drifted {
		r.Diffs = []drift.Diff{
			{Field: "replicas", Want: "3", Got: "1"},
		}
	}
	return r
}

func TestIsSuppressed_NoDrift_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.New(dir)
	checker := baseline.NewChecker(store)

	ok, err := checker.IsSuppressed(makeResult("api", false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("clean result should never be suppressed")
	}
}

func TestIsSuppressed_NoBaseline_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.New(dir)
	checker := baseline.NewChecker(store)

	ok, err := checker.IsSuppressed(makeResult("api", true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("should not be suppressed without a stored baseline")
	}
}

func TestApprove_Then_IsSuppressed_ReturnsTrue(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.New(dir)
	checker := baseline.NewChecker(store)
	r := makeResult("api", true)

	if err := checker.Approve(r, "known issue"); err != nil {
		t.Fatalf("Approve: %v", err)
	}

	ok, err := checker.IsSuppressed(r)
	if err != nil {
		t.Fatalf("IsSuppressed: %v", err)
	}
	if !ok {
		t.Error("result should be suppressed after approval")
	}
}

func TestIsSuppressed_DifferentChecksum_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.New(dir)
	checker := baseline.NewChecker(store)

	r1 := makeResult("api", true)
	checker.Approve(r1, "")

	r2 := makeResult("api", true)
	r2.Diffs = append(r2.Diffs, drift.Diff{Field: "image", Want: "v1", Got: "v2"})

	ok, err := checker.IsSuppressed(r2)
	if err != nil {
		t.Fatalf("IsSuppressed: %v", err)
	}
	if ok {
		t.Error("changed drift should not match old baseline")
	}
}

func TestRevoke_ClearsBaseline(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.New(dir)
	checker := baseline.NewChecker(store)
	r := makeResult("api", true)

	checker.Approve(r, "")
	if err := checker.Revoke("api"); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	ok, _ := checker.IsSuppressed(r)
	if ok {
		t.Error("result should not be suppressed after revocation")
	}
}

// Ensure Checksum is stable for the same result.
func TestChecksum_Deterministic(t *testing.T) {
	r := makeResult("svc", true)
	if cache.Checksum(r) != cache.Checksum(r) {
		t.Error("checksum should be deterministic")
	}
}
