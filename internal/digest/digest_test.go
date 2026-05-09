package digest_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/digest"
)

func TestCompute_DeterministicForSameInput(t *testing.T) {
	v := map[string]any{"name": "api", "replicas": 3, "image": "nginx:1.25"}
	a, err := digest.Compute(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := digest.Compute(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a != b {
		t.Errorf("expected identical digests, got %q and %q", a, b)
	}
}

func TestCompute_DifferentForDifferentInput(t *testing.T) {
	a, _ := digest.Compute(map[string]any{"replicas": 1})
	b, _ := digest.Compute(map[string]any{"replicas": 2})
	if a == b {
		t.Error("expected different digests for different inputs")
	}
}

func TestCompute_StableAcrossKeyOrder(t *testing.T) {
	// Go map iteration is random; marshalling the same logical object twice
	// should still produce the same digest.
	v1 := map[string]any{"z": "last", "a": "first", "m": 42}
	v2 := map[string]any{"m": 42, "z": "last", "a": "first"}
	a, err := digest.Compute(v1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := digest.Compute(v2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a != b {
		t.Errorf("key order should not affect digest: %q vs %q", a, b)
	}
}

func TestCompute_NonEmptyHexString(t *testing.T) {
	sum, err := digest.Compute("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sum) != 64 {
		t.Errorf("expected 64-char hex string, got %d chars: %q", len(sum), sum)
	}
}

func TestEqual_MatchingDigests(t *testing.T) {
	if !digest.Equal("abc123", "abc123") {
		t.Error("expected Equal to return true for matching digests")
	}
}

func TestEqual_DifferentDigests(t *testing.T) {
	if digest.Equal("abc123", "xyz789") {
		t.Error("expected Equal to return false for different digests")
	}
}

func TestEqual_EmptyDigest_ReturnsFalse(t *testing.T) {
	if digest.Equal("", "") {
		t.Error("expected Equal to return false when both digests are empty")
	}
}
