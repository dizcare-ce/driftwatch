// Package driftbudget enforces a configurable cap on the total number of
// drifted fields observed within a rolling time window.
//
// Use New or NewFromPolicy to create a Budget, then call Record after each
// detection run. If the cumulative drift count would exceed the configured
// limit, Record returns a non-nil error that the caller can use to trigger
// escalation, pause automated remediation, or raise an alert.
//
// Example:
//
//	b, err := driftbudget.NewFromPolicy(driftbudget.DefaultPolicy())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := b.Record(results); err != nil {
//		log.Printf("drift budget exhausted: %v", err)
//	}
package driftbudget
