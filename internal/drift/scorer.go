package drift

import "sort"

// DriftScore represents a weighted drift score for a release.
type DriftScore struct {
	Release   string
	Namespace string
	Score     int
	Critical  int
	High      int
	Medium    int
	Low       int
}

// severityWeight maps severity levels to numeric weights.
var severityWeight = map[Severity]int{
	SeverityCritical: 100,
	SeverityHigh:     50,
	SeverityMedium:   10,
	SeverityLow:      1,
}

// ScoreRelease computes a weighted DriftScore from a slice of DriftResults
// belonging to a single release/namespace pair.
func ScoreRelease(release, namespace string, results []DriftResult) DriftScore {
	ds := DriftScore{
		Release:   release,
		Namespace: namespace,
	}
	for _, r := range results {
		switch r.Severity {
		case SeverityCritical:
			ds.Critical++
		case SeverityHigh:
			ds.High++
		case SeverityMedium:
			ds.Medium++
		case SeverityLow:
			ds.Low++
		}
		ds.Score += severityWeight[r.Severity]
	}
	return ds
}

// RankScores returns a slice of DriftScores sorted descending by Score.
func RankScores(scores []DriftScore) []DriftScore {
	out := make([]DriftScore, len(scores))
	copy(out, scores)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if out[i].Namespace != out[j].Namespace {
			return out[i].Namespace < out[j].Namespace
		}
		return out[i].Release < out[j].Release
	})
	return out
}
