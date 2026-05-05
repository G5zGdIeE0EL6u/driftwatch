package helm

import (
	"context"
	"errors"
	"testing"
)

type stubHelmClient struct {
	releases map[string]*Release
	err      error
}

func (s *stubHelmClient) GetRelease(_ context.Context, namespace, name string) (*Release, error) {
	if s.err != nil {
		return nil, s.err
	}
	key := namespace + "/" + name
	if r, ok := s.releases[key]; ok {
		return r, nil
	}
	return nil, errors.New("not found")
}

func (s *stubHelmClient) GetValues(_ context.Context, namespace, name string) (map[string]interface{}, error) {
	return nil, nil
}

func buildMultiNSLister(refs []ReleaseRef) *NamespaceLister {
	return &NamespaceLister{staticRefs: refs}
}

func TestFetchAllNamespaceReleases_Success(t *testing.T) {
	refs := []ReleaseRef{
		{Namespace: "ns1", Name: "app-a"},
		{Namespace: "ns2", Name: "app-b"},
	}
	client := &stubHelmClient{
		releases: map[string]*Release{
			"ns1/app-a": {Name: "app-a", Namespace: "ns1"},
			"ns2/app-b": {Name: "app-b", Namespace: "ns2"},
		},
	}
	lister := buildMultiNSLister(refs)

	results, err := FetchAllNamespaceReleases(context.Background(), client, lister, []string{"ns1", "ns2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected per-result error for %s: %v", r.Ref.Key(), r.Err)
		}
		if r.Release == nil {
			t.Errorf("expected release for %s, got nil", r.Ref.Key())
		}
	}
}

func TestFetchAllNamespaceReleases_PartialError(t *testing.T) {
	refs := []ReleaseRef{
		{Namespace: "ns1", Name: "app-a"},
		{Namespace: "ns1", Name: "missing"},
	}
	client := &stubHelmClient{
		releases: map[string]*Release{
			"ns1/app-a": {Name: "app-a", Namespace: "ns1"},
		},
	}
	lister := buildMultiNSLister(refs)

	results, err := FetchAllNamespaceReleases(context.Background(), client, lister, []string{"ns1"})
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	var errCount int
	for _, r := range results {
		if r.Err != nil {
			errCount++
		}
	}
	if errCount != 1 {
		t.Errorf("expected 1 per-result error, got %d", errCount)
	}
}
