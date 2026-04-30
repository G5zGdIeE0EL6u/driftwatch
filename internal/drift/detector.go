package drift

import (
	"fmt"
	"reflect"

	"helm.sh/helm/v3/pkg/release"
)

// DriftResult holds the result of a drift detection comparison.
type DriftResult struct {
	ReleaseName string
	Namespace   string
	HasDrift    bool
	Changes     []Change
}

// Change describes a single configuration difference.
type Change struct {
	Key      string
	OldValue interface{}
	NewValue interface{}
}

// Detector compares Helm release values against a desired state.
type Detector struct{}

// NewDetector creates a new Detector instance.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect compares the live release values with the desired values map.
// It returns a DriftResult describing any differences found.
func (d *Detector) Detect(rel *release.Release, desiredValues map[string]interface{}) (*DriftResult, error) {
	if rel == nil {
		return nil, fmt.Errorf("release must not be nil")
	}

	result := &DriftResult{
		ReleaseName: rel.Name,
		Namespace:   rel.Namespace,
	}

	liveValues := rel.Config
	if liveValues == nil {
		liveValues = map[string]interface{}{}
	}

	changes := diffValues("", liveValues, desiredValues)
	result.Changes = changes
	result.HasDrift = len(changes) > 0

	return result, nil
}

// diffValues recursively compares two value maps and returns a list of changes.
func diffValues(prefix string, live, desired map[string]interface{}) []Change {
	var changes []Change

	visited := map[string]bool{}

	for k, desiredVal := range desired {
		visited[k] = true
		fullKey := key(prefix, k)
		liveVal, exists := live[k]
		if !exists {
			changes = append(changes, Change{Key: fullKey, OldValue: nil, NewValue: desiredVal})
			continue
		}
		desiredMap, dIsMap := desiredVal.(map[string]interface{})
		liveMap, lIsMap := liveVal.(map[string]interface{})
		if dIsMap && lIsMap {
			changes = append(changes, diffValues(fullKey, liveMap, desiredMap)...)
		} else if !reflect.DeepEqual(liveVal, desiredVal) {
			changes = append(changes, Change{Key: fullKey, OldValue: liveVal, NewValue: desiredVal})
		}
	}

	for k, liveVal := range live {
		if visited[k] {
			continue
		}
		changes = append(changes, Change{Key: key(prefix, k), OldValue: liveVal, NewValue: nil})
	}

	return changes
}

func key(prefix, k string) string {
	if prefix == "" {
		return k
	}
	return prefix + "." + k
}
