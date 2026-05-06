package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeResults(keys ...string) []DriftResult {
	var out []DriftResult
	for _, k := range keys {
		out = append(out, DriftResult{Key: k, Path: k, LiveVal: "a", ChartVal: "b", Severity: SeverityMedium})
	}
	return out
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	results := map[string][]DriftResult{
		"ns/app": makeResults("image.tag", "replicas"),
	}
	b := NewBaseline(results)
	if err := SaveBaseline(path, b); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if len(loaded.Results["ns/app"]) != 2 {
		t.Errorf("expected 2 results, got %d", len(loaded.Results["ns/app"]))
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewDriftsSince_AllNew(t *testing.T) {
	b := &Baseline{
		CapturedAt: time.Now(),
		Results:    map[string][]DriftResult{},
	}
	current := makeResults("image.tag", "replicas")
	novel := b.NewDriftsSince("ns/app", current)
	if len(novel) != 2 {
		t.Errorf("expected 2 novel drifts, got %d", len(novel))
	}
}

func TestNewDriftsSince_NoneNew(t *testing.T) {
	b := &Baseline{
		CapturedAt: time.Now(),
		Results: map[string][]DriftResult{
			"ns/app": makeResults("image.tag", "replicas"),
		},
	}
	current := makeResults("image.tag", "replicas")
	novel := b.NewDriftsSince("ns/app", current)
	if len(novel) != 0 {
		t.Errorf("expected 0 novel drifts, got %d", len(novel))
	}
}

func TestNewDriftsSince_PartiallyNew(t *testing.T) {
	b := &Baseline{
		CapturedAt: time.Now(),
		Results: map[string][]DriftResult{
			"ns/app": makeResults("image.tag"),
		},
	}
	current := makeResults("image.tag", "replicas", "resources.limits.cpu")
	novel := b.NewDriftsSince("ns/app", current)
	if len(novel) != 2 {
		t.Errorf("expected 2 novel drifts, got %d", len(novel))
	}
}

func TestSaveBaseline_BadPath(t *testing.T) {
	b := NewBaseline(nil)
	err := SaveBaseline("/no/such/dir/baseline.json", b)
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}

func TestNewBaseline_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	b := NewBaseline(nil)
	after := time.Now().UTC()
	if b.CapturedAt.Before(before) || b.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not in expected range", b.CapturedAt)
	}
	_ = os.Getenv("") // suppress unused import warning
}
