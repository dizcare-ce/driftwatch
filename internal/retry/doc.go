// Package retry implements a configurable exponential-backoff retry helper
// for use throughout driftwatch whenever operations may fail transiently.
//
// Basic usage:
//
//	p := retry.DefaultPolicy()
//	err := p.Do(ctx, func() error {
//		return doSomethingThatMightFail()
//	})
//
// To prevent a specific error from being retried, wrap it with Permanent:
//
//	return retry.Permanent(fmt.Errorf("config invalid: %w", err))
//
// The policy is safe for concurrent use; each call to Do maintains its
// own state.
package retry
