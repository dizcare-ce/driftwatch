// Package digest provides content fingerprinting for service definitions.
//
// # Overview
//
// Compute produces a stable SHA-256 hex digest for any JSON-serialisable
// value. The digest is canonical: map keys are sorted before hashing, so
// two structurally identical definitions always yield the same fingerprint
// regardless of the order in which their fields were populated.
//
// # Usage
//
//	sum, err := digest.Compute(definition)
//	if err != nil { ... }
//
//	if digest.Equal(cached, sum) {
//	    // definition has not changed — skip expensive comparison
//	}
//
// # Integration
//
// The digest package is intentionally dependency-free. Callers such as the
// cache and snapshot packages use it to avoid redundant drift checks when
// a source definition has not changed since the last run.
package digest
