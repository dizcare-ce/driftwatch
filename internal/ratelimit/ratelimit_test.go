package ratelimit_test

import (
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/ratelimit"
)

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	l := ratelimit.New(5 * time.Minute)
	if !l.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	l := ratelimit.New(5 * time.Minute)
	l.Allow("svc-a")
	if l.Allow("svc-a") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_SecondCallAfterCooldownPasses(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("svc-a")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("svc-a") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_IndependentServicesDoNotInterfere(t *testing.T) {
	l := ratelimit.New(5 * time.Minute)
	l.Allow("svc-a")
	if !l.Allow("svc-b") {
		t.Fatal("expected independent service to be allowed")
	}
}

func TestReset_ClearsTimestampForService(t *testing.T) {
	l := ratelimit.New(5 * time.Minute)
	l.Allow("svc-a")
	l.Reset("svc-a")
	if !l.Allow("svc-a") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAllServices(t *testing.T) {
	l := ratelimit.New(5 * time.Minute)
	l.Allow("svc-a")
	l.Allow("svc-b")
	l.ResetAll()
	if !l.Allow("svc-a") || !l.Allow("svc-b") {
		t.Fatal("expected both services to be allowed after ResetAll")
	}
}

func TestAllow_ZeroCooldownAlwaysPasses(t *testing.T) {
	l := ratelimit.New(0)
	l.Allow("svc-a")
	if !l.Allow("svc-a") {
		t.Fatal("expected zero cooldown to always allow")
	}
}
