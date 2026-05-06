package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/filter"
)

func TestAllow_NoPatterns_AllowsEverything(t *testing.T) {
	f := filter.New(nil, nil)
	for _, name := range []string{"api", "worker", "db-proxy"} {
		if !f.Allow(name) {
			t.Errorf("expected %q to be allowed", name)
		}
	}
}

func TestAllow_IncludePattern_FiltersOthers(t *testing.T) {
	f := filter.New([]string{"api-*"}, nil)

	if !f.Allow("api-gateway") {
		t.Error("expected api-gateway to be allowed")
	}
	if f.Allow("worker") {
		t.Error("expected worker to be blocked")
	}
}

func TestAllow_ExcludePattern_BlocksMatch(t *testing.T) {
	f := filter.New(nil, []string{"legacy-*"})

	if f.Allow("legacy-auth") {
		t.Error("expected legacy-auth to be blocked")
	}
	if !f.Allow("api-gateway") {
		t.Error("expected api-gateway to be allowed")
	}
}

func TestAllow_ExcludeTakesPrecedenceOverInclude(t *testing.T) {
	f := filter.New([]string{"api-*"}, []string{"api-legacy"})

	if !f.Allow("api-gateway") {
		t.Error("expected api-gateway to be allowed")
	}
	if f.Allow("api-legacy") {
		t.Error("expected api-legacy to be blocked by exclude")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f := filter.New([]string{"API-*"}, nil)

	if !f.Allow("api-gateway") {
		t.Error("expected case-insensitive match for api-gateway")
	}
}

func TestApply_ReturnsFilteredSlice(t *testing.T) {
	f := filter.New([]string{"svc-*"}, []string{"svc-bad"})

	input := []string{"svc-alpha", "svc-bad", "other", "svc-beta"}
	got := f.Apply(input)

	expected := []string{"svc-alpha", "svc-beta"}
	if len(got) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
	for i, v := range expected {
		if got[i] != v {
			t.Errorf("index %d: expected %q, got %q", i, v, got[i])
		}
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	f := filter.New(nil, nil)
	got := f.Apply([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
