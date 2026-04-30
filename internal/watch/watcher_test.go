package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/output"
	"github.com/example/driftwatch/internal/watch"
	"helm.sh/helm/v3/pkg/release"
)

type mockClient struct {
	rel *release.Release
	err error
	calls int
}

func (m *mockClient) GetRelease(name, ns string) (*release.Release, error) {
	m.calls++
	return m.rel, m.err
}

func (m *mockClient) GetValues(name, ns string) (map[string]interface{}, error) {
	return nil, nil
}

func TestWatcher_RunsTicks(t *testing.T) {
	mc := &mockClient{rel: &release.Release{Name: "myapp"}}
	det := drift.NewDetector(map[string]interface{}{})
	fmt := output.NewFormatter("text", nil)

	cfg := watch.Config{
		Release:   "myapp",
		Namespace: "default",
		Interval:  20 * time.Millisecond,
		Formatter: fmt,
	}

	w := watch.New(cfg, mc, det)
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	if mc.calls < 2 {
		t.Errorf("expected at least 2 checks, got %d", mc.calls)
	}
}

func TestWatcher_StopsOnCancel(t *testing.T) {
	mc := &mockClient{rel: &release.Release{Name: "myapp"}}
	det := drift.NewDetector(map[string]interface{}{})
	fmt := output.NewFormatter("text", nil)

	cfg := watch.Config{
		Release:   "myapp",
		Namespace: "default",
		Interval:  1 * time.Second,
		Formatter: fmt,
	}

	w := watch.New(cfg, mc, det)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := w.Run(ctx)
	if err == nil {
		t.Error("expected context cancellation error")
	}
}
