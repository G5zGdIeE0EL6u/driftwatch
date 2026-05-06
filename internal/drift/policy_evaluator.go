package drift

import "sort"

// EvaluationSummary provides aggregate statistics from a policy evaluation.
type EvaluationSummary struct {
	Total      int `json:"total"`
	Violations int `json:"violations"`
	Clean      int `json:"clean"`
}

// Summarize computes aggregate statistics over a set of PolicyResults.
func Summarize(results []PolicyResult) EvaluationSummary {
	s := EvaluationSummary{Total: len(results)}
	for _, r := range results {
		if r.Violation {
			s.Violations++
		} else {
			s.Clean++
		}
	}
	return s
}

// ViolationsOnly returns only the results that are policy violations.
func ViolationsOnly(results []PolicyResult) []PolicyResult {
	out := make([]PolicyResult, 0)
	for _, r := range results {
		if r.Violation {
			out = append(out, r)
		}
	}
	return out
}

// SortPolicyResults sorts PolicyResults with violations first, then by key.
func SortPolicyResults(results []PolicyResult) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].Violation != results[j].Violation {
			return results[i].Violation
		}
		return results[i].Key < results[j].Key
	})
}

// PolicyEvaluator combines a Policy with helper methods for batch evaluation.
type PolicyEvaluator struct {
	policy *Policy
}

// NewPolicyEvaluator creates a PolicyEvaluator wrapping the given Policy.
func NewPolicyEvaluator(p *Policy) *PolicyEvaluator {
	return &PolicyEvaluator{policy: p}
}

// Run evaluates results, sorts them, and returns the summary alongside results.
func (pe *PolicyEvaluator) Run(results []DriftResult) ([]PolicyResult, EvaluationSummary) {
	pr := pe.policy.Evaluate(results)
	SortPolicyResults(pr)
	return pr, Summarize(pr)
}
