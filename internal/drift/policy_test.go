package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func buildPolicyResults() []DriftResult {
	return []DriftResult{
		{Release: "app", Namespace: "default", Key: "image.tag", LiveValue: "v1", ChartValue: "v2", Severity: SeverityHigh},
		{Release: "app", Namespace: "default", Key: "replicas", LiveValue: "3", ChartValue: "1", Severity: SeverityLow},
		{Release: "db", Namespace: "prod", Key: "resources.limits.cpu", LiveValue: "500m", ChartValue: "250m", Severity: SeverityMedium},
	}
}

func TestPolicy_NoRules_NoViolations(t *testing.T) {
	p := &Policy{}
	results := p.Evaluate(buildPolicyResults())
	for _, r := range results {
		if r.Violation {
			t.Errorf("expected no violation for key %s", r.Key)
		}
	}
}

func TestPolicy_DenyRule_MatchesKeyPrefix(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{KeyPrefix: "image.", Action: "deny", Reason: "image changes require approval"},
		},
	}
	results := p.Evaluate(buildPolicyResults())
	var violations []PolicyResult
	for _, r := range results {
		if r.Violation {
			violations = append(violations, r)
		}
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "image.tag" {
		t.Errorf("unexpected violation key: %s", violations[0].Key)
	}
	if violations[0].Reason != "image changes require approval" {
		t.Errorf("unexpected reason: %s", violations[0].Reason)
	}
}

func TestPolicy_DenyRule_ScopedToRelease(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{Release: "db", KeyPrefix: "resources.", Action: "deny", Reason: "resource changes locked"},
		},
	}
	results := p.Evaluate(buildPolicyResults())
	var violations []PolicyResult
	for _, r := range results {
		if r.Violation {
			violations = append(violations, r)
		}
	}
	if len(violations) != 1 || violations[0].Release != "db" {
		t.Errorf("expected 1 violation for release db, got %+v", violations)
	}
}

func TestSaveLoad_Policy_RoundTrip(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{KeyPrefix: "image.", Action: "deny", Reason: "approval required"},
			{Release: "app", KeyPrefix: "replicas", Action: "allow"},
		},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	if err := SavePolicy(path, p); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(loaded.Rules))
	}
}

func TestLoadPolicy_MissingFile(t *testing.T) {
	_, err := LoadPolicy("/nonexistent/policy.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}
