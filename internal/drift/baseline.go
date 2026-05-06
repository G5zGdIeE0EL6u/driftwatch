package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline holds a snapshot of drift results for a set of releases
// at a particular point in time. It is used to compare current drift
// against a previously recorded state.
type Baseline struct {
	CapturedAt time.Time               `json:"captured_at"`
	Results    map[string][]DriftResult `json:"results"`
}

// NewBaseline creates a Baseline from a map of release key → drift results.
func NewBaseline(results map[string][]DriftResult) *Baseline {
	return &Baseline{
		CapturedAt: time.Now().UTC(),
		Results:    results,
	}
}

// SaveBaseline serialises b to the file at path (creates or truncates).
func SaveBaseline(path string, b *Baseline) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create %q: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("baseline: encode: %w", err)
	}
	return nil
}

// LoadBaseline reads and deserialises a Baseline from path.
func LoadBaseline(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: open %q: %w", path, err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode: %w", err)
	}
	return &b, nil
}

// NewDriftsSince returns only the DriftResults that are present in current
// but were absent (by Key) in the baseline for the given release key.
func (b *Baseline) NewDriftsSince(releaseKey string, current []DriftResult) []DriftResult {
	known := make(map[string]struct{})
	for _, r := range b.Results[releaseKey] {
		known[r.Key] = struct{}{}
	}
	var novel []DriftResult
	for _, r := range current {
		if _, seen := known[r.Key]; !seen {
			novel = append(novel, r)
		}
	}
	return novel
}
