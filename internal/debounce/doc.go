// Package debounce provides a service-keyed debouncer that coalesces rapid
// successive triggers into a single delayed callback invocation.
//
// This is useful in driftwatch when a file watcher emits multiple change
// events in quick succession for the same service definition (e.g. due to
// editor save behaviour). Rather than triggering a full drift check on every
// event, the Debouncer waits for a quiet period before acting.
//
// Basic usage:
//
//	d := debounce.New(500*time.Millisecond, func(service string) {
//		// run drift check for service
//	})
//
//	// Called from a watcher callback:
//	d.Trigger(serviceName)
package debounce
