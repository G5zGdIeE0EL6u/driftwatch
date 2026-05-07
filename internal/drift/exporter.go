package drift

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"time"
)

// ExportFormat identifies the output format for drift exports.
type ExportFormat string

const (
	ExportCSV  ExportFormat = "csv"
	ExportTSV  ExportFormat = "tsv"
)

// ExportRecord is a flattened row representation of a DriftResult.
type ExportRecord struct {
	Timestamp string
	Namespace string
	Release   string
	Key       string
	Severity  string
	Live      string
	Chart     string
}

// Exporter writes drift results to a tabular format (CSV or TSV).
type Exporter struct {
	format    ExportFormat
	timestamp time.Time
}

// NewExporter creates an Exporter for the given format.
func NewExporter(format ExportFormat) *Exporter {
	return &Exporter{
		format:    format,
		timestamp: time.Now().UTC(),
	}
}

// Write encodes results to the writer in the configured format.
func (e *Exporter) Write(w io.Writer, results []DriftResult) error {
	delimiter := ','
	if e.format == ExportTSV {
		delimiter = '\t'
	}

	cw := csv.NewWriter(w)
	cw.Comma = rune(delimiter)

	// header
	if err := cw.Write([]string{"timestamp", "namespace", "release", "key", "severity", "live_value", "chart_value"}); err != nil {
		return fmt.Errorf("exporter: write header: %w", err)
	}

	sorted := make([]DriftResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Namespace != sorted[j].Namespace {
			return sorted[i].Namespace < sorted[j].Namespace
		}
		if sorted[i].Release != sorted[j].Release {
			return sorted[i].Release < sorted[j].Release
		}
		return sorted[i].Key < sorted[j].Key
	})

	ts := e.timestamp.Format(time.RFC3339)
	for _, r := range sorted {
		row := []string{
			ts,
			r.Namespace,
			r.Release,
			r.Key,
			r.Severity.String(),
			fmt.Sprintf("%v", r.LiveValue),
			fmt.Sprintf("%v", r.ChartValue),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("exporter: write row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
