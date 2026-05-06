// Package history provides persistent storage of drift check runs.
//
// Each time the runner completes a check cycle, the results can be
// appended to the history store. Entries are written as timestamped
// JSON files under a configurable directory, making them easy to
// inspect, archive, or feed into external tooling.
//
// Usage:
//
//	h, err := history.New("/var/lib/driftwatch/history")
//	if err != nil { ... }
//
//	// after a check run:
//	if err := h.Record(results); err != nil { ... }
//
//	// retrieve the last 10 runs:
//	entries, err := h.Last(10)
package history
