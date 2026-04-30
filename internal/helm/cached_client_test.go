package helm

import (
	"testing"
	"time"

	"helm.sh/helm/v3/pkg/release"
)

func buildCachedClient(t *testing.T, rel *release.Release) *CachedClient {
	t.Helper()
	inner := newTestClientWithRelease(t, rel)
	return NewCachedClient(inner, 30*time.Second)
}

func TestCachedClient_GetRelease_Hit(t *testing.T) {
	rel := &release.Release{Name: "myapp", Namespace: "default"}
	cc := buildCachedClient(t, rel)

	first, err := cc.GetRelease("default", "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.Name != "myapp" {
		t.Fatalf("expected 'myapp', got %s", first.Name)
	}

	// Second call should hit cache (no panic / no error).
	second, err := cc.GetRelease("default", "myapp")
	if err != nil {
		t.Fatalf("unexpected error on cached call: %v", err)
	}
	if first != second {
		t.Fatal("expected same pointer from cache")
	}
}

func TestCachedClient_Invalidate(t *testing.T) {
	rel := &release.Release{Name: "myapp", Namespace: "default"}
	cc := buildCachedClient(t, rel)

	_, _ = cc.GetRelease("default", "myapp")
	cc.Invalidate("default", "myapp")

	// After invalidation the cache should be empty for this key.
	_, ok := cc.cache.Get("release::default::myapp")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestCachedClient_GetValues_Cached(t *testing.T) {
	rel := &release.Release{
		Name:      "myapp",
		Namespace: "default",
		Config:    map[string]interface{}{"replicas": 3},
	}
	cc := buildCachedClient(t, rel)

	v1, err := cc.GetValues("default", "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v2, err := cc.GetValues("default", "myapp")
	if err != nil {
		t.Fatalf("unexpected error on cached call: %v", err)
	}
	if len(v1) != len(v2) {
		t.Fatal("cached values differ in length")
	}
}
