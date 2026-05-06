package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot captures drift results at a point in time for later comparison.
type Snapshot struct {
	CapturedAt time.Time    `json:"captured_at"`
	Release    string       `json:"release"`
	Namespace  string       `json:"namespace"`
	Results    []DriftResult `json:"results"`
}

// NewSnapshot creates a Snapshot from the given drift results.
func NewSnapshot(release, namespace string, results []DriftResult) *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Release:    release,
		Namespace:  namespace,
		Results:    results,
	}
}

// SaveSnapshot writes a snapshot to disk as JSON.
func SaveSnapshot(path string, s *Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create snapshot dir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create snapshot file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadSnapshot reads a snapshot from disk.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot not found: %s", path)
		}
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("decode snapshot: %w", err)
	}
	return &s, nil
}

// DiffSnapshots returns results present in current but not in previous (by key).
func DiffSnapshots(previous, current *Snapshot) []DriftResult {
	prevKeys := make(map[string]struct{}, len(previous.Results))
	for _, r := range previous.Results {
		prevKeys[r.Key] = struct{}{}
	}
	var added []DriftResult
	for _, r := range current.Results {
		if _, seen := prevKeys[r.Key]; !seen {
			added = append(added, r)
		}
	}
	return added
}
