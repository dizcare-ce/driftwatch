package drift

import (
	"testing"

	"github.com/user/driftwatch/internal/source"
)

func makeDefinition(name string, spec map[string]interface{}) source.ServiceDefinition {
	return source.ServiceDefinition{
		Name: name,
		Spec: spec,
	}
}

func TestCompare_NoDrift(t *testing.T) {
	d := NewDetector()
	def := makeDefinition("api", map[string]interface{}{
		"image":    "nginx:1.25",
		"replicas": 3,
	})
	live := map[string]interface{}{
		"image":    "nginx:1.25",
		"replicas": 3,
	}

	result := d.Compare(def, live)

	if result.Drifted {
		t.Errorf("expected no drift, got drifted=true")
	}
	for _, diff := range result.Diffs {
		if diff.Status != StatusMatch {
			t.Errorf("field %q: expected status %q, got %q", diff.Field, StatusMatch, diff.Status)
		}
	}
}

func TestCompare_DriftedValue(t *testing.T) {
	d := NewDetector()
	def := makeDefinition("api", map[string]interface{}{
		"image": "nginx:1.25",
	})
	live := map[string]interface{}{
		"image": "nginx:1.19",
	}

	result := d.Compare(def, live)

	if !result.Drifted {
		t.Fatal("expected drift, got drifted=false")
	}
	if len(result.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(result.Diffs))
	}
	if result.Diffs[0].Status != StatusDrifted {
		t.Errorf("expected status %q, got %q", StatusDrifted, result.Diffs[0].Status)
	}
}

func TestCompare_MissingField(t *testing.T) {
	d := NewDetector()
	def := makeDefinition("api", map[string]interface{}{
		"image":    "nginx:1.25",
		"replicas": 3,
	})
	live := map[string]interface{}{
		"image": "nginx:1.25",
		// replicas is missing
	}

	result := d.Compare(def, live)

	if !result.Drifted {
		t.Fatal("expected drift due to missing field")
	}

	var missingDiff *FieldDiff
	for i := range result.Diffs {
		if result.Diffs[i].Field == "replicas" {
			missingDiff = &result.Diffs[i]
		}
	}
	if missingDiff == nil {
		t.Fatal("expected diff for field 'replicas'")
	}
	if missingDiff.Status != StatusMissing {
		t.Errorf("expected status %q, got %q", StatusMissing, missingDiff.Status)
	}
}

func TestCompare_ServiceName(t *testing.T) {
	d := NewDetector()
	def := makeDefinition("my-service", map[string]interface{}{})
	result := d.Compare(def, map[string]interface{}{})
	if result.ServiceName != "my-service" {
		t.Errorf("expected service name 'my-service', got %q", result.ServiceName)
	}
}
