package driftbudget_test

import (
	"testing"
	"time"

	"driftwatch/internal/driftbudget"
)

func TestDefaultPolicy_IsValid(t *testing.T) {
	p := driftbudget.DefaultPolicy()
	if err := p.Validate(); err != nil {
		t.Fatalf("default policy should be valid: %v", err)
	}
}

func TestValidate_ZeroLimit_ReturnsError(t *testing.T) {
	p := driftbudget.Policy{Limit: 0, Window: time.Minute}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestValidate_ZeroWindow_ReturnsError(t *testing.T) {
	p := driftbudget.Policy{Limit: 10, Window: 0}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewFromPolicy_ValidPolicy_ReturnsBudget(t *testing.T) {
	p := driftbudget.Policy{Limit: 20, Window: 30 * time.Minute}
	b, err := driftbudget.NewFromPolicy(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil budget")
	}
}

func TestNewFromPolicy_InvalidPolicy_ReturnsError(t *testing.T) {
	p := driftbudget.Policy{Limit: -1, Window: time.Minute}
	_, err := driftbudget.NewFromPolicy(p)
	if err == nil {
		t.Fatal("expected error for invalid policy")
	}
}

func TestNewFromPolicy_RemainingMatchesLimit(t *testing.T) {
	p := driftbudget.Policy{Limit: 15, Window: time.Hour}
	b, _ := driftbudget.NewFromPolicy(p)
	if got := b.Remaining(); got != 15 {
		t.Fatalf("expected 15 remaining, got %d", got)
	}
}
