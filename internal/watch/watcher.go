package watch

import (
	"context"
	"log"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/helm"
	"github.com/example/driftwatch/internal/output"
)

// Config holds configuration for the watcher.
type Config struct {
	Release   string
	Namespace string
	Interval  time.Duration
	Formatter *output.Formatter
}

// Watcher periodically checks for configuration drift.
type Watcher struct {
	cfg      Config
	client   helm.ClientIface
	detector *drift.Detector
}

// New creates a new Watcher.
func New(cfg Config, client helm.ClientIface, detector *drift.Detector) *Watcher {
	return &Watcher{
		cfg:      cfg,
		client:   client,
		detector: detector,
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	log.Printf("[driftwatch] starting watch for release %q every %s", w.cfg.Release, w.cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("[driftwatch] watch stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := w.check(ctx); err != nil {
				log.Printf("[driftwatch] check error: %v", err)
			}
		}
	}
}

func (w *Watcher) check(ctx context.Context) error {
	rel, err := w.client.GetRelease(w.cfg.Release, w.cfg.Namespace)
	if err != nil {
		return err
	}
	results := w.detector.Detect(rel)
	return w.cfg.Formatter.Write(results)
}
