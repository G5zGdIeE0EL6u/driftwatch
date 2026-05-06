package drift

import (
	"testing"
	"time"
)

func makeTrendResults(keys ...string) []DriftResult {
	var out []DriftResult
	for _, k := range keys {
		out = append(out, DriftResult{Key: k, Live: "a", Chart: "b"})
	}
	return out
}

func TestTrend_EmptyOnCreation(t *testing.T) {
	tr := NewTrend("myapp", "default")
	if len(tr.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(tr.Entries))
	}
	if tr.Latest() != nil {
		t.Fatal("expected nil Latest on empty trend")
	}
}

func TestTrend_RecordAddsEntry(t *testing.T) {
	tr := NewTrend("myapp", "default")
	tr.Record(makeTrendResults("image.tag", "replicas"))

	if len(tr.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tr.Entries))
	}
	if tr.Entries[0].DriftCount != 2 {
		t.Errorf("expected DriftCount 2, got %d", tr.Entries[0].DriftCount)
	}
}

func TestTrend_Latest(t *testing.T) {
	tr := NewTrend("myapp", "default")
	tr.Record(makeTrendResults("a"))
	tr.Record(makeTrendResults("a", "b", "c"))

	latest := tr.Latest()
	if latest == nil {
		t.Fatal("expected non-nil Latest")
	}
	if latest.DriftCount != 3 {
		t.Errorf("expected DriftCount 3, got %d", latest.DriftCount)
	}
}

func TestTrend_Sorted(t *testing.T) {
	tr := NewTrend("myapp", "default")
	now := time.Now().UTC()
	tr.Entries = []TrendEntry{
		{Timestamp: now.Add(2 * time.Minute), DriftCount: 3},
		{Timestamp: now.Add(1 * time.Minute), DriftCount: 1},
		{Timestamp: now, DriftCount: 2},
	}

	sorted := tr.Sorted()
	if sorted[0].DriftCount != 2 || sorted[1].DriftCount != 1 || sorted[2].DriftCount != 3 {
		t.Errorf("unexpected sort order: %v", sorted)
	}
	// original should be unchanged
	if tr.Entries[0].DriftCount != 3 {
		t.Error("Sorted mutated original entries")
	}
}

func TestTrend_Growing_True(t *testing.T) {
	tr := NewTrend("myapp", "default")
	tr.Record(makeTrendResults("a"))
	tr.Record(makeTrendResults("a", "b"))

	if !tr.Growing() {
		t.Error("expected Growing to be true")
	}
}

func TestTrend_Growing_False(t *testing.T) {
	tr := NewTrend("myapp", "default")
	tr.Record(makeTrendResults("a", "b"))
	tr.Record(makeTrendResults("a"))

	if tr.Growing() {
		t.Error("expected Growing to be false")
	}
}

func TestTrend_Growing_SingleEntry(t *testing.T) {
	tr := NewTrend("myapp", "default")
	tr.Record(makeTrendResults("a"))

	if tr.Growing() {
		t.Error("expected Growing false with single entry")
	}
}
