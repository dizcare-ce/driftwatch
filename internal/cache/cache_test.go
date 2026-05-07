package cache_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/cache"
	"github.com/driftwatch/internal/drift"
)

func makeResult(name string, drifted bool) drift.Result {
	return drift.Result{
		Service: name,
		Drifted: drifted,
	}
}

func TestGet_MissingEntry_ReturnsFalse(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("svc", "abc")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	c := cache.New(time.Minute)
	r := makeResult("api", false)
	c.Set("api", "sum1", r)

	e, ok := c.Get("api", "sum1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if e.Result.Service != "api" {
		t.Errorf("unexpected service: %s", e.Result.Service)
	}
}

func TestGet_WrongChecksum_ReturnsFalse(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("api", "sum1", makeResult("api", false))

	_, ok := c.Get("api", "sum2")
	if ok {
		t.Fatal("expected cache miss on checksum mismatch")
	}
}

func TestGet_ExpiredEntry_ReturnsFalse(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("api", "sum1", makeResult("api", false))

	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("api", "sum1")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestGet_ZeroTTL_NeverExpires(t *testing.T) {
	c := cache.New(0)
	c.Set("api", "sum1", makeResult("api", false))

	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("api", "sum1")
	if !ok {
		t.Fatal("expected cache hit with zero TTL")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("api", "sum1", makeResult("api", false))
	c.Invalidate("api")

	_, ok := c.Get("api", "sum1")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestLen_ReflectsEntryCount(t *testing.T) {
	c := cache.New(time.Minute)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	c.Set("a", "s", makeResult("a", false))
	c.Set("b", "s", makeResult("b", true))
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	c.Invalidate("a")
	if c.Len() != 1 {
		t.Fatalf("expected 1, got %d", c.Len())
	}
}
