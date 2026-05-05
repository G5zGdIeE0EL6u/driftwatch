package drift

import (
	"sort"

	"github.com/yourusername/driftwatch/internal/helm"
)

// AggregateResult holds drift results for multiple releases.
type AggregateResult struct {
	Results map[string]*Result
	Errors  map[string]error
}

// TotalDrifted returns the number of releases with drift detected.
func (a *AggregateResult) TotalDrifted() int {
	count := 0
	for _, r := range a.Results {
		if r.HasDrift() {
			count++
		}
	}
	return count
}

// TotalClean returns the number of releases with no drift.
func (a *AggregateResult) TotalClean() int {
	return len(a.Results) - a.TotalDrifted()
}

// SortedKeys returns release keys in deterministic order.
func (a *AggregateResult) SortedKeys() []string {
	keys := make([]string, 0, len(a.Results))
	for k := range a.Results {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Aggregator runs drift detection across multiple releases.
type Aggregator struct {
	detector *Detector
}

// NewAggregator creates an Aggregator using the provided Detector.
func NewAggregator(d *Detector) *Aggregator {
	return &Aggregator{detector: d}
}

// Run detects drift for each release reference and returns an AggregateResult.
func (a *Aggregator) Run(releases map[string]*helm.Release, overrides map[string]map[string]interface{}) *AggregateResult {
	out := &AggregateResult{
		Results: make(map[string]*Result),
		Errors:  make(map[string]error),
	}
	for key, rel := range releases {
		desired := overrides[key]
		result := a.detector.Detect(rel, desired)
		out.Results[key] = result
	}
	return out
}
