// Package metrics provides lightweight in-process counters that track
// drift-detection run outcomes over the daemon lifetime.
//
// Usage:
//
//	counters := metrics.New()
//
//	// after each detection run:
//	counters.RecordRun(hasDrift, err)
//
//	// emit a human-readable summary:
//	counters.Write(os.Stdout)
//
// Counters are safe for concurrent use. All fields use atomic operations
// so callers do not need additional synchronisation.
package metrics
