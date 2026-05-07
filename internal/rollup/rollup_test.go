package rollup_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/rollup"
)

func makeResult(service string, diffs []drift.Diff) drift.Result {
	return drift.Result{
		Service: service,
		Diffs:   diffs,
	}
}

func TestAggregate_Empty(t *testing.T) {
	s := rollup.Aggregate(nil)
	if s.Total != 0 || s.Drifted != 0 || s.Clean != 0 {
		t.Fatalf("expected zero counts, got %+v", s)
	}
}

func TestAggregate_AllClean(t *testing.T) {
	results := []drift.Result{
		makeResult("alpha", nil),
		makeResult("beta", nil),
	}
	s := rollup.Aggregate(results)
	if s.Total != 2 || s.Clean != 2 || s.Drifted != 0 {
		t.Fatalf("unexpected counts: %+v", s)
	}
}

func TestAggregate_SomeDrifted(t *testing.T) {
	results := []drift.Result{
		makeResult("alpha", []drift.Diff{{Field: "replicas", Want: 3, Got: 1}}),
		makeResult("beta", nil),
	}
	s := rollup.Aggregate(results)
	if s.Drifted != 1 || s.Clean != 1 {
		t.Fatalf("expected 1 drifted 1 clean, got %+v", s)
	}
	if s.Services[0].Name != "alpha" {
		t.Fatalf("expected alpha first (sorted), got %s", s.Services[0].Name)
	}
}

func TestAggregate_SortedByName(t *testing.T) {
	results := []drift.Result{
		makeResult("zebra", nil),
		makeResult("apple", nil),
		makeResult("mango", nil),
	}
	s := rollup.Aggregate(results)
	names := []string{s.Services[0].Name, s.Services[1].Name, s.Services[2].Name}
	if names[0] != "apple" || names[1] != "mango" || names[2] != "zebra" {
		t.Fatalf("unexpected order: %v", names)
	}
}

func TestAggregate_DiffDetailsPopulated(t *testing.T) {
	results := []drift.Result{
		makeResult("svc", []drift.Diff{
			{Field: "image", Want: "v2", Got: "v1"},
		}),
	}
	s := rollup.Aggregate(results)
	if len(s.Services[0].Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(s.Services[0].Details))
	}
}

func TestSummary_String_ContainsServiceName(t *testing.T) {
	results := []drift.Result{
		makeResult("my-service", []drift.Diff{{Field: "port", Want: 8080, Got: 9090}}),
	}
	s := rollup.Aggregate(results)
	out := s.String()
	if !strings.Contains(out, "my-service") {
		t.Fatalf("expected service name in output:\n%s", out)
	}
	if !strings.Contains(out, "DRIFTED") {
		t.Fatalf("expected DRIFTED label in output:\n%s", out)
	}
}
