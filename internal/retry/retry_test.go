package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"driftwatch/internal/retry"
)

func fastPolicy(attempts int) retry.Policy {
	return retry.Policy{
		MaxAttempts:  attempts,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := fastPolicy(3).Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnTransientError(t *testing.T) {
	calls := 0
	transient := errors.New("transient")
	err := fastPolicy(3).Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return transient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsMaxAttempts(t *testing.T) {
	calls := 0
	sentinel := errors.New("always fails")
	err := fastPolicy(3).Do(context.Background(), func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_PermanentError_StopsImmediately(t *testing.T) {
	calls := 0
	perm := errors.New("permanent")
	err := fastPolicy(5).Do(context.Background(), func() error {
		calls++
		return retry.Permanent(perm)
	})
	if !errors.Is(err, perm) {
		t.Fatalf("expected permanent error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_ContextCancelled_StopsRetrying(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := fastPolicy(5).Do(ctx, func() error {
		calls++
		return errors.New("transient")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_DefaultPolicy_ReturnsPolicy(t *testing.T) {
	p := retry.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Fatalf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
}
