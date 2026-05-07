package throttle_test

import (
	"testing"
	"time"

	"driftwatch/internal/throttle"
)

func TestAllow_UnderLimit_Passes(t *testing.T) {
	th := throttle.New(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !th.Allow() {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
}

func TestAllow_AtLimit_Blocked(t *testing.T) {
	th := throttle.New(time.Minute, 2)
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("third call should be blocked")
	}
}

func TestAllow_WindowExpiry_ResetsQuota(t *testing.T) {
	now := time.Now()
	th := throttle.New(100*time.Millisecond, 1)
	th.Allow() // consume the single slot

	// Advance internal clock past the window by swapping the now func.
	// We rebuild a fresh throttle with a stubbed clock instead.
	_ = now

	// Use a real short window and sleep past it.
	th2 := throttle.New(50*time.Millisecond, 1)
	if !th2.Allow() {
		t.Fatal("first call should pass")
	}
	if th2.Allow() {
		t.Fatal("second call within window should be blocked")
	}
	time.Sleep(60 * time.Millisecond)
	if !th2.Allow() {
		t.Fatal("call after window expiry should pass")
	}
}

func TestRemaining_ReflectsUsage(t *testing.T) {
	th := throttle.New(time.Minute, 5)
	if got := th.Remaining(); got != 5 {
		t.Fatalf("expected 5 remaining, got %d", got)
	}
	th.Allow()
	th.Allow()
	if got := th.Remaining(); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestRemaining_NeverNegative(t *testing.T) {
	th := throttle.New(time.Minute, 1)
	th.Allow()
	th.Allow() // blocked but should not corrupt state
	if got := th.Remaining(); got != 0 {
		t.Fatalf("expected 0 remaining, got %d", got)
	}
}

func TestReset_ClearsQuota(t *testing.T) {
	th := throttle.New(time.Minute, 2)
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("should be blocked before reset")
	}
	th.Reset()
	if !th.Allow() {
		t.Fatal("should be allowed after reset")
	}
}
