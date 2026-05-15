package patchplan_test

import (
	"strings"
	"testing"

	"driftwatch/internal/drift"
	"driftwatch/internal/patchplan"
)

func makeResult(service string, diffs []drift.Diff) drift.Result {
	return drift.Result{Service: service, Diffs: diffs}
}

func TestBuild_NoDrift_ReturnsEmpty(t *testing.T) {
	results := []drift.Result{makeResult("api", nil)}
	plans := patchplan.Build(results)
	if len(plans) != 0 {
		t.Fatalf("expected no plans, got %d", len(plans))
	}
}

func TestBuild_SingleDrift_ProducesAction(t *testing.T) {
	results := []drift.Result{
		makeResult("api", []drift.Diff{
			{Field: "replicas", Expected: 3, Actual: 1},
		}),
	}
	plans := patchplan.Build(results)
	if len(plans) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(plans))
	}
	if plans[0].Service != "api" {
		t.Errorf("unexpected service name: %s", plans[0].Service)
	}
	if len(plans[0].Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(plans[0].Actions))
	}
	if plans[0].Actions[0].Field != "replicas" {
		t.Errorf("unexpected field: %s", plans[0].Actions[0].Field)
	}
}

func TestBuild_SortedByServiceAndField(t *testing.T) {
	results := []drift.Result{
		makeResult("zebra", []drift.Diff{
			{Field: "mem", Expected: "512Mi", Actual: "256Mi"},
			{Field: "cpu", Expected: "200m", Actual: "100m"},
		}),
		makeResult("alpha", []drift.Diff{
			{Field: "image", Expected: "v2", Actual: "v1"},
		}),
	}
	plans := patchplan.Build(results)
	if plans[0].Service != "alpha" {
		t.Errorf("expected alpha first, got %s", plans[0].Service)
	}
	if plans[1].Actions[0].Field != "cpu" {
		t.Errorf("expected cpu before mem, got %s", plans[1].Actions[0].Field)
	}
}

func TestBuild_MissingExpected_DescribesRemoval(t *testing.T) {
	results := []drift.Result{
		makeResult("svc", []drift.Diff{
			{Field: "debug", Expected: nil, Actual: true},
		}),
	}
	plans := patchplan.Build(results)
	desc := plans[0].Actions[0].Description
	if !strings.Contains(desc, "remove") {
		t.Errorf("expected removal description, got: %s", desc)
	}
}

func TestBuild_MissingActual_DescribesAddition(t *testing.T) {
	results := []drift.Result{
		makeResult("svc", []drift.Diff{
			{Field: "timeout", Expected: "30s", Actual: nil},
		}),
	}
	plans := patchplan.Build(results)
	desc := plans[0].Actions[0].Description
	if !strings.Contains(desc, "absent") {
		t.Errorf("expected absent description, got: %s", desc)
	}
}

func TestFormat_NoDrift_ReturnsSyncMessage(t *testing.T) {
	out := patchplan.Format(nil)
	if !strings.Contains(out, "in sync") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormat_WithPlans_ContainsServiceAndAction(t *testing.T) {
	plans := []patchplan.Plan{
		{
			Service: "api",
			Actions: []patchplan.Action{
				{Field: "replicas", Description: "update field \"replicas\" from 1 to 3"},
			},
		},
	}
	out := patchplan.Format(plans)
	if !strings.Contains(out, "api") {
		t.Errorf("expected service name in output")
	}
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected field name in output")
	}
}
