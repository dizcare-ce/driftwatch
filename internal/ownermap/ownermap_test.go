package ownermap_test

import (
	"testing"

	"driftwatch/internal/ownermap"
)

func TestSet_And_Get_RoundTrip(t *testing.T) {
	m := ownermap.New()
	if err := m.Set("api", "platform", "platform@example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	o, ok := m.Get("api")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if o.Team != "platform" || o.Contact != "platform@example.com" {
		t.Errorf("unexpected owner: %+v", o)
	}
}

func TestGet_CaseInsensitive(t *testing.T) {
	m := ownermap.New()
	_ = m.Set("API", "team-a", "a@example.com")
	_, ok := m.Get("api")
	if !ok {
		t.Error("expected case-insensitive lookup to succeed")
	}
}

func TestGet_MissingEntry_ReturnsFalse(t *testing.T) {
	m := ownermap.New()
	_, ok := m.Get("unknown")
	if ok {
		t.Error("expected false for unknown service")
	}
}

func TestSet_EmptyService_ReturnsError(t *testing.T) {
	m := ownermap.New()
	if err := m.Set("", "team", "x@example.com"); err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestSet_OverwritesPreviousEntry(t *testing.T) {
	m := ownermap.New()
	_ = m.Set("svc", "old-team", "old@example.com")
	_ = m.Set("svc", "new-team", "new@example.com")
	o, _ := m.Get("svc")
	if o.Team != "new-team" {
		t.Errorf("expected new-team, got %s", o.Team)
	}
}

func TestRemove_ClearsEntry(t *testing.T) {
	m := ownermap.New()
	_ = m.Set("svc", "team", "t@example.com")
	m.Remove("svc")
	_, ok := m.Get("svc")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestRemove_UnknownService_NoError(t *testing.T) {
	m := ownermap.New()
	m.Remove("ghost") // must not panic
}

func TestAll_ReturnsSortedOwners(t *testing.T) {
	m := ownermap.New()
	_ = m.Set("zebra", "z-team", "z@example.com")
	_ = m.Set("alpha", "a-team", "a@example.com")
	_ = m.Set("mango", "m-team", "m@example.com")
	all := m.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
	names := []string{all[0].Service, all[1].Service, all[2].Service}
	expected := []string{"alpha", "mango", "zebra"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("position %d: got %s, want %s", i, n, expected[i])
		}
	}
}
