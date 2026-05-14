// Package dependency provides a directed dependency graph for driftwatch
// services. It allows the runner to understand which services are affected
// when an upstream dependency is found to have drifted, enabling targeted
// re-evaluation of dependant services without a full scan.
//
// Typical usage:
//
//	g := dependency.New()
//	g.Add("api", "auth")      // api depends on auth
//	g.Add("worker", "auth")   // worker also depends on auth
//
//	// When auth drifts, find everything that might be affected:
//	affected := g.Dependants("auth") // ["api", "worker"]
package dependency
