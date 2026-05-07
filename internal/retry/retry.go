// Package retry provides a simple exponential-backoff retry mechanism
// used when transient errors occur during drift checks or source loading.
package retry

import (
	"context"
	"errors"
	"time"
)

// Policy defines how retries are attempted.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// MaxDelay caps the exponential back-off.
	MaxDelay time.Duration
	// Multiplier scales the delay after each failure (default 2.0).
	Multiplier float64
}

// DefaultPolicy returns a sensible out-of-the-box retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  3,
		InitialDelay: 250 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// Do calls fn up to p.MaxAttempts times, backing off between failures.
// It stops early if ctx is cancelled or fn returns a non-retryable error
// wrapped with Permanent.
func (p Policy) Do(ctx context.Context, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	mul := p.Multiplier
	if mul <= 0 {
		mul = 2.0
	}

	delay := p.InitialDelay
	var last error
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		last = fn()
		if last == nil {
			return nil
		}
		var pe *permanentError
		if errors.As(last, &pe) {
			return pe.cause
		}
		if attempt < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * mul)
			if delay > p.MaxDelay {
				delay = p.MaxDelay
			}
		}
	}
	return last
}

// permanentError wraps an error that should not be retried.
type permanentError struct{ cause error }

func (e *permanentError) Error() string { return e.cause.Error() }
func (e *permanentError) Unwrap() error { return e.cause }

// Permanent marks err so that Do stops retrying immediately.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &permanentError{cause: err}
}
