package drift

import (
	"sort"
	"time"
)

// TrendEntry records drift results at a point in time.
type TrendEntry struct {
	Timestamp  time.Time     `json:"timestamp"`
	Results    []DriftResult `json:"results"`
	DriftCount int           `json:"drift_count"`
}

// Trend holds a time-ordered series of drift snapshots for a release.
type Trend struct {
	Release   string       `json:"release"`
	Namespace string       `json:"namespace"`
	Entries   []TrendEntry `json:"entries"`
}

// NewTrend creates an empty Trend for the given release and namespace.
func NewTrend(release, namespace string) *Trend {
	return &Trend{
		Release:   release,
		Namespace: namespace,
	}
}

// Record appends a new entry to the trend using the current time.
func (t *Trend) Record(results []DriftResult) {
	t.Entries = append(t.Entries, TrendEntry{
		Timestamp:  time.Now().UTC(),
		Results:    results,
		DriftCount: len(results),
	})
}

// Latest returns the most recent TrendEntry, or nil if there are none.
func (t *Trend) Latest() *TrendEntry {
	if len(t.Entries) == 0 {
		return nil
	}
	return &t.Entries[len(t.Entries)-1]
}

// Sorted returns entries ordered by timestamp ascending.
func (t *Trend) Sorted() []TrendEntry {
	copy := make([]TrendEntry, len(t.Entries))
	for i, e := range t.Entries {
		copy[i] = e
	}
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Timestamp.Before(copy[j].Timestamp)
	})
	return copy
}

// Growing reports whether drift count has increased between the last two entries.
func (t *Trend) Growing() bool {
	if len(t.Entries) < 2 {
		return false
	}
	last := t.Entries[len(t.Entries)-1]
	prev := t.Entries[len(t.Entries)-2]
	return last.DriftCount > prev.DriftCount
}
