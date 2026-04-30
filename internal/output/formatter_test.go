package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/driftwatch/internal/drift"
	"github.com/yourusername/driftwatch/internal/output"
)

func makeDriftResults() []drift.DriftResult {
	return []drift.DriftResult{
		{Key: "replicaCount", LiveValue: 3, ChartValue: 1, Severity: drift.SeverityHigh},
		{Key: "image.tag", LiveValue: "v1.1.0", ChartValue: "v1.0.0", Severity: drift.SeverityMedium},
	}
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText)
	if err := f.Write(makeDriftResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "replicaCount") {
		t.Errorf("expected key replicaCount in output, got: %q", out)
	}
	if !strings.Contains(out, "HIGH") {
		t.Errorf("expected severity HIGH in output, got: %q", out)
	}
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var results []drift.DriftResult
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty slice, got %d items", len(results))
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON)
	if err := f.Write(makeDriftResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var results []drift.DriftResult
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if results[0].Key != "replicaCount" {
		t.Errorf("unexpected first key: %q", results[0].Key)
	}
}
