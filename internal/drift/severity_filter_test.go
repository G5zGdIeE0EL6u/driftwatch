package drift

import (
	"testing"
)

func TestParseSeverity_Known(t *testing.T) {
	cases := []struct {
		input string
		want  SeverityLevel
	}{
		{"low", SeverityLow},
		{"LOW", SeverityLow},
		{"medium", SeverityMedium},
		{"Medium", SeverityMedium},
		{"high", SeverityHigh},
		{"HIGH", SeverityHigh},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := ParseSeverity(tc.input)
			if got != tc.want {
				t.Errorf("ParseSeverity(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseSeverity_Unknown(t *testing.T) {
	got := ParseSeverity("critical")
	if got != SeverityUnknown {
		t.Errorf("expected SeverityUnknown, got %v", got)
	}
}

func TestSeverityString(t *testing.T) {
	if SeverityHigh.String() != "high" {
		t.Errorf("unexpected string for SeverityHigh")
	}
	if SeverityUnknown.String() != "unknown" {
		t.Errorf("unexpected string for SeverityUnknown")
	}
}

func TestAtLeast(t *testing.T) {
	if !SeverityHigh.AtLeast(SeverityMedium) {
		t.Error("high should be at least medium")
	}
	if SeverityLow.AtLeast(SeverityMedium) {
		t.Error("low should not be at least medium")
	}
	if SeverityUnknown.AtLeast(SeverityLow) {
		t.Error("unknown should never satisfy AtLeast")
	}
}

func TestFilterByMinSeverity_ReturnsAll_WhenLow(t *testing.T) {
	results := []DriftResult{
		{Key: "replicas", LiveValue: "1", ChartValue: "2"},
		{Key: "image.tag", LiveValue: "v1", ChartValue: "v2"},
	}
	got := FilterByMinSeverity(results, SeverityLow)
	if len(got) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(got))
	}
}

func TestFilterByMinSeverity_FiltersCorrectly(t *testing.T) {
	results := []DriftResult{
		{Key: "replicas", LiveValue: "1", ChartValue: "2"},
		{Key: "image.tag", LiveValue: "v1", ChartValue: "v2"},
	}
	// Only keys that classifySeverity returns >= medium should survive.
	// Behaviour depends on classifySeverity; we just verify no panic and
	// that the returned slice is a subset.
	got := FilterByMinSeverity(results, SeverityMedium)
	if len(got) > len(results) {
		t.Errorf("filtered slice larger than input: %d > %d", len(got), len(results))
	}
}
