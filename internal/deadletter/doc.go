// Package deadletter provides a persistent dead-letter queue for drift check
// failures that could not be delivered to a downstream sink (e.g. a notifier
// or audit log that was temporarily unavailable).
//
// Each failed event is written as a JSON file under a configurable directory.
// Entries can be inspected with All and cleared with Purge once the root cause
// has been resolved.
//
// Usage:
//
//	q, err := deadletter.New("/var/lib/driftwatch/dlq")
//	if err != nil { ... }
//
//	// record a failure
//	q.Push(deadletter.Entry{
//		Service:  "api-gateway",
//		Reason:   "notifier unavailable",
//		Payload:  resultJSON,
//		Attempts: 3,
//	})
//
//	// inspect later
//	entries, _ := q.All()
//
//	// clear once resolved
//	q.Purge()
package deadletter
