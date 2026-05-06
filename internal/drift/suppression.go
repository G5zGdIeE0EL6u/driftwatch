package drift

import (
	"encoding/json"
	"os"
	"strings"
)

// SuppressionRule defines a pattern-based rule to ignore specific drift keys.
type SuppressionRule struct {
	Release   string `json:"release"`
	Namespace string `json:"namespace"`
	KeyPrefix string `json:"keyPrefix"`
	Reason    string `json:"reason,omitempty"`
}

// Suppressor holds a set of suppression rules and applies them to drift results.
type Suppressor struct {
	rules []SuppressionRule
}

// NewSuppressor creates a Suppressor from the given rules.
func NewSuppressor(rules []SuppressionRule) *Suppressor {
	return &Suppressor{rules: rules}
}

// LoadSuppressor reads suppression rules from a JSON file.
func LoadSuppressor(path string) (*Suppressor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules []SuppressionRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	return NewSuppressor(rules), nil
}

// Apply filters out DriftResult entries whose key matches a suppression rule
// for the given release and namespace.
func (s *Suppressor) Apply(results []DriftResult, release, namespace string) []DriftResult {
	if len(s.rules) == 0 {
		return results
	}
	out := results[:0:0]
	for _, r := range results {
		if !s.isSuppressed(r.Key, release, namespace) {
			out = append(out, r)
		}
	}
	return out
}

// isSuppressed returns true when any rule matches the given key/release/namespace.
func (s *Suppressor) isSuppressed(key, release, namespace string) bool {
	for _, rule := range s.rules {
		if rule.Release != "" && rule.Release != release {
			continue
		}
		if rule.Namespace != "" && rule.Namespace != namespace {
			continue
		}
		if rule.KeyPrefix != "" && !strings.HasPrefix(key, rule.KeyPrefix) {
			continue
		}
		return true
	}
	return false
}
