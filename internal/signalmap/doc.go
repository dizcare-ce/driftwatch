// Package signalmap translates drift detection results into named severity
// signals suitable for forwarding to external alerting or observability
// systems.
//
// A Mapper is constructed with a signal name prefix and a priority.Scorer.
// Calling Map on a set of drift.Result values produces a []Signal slice
// containing only those results whose computed priority level meets or exceeds
// the caller-supplied minimum level.
//
// Example usage:
//
//	scorer := priority.New()
//	mapper := signalmap.New("driftwatch", scorer)
//	signals := mapper.Map(results, priority.LevelMedium)
//	for _, s := range signals {
//		fmt.Printf("[%s] %s: %s\n", s.Level, s.Name, s.Message)
//	}
package signalmap
