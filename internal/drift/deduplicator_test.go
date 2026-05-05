package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeDR(ns, release, field string, sev Severity) DriftResult {
	return DriftResult{
		Namespace: ns,
		Release:   release,
		FieldPath: field,
		Severity:  sev,
		WantVal:   "a",
		GotVal:    "b",
	}
}

func TestDeduplicator_NoDuplicates(t *testing.T) {
	d := NewDeduplicator()
	d.Add(makeDR("default", "app", "replicas", SeverityHigh))
	d.Add(makeDR("default", "app", "image", SeverityLow))

	results := d.Results()
	require.Len(t, results, 2)
	assert.Equal(t, "image", results[0].FieldPath)
	assert.Equal(t, "replicas", results[1].FieldPath)
}

func TestDeduplicator_KeepsHigherSeverity(t *testing.T) {
	d := NewDeduplicator()
	d.Add(makeDR("default", "app", "replicas", SeverityLow))
	d.Add(makeDR("default", "app", "replicas", SeverityHigh))
	d.Add(makeDR("default", "app", "replicas", SeverityMedium))

	results := d.Results()
	require.Len(t, results, 1)
	assert.Equal(t, SeverityHigh, results[0].Severity)
}

func TestDeduplicator_AddAll(t *testing.T) {
	d := NewDeduplicator()
	batch := []DriftResult{
		makeDR("prod", "svc", "cpu", SeverityMedium),
		makeDR("prod", "svc", "mem", SeverityLow),
		makeDR("staging", "svc", "cpu", SeverityHigh),
	}
	d.AddAll(batch)

	results := d.Results()
	require.Len(t, results, 3)
}

func TestDeduplicator_SortOrder(t *testing.T) {
	d := NewDeduplicator()
	d.Add(makeDR("z-ns", "app", "field", SeverityLow))
	d.Add(makeDR("a-ns", "app", "field", SeverityLow))
	d.Add(makeDR("a-ns", "app", "aaa", SeverityLow))

	results := d.Results()
	require.Len(t, results, 3)
	assert.Equal(t, "a-ns", results[0].Namespace)
	assert.Equal(t, "aaa", results[0].FieldPath)
	assert.Equal(t, "a-ns", results[1].Namespace)
	assert.Equal(t, "field", results[1].FieldPath)
	assert.Equal(t, "z-ns", results[2].Namespace)
}

func TestDeduplicator_Empty(t *testing.T) {
	d := NewDeduplicator()
	results := d.Results()
	assert.Empty(t, results)
}
