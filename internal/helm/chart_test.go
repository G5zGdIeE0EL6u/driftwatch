package helm

import (
	"testing"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

func buildReleaseWithChart(name string, chartValues map[string]interface{}) *release.Release {
	ch := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:        "mychart",
			Version:     "1.2.3",
			Description: "A test chart",
		},
		Values: chartValues,
	}
	return &release.Release{
		Name:      name,
		Namespace: "default",
		Chart:     ch,
		Config:    map[string]interface{}{},
	}
}

func TestChartInfoFromChart_PopulatesFields(t *testing.T) {
	ch := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:        "nginx",
			Version:     "0.9.0",
			Description: "nginx chart",
		},
		Values: map[string]interface{}{"replicaCount": 1},
	}
	info := chartInfoFromChart(ch)
	if info.Name != "nginx" {
		t.Errorf("expected name nginx, got %s", info.Name)
	}
	if info.Version != "0.9.0" {
		t.Errorf("expected version 0.9.0, got %s", info.Version)
	}
	if info.DefaultValues["replicaCount"] != 1 {
		t.Errorf("expected replicaCount=1, got %v", info.DefaultValues["replicaCount"])
	}
}

func TestChartInfoFromChart_NilValues(t *testing.T) {
	ch := &chart.Chart{
		Metadata: &chart.Metadata{Name: "empty", Version: "0.0.1"},
		Values:   nil,
	}
	info := chartInfoFromChart(ch)
	if len(info.DefaultValues) != 0 {
		t.Errorf("expected empty default values, got %v", info.DefaultValues)
	}
}

func TestGetChartFromRelease_NoChart(t *testing.T) {
	client := newTestClient(t)
	// Store a release without a chart
	rel := &release.Release{
		Name:      "nochart",
		Namespace: "default",
		Chart:     nil,
		Config:    map[string]interface{}{},
	}
	client.cfg.Releases.Create(rel) //nolint:errcheck
	_, err := client.GetChartFromRelease("default", "nochart")
	if err == nil {
		t.Fatal("expected error for release with no chart")
	}
}
