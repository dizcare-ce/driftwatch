// Package slatracker tracks service-level agreement (SLA) compliance for
// driftwatch-monitored services.
//
// A Tracker accumulates drift check results over a configurable rolling time
// window and computes the percentage of checks during which each service was
// drift-free (compliant). A service is considered to meet its SLA when its
// compliance percentage is at or above a configurable threshold.
//
// Typical usage:
//
//	tr := slatracker.New(24*time.Hour, 99.0)
//
//	// After each detector run:
//	tr.Record(result)
//
//	// Query compliance:
//	rep, err := tr.Report("api-gateway")
//	if err == nil && !rep.MeetsSLA {
//		log.Printf("SLA breach: %.1f%% compliant", rep.CompliancePct)
//	}
//
//	// Periodically discard entries older than the window:
//	tr.Prune()
package slatracker
