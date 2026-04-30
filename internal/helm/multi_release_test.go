package helm

import (
	"errors"
	"testing"
)

// stubClient is a minimal Client implementation for testing.
type stubClient struct {
	releases map[string]*Release
	values   map[string]map[string]interface{}
	errOn    string
}

func (s *stubClient) GetRelease(name string) (*Release, error) {
	if s.errOn == name {
		return nil, errors.New("injected error")
	}
	if r, ok := s.releases[name]; ok {
		return r, nil
	}
	return nil, errors.New("not found")
}

func (s *stubClient) GetValues(name string) (map[string]interface{}, error) {
	if v, ok := s.values[name]; ok {
		return v, nil
	}
	return nil, nil
}

func TestFetchMultipleReleases_AllSuccess(t *testing.T) {
	client := &stubClient{
		releases: map[string]*Release{
			"app-a": {Name: "app-a"},
			"app-b": {Name: "app-b"},
		},
		values: map[string]map[string]interface{}{
			"app-a": {"replicas": 2},
			"app-b": {"replicas": 3},
		},
	}

	results := FetchMultipleReleases(client, []string{"app-a", "app-b"})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %q: %v", r.Name, r.Err)
		}
		if r.Release == nil {
			t.Errorf("expected release for %q, got nil", r.Name)
		}
	}
}

func TestFetchMultipleReleases_PartialError(t *testing.T) {
	client := &stubClient{
		releases: map[string]*Release{
			"app-a": {Name: "app-a"},
		},
		values:  map[string]map[string]interface{}{},
		errOn: "app-b",
	}

	results := FetchMultipleReleases(client, []string{"app-a", "app-b"})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	var errCount int
	for _, r := range results {
		if r.Err != nil {
			errCount++
		}
	}
	if errCount != 1 {
		t.Errorf("expected 1 error, got %d", errCount)
	}
}

func TestFetchMultipleReleases_PreservesOrder(t *testing.T) {
	names := []string{"z-release", "a-release", "m-release"}
	client := &stubClient{
		releases: map[string]*Release{
			"z-release": {Name: "z-release"},
			"a-release": {Name: "a-release"},
			"m-release": {Name: "m-release"},
		},
		values: map[string]map[string]interface{}{},
	}

	results := FetchMultipleReleases(client, names)

	for i, r := range results {
		if r.Name != names[i] {
			t.Errorf("index %d: expected %q, got %q", i, names[i], r.Name)
		}
	}
}
