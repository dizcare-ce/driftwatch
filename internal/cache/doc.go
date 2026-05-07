// Package cache provides a lightweight in-memory result cache for driftwatch.
//
// # Overview
//
// When a service definition has not changed between scheduler ticks, repeating
// a full drift comparison is wasteful. The cache stores the most recent
// [drift.Result] for each service, keyed by a SHA-256 checksum of its
// definition. A configurable TTL ensures stale entries are not served
// indefinitely even if the checksum happens to remain stable.
//
// # Usage
//
//	c := cache.New(5 * time.Minute)
//
//	sum, err := cache.Checksum(definition)
//	if err != nil { ... }
//
//	if entry, ok := c.Get(service, sum); ok {
//	    return entry.Result, nil
//	}
//
//	result, err := detector.Compare(definition, live)
//	c.Set(service, sum, result)
//	return result, nil
package cache
