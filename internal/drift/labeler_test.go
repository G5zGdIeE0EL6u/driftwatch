package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildLabelerResults() []DriftResult {
	return []DriftResult{
		{Key: "image.tag", Severity: SeverityHigh},
		{Key: "resources.limits.cpu", Severity: SeverityMedium},
		{Key: "replicas", Severity: SeverityLow},
	}
}

func TestLabeler_NoRules_NoLabels(t *testing.T) {
	l := NewLabeler(nil)
	results := buildLabelerResults()
	out := l.Label(results)
	require.Len(t, out, len(results))
	for _, r := range out {
		assert.Nil(t, r.Labels)
	}
}

func TestLabeler_PrefixMatch_AttachesLabels(t *testing.T) {
	rules := []LabelRule{
		{
			KeyPrefix: "image.",
			Labels:    map[string]string{"team": "platform", "tier": "runtime"},
		},
	}
	l := NewLabeler(rules)
	out := l.Label(buildLabelerResults())

	require.Len(t, out, 3)
	assert.Equal(t, map[string]string{"team": "platform", "tier": "runtime"}, out[0].Labels)
	assert.Nil(t, out[1].Labels)
	assert.Nil(t, out[2].Labels)
}

func TestLabeler_FirstRuleWins(t *testing.T) {
	rules := []LabelRule{
		{KeyPrefix: "image.", Labels: map[string]string{"owner": "first"}},
		{KeyPrefix: "image.", Labels: map[string]string{"owner": "second"}},
	}
	l := NewLabeler(rules)
	out := l.Label(buildLabelerResults())
	assert.Equal(t, "first", out[0].Labels["owner"])
}

func TestLabeler_PreservesExistingLabels(t *testing.T) {
	rules := []LabelRule{
		{KeyPrefix: "replicas", Labels: map[string]string{"new": "label"}},
	}
	l := NewLabeler(rules)
	input := []DriftResult{
		{Key: "replicas", Labels: map[string]string{"existing": "yes"}},
	}
	out := l.Label(input)
	assert.Equal(t, "yes", out[0].Labels["existing"])
	assert.Equal(t, "label", out[0].Labels["new"])
}

func TestMergeLabels_BothNil(t *testing.T) {
	result := mergeLabels(nil, nil)
	assert.Nil(t, result)
}

func TestMergeLabels_OverrideWins(t *testing.T) {
	base := map[string]string{"k": "base", "only": "base"}
	over := map[string]string{"k": "override"}
	out := mergeLabels(base, over)
	assert.Equal(t, "override", out["k"])
	assert.Equal(t, "base", out["only"])
}
