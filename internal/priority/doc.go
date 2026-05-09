// Package priority provides severity scoring for drift results.
//
// A Scorer assigns a Level (low, medium, high, critical) to each
// drift.Result based on the number of detected diffs. Operators can
// use these levels to route alerts, suppress low-priority noise, or
// integrate with incident-management systems.
//
// Usage:
//
//	scorer := priority.New()
//	level  := scorer.Score(result)
//
//	// Sort a batch of results, highest priority first:
//	ranked := priority.Sort(results, scorer)
//	for _, r := range ranked {
//		fmt.Printf("%s [%s]\n", r.Result.Service, r.Level)
//	}
package priority
