// Package groupby partitions a slice of drift.Result values into named
// buckets using a caller-supplied KeyFunc.
//
// A KeyFunc is any function that maps a drift.Result to a string key.
// Two built-in helpers are provided:
//
//   - ByService — groups by the result's Service field (the default for
//     most dashboards and reports).
//
// Results whose KeyFunc returns an empty string are placed in a
// configurable fallback bucket (default: "(ungrouped)").
//
// Groups are returned sorted alphabetically by key so that callers
// receive a deterministic ordering regardless of input order.
//
// Usage:
//
//	g, err := groupby.New(groupby.ByService, "")
//	groups := g.Apply(results)
//	for _, grp := range groups {
//		fmt.Println(grp.Key, len(grp.Results))
//	}
package groupby
