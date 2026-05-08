package backoff_test

import (
	"testing"
	"time"

	"driftwatch/internal/backoff"
)

func TestNext_AttemptZero_ReturnsInitialInterval(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 100 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     10 * time.Second,
		Jitter:          false,
	}
	got := p.Next(0)
	if got != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", got)
	}
}

func TestNext_Doubles_EachAttempt(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 100 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     10 * time.Second,
		Jitter:          false,
	}
	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}
	for i, want := range expected {
		if got := p.Next(i); got != want {
			t.Errorf("attempt %d: expected %v, got %v", i, want, got)
		}
	}
}

func TestNext_CapsAtMaxInterval(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		Multiplier:      10.0,
		MaxInterval:     5 * time.Second,
		Jitter:          false,
	}
	for attempt := 0; attempt < 10; attempt++ {
		got := p.Next(attempt)
		if got > p.MaxInterval {
			t.Errorf("attempt %d: %v exceeds MaxInterval %v", attempt, got, p.MaxInterval)
		}
	}
}

func TestNext_NegativeAttempt_TreatedAsZero(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 200 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     1 * time.Minute,
		Jitter:          false,
	}
	if got := p.Next(-5); got != 200*time.Millisecond {
		t.Fatalf("expected 200ms for negative attempt, got %v", got)
	}
}

func TestNext_Jitter_DoesNotExceedCap(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		Multiplier:      2.0,
		MaxInterval:     4 * time.Second,
		Jitter:          true,
	}
	for i := 0; i < 100; i++ {
		for attempt := 0; attempt < 5; attempt++ {
			got := p.Next(attempt)
			// Allow up to 20 % jitter on top of MaxInterval.
			ceiling := time.Duration(float64(p.MaxInterval) * 1.21)
			if got > ceiling {
				t.Errorf("attempt %d: %v exceeds jitter ceiling %v", attempt, got, ceiling)
			}
		}
	}
}

func TestSequence_ReturnsCorrectLength(t *testing.T) {
	p := backoff.DefaultPolicy()
	seq := p.Sequence(5)
	if len(seq) != 5 {
		t.Fatalf("expected 5 intervals, got %d", len(seq))
	}
}

func TestDefaultPolicy_SanityCheck(t *testing.T) {
	p := backoff.DefaultPolicy()
	if p.InitialInterval <= 0 {
		t.Error("InitialInterval must be positive")
	}
	if p.Multiplier <= 1 {
		t.Error("Multiplier must be greater than 1")
	}
	if p.MaxInterval < p.InitialInterval {
		t.Error("MaxInterval must be >= InitialInterval")
	}
}
