package drift

import (
	"strings"
	"testing"
)

func makeDriftResult(drifted bool, changes []Change) *DriftResult {
	return &DriftResult{
		Release:   "my-app",
		Namespace: "production",
		Drifted:   drifted,
		Changes:   changes,
	}
}

func TestSummary_NoDrift(t *testing.T) {
	r := makeDriftResult(false, nil)
	out := r.Summary()

	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected 'No drift detected' in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "my-app") {
		t.Errorf("expected release name in summary, got:\n%s", out)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	changes := []Change{
		{Key: "replicaCount", LiveVal: 3, ChartVal: 1},
		{Key: "image.tag", LiveVal: "v1.2.0", ChartVal: "v1.1.0"},
	}
	r := makeDriftResult(true, changes)
	out := r.Summary()

	if !strings.Contains(out, "Drift detected") {
		t.Errorf("expected 'Drift detected' in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "replicaCount") {
		t.Errorf("expected key 'replicaCount' in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "image.tag") {
		t.Errorf("expected key 'image.tag' in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "2 change(s)") {
		t.Errorf("expected change count in summary, got:\n%s", out)
	}
}

func TestExitCode_NoDrift(t *testing.T) {
	r := makeDriftResult(false, nil)
	if code := r.ExitCode(); code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestExitCode_WithDrift(t *testing.T) {
	r := makeDriftResult(true, []Change{{Key: "x", LiveVal: 1, ChartVal: 2}})
	if code := r.ExitCode(); code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}
