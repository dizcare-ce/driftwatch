package groupby_test

import (
	"testing"

	"driftwatch/internal/drift"
	"driftwatch/internal/groupby"
)

func makeResult(service string, diffs int) drift.Result {
	d := make([]drift.Diff, diffs)
	for i := range d {
		d[i] = drift.Diff{Field: "f", Expected: "a", Actual: "b"}
	}
	return drift.Result{Service: service, Diffs: d}
}

func TestNew_NilKeyFunc_ReturnsError(t *testing.T) {
	_, err := groupby.New(nil, "")
	if err == nil {
		t.Fatal("expected error for nil keyFn")
	}
}

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	g, _ := groupby.New(groupby.ByService, "")
	if got := g.Apply(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestApply_GroupsByKey(t *testing.T) {
	results := []drift.Result{
		makeResult("alpha", 1),
		makeResult("beta", 0),
		makeResult("alpha", 2),
	}
	g, _ := groupby.New(groupby.ByService, "")
	groups := g.Apply(results)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "alpha" {
		t.Errorf("expected first key 'alpha', got %q", groups[0].Key)
	}
	if len(groups[0].Results) != 2 {
		t.Errorf("expected 2 results in 'alpha', got %d", len(groups[0].Results))
	}
}

func TestApply_SortedByKeyAscending(t *testing.T) {
	results := []drift.Result{
		makeResult("zebra", 0),
		makeResult("ant", 1),
		makeResult("mango", 0),
	}
	g, _ := groupby.New(groupby.ByService, "")
	groups := g.Apply(results)

	keys := []string{groups[0].Key, groups[1].Key, groups[2].Key}
	want := []string{"ant", "mango", "zebra"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestApply_EmptyKeyFallback(t *testing.T) {
	keyFn := func(r drift.Result) string { return "" }
	g, _ := groupby.New(keyFn, "misc")
	groups := g.Apply([]drift.Result{makeResult("svc", 0)})

	if len(groups) != 1 || groups[0].Key != "misc" {
		t.Errorf("expected fallback key 'misc', got %v", groups)
	}
}

func TestApply_DefaultFallbackLabel(t *testing.T) {
	keyFn := func(r drift.Result) string { return "" }
	g, _ := groupby.New(keyFn, "")
	groups := g.Apply([]drift.Result{makeResult("svc", 0)})

	if groups[0].Key != "(ungrouped)" {
		t.Errorf("expected '(ungrouped)', got %q", groups[0].Key)
	}
}

func TestApply_CustomKeyFunc(t *testing.T) {
	results := []drift.Result{
		makeResult("payment-api", 1),
		makeResult("payment-worker", 0),
		makeResult("auth-service", 2),
	}
	prefixFn := func(r drift.Result) string {
		if len(r.Service) >= 4 {
			return r.Service[:4]
		}
		return r.Service
	}
	g, _ := groupby.New(prefixFn, "")
	groups := g.Apply(results)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}
