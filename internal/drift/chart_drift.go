package drift

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/helm"
	helmrelease "helm.sh/helm/v3/pkg/release"
)

// ChartDriftResult represents a single value that differs between
// the chart's default and the user-supplied override.
type ChartDriftResult struct {
	Key          string
	ChartDefault interface{}
	UserValue    interface{}
	Severity     string
}

// DetectChartDrift compares a release's chart defaults against the
// user-supplied values and returns any overridden keys.
func DetectChartDrift(rel *helmrelease.Release, info helm.ChartInfo) []ChartDriftResult {
	if rel == nil || rel.Chart == nil {
		return nil
	}

	chartDefaults := rel.Chart.Values
	userValues := rel.Config

	var results []ChartDriftResult
	diffChartDefaults("", chartDefaults, userValues, &results)
	return results
}

// diffChartDefaults recursively walks chart defaults and compares
// them against user-supplied values, recording any differences.
func diffChartDefaults(prefix string, defaults, user map[string]interface{}, out *[]ChartDriftResult) {
	for k, defaultVal := range defaults {
		fullKey := k
		if prefix != "" {
			fullKey = fmt.Sprintf("%s.%s", prefix, k)
		}

		userVal, exists := user[k]
		if !exists {
			// Key not overridden — no drift from chart perspective.
			continue
		}

		defaultMap, defaultIsMap := defaultVal.(map[string]interface{})
		userMap, userIsMap := userVal.(map[string]interface{})

		if defaultIsMap && userIsMap {
			diffChartDefaults(fullKey, defaultMap, userMap, out)
			continue
		}

		if fmt.Sprintf("%v", defaultVal) != fmt.Sprintf("%v", userVal) {
			*out = append(*out, ChartDriftResult{
				Key:          fullKey,
				ChartDefault: defaultVal,
				UserValue:    userVal,
				Severity:     classifySeverity(fullKey),
			})
		}
	}
}
