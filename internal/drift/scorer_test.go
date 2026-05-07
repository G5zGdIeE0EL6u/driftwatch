package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeScorerResults(svs ...Severity) []DriftResult {
	results := make([]DriftResult, len(svs))
	for i, s := range svs {
		results[i] = DriftResult{Key: "k", Severity: s}
	}
	return results
}

func TestScoreRelease_Empty(t *testing.T) {
	ds := ScoreRelease("rel", "ns", nil)
	assert.Equal(t, "rel", ds.Release)
	assert.Equal(t, "ns", ds.Namespace)
	assert.Equal(t, 0, ds.Score)
}

func TestScoreRelease_Weights(t *testing.T) {
	results := makeScorerResults(
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
	)
	ds := ScoreRelease("rel", "default", results)
	assert.Equal(t, 161, ds.Score) // 100+50+10+1
	assert.Equal(t, 1, ds.Critical)
	assert.Equal(t, 1, ds.High)
	assert.Equal(t, 1, ds.Medium)
	assert.Equal(t, 1, ds.Low)
}

func TestScoreRelease_OnlyCritical(t *testing.T) {
	results := makeScorerResults(SeverityCritical, SeverityCritical)
	ds := ScoreRelease("rel", "ns", results)
	assert.Equal(t, 200, ds.Score)
	assert.Equal(t, 2, ds.Critical)
	assert.Equal(t, 0, ds.High)
}

func TestRankScores_Descending(t *testing.T) {
	scores := []DriftScore{
		{Release: "a", Namespace: "ns", Score: 10},
		{Release: "b", Namespace: "ns", Score: 200},
		{Release: "c", Namespace: "ns", Score: 50},
	}
	ranked := RankScores(scores)
	assert.Equal(t, "b", ranked[0].Release)
	assert.Equal(t, "c", ranked[1].Release)
	assert.Equal(t, "a", ranked[2].Release)
}

func TestRankScores_TieBreakByNamespaceThenRelease(t *testing.T) {
	scores := []DriftScore{
		{Release: "z", Namespace: "ns-b", Score: 50},
		{Release: "a", Namespace: "ns-a", Score: 50},
		{Release: "m", Namespace: "ns-a", Score: 50},
	}
	ranked := RankScores(scores)
	assert.Equal(t, "ns-a", ranked[0].Namespace)
	assert.Equal(t, "a", ranked[0].Release)
	assert.Equal(t, "m", ranked[1].Release)
	assert.Equal(t, "z", ranked[2].Release)
}

func TestRankScores_DoesNotMutateInput(t *testing.T) {
	original := []DriftScore{
		{Release: "a", Score: 5},
		{Release: "b", Score: 100},
	}
	_ = RankScores(original)
	assert.Equal(t, "a", original[0].Release)
}
