// Package patchplan analyses drift.Result values and produces an ordered,
// human-readable remediation plan that operators can follow to bring deployed
// services back into alignment with their source definitions.
//
// Usage:
//
//	results := detector.Compare(definitions, live)
//	plans   := patchplan.Build(results)
//	fmt.Println(patchplan.Format(plans))
//
// Build returns one Plan per drifted service, each containing an ordered list
// of Actions. Actions are sorted alphabetically by field name so that output
// is deterministic across runs. Services with no drift are omitted entirely.
//
// Format renders the plans as a plain-text string. When there is nothing to
// remediate it returns a short "all services are in sync" message instead.
package patchplan
