package drift

import (
	"fmt"

	"github.com/your-org/driftwatch/internal/helm"
)

// ChartDriftResult captures differences between a release's deployed chart
// version and an expected chart version.
type ChartDriftResult struct {
	ReleaseName     string
	DeployedVersion string
	ExpectedVersion string
	VersionDrifted  bool
	DefaultsChanged []DiffEntry
}

// DiffEntry represents a single key that differs between two value maps.
type DiffEntry struct {
	Key      string
	Deployed interface{}
	Expected interface{}
}

// DetectChartDrift compares the deployed chart info against an expected ChartInfo
// and returns a ChartDriftResult describing any version or default-value drift.
func DetectChartDrift(releaseName string, deployed, expected *helm.ChartInfo) (*ChartDriftResult, error) {
	if deployed == nil {
		return nil, fmt.Errorf("deployed chart info is nil for release %q", releaseName)
	}
	if expected == nil {
		return nil, fmt.Errorf("expected chart info is nil for release %q", releaseName)
	}

	result := &ChartDriftResult{
		ReleaseName:     releaseName,
		DeployedVersion: deployed.Version,
		ExpectedVersion: expected.Version,
		VersionDrifted:  deployed.Version != expected.Version,
	}

	result.DefaultsChanged = diffChartDefaults(deployed.DefaultValues, expected.DefaultValues)
	return result, nil
}

// HasDrift returns true if either the version or any default values have drifted.
func (r *ChartDriftResult) HasDrift() bool {
	return r.VersionDrifted || len(r.DefaultsChanged) > 0
}

func diffChartDefaults(deployed, expected map[string]interface{}) []DiffEntry {
	var diffs []DiffEntry
	for k, expVal := range expected {
		depVal, exists := deployed[k]
		if !exists || fmt.Sprintf("%v", depVal) != fmt.Sprintf("%v", expVal) {
			diffs = append(diffs, DiffEntry{Key: k, Deployed: depVal, Expected: expVal})
		}
	}
	return diffs
}
