// Package ratelimit implements a per-service cooldown limiter used to
// suppress repeated drift notifications for the same service within a
// configurable time window.
//
// Usage:
//
//	limiter := ratelimit.New(10 * time.Minute)
//
//	if limiter.Allow(result.Service) {
//		notifier.Notify(result)
//	}
//
// The cooldown window is measured from the most recent allowed call.
// Calls arriving before the window expires are silently dropped.
// Use Reset or ResetAll to clear recorded timestamps, for example
// after a configuration reload.
package ratelimit
