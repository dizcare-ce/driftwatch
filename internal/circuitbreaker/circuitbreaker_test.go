package circuitbreaker_test

import (
	"testing"
	"time"

	"driftwatch/internal/circuitbreaker"
)

func TestAllow_InitiallyPermits(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if !b.Allow("svc") {
		t.Fatal("expected Allow to return true before any failures")
	}
}

func TestAllow_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure("svc")
	}
	if b.Allow("svc") {
		t.Fatal("expected Allow to return false after threshold reached")
	}
}

func TestAllow_ReopensAfterTimeout(t *testing.T) {
	b := circuitbreaker.New(2, 20*time.Millisecond)
	b.RecordFailure("svc")
	b.RecordFailure("svc")

	time.Sleep(30 * time.Millisecond)

	if !b.Allow("svc") {
		t.Fatal("expected Allow to return true after reset timeout")
	}
}

func TestRecordSuccess_ResetsCounts(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	b.RecordFailure("svc")
	b.RecordFailure("svc")
	b.RecordSuccess("svc")
	// one more failure should not open the circuit
	b.RecordFailure("svc")
	if !b.Allow("svc") {
		t.Fatal("expected circuit to remain closed after success reset")
	}
}

func TestStateOf_ReflectsState(t *testing.T) {
	b := circuitbreaker.New(2, time.Second)
	if b.StateOf("svc") != circuitbreaker.StateClosed {
		t.Fatal("expected StateClosed initially")
	}
	b.RecordFailure("svc")
	b.RecordFailure("svc")
	if b.StateOf("svc") != circuitbreaker.StateOpen {
		t.Fatal("expected StateOpen after threshold")
	}
}

func TestIndependentServices_DoNotInterfere(t *testing.T) {
	b := circuitbreaker.New(2, time.Second)
	b.RecordFailure("alpha")
	b.RecordFailure("alpha")

	if !b.Allow("beta") {
		t.Fatal("expected beta to be unaffected by alpha failures")
	}
}

func TestAllow_BelowThreshold_StillPermits(t *testing.T) {
	b := circuitbreaker.New(5, time.Second)
	for i := 0; i < 4; i++ {
		b.RecordFailure("svc")
	}
	if !b.Allow("svc") {
		t.Fatal("expected Allow to return true when below threshold")
	}
}
