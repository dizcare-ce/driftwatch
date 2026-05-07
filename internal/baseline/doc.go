// Package baseline provides approved-drift suppression for driftwatch.
//
// When an operator determines that a detected drift is intentional or
// acceptable, they can approve it. Subsequent runs will suppress
// notifications for that service until the drift changes.
//
// # Usage
//
//	store, err := baseline.New("/var/lib/driftwatch/baselines")
//	checker := baseline.NewChecker(store)
//
//	// Suppress check:
//	suppressed, err := checker.IsSuppressed(result)
//
//	// Approve current drift:
//	err = checker.Approve(result, "scaled down for maintenance")
//
//	// Revoke approval:
//	err = checker.Revoke("api")
//
// Baselines are keyed by service name and stored as JSON files under
// the configured directory. Each file contains the checksum of the
// drift state at approval time; if the drift changes the checksum
// will no longer match and notifications resume automatically.
package baseline
