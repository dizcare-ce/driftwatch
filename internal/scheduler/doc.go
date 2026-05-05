// Package scheduler provides a simple interval-based job runner for driftwatch.
//
// A Scheduler executes a user-supplied job function immediately upon start,
// then repeats on a fixed interval until the provided context is cancelled.
// Transient job errors are logged but do not halt the scheduler, allowing
// drift checks to recover on the next tick without operator intervention.
//
// Example usage:
//
//	s := scheduler.New(5*time.Minute, func(ctx context.Context) error {
//		// run drift detection
//		return nil
//	}, nil)
//
//	if err := s.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package scheduler
