package drift

import (
	"strings"
	"testing"
)

func makeResults(keys ...string) []DriftResult {
	var out []DriftResult
	for _, k := range keys {
		out = append(out, DriftResult{
			Key:      k,
			Expected: "default",
			Actual:   "override",
			Severity: SeverityMedium,
		})
	}
	return out
}

func TestChangeLog_EmptyOnCreation(t *testing.T) {
	cl := NewChangeLog()
	if cl.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", cl.Len())
	}
}

func TestChangeLog_RecordAddsEntries(t *testing.T) {
	cl := NewChangeLog()
	cl.Record("default", "my-release", makeResults("image.tag", "replicas"))
	if cl.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", cl.Len())
	}
}

func TestChangeLog_EntriesAreSortedByKey(t *testing.T) {
	cl := NewChangeLog()
	cl.Record("default", "my-release", makeResults("z.key", "a.key", "m.key"))
	entries := cl.Entries()
	if entries[0].Key != "a.key" {
		t.Errorf("expected first key to be a.key, got %s", entries[0].Key)
	}
	if entries[2].Key != "z.key" {
		t.Errorf("expected last key to be z.key, got %s", entries[2].Key)
	}
}

func TestChangeLog_EntriesPreservesNamespaceAndRelease(t *testing.T) {
	cl := NewChangeLog()
	cl.Record("staging", "nginx", makeResults("image.tag"))
	e := cl.Entries()[0]
	if e.Namespace != "staging" {
		t.Errorf("expected namespace staging, got %s", e.Namespace)
	}
	if e.Release != "nginx" {
		t.Errorf("expected release nginx, got %s", e.Release)
	}
}

func TestChangeLog_SummaryNoDrift(t *testing.T) {
	cl := NewChangeLog()
	s := cl.Summary()
	if s != "no drift changes recorded" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestChangeLog_SummaryWithDrift(t *testing.T) {
	cl := NewChangeLog()
	cl.Record("prod", "api", makeResults("image.tag"))
	s := cl.Summary()
	if !strings.Contains(s, "1 drift change(s)") {
		t.Errorf("summary missing count: %s", s)
	}
	if !strings.Contains(s, "prod/api") {
		t.Errorf("summary missing release ref: %s", s)
	}
}

func TestChangeLog_MultipleRecordCalls(t *testing.T) {
	cl := NewChangeLog()
	cl.Record("ns1", "rel1", makeResults("k1"))
	cl.Record("ns2", "rel2", makeResults("k2", "k3"))
	if cl.Len() != 3 {
		t.Fatalf("expected 3 total entries, got %d", cl.Len())
	}
}
