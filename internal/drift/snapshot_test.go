package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeSnapshotResults() []DriftResult {
	return []DriftResult{
		{Key: "image.tag", LiveValue: "v1.0", ChartValue: "v1.1", Severity: SeverityHigh},
		{Key: "replicas", LiveValue: "2", ChartValue: "3", Severity: SeverityLow},
	}
}

func TestSaveLoad_Snapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	s := NewSnapshot("myapp", "default", makeSnapshotResults())
	if err := SaveSnapshot(path, s); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if loaded.Release != "myapp" {
		t.Errorf("release: got %q want %q", loaded.Release, "myapp")
	}
	if loaded.Namespace != "default" {
		t.Errorf("namespace: got %q want %q", loaded.Namespace, "default")
	}
	if len(loaded.Results) != 2 {
		t.Errorf("results len: got %d want 2", len(loaded.Results))
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/tmp/does-not-exist-snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDiffSnapshots_NewEntries(t *testing.T) {
	prev := &Snapshot{CapturedAt: time.Now(), Results: makeSnapshotResults()}
	newResult := DriftResult{Key: "resources.limits.cpu", LiveValue: "100m", ChartValue: "200m", Severity: SeverityMedium}
	curr := &Snapshot{CapturedAt: time.Now(), Results: append(makeSnapshotResults(), newResult)}

	added := DiffSnapshots(prev, curr)
	if len(added) != 1 {
		t.Fatalf("expected 1 new drift, got %d", len(added))
	}
	if added[0].Key != "resources.limits.cpu" {
		t.Errorf("unexpected key: %s", added[0].Key)
	}
}

func TestDiffSnapshots_NoDiff(t *testing.T) {
	results := makeSnapshotResults()
	prev := &Snapshot{Results: results}
	curr := &Snapshot{Results: results}

	added := DiffSnapshots(prev, curr)
	if len(added) != 0 {
		t.Errorf("expected no new drifts, got %d", len(added))
	}
}

func TestNewSnapshot_SetsFields(t *testing.T) {
	s := NewSnapshot("svc", "prod", nil)
	if s.Release != "svc" || s.Namespace != "prod" {
		t.Errorf("unexpected fields: %+v", s)
	}
	if s.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestSaveSnapshot_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "snap.json")
	s := NewSnapshot("app", "ns", nil)
	if err := SaveSnapshot(path, s); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
