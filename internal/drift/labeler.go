package drift

import (
	"strings"
)

// LabelRule maps a key prefix to a set of labels to attach.
type LabelRule struct {
	KeyPrefix string            `json:"key_prefix" yaml:"key_prefix"`
	Labels    map[string]string `json:"labels"     yaml:"labels"`
}

// Labeler attaches metadata labels to DriftResults based on configurable rules.
type Labeler struct {
	rules []LabelRule
}

// NewLabeler constructs a Labeler with the provided rules.
func NewLabeler(rules []LabelRule) *Labeler {
	return &Labeler{rules: rules}
}

// Label returns a copy of results with labels applied according to matching rules.
// The first matching rule for each result wins.
func (l *Labeler) Label(results []DriftResult) []DriftResult {
	out := make([]DriftResult, len(results))
	for i, r := range results {
		out[i] = r
		for _, rule := range l.rules {
			if strings.HasPrefix(r.Key, rule.KeyPrefix) {
				out[i].Labels = mergeLabels(r.Labels, rule.Labels)
				break
			}
		}
	}
	return out
}

// MergeLabels combines two label maps, with overrides taking precedence.
func mergeLabels(base, overrides map[string]string) map[string]string {
	if len(base) == 0 && len(overrides) == 0 {
		return nil
	}
	out := make(map[string]string, len(base)+len(overrides))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
	}
	return out
}
