package helm

import (
	"testing"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func newTestClientWithRelease(t *testing.T, rel *release.Release) *Client {
	t.Helper()
	mem := driver.NewMemory()
	store := storage.Init(mem)
	if err := store.Create(rel); err != nil {
		t.Fatalf("store.Create: %v", err)
	}
	return &Client{getter: &fakeRESTClientGetter{store: store}}
}

func TestGetValues_UserSupplied(t *testing.T) {
	rel := &release.Release{
		Name:      "myapp",
		Namespace: "default",
		Version:   1,
		Info:      &release.Info{Status: release.StatusDeployed},
		Chart:     &chart.Chart{Metadata: &chart.Metadata{Name: "myapp"}},
		Config:    map[string]interface{}{"replicaCount": 3},
	}

	c := newTestClientWithRelease(t, rel)
	vals, err := c.GetValues("myapp", "default", UserSupplied)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v, ok := vals["replicaCount"]; !ok || v != 3 {
		t.Errorf("expected replicaCount=3, got %v", vals)
	}
}

func TestGetValues_NotFound(t *testing.T) {
	rel := &release.Release{
		Name:      "other",
		Namespace: "default",
		Version:   1,
		Info:      &release.Info{Status: release.StatusDeployed},
		Chart:     &chart.Chart{Metadata: &chart.Metadata{Name: "other"}},
		Config:    map[string]interface{}{},
	}

	c := newTestClientWithRelease(t, rel)
	_, err := c.GetValues("missing", "default", UserSupplied)
	if err == nil {
		t.Fatal("expected error for missing release, got nil")
	}
}

func TestGetValues_EmptyConfig(t *testing.T) {
	rel := &release.Release{
		Name:      "bare",
		Namespace: "default",
		Version:   1,
		Info:      &release.Info{Status: release.StatusDeployed},
		Chart:     &chart.Chart{Metadata: &chart.Metadata{Name: "bare"}},
		Config:    nil,
	}

	c := newTestClientWithRelease(t, rel)
	vals, err := c.GetValues("bare", "default", UserSupplied)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 0 {
		t.Errorf("expected empty map, got %v", vals)
	}
}
