package labelfilter_test

import (
	"testing"

	"driftwatch/internal/labelfilter"
)

func TestMatch_NoSelectors_AcceptsAll(t *testing.T) {
	f := labelfilter.New(nil)
	if !f.Match(map[string]string{"env": "prod"}) {
		t.Fatal("expected match with no selectors")
	}
	if !f.Match(map[string]string{}) {
		t.Fatal("expected match for empty labels with no selectors")
	}
}

func TestMatch_SingleSelector_Matches(t *testing.T) {
	f := labelfilter.New([]string{"env=prod"})
	if !f.Match(map[string]string{"env": "prod", "team": "platform"}) {
		t.Fatal("expected match")
	}
}

func TestMatch_SingleSelector_NoMatch(t *testing.T) {
	f := labelfilter.New([]string{"env=prod"})
	if f.Match(map[string]string{"env": "staging"}) {
		t.Fatal("expected no match")
	}
}

func TestMatch_MissingLabel_NoMatch(t *testing.T) {
	f := labelfilter.New([]string{"env=prod"})
	if f.Match(map[string]string{"team": "platform"}) {
		t.Fatal("expected no match when label key absent")
	}
}

func TestMatch_MultipleSelectors_AllMustMatch(t *testing.T) {
	f := labelfilter.New([]string{"env=prod", "tier=api"})
	if !f.Match(map[string]string{"env": "prod", "tier": "api"}) {
		t.Fatal("expected match")
	}
	if f.Match(map[string]string{"env": "prod", "tier": "worker"}) {
		t.Fatal("expected no match when one selector fails")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	f := labelfilter.New([]string{"ENV=PROD"})
	if !f.Match(map[string]string{"env": "prod"}) {
		t.Fatal("expected case-insensitive match")
	}
}

func TestNew_MalformedSelector_Ignored(t *testing.T) {
	f := labelfilter.New([]string{"no-equals-sign", "env=prod"})
	if !f.Match(map[string]string{"env": "prod"}) {
		t.Fatal("expected match; malformed selector should be ignored")
	}
}

func TestMatchAll_FiltersEntries(t *testing.T) {
	f := labelfilter.New([]string{"env=prod"})
	entries := []labelfilter.Entry{
		{Name: "api", Labels: map[string]string{"env": "prod"}},
		{Name: "worker", Labels: map[string]string{"env": "staging"}},
		{Name: "gateway", Labels: map[string]string{"env": "prod"}},
	}
	got := f.MatchAll(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0] != "api" || got[1] != "gateway" {
		t.Fatalf("unexpected results: %v", got)
	}
}

func TestMatchAll_NoSelectors_ReturnsAll(t *testing.T) {
	f := labelfilter.New(nil)
	entries := []labelfilter.Entry{
		{Name: "a", Labels: map[string]string{}},
		{Name: "b", Labels: map[string]string{"env": "prod"}},
	}
	got := f.MatchAll(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
