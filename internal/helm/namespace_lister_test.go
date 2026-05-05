package helm

import (
	"context"
	"testing"

	helm "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func buildNamespaceLister(releases []*release.Release) *NamespaceLister {
	return &NamespaceLister{
		cfgFn: func(namespace string) (*helm.Configuration, error) {
			cfg := &helm.Configuration{}
			_ = cfg.Init(&fakeGetter{}, namespace, "memory", func(string, ...interface{}) {})
			for _, r := range releases {
				_ = cfg.Releases.Create(r)
			}
			return cfg, nil
		},
	}
}

func makeRelRef(name, ns string) *release.Release {
	return &release.Release{
		Name:      name,
		Namespace: ns,
		Info:      &release.Info{Status: release.StatusDeployed},
		Chart:     buildReleaseWithChart(name).Chart,
	}
}

func TestListReleaseNames_ReturnsAll(t *testing.T) {
	lister := buildNamespaceLister([]*release.Release{
		makeRelRef("alpha", "default"),
		makeRelRef("beta", "default"),
	})

	refs, err := lister.ListReleaseNames(context.Background(), []string{"default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d", len(refs))
	}
}

func TestListReleaseNames_EmptyNamespacesUsesDefault(t *testing.T) {
	lister := buildNamespaceLister([]*release.Release{
		makeRelRef("gamma", ""),
	})

	refs, err := lister.ListReleaseNames(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("expected 1 ref, got %d", len(refs))
	}
}

func TestReleaseRef_KeyUnique(t *testing.T) {
	a := ReleaseRef{Name: "svc", Namespace: "ns1"}
	b := ReleaseRef{Name: "svc", Namespace: "ns2"}
	if a.Key() == b.Key() {
		t.Errorf("expected different keys for different namespaces, got %q", a.Key())
	}
}
