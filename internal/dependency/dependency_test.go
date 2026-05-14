package dependency_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/dependency"
)

func TestAdd_And_Dependencies(t *testing.T) {
	g := dependency.New()
	g.Add("api", "auth")
	g.Add("api", "db")

	deps := g.Dependencies("api")
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(deps))
	}
}

func TestAdd_DeduplicatesEdges(t *testing.T) {
	g := dependency.New()
	g.Add("api", "auth")
	g.Add("api", "auth") // duplicate

	if got := len(g.Dependencies("api")); got != 1 {
		t.Fatalf("expected 1 dependency after dedup, got %d", got)
	}
}

func TestDependants_ReturnsUpstreamCallers(t *testing.T) {
	g := dependency.New()
	g.Add("api", "auth")
	g.Add("worker", "auth")
	g.Add("api", "db")

	deps := g.Dependants("auth")
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependants of auth, got %d: %v", len(deps), deps)
	}
}

func TestDependants_UnknownService_ReturnsEmpty(t *testing.T) {
	g := dependency.New()
	if got := g.Dependants("ghost"); len(got) != 0 {
		t.Fatalf("expected no dependants, got %v", got)
	}
}

func TestRemove_ClearsEdges(t *testing.T) {
	g := dependency.New()
	g.Add("api", "auth")
	g.Remove("api")

	if got := len(g.Dependencies("api")); got != 0 {
		t.Fatalf("expected 0 dependencies after remove, got %d", got)
	}
}

func TestValidate_SelfReferential_ReturnsError(t *testing.T) {
	g := dependency.New()
	g.Add("api", "api")

	if err := g.Validate(); err == nil {
		t.Fatal("expected error for self-referential edge, got nil")
	}
}

func TestValidate_ValidGraph_ReturnsNil(t *testing.T) {
	g := dependency.New()
	g.Add("api", "auth")
	g.Add("worker", "db")

	if err := g.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
