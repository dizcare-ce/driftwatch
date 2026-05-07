// Package rollup provides aggregation of drift detection results across
// multiple services into a unified summary.
//
// Use Aggregate to combine a []drift.Result slice into a Summary value that
// reports total, drifted, and clean counts together with per-service details.
// The resulting Summary can be rendered as a human-readable string or
// inspected programmatically by downstream reporters and notifiers.
//
// Example:
//
//	summary := rollup.Aggregate(results)
//	fmt.Println(summary)
package rollup
