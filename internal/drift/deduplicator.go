package drift

import "sort"

// DriftKey uniquely identifies a drift result by release and field path.
type DriftKey struct {
	Release   string
	Namespace string
	FieldPath string
}

// Deduplicator removes duplicate DriftResult entries across multiple scans
// or aggregated results, keeping the entry with the highest severity.
type Deduplicator struct {
	seen map[DriftKey]DriftResult
}

// NewDeduplicator creates a new Deduplicator.
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{
		seen: make(map[DriftKey]DriftResult),
	}
}

// Add inserts a DriftResult into the deduplicator. If an entry with the same
// key already exists, the one with the higher severity is retained.
func (d *Deduplicator) Add(r DriftResult) {
	k := DriftKey{
		Release:   r.Release,
		Namespace: r.Namespace,
		FieldPath: r.FieldPath,
	}
	existing, ok := d.seen[k]
	if !ok || severityRank(r.Severity) > severityRank(existing.Severity) {
		d.seen[k] = r
	}
}

// AddAll adds multiple DriftResult entries.
func (d *Deduplicator) AddAll(results []DriftResult) {
	for _, r := range results {
		d.Add(r)
	}
}

// Results returns deduplicated results sorted by namespace, release, then field path.
func (d *Deduplicator) Results() []DriftResult {
	out := make([]DriftResult, 0, len(d.seen))
	for _, v := range d.seen {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Namespace != out[j].Namespace {
			return out[i].Namespace < out[j].Namespace
		}
		if out[i].Release != out[j].Release {
			return out[i].Release < out[j].Release
		}
		return out[i].FieldPath < out[j].FieldPath
	})
	return out
}

// severityRank maps a Severity to a numeric rank for comparison.
func severityRank(s Severity) int {
	switch s {
	case SeverityLow:
		return 1
	case SeverityMedium:
		return 2
	case SeverityHigh:
		return 3
	default:
		return 0
	}
}
