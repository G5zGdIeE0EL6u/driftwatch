package drift

import "strings"

// DriftResult represents a single detected configuration drift.
type DriftResult struct {
	Release   string            `json:"release"`
	Namespace string            `json:"namespace"`
	Key       string            `json:"key"`
	LiveVal   interface{}       `json:"live_value"`
	ChartVal  interface{}       `json:"chart_value"`
	Severity  Severity          `json:"severity"`
	Labels    map[string]string `json:"labels,omitempty"`
	Note      string            `json:"note,omitempty"`
}

// classifySeverity assigns a severity level based on the key path.
func classifySeverity(key string) Severity {
	switch {
	case hasPrefix(key, "image.", "securityContext.", "serviceAccountName"):
		return SeverityHigh
	case hasPrefix(key, "resources.", "replicas", "env."):
		return SeverityMedium
	default:
		return SeverityLow
	}
}

func hasPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) || s == p {
			return true
		}
	}
	return false
}

// FilterBySeverity returns only results at or above the given severity.
func FilterBySeverity(results []DriftResult, min Severity) []DriftResult {
	var out []DriftResult
	for _, r := range results {
		if r.Severity >= min {
			out = append(out, r)
		}
	}
	return out
}
