package drift

import (
	"fmt"

	"github.com/user/driftwatch/internal/source"
)

// Status represents the drift state of a single field.
type Status string

const (
	StatusMatch   Status = "match"
	StatusDrifted Status = "drifted"
	StatusMissing Status = "missing"
)

// FieldDiff describes a single field-level difference.
type FieldDiff struct {
	Field    string
	Expected interface{}
	Actual   interface{}
	Status   Status
}

// Result holds the drift report for one service.
type Result struct {
	ServiceName string
	Drifted     bool
	Diffs       []FieldDiff
}

// Detector compares live service state against source definitions.
type Detector struct{}

// NewDetector creates a new Detector.
func NewDetector() *Detector {
	return &Detector{}
}

// Compare checks a live metadata map against a source ServiceDefinition.
// live represents the currently deployed service's key/value metadata.
func (d *Detector) Compare(def source.ServiceDefinition, live map[string]interface{}) Result {
	result := Result{
		ServiceName: def.Name,
		Drifted:     false,
	}

	for key, expectedVal := range def.Spec {
		actualVal, exists := live[key]
		if !exists {
			result.Diffs = append(result.Diffs, FieldDiff{
				Field:    key,
				Expected: expectedVal,
				Actual:   nil,
				Status:   StatusMissing,
			})
			result.Drifted = true
			continue
		}

		expectedStr := fmt.Sprintf("%v", expectedVal)
		actualStr := fmt.Sprintf("%v", actualVal)

		if expectedStr != actualStr {
			result.Diffs = append(result.Diffs, FieldDiff{
				Field:    key,
				Expected: expectedVal,
				Actual:   actualVal,
				Status:   StatusDrifted,
			})
			result.Drifted = true
		} else {
			result.Diffs = append(result.Diffs, FieldDiff{
				Field:    key,
				Expected: expectedVal,
				Actual:   actualVal,
				Status:   StatusMatch,
			})
		}
	}

	return result
}
