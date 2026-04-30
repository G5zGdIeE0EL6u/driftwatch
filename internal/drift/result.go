package drift

// Severity indicates how critical a drift item is.
type Severity string

const (
	SeverityHigh   Severity = "high"
	SeverityMedium Severity = "medium"
	SeverityLow    Severity = "low"
)

// DriftResult represents a single value that has drifted between
// the live release and the chart default.
type DriftResult struct {
	// Key is the dot-separated path of the value, e.g. "image.tag".
	Key string `json:"key"`

	// LiveValue is the value observed in the running release.
	LiveValue interface{} `json:"liveValue"`

	// ChartValue is the value defined in the chart defaults.
	ChartValue interface{} `json:"chartValue"`

	// Severity classifies the importance of the drift.
	Severity Severity `json:"severity"`
}

// classifySeverity returns a Severity based on the key path.
// Keys under "image" or "replicas" are considered high severity;
// keys under "resources" are medium; everything else is low.
func classifySeverity(key string) Severity {
	switch {
	case hasPrefix(key, "image"), hasPrefix(key, "replicas"), hasPrefix(key, "replicaCount"):
		return SeverityHigh
	case hasPrefix(key, "resources"), hasPrefix(key, "limits"), hasPrefix(key, "requests"):
		return SeverityMedium
	default:
		return SeverityLow
	}
}

func hasPrefix(key, prefix string) bool {
	return key == prefix || len(key) > len(prefix) && key[:len(prefix)+1] == prefix+"."
}
