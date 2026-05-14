package tagindex_test

import (
	"sort"
	"testing"

	"driftwatch/internal/tagindex"
)

func TestAdd_And_Lookup(t *testing.T) {
	idx := tagindex.New()
	idx.Add("api", map[string]string{"env": "prod", "team": "platform"})
	idx.Add("worker", map[string]string{"env": "prod", "team": "data"})

	got := idx.Lookup("env", "prod")
	sort.Strings(got)
	if len(got) != 2 || got[0] != "api" || got[1] != "worker" {
		t.Fatalf("expected [api worker], got %v", got)
	}
}

func TestLookup_MissingKey_ReturnsNil(t *testing.T) {
	idx := tagindex.New()
	if got := idx.Lookup("env", "staging"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestAdd_ReplacesLabels(t *testing.T) {
	idx := tagindex.New()
	idx.Add("api", map[string]string{"env": "prod"})
	idx.Add("api", map[string]string{"env": "staging"})

	if got := idx.Lookup("env", "prod"); len(got) != 0 {
		t.Fatalf("old label should be gone, got %v", got)
	}
	got := idx.Lookup("env", "staging")
	if len(got) != 1 || got[0] != "api" {
		t.Fatalf("expected [api], got %v", got)
	}
}

func TestRemove_ClearsService(t *testing.T) {
	idx := tagindex.New()
	idx.Add("api", map[string]string{"env": "prod"})
	idx.Remove("api")

	if got := idx.Lookup("env", "prod"); len(got) != 0 {
		t.Fatalf("expected empty after remove, got %v", got)
	}
}

func TestRemove_UnknownService_NoError(t *testing.T) {
	idx := tagindex.New()
	idx.Remove("ghost") // must not panic
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	idx := tagindex.New()
	idx.Add("api", map[string]string{"env": "prod"})
	idx.Add("worker", map[string]string{"env": "staging"})

	entries := idx.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestAll_Empty_ReturnsEmptySlice(t *testing.T) {
	idx := tagindex.New()
	if entries := idx.All(); len(entries) != 0 {
		t.Fatalf("expected empty, got %v", entries)
	}
}

func TestLookup_MultipleServicesShareLabel(t *testing.T) {
	idx := tagindex.New()
	for _, svc := range []string{"a", "b", "c"} {
		idx.Add(svc, map[string]string{"tier": "backend"})
	}
	got := idx.Lookup("tier", "backend")
	if len(got) != 3 {
		t.Fatalf("expected 3 services, got %d", len(got))
	}
}
