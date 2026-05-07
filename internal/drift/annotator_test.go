package drift

import (
	"testing"
)

func buildAnnotatorResults() []DriftResult {
	return []DriftResult{
		{Key: "image.tag", LiveVal: "v1.0", ChartVal: "v1.1", Severity: SeverityHigh},
		{Key: "image.repository", LiveVal: "myrepo", ChartVal: "myrepo", Severity: SeverityLow},
		{Key: "resources.limits.cpu", LiveVal: "500m", ChartVal: "1000m", Severity: SeverityCritical},
		{Key: "replicas", LiveVal: "2", ChartVal: "3", Severity: SeverityMedium},
	}
}

func TestAnnotator_NoRules_NoAnnotations(t *testing.T) {
	a := NewAnnotator()
	results := buildAnnotatorResults()
	anns := a.Annotate(results)
	if len(anns) != 0 {
		t.Fatalf("expected 0 annotations, got %d", len(anns))
	}
}

func TestAnnotator_PrefixMatch(t *testing.T) {
	a := NewAnnotator()
	a.AddRule("image.", "image field change — verify rollout")

	anns := a.Annotate(buildAnnotatorResults())
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
	for _, key := range []string{"image.tag", "image.repository"} {
		if _, ok := anns[key]; !ok {
			t.Errorf("expected annotation for key %q", key)
		}
	}
}

func TestAnnotator_FirstRuleWins(t *testing.T) {
	a := NewAnnotator()
	a.AddRule("resources.", "resource limit change")
	a.AddRule("resources.limits.", "should not appear")

	anns := a.Annotate(buildAnnotatorResults())
	ann, ok := anns["resources.limits.cpu"]
	if !ok {
		t.Fatal("expected annotation for resources.limits.cpu")
	}
	if ann.Message == "" {
		t.Error("annotation message should not be empty")
	}
	// first rule prefix should appear in message
	if !containsStr(ann.Message, "resources.") {
		t.Errorf("expected message to reference matched prefix, got %q", ann.Message)
	}
}

func TestMerge_AttachesAnnotations(t *testing.T) {
	a := NewAnnotator()
	a.AddRule("image.", "check image policy")
	results := buildAnnotatorResults()
	anns := a.Annotate(results)
	merged := Merge(results, anns)

	if len(merged) != len(results) {
		t.Fatalf("expected %d merged results, got %d", len(results), len(merged))
	}
	for _, m := range merged {
		if m.Key == "image.tag" || m.Key == "image.repository" {
			if !m.HasAnnotation {
				t.Errorf("expected HasAnnotation=true for key %q", m.Key)
			}
		} else {
			if m.HasAnnotation {
				t.Errorf("expected HasAnnotation=false for key %q", m.Key)
			}
		}
	}
}

func TestMerge_EmptyAnnotations(t *testing.T) {
	results := buildAnnotatorResults()
	merged := Merge(results, map[string]Annotation{})
	for _, m := range merged {
		if m.HasAnnotation {
			t.Errorf("expected no annotation for key %q", m.Key)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
