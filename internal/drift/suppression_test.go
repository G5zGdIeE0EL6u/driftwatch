package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func buildSuppressor(rules []SuppressionRule) *Suppressor {
	return NewSuppressor(rules)
}

func makeDriftResults(keys ...string) []DriftResult {
	out := make([]DriftResult, len(keys))
	for i, k := range keys {
		out[i] = DriftResult{Key: k, Live: "a", Chart: "b", Severity: SeverityMedium}
	}
	return out
}

func TestSuppressor_NoRules_PassesAll(t *testing.T) {
	s := buildSuppressor(nil)
	res := makeDriftResults("image.tag", "replicas")
	got := s.Apply(res, "myapp", "default")
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestSuppressor_MatchesKeyPrefix(t *testing.T) {
	s := buildSuppressor([]SuppressionRule{
		{KeyPrefix: "annotations."},
	})
	res := makeDriftResults("annotations.checksum", "image.tag")
	got := s.Apply(res, "myapp", "default")
	if len(got) != 1 || got[0].Key != "image.tag" {
		t.Fatalf("unexpected results: %+v", got)
	}
}

func TestSuppressor_MatchesReleaseAndNamespace(t *testing.T) {
	s := buildSuppressor([]SuppressionRule{
		{Release: "myapp", Namespace: "staging", KeyPrefix: "replicas"},
	})
	res := makeDriftResults("replicas")
	// should suppress in staging
	if got := s.Apply(res, "myapp", "staging"); len(got) != 0 {
		t.Fatalf("expected suppression, got %+v", got)
	}
	// should NOT suppress in prod
	if got := s.Apply(res, "myapp", "prod"); len(got) != 1 {
		t.Fatalf("expected pass-through, got %+v", got)
	}
}

func TestSuppressor_WildcardRelease(t *testing.T) {
	// Empty release means "any release"
	s := buildSuppressor([]SuppressionRule{
		{KeyPrefix: "debug."},
	})
	res := makeDriftResults("debug.enabled", "image.tag")
	got := s.Apply(res, "anything", "default")
	if len(got) != 1 || got[0].Key != "image.tag" {
		t.Fatalf("unexpected results: %+v", got)
	}
}

func TestLoadSuppressor_RoundTrip(t *testing.T) {
	rules := []SuppressionRule{
		{Release: "app", Namespace: "prod", KeyPrefix: "labels.", Reason: "managed externally"},
	}
	data, _ := json.Marshal(rules)
	tmp := filepath.Join(t.TempDir(), "rules.json")
	_ = os.WriteFile(tmp, data, 0o644)

	s, err := LoadSuppressor(tmp)
	if err != nil {
		t.Fatalf("LoadSuppressor: %v", err)
	}
	if len(s.rules) != 1 || s.rules[0].Reason != "managed externally" {
		t.Fatalf("unexpected rules: %+v", s.rules)
	}
}

func TestLoadSuppressor_MissingFile(t *testing.T) {
	_, err := LoadSuppressor("/nonexistent/rules.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
