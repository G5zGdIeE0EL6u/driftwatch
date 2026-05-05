package drift

import "strings"

// SeverityLevel represents the importance of a detected drift.
type SeverityLevel int

const (
	SeverityLow SeverityLevel = iota
	SeverityMedium
	SeverityHigh
	SeverityUnknown
)

// String returns the string representation of a SeverityLevel.
func (s SeverityLevel) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	default:
		return "unknown"
	}
}

// ParseSeverity converts a string to a SeverityLevel.
// Returns SeverityUnknown if the value is unrecognised.
func ParseSeverity(s string) SeverityLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "low":
		return SeverityLow
	case "medium":
		return SeverityMedium
	case "high":
		return SeverityHigh
	default:
		return SeverityUnknown
	}
}

// AtLeast reports whether s is at or above the given minimum severity.
func (s SeverityLevel) AtLeast(min SeverityLevel) bool {
	return s >= min && s != SeverityUnknown
}

// FilterByMinSeverity returns only those DriftResult entries whose
// severity is at or above minSeverity.
func FilterByMinSeverity(results []DriftResult, minSeverity SeverityLevel) []DriftResult {
	if minSeverity == SeverityLow {
		return results
	}
	out := make([]DriftResult, 0, len(results))
	for _, r := range results {
		if classifySeverity(r.Key).AtLeast(minSeverity) {
			out = append(out, r)
		}
	}
	return out
}
