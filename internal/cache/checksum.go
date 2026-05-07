package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Checksum returns a stable SHA-256 hex digest of v by marshalling it to JSON.
// It is used to detect whether a service definition has changed between runs so
// that the cache can decide whether a stored result is still valid.
func Checksum(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("cache: checksum marshal: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
