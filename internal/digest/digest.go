// Package digest computes and compares content fingerprints for service
// definitions, enabling driftwatch to detect when a source file has changed
// without performing a full drift comparison.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// Entry holds the SHA-256 fingerprint of a single service definition.
type Entry struct {
	Service  string `json:"service"`
	Checksum string `json:"checksum"`
}

// Compute returns a stable SHA-256 hex digest for the given value by
// marshalling it to canonical JSON first. The sort order of map keys is
// normalised by the JSON encoder, so structurally identical definitions
// always produce the same digest regardless of field insertion order.
func Compute(v any) (string, error) {
	b, err := stableJSON(v)
	if err != nil {
		return "", fmt.Errorf("digest: marshal: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// Equal reports whether two digests match.
func Equal(a, b string) bool { return a != "" && a == b }

// stableJSON marshals v to JSON with sorted map keys.
func stableJSON(v any) ([]byte, error) {
	// Marshal then unmarshal into an ordered structure so that map keys are
	// sorted deterministically by the standard library encoder.
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var ordered any
	if err := json.Unmarshal(raw, &ordered); err != nil {
		return nil, err
	}
	sortKeys(ordered)
	return json.Marshal(ordered)
}

// sortKeys recursively sorts map keys within an unmarshalled JSON value.
func sortKeys(v any) {
	switch val := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sortKeys(val[k])
		}
	case []any:
		for _, item := range val {
			sortKeys(item)
		}
	}
}
