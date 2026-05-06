package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// PolicyRule defines a rule that marks certain drift results as expected or forbidden.
type PolicyRule struct {
	Release   string `json:"release,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	KeyPrefix string `json:"key_prefix"`
	Action    string `json:"action"` // "allow" or "deny"
	Reason    string `json:"reason,omitempty"`
}

// Policy holds a collection of rules.
type Policy struct {
	Rules []PolicyRule `json:"rules"`
}

// PolicyResult is the outcome of evaluating a DriftResult against a Policy.
type PolicyResult struct {
	DriftResult
	Violation bool   `json:"violation"`
	Reason    string `json:"reason,omitempty"`
}

// LoadPolicy reads a policy file from disk.
func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read policy %s: %w", path, err)
	}
	p := &Policy{}
	if err := json.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("parse policy %s: %w", path, err)
	}
	return p, nil
}

// SavePolicy writes a policy to disk.
func SavePolicy(path string, p *Policy) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal policy: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Evaluate applies the policy rules to a slice of DriftResults and returns PolicyResults.
func (p *Policy) Evaluate(results []DriftResult) []PolicyResult {
	out := make([]PolicyResult, 0, len(results))
	for _, dr := range results {
		pr := PolicyResult{DriftResult: dr}
		for _, rule := range p.Rules {
			if !matchesRule(rule, dr) {
				continue
			}
			if rule.Action == "deny" {
				pr.Violation = true
				pr.Reason = rule.Reason
			}
			break
		}
		out = append(out, pr)
	}
	return out
}

func matchesRule(rule PolicyRule, dr DriftResult) bool {
	if rule.Release != "" && rule.Release != dr.Release {
		return false
	}
	if rule.Namespace != "" && rule.Namespace != dr.Namespace {
		return false
	}
	return strings.HasPrefix(dr.Key, rule.KeyPrefix)
}
