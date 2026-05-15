// Package shadowmode provides a shadow-run facility for driftwatch.
//
// A shadow run executes the full drift-detection pipeline but suppresses
// all side-effects (notifications, remediations, audit writes). Results
// are buffered in memory so operators can compare shadow output against
// live output before promoting a configuration change.
//
// Typical usage:
//
//	sr := shadowmode.New(20)        // keep last 20 runs
//	sr.Run(ctx, myDetectFunc)       // non-blocking; safe to call concurrently
//	records := sr.Last(5)           // inspect most-recent 5 outcomes
//
// Panics inside the supplied function are recovered and surfaced as
// errors so the host process is never destabilised by a shadow run.
package shadowmode
