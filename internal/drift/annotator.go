package drift

import (
	"fmt"
	"strings"
)

// Annotation holds a human-readable explanation attached to a DriftResult.
type Annotation struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

// Annotator enriches DriftResults with explanatory annotations based on
// configurable rules that match key prefixes or exact keys.
type Annotator struct {
	rules []annotationRule
}

type annotationRule struct {
	prefix  string
	message string
}

// NewAnnotator constructs an Annotator with no rules.
func NewAnnotator() *Annotator {
	return &Annotator{}
}

// AddRule registers a rule: any DriftResult whose key starts with prefix will
// receive the given message as an annotation.
func (a *Annotator) AddRule(prefix, message string) {
	a.rules = append(a.rules, annotationRule{prefix: prefix, message: message})
}

// Annotate returns a map from DriftResult key to Annotation for all results
// that match at least one rule. Results with no matching rule are omitted.
func (a *Annotator) Annotate(results []DriftResult) map[string]Annotation {
	out := make(map[string]Annotation)
	for _, r := range results {
		for _, rule := range a.rules {
			if strings.HasPrefix(r.Key, rule.prefix) {
				out[r.Key] = Annotation{
					Key:     r.Key,
					Message: fmt.Sprintf("%s (matched prefix %q)", rule.message, rule.prefix),
				}
				break // first matching rule wins
			}
		}
	}
	return out
}

// Merge attaches annotations from the provided map directly onto a copy of
// each DriftResult, returning a new slice of AnnotatedDriftResult.
func Merge(results []DriftResult, annotations map[string]Annotation) []AnnotatedDriftResult {
	out := make([]AnnotatedDriftResult, 0, len(results))
	for _, r := range results {
		ann, ok := annotations[r.Key]
		out = append(out, AnnotatedDriftResult{DriftResult: r, Annotation: ann, HasAnnotation: ok})
	}
	return out
}

// AnnotatedDriftResult wraps a DriftResult with an optional Annotation.
type AnnotatedDriftResult struct {
	DriftResult
	Annotation    Annotation `json:"annotation,omitempty"`
	HasAnnotation bool       `json:"has_annotation"`
}
