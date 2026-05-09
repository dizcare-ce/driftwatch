package priority_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/priority"
)

func makeResult(service string, diffs int) drift.Result {
	d := make([]drift.Diff, diffs)
	for i := range d {
		d[i] = drift.Diff{Field: "field", Want: "a", Got: "b"}
	}
	return drift.Result{Service: service, Diffs: d}
}

func TestScore_NoDrift_IsLow(t *testing.T) {
	s := priority.New()
	got := s.Score(makeResult("svc", 0))
	if got != priority.LevelLow {
		t.Fatalf("expected low, got %s", got)
	}
}

func TestScore_OneDiff_IsMedium(t *testing.T) {
	s := priority.New()
	got := s.Score(makeResult("svc", 1))
	if got != priority.LevelMedium {
		t.Fatalf("expected medium, got %s", got)
	}
}

func TestScore_ThreeDiffs_IsHigh(t *testing.T) {
	s := priority.New()
	got := s.Score(makeResult("svc", 3))
	if got != priority.LevelHigh {
		t.Fatalf("expected high, got %s", got)
	}
}

func TestScore_SixDiffs_IsCritical(t *testing.T) {
	s := priority.New()
	got := s.Score(makeResult("svc", 6))
	if got != priority.LevelCritical {
		t.Fatalf("expected critical, got %s", got)
	}
}

func TestScoreAll_ReturnsMappedLevels(t *testing.T) {
	s := priority.New()
	results := []drift.Result{
		makeResult("alpha", 0),
		makeResult("beta", 4),
	}
	got := s.ScoreAll(results)
	if got["alpha"] != priority.LevelLow {
		t.Errorf("alpha: expected low, got %s", got["alpha"])
	}
	if got["beta"] != priority.LevelHigh {
		t.Errorf("beta: expected high, got %s", got["beta"])
	}
}

func TestParseLevel_Valid(t *testing.T) {
	cases := []struct {
		in   string
		want priority.Level
	}{
		{"low", priority.LevelLow},
		{"Medium", priority.LevelMedium},
		{"HIGH", priority.LevelHigh},
		{"critical", priority.LevelCritical},
	}
	for _, c := range cases {
		got, err := priority.ParseLevel(c.in)
		if err != nil {
			t.Fatalf("ParseLevel(%q): unexpected error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("ParseLevel(%q) = %s, want %s", c.in, got, c.want)
		}
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	_, err := priority.ParseLevel("extreme")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestLevel_String(t *testing.T) {
	if priority.LevelCritical.String() != "critical" {
		t.Errorf("unexpected string for LevelCritical")
	}
}
