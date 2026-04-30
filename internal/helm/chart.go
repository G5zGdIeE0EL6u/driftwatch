package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/repo"
)

// ChartInfo holds metadata about a Helm chart fetched from a release or repository.
type ChartInfo struct {
	Name        string
	Version     string
	Description string
	DefaultValues map[string]interface{}
}

// GetChartFromRelease extracts chart metadata and default values from a deployed release.
func (c *Client) GetChartFromRelease(namespace, releaseName string) (*ChartInfo, error) {
	rel, err := c.GetRelease(namespace, releaseName)
	if err != nil {
		return nil, fmt.Errorf("get release: %w", err)
	}
	if rel.Chart == nil {
		return nil, fmt.Errorf("release %q has no embedded chart", releaseName)
	}
	return chartInfoFromChart(rel.Chart), nil
}

// LoadChartFromPath loads a chart from a local filesystem path.
func LoadChartFromPath(path string) (*ChartInfo, error) {
	ch, err := loader.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load chart from %q: %w", path, err)
	}
	return chartInfoFromChart(ch), nil
}

// IndexEntryToChartInfo converts a repo index entry to a lightweight ChartInfo.
func IndexEntryToChartInfo(entry *repo.ChartVersion) *ChartInfo {
	return &ChartInfo{
		Name:          entry.Name,
		Version:       entry.Version,
		Description:   entry.Description,
		DefaultValues: nil,
	}
}

func chartInfoFromChart(ch *chart.Chart) *ChartInfo {
	defaults := make(map[string]interface{})
	if ch.Values != nil {
		for k, v := range ch.Values {
			defaults[k] = v
		}
	}
	return &ChartInfo{
		Name:          ch.Metadata.Name,
		Version:       ch.Metadata.Version,
		Description:   ch.Metadata.Description,
		DefaultValues: defaults,
	}
}
