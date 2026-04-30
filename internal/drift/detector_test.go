package drift

import (
	"testing"

	"helm.sh/helm/v3/pkg/release"
)

func makeRelease(name, ns string, config map[string]interface{}) *release.Release {
	return &release.Release{
		Name:      name,
		Namespace: ns,
		Config:    config,
	}
}

func TestDetect_NoDrift(t *testing.T) {
	d := NewDetector()
	rel := makeRelease("myapp", "default", map[string]interface{}{
		"replicaCount": 2,
		"image":        map[string]interface{}{"tag": "v1.0"},
	})
	desired := map[string]interface{}{
		"replicaCount": 2,
		"image":        map[string]interface{}{"tag": "v1.0"},
	}
	result, err := d.Detect(rel, desired)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift {
		t.Errorf("expected no drift, got changes: %v", result.Changes)
	}
}

func TestDetect_ValueChanged(t *testing.T) {
	d := NewDetector()
	rel := makeRelease("myapp", "default", map[string]interface{}{
		"replicaCount": 1,
	})
	desired := map[string]interface{}{
		"replicaCount": 3,
	}
	result, err := d.Detect(rel, desired)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDrift {
		t.Fatal("expected drift but none detected")
	}
	if len(result.Changes) != 1 || result.Changes[0].Key != "replicaCount" {
		t.Errorf("unexpected changes: %v", result.Changes)
	}
}

func TestDetect_NestedDrift(t *testing.T) {
	d := NewDetector()
	rel := makeRelease("myapp", "default", map[string]interface{}{
		"image": map[string]interface{}{"tag": "v1.0", "pullPolicy": "Always"},
	})
	desired := map[string]interface{}{
		"image": map[string]interface{}{"tag": "v2.0", "pullPolicy": "Always"},
	}
	result, err := d.Detect(rel, desired)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDrift {
		t.Fatal("expected drift but none detected")
	}
	if result.Changes[0].Key != "image.tag" {
		t.Errorf("expected key 'image.tag', got '%s'", result.Changes[0].Key)
	}
}

func TestDetect_NilRelease(t *testing.T) {
	d := NewDetector()
	_, err := d.Detect(nil, map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for nil release")
	}
}

func TestDetect_MissingKey(t *testing.T) {
	d := NewDetector()
	rel := makeRelease("myapp", "default", map[string]interface{}{})
	desired := map[string]interface{}{"replicaCount": 2}
	result, err := d.Detect(rel, desired)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDrift || len(result.Changes) != 1 {
		t.Errorf("expected 1 drift change, got %v", result.Changes)
	}
	if result.Changes[0].OldValue != nil {
		t.Errorf("expected OldValue nil for missing key")
	}
}
