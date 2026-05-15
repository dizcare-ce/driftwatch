// Package ruleengine evaluates named drift rules against detection results,
// allowing operators to define custom pass/fail conditions beyond the default
// threshold-based checks.
package ruleengine

import (
	"fmt"
	"regexp"

	"github.com/driftwatch/internal/drift"
)

// Rule describes a single evaluation rule.
type Rule struct {
	// Name is a human-readable identifier for the rule.
	Name string
	// ServicePattern is a regex matched against the service name.
	// An empty string matches all services.
	ServicePattern string
	// MaxDiffs is the maximum number of diffs allowed before the rule fails.
	// A value of -1 disables the diff-count check.
	MaxDiffs int
	// RequireClean demands zero diffs when true.
	RequireClean bool
}

// Violation records a rule that was not satisfied.
type Violation struct {
	Rule    string
	Service string
	Reason  string
}

// Engine holds a set of rules and evaluates them against drift results.
type Engine struct {
	rules   []Rule
	compiled []*regexp.Regexp
}

// New creates an Engine from the provided rules.
// Returns an error if any ServicePattern fails to compile.
func New(rules []Rule) (*Engine, error) {
	compiled := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		pattern := r.ServicePattern
		if pattern == "" {
			pattern = ".*"
		}
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return nil, fmt.Errorf("ruleengine: rule %q has invalid pattern: %w", r.Name, err)
		}
		compiled[i] = re
	}
	return &Engine{rules: rules, compiled: compiled}, nil
}

// Evaluate checks every rule against the provided result.
// It returns all violations found; a nil slice means all rules passed.
func (e *Engine) Evaluate(result drift.Result) []Violation {
	var violations []Violation
	for i, rule := range e.rules {
		if !e.compiled[i].MatchString(result.Service) {
			continue
		}
		ndiffs := len(result.Diffs)
		if rule.RequireClean && ndiffs > 0 {
			violations = append(violations, Violation{
				Rule:    rule.Name,
				Service: result.Service,
				Reason:  fmt.Sprintf("requires clean state but found %d diff(s)", ndiffs),
			})
			continue
		}
		if rule.MaxDiffs >= 0 && ndiffs > rule.MaxDiffs {
			violations = append(violations, Violation{
				Rule:    rule.Name,
				Service: result.Service,
				Reason:  fmt.Sprintf("found %d diff(s), max allowed is %d", ndiffs, rule.MaxDiffs),
			})
		}
	}
	return violations
}
