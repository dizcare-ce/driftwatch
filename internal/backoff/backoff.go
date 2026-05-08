// Package backoff provides configurable exponential back-off helpers used
// when retrying failed checks or notification deliveries.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Policy defines the parameters for an exponential back-off strategy.
type Policy struct {
	// InitialInterval is the wait time after the first failure.
	InitialInterval time.Duration
	// Multiplier is applied to the current interval on each retry.
	Multiplier float64
	// MaxInterval caps the computed interval.
	MaxInterval time.Duration
	// Jitter adds a random fraction of the current interval to spread load.
	Jitter bool
}

// DefaultPolicy returns a Policy suitable for most drift-check retries.
func DefaultPolicy() Policy {
	return Policy{
		InitialInterval: 500 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     30 * time.Second,
		Jitter:          true,
	}
}

// Next returns the wait duration for the given attempt number (0-based).
// Attempt 0 returns InitialInterval; each subsequent attempt multiplies by
// Multiplier and is capped at MaxInterval.
func (p Policy) Next(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	interval := float64(p.InitialInterval) * math.Pow(p.Multiplier, float64(attempt))
	if max := float64(p.MaxInterval); interval > max {
		interval = max
	}

	if p.Jitter {
		// Add up to 20 % random jitter.
		interval += interval * 0.2 * rand.Float64() //nolint:gosec
	}

	return time.Duration(interval)
}

// Sequence returns the first n intervals for a policy, useful in tests and
// pre-flight validation.
func (p Policy) Sequence(n int) []time.Duration {
	out := make([]time.Duration, n)
	for i := range out {
		out[i] = p.Next(i)
	}
	return out
}
