package drift

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ChangeEntry records a single observed drift event with metadata.
type ChangeEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Release    string    `json:"release"`
	Namespace  string    `json:"namespace"`
	Key        string    `json:"key"`
	OldValue   string    `json:"old_value"`
	NewValue   string    `json:"new_value"`
	Severity   string    `json:"severity"`
}

// ChangeLog accumulates ChangeEntry records produced during a scan.
type ChangeLog struct {
	entries []ChangeEntry
}

// NewChangeLog returns an empty ChangeLog.
func NewChangeLog() *ChangeLog {
	return &ChangeLog{}
}

// Record appends a new entry derived from a DriftResult.
func (cl *ChangeLog) Record(namespace, release string, results []DriftResult) {
	now := time.Now().UTC()
	for _, r := range results {
		cl.entries = append(cl.entries, ChangeEntry{
			Timestamp: now,
			Release:   release,
			Namespace: namespace,
			Key:       r.Key,
			OldValue:  fmt.Sprintf("%v", r.Expected),
			NewValue:  fmt.Sprintf("%v", r.Actual),
			Severity:  string(r.Severity),
		})
	}
}

// Entries returns a copy of all recorded entries, sorted by timestamp then key.
func (cl *ChangeLog) Entries() []ChangeEntry {
	out := make([]ChangeEntry, len(cl.entries))
	copy(out, cl.entries)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Timestamp.Equal(out[j].Timestamp) {
			return out[i].Key < out[j].Key
		}
		return out[i].Timestamp.Before(out[j].Timestamp)
	})
	return out
}

// Len returns the number of recorded entries.
func (cl *ChangeLog) Len() int { return len(cl.entries) }

// Summary returns a human-readable summary of the change log.
func (cl *ChangeLog) Summary() string {
	if len(cl.entries) == 0 {
		return "no drift changes recorded"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d drift change(s) recorded:\n", len(cl.entries))
	for _, e := range cl.Entries() {
		fmt.Fprintf(&sb, "  [%s] %s/%s %s: %q -> %q (%s)\n",
			e.Timestamp.Format(time.RFC3339),
			e.Namespace, e.Release,
			e.Key, e.OldValue, e.NewValue, e.Severity)
	}
	return sb.String()
}
