package tagindex_test

import (
	"testing"

	"driftwatch/internal/source"
	"driftwatch/internal/tagindex"
)

func defs(pairs ...any) []source.Definition {
	var out []source.Definition
	for i := 0; i+1 < len(pairs); i += 2 {
		name := pairs[i].(string)
		labels := pairs[i+1].(map[string]string)
		out = append(out, source.Definition{Name: name, Labels: labels})
	}
	return out
}

func TestBuild_PopulatesIndex(t *testing.T) {
	idx := tagindex.New()
	b := tagindex.NewBuilder(idx)

	b.Build(defs(
		"api", map[string]string{"env": "prod"},
		"worker", map[string]string{"env": "prod"},
	))

	got := idx.Lookup("env", "prod")
	if len(got) != 2 {
		t.Fatalf("expected 2 services, got %d", len(got))
	}
}

func TestBuild_RemovesStaleServices(t *testing.T) {
	idx := tagindex.New()
	b := tagindex.NewBuilder(idx)

	b.Build(defs("api", map[string]string{"env": "prod"}))
	b.Build(defs("worker", map[string]string{"env": "prod"}))

	if got := idx.Lookup("env", "prod"); len(got) != 1 || got[0] != "worker" {
		t.Fatalf("expected only worker, got %v", got)
	}
}

func TestBuild_NilLabels_RegistersService(t *testing.T) {
	idx := tagindex.New()
	b := tagindex.NewBuilder(idx)

	b.Build([]source.Definition{{Name: "bare"}})

	entries := idx.All()
	if len(entries) != 1 || entries[0].Service != "bare" {
		t.Fatalf("expected bare service registered, got %v", entries)
	}
}

func TestBuild_IdempotentOnSameInput(t *testing.T) {
	idx := tagindex.New()
	b := tagindex.NewBuilder(idx)

	input := defs("api", map[string]string{"tier": "frontend"})
	b.Build(input)
	b.Build(input)

	if entries := idx.All(); len(entries) != 1 {
		t.Fatalf("expected 1 entry after double build, got %d", len(entries))
	}
}
