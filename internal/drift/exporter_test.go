package drift

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildExporterResults() []DriftResult {
	return []DriftResult{
		{
			Namespace:  "prod",
			Release:    "api",
			Key:        "replicaCount",
			Severity:   SeverityHigh,
			LiveValue:  3,
			ChartValue: 1,
		},
		{
			Namespace:  "prod",
			Release:    "api",
			Key:        "image.tag",
			Severity:   SeverityLow,
			LiveValue:  "v2.1.0",
			ChartValue: "v2.0.0",
		},
		{
			Namespace:  "staging",
			Release:    "worker",
			Key:        "resources.limits.cpu",
			Severity:   SeverityMedium,
			LiveValue:  "500m",
			ChartValue: "250m",
		},
	}
}

func TestExporter_CSV_Header(t *testing.T) {
	e := NewExporter(ExportCSV)
	e.timestamp = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	var buf bytes.Buffer
	if err := e.Write(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line (header only), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp,namespace,release") {
		t.Errorf("unexpected header: %s", lines[0])
	}
}

func TestExporter_CSV_RowCount(t *testing.T) {
	e := NewExporter(ExportCSV)
	e.timestamp = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	var buf bytes.Buffer
	if err := e.Write(&buf, buildExporterResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// 1 header + 3 data rows
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
}

func TestExporter_CSV_SortedOutput(t *testing.T) {
	e := NewExporter(ExportCSV)
	e.timestamp = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	var buf bytes.Buffer
	_ = e.Write(&buf, buildExporterResults())

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// rows sorted by namespace then release then key
	// prod/api/image.tag < prod/api/replicaCount < staging/worker/...
	if !strings.Contains(lines[1], "image.tag") {
		t.Errorf("expected first data row to be image.tag, got: %s", lines[1])
	}
	if !strings.Contains(lines[3], "staging") {
		t.Errorf("expected last row to be staging, got: %s", lines[3])
	}
}

func TestExporter_TSV_Delimiter(t *testing.T) {
	e := NewExporter(ExportTSV)
	e.timestamp = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	var buf bytes.Buffer
	_ = e.Write(&buf, buildExporterResults())

	if !strings.Contains(buf.String(), "\t") {
		t.Error("expected tab delimiter in TSV output")
	}
	if strings.Contains(strings.SplitN(buf.String(), "\n", 2)[0], ",") {
		t.Error("unexpected comma in TSV header")
	}
}
