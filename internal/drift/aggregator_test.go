package drift

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/helm"
	"helm.sh/helm/v3/pkg/release"
)

func buildRelease(name string, vals map[string]interface{}) *helm.Release {
	return &helm.Release{
		Release: &release.Release{
			Name:   name,
			Config: vals,
		},
	}
}

func TestAggregator_NoDrift(t *testing.T) {
	detector := NewDetector()
	agg := NewAggregator(detector)

	releases := map[string]*helm.Release{
		"ns/app": buildRelease("app", map[string]interface{}{"replicas": 2}),
	}
	overrides := map[string]map[string]interface{}{
		"ns/app": {"replicas": 2},
	}

	res := agg.Run(releases, overrides)

	if res.TotalDrifted() != 0 {
		t.Errorf("expected 0 drifted, got %d", res.TotalDrifted())
	}
	if res.TotalClean() != 1 {
		t.Errorf("expected 1 clean, got %d", res.TotalClean())
	}
}

func TestAggregator_WithDrift(t *testing.T) {
	detector := NewDetector()
	agg := NewAggregator(detector)

	releases := map[string]*helm.Release{
		"ns/app": buildRelease("app", map[string]interface{}{"replicas": 3}),
		"ns/db":  buildRelease("db", map[string]interface{}{"storage": "10Gi"}),
	}
	overrides := map[string]map[string]interface{}{
		"ns/app": {"replicas": 2},
		"ns/db":  {"storage": "10Gi"},
	}

	res := agg.Run(releases, overrides)

	if res.TotalDrifted() != 1 {
		t.Errorf("expected 1 drifted, got %d", res.TotalDrifted())
	}
	if res.TotalClean() != 1 {
		t.Errorf("expected 1 clean, got %d", res.TotalClean())
	}
}

func TestAggregator_SortedKeys(t *testing.T) {
	detector := NewDetector()
	agg := NewAggregator(detector)

	releases := map[string]*helm.Release{
		"ns/zebra": buildRelease("zebra", nil),
		"ns/alpha": buildRelease("alpha", nil),
		"ns/mango": buildRelease("mango", nil),
	}

	res := agg.Run(releases, nil)
	keys := res.SortedKeys()

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "ns/alpha" || keys[1] != "ns/mango" || keys[2] != "ns/zebra" {
		t.Errorf("unexpected order: %v", keys)
	}
}
