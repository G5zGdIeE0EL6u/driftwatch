package drift_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/helm"
	helmchart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

func buildReleaseWithDefaults(name string, chartDefaults, userValues map[string]interface{}) *release.Release {
	return &release.Release{
		Name: name,
		Chart: &helmchart.Chart{
			Metadata: &helmchart.Metadata{Name: name, Version: "1.0.0"},
			Values:   chartDefaults,
		},
		Config: userValues,
	}
}

func TestDetectChartDrift_NoDrift(t *testing.T) {
	rel := buildReleaseWithDefaults("myapp", map[string]interface{}{
		"replicaCount": float64(2),
	}, map[string]interface{}{
		"replicaCount": float64(2),
	})

	info := helm.ChartInfo{Name: "myapp", Version: "1.0.0"}
	results := drift.DetectChartDrift(rel, info)
	if len(results) != 0 {
		t.Fatalf("expected no drift, got %d results", len(results))
	}
}

func TestDetectChartDrift_ValueOverridden(t *testing.T) {
	rel := buildReleaseWithDefaults("myapp", map[string]interface{}{
		"replicaCount": float64(1),
	}, map[string]interface{}{
		"replicaCount": float64(3),
	})

	info := helm.ChartInfo{Name: "myapp", Version: "1.0.0"}
	results := drift.DetectChartDrift(rel, info)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Key != "replicaCount" {
		t.Errorf("unexpected key: %s", results[0].Key)
	}
}

func TestDetectChartDrift_NilChart(t *testing.T) {
	rel := &release.Release{Name: "empty", Chart: nil}
	info := helm.ChartInfo{Name: "empty", Version: "0.0.0"}
	results := drift.DetectChartDrift(rel, info)
	if results != nil {
		t.Errorf("expected nil results for nil chart, got %v", results)
	}
}

func TestDetectChartDrift_NestedOverride(t *testing.T) {
	rel := buildReleaseWithDefaults("myapp", map[string]interface{}{
		"image": map[string]interface{}{"tag": "stable", "pullPolicy": "IfNotPresent"},
	}, map[string]interface{}{
		"image": map[string]interface{}{"tag": "latest", "pullPolicy": "IfNotPresent"},
	})

	info := helm.ChartInfo{Name: "myapp", Version: "1.0.0"}
	results := drift.DetectChartDrift(rel, info)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Key != "image.tag" {
		t.Errorf("unexpected key: %s", results[0].Key)
	}
}
