package helm

import (
	"testing"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"helm.sh/helm/v3/pkg/action"
)

// newTestClient builds a Helm Client backed by an in-memory store.
func newTestClient(t *testing.T, releases ...*release.Release) *Client {
	t.Helper()
	store := storage.Init(driver.NewMemory())
	for _, rel := range releases {
		if err := store.Create(rel); err != nil {
			t.Fatalf("seeding test release: %v", err)
		}
	}
	cfg := &action.Configuration{Releases: store}
	return &Client{cfg: cfg}
}

func TestGetRelease_Found(t *testing.T) {
	rel := &release.Release{
		Name:    "my-app",
		Version: 1,
		Info:    &release.Info{Status: release.StatusDeployed},
	}
	client := newTestClient(t, rel)

	got, err := client.GetRelease("my-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "my-app" {
		t.Errorf("expected release name %q, got %q", "my-app", got.Name)
	}
}

func TestGetRelease_NotFound(t *testing.T) {
	client := newTestClient(t)

	_, err := client.GetRelease("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing release, got nil")
	}
}

func TestGetReleaseValues(t *testing.T) {
	rel := &release.Release{
		Name:    "my-app",
		Version: 1,
		Info:    &release.Info{Status: release.StatusDeployed},
		Config:  map[string]interface{}{"replicaCount": 3},
	}
	client := newTestClient(t, rel)

	vals, err := client.GetReleaseValues("my-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vals["replicaCount"] != 3 {
		t.Errorf("expected replicaCount=3, got %v", vals["replicaCount"])
	}
}
