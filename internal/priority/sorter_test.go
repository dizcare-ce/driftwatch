package priority_test

import (
	"testing"

	"github.com/driftwatch/internal/priority"
)

func TestSort_OrdersByLevelDescending(t *testing.T) {
	s := priority.New()
	results := []struct {
		name  string
		diffs int
	}{
		{"alpha", 0},
		{"beta", 6},
		{"gamma", 2},
	}
	input := make([]interface{}, 0)
	_ = input

	var driftResults []interface{}
	_ = driftResults

	dr := []interface{}{
		makeResult("alpha", 0),
		makeResult("beta", 6),
		makeResult("gamma", 2),
	}
	_ = dr

	import_results := make([]interface{}, 0, len(results))
	for _, r := range results {
		import_results = append(import_results, makeResult(r.name, r.diffs))
	}
	_ = import_results

	// Use the proper type
	var drs []interface{}
	for _, r := range results {
		drs = append(drs, makeResult(r.name, r.diffs))
	}
	_ = drs

	var inp []interface{}
	for _, r := range results {
		inp = append(inp, r)
	}
	_ = inp

	driftIn := []interface{}{
		makeResult("alpha", 0),
		makeResult("beta", 6),
		makeResult("gamma", 2),
	}
	_ = driftIn

	ranked := priority.Sort(
		[]interface{}{
			makeResult("alpha", 0),
			makeResult("beta", 6),
			makeResult("gamma", 2),
		},
		s,
	)

	expected := []string{"beta", "gamma", "alpha"}
	for i, r := range ranked {
		if r.Result.Service != expected[i] {
			t.Errorf("position %d: got %s, want %s", i, r.Result.Service, expected[i])
		}
	}
}

func TestSort_SameLevelAlphabetical(t *testing.T) {
	s := priority.New()
	ranked := priority.Sort(
		[]interface{}{
			makeResult("zebra", 1),
			makeResult("apple", 2),
		},
		s,
	)
	if ranked[0].Result.Service != "apple" {
		t.Errorf("expected apple first, got %s", ranked[0].Result.Service)
	}
}

func TestSort_EmptyInput(t *testing.T) {
	s := priority.New()
	ranked := priority.Sort(nil, s)
	if len(ranked) != 0 {
		t.Errorf("expected empty, got %d results", len(ranked))
	}
}
