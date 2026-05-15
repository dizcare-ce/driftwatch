// Package costestimator assigns a weighted cost score to detected drift,
// enabling operators to prioritise remediation of the most impactful
// configuration discrepancies.
//
// # Overview
//
// Each drifted field contributes a cost determined by its entry in a
// FieldWeight map. Fields not explicitly listed fall back to a configurable
// default cost. The Estimator aggregates per-field costs into a per-service
// total and returns estimates sorted by descending cost.
//
// # Usage
//
//	e := costestimator.New(
//		costestimator.FieldWeight{"replicas": 5.0, "image": 3.0},
//		1.0, // default cost for unlisted fields
//	)
//	estimates := e.Compute(results)
//	fmt.Print(costestimator.Format(estimates))
//
// # Configuration
//
// Use DefaultConfig and NewFromConfig to drive the estimator from the
// application's YAML configuration file.
package costestimator
