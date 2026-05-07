package debounce_test

import (
	"sync"
	"testing"
	"time"

	"driftwatch/internal/debounce"
)

func TestTrigger_FiresAfterWait(t *testing.T) {
	var mu sync.Mutex
	fired := []string{}

	d := debounce.New(30*time.Millisecond, func(svc string) {
		mu.Lock()
		fired = append(fired, svc)
		mu.Unlock()
	})

	d.Trigger("api")
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 || fired[0] != "api" {
		t.Fatalf("expected [api], got %v", fired)
	}
}

func TestTrigger_ResetsTimer_OnRapidCalls(t *testing.T) {
	var mu sync.Mutex
	count := 0

	d := debounce.New(50*time.Millisecond, func(_ string) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Trigger three times quickly; only one callback should fire.
	d.Trigger("svc")
	time.Sleep(20 * time.Millisecond)
	d.Trigger("svc")
	time.Sleep(20 * time.Millisecond)
	d.Trigger("svc")
	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected 1 callback, got %d", count)
	}
}

func TestCancel_PreventsFiring(t *testing.T) {
	fired := false
	d := debounce.New(40*time.Millisecond, func(_ string) { fired = true })

	d.Trigger("svc")
	d.Cancel("svc")
	time.Sleep(70 * time.Millisecond)

	if fired {
		t.Fatal("callback should not have fired after Cancel")
	}
}

func TestPending_ReflectsState(t *testing.T) {
	d := debounce.New(50*time.Millisecond, func(_ string) {})

	if d.Pending("svc") {
		t.Fatal("expected no pending timer before Trigger")
	}

	d.Trigger("svc")
	if !d.Pending("svc") {
		t.Fatal("expected pending timer after Trigger")
	}

	time.Sleep(80 * time.Millisecond)
	if d.Pending("svc") {
		t.Fatal("expected no pending timer after callback fired")
	}
}

func TestTrigger_IndependentServices(t *testing.T) {
	var mu sync.Mutex
	fired := map[string]int{}

	d := debounce.New(30*time.Millisecond, func(svc string) {
		mu.Lock()
		fired[svc]++
		mu.Unlock()
	})

	d.Trigger("alpha")
	d.Trigger("beta")
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if fired["alpha"] != 1 || fired["beta"] != 1 {
		t.Fatalf("expected one fire each, got %v", fired)
	}
}
