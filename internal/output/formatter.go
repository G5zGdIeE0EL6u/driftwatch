package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/driftwatch/internal/drift"
)

// Format represents the output format for drift results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift results to an io.Writer in a given format.
type Formatter struct {
	w      io.Writer
	format Format
}

// NewFormatter creates a new Formatter.
func NewFormatter(w io.Writer, format Format) *Formatter {
	return &Formatter{w: w, format: format}
}

// Write outputs the drift results according to the configured format.
func (f *Formatter) Write(results []drift.DriftResult) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(results)
	default:
		return f.writeText(results)
	}
}

func (f *Formatter) writeText(results []drift.DriftResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.w, "No drift detected.")
		return err
	}
	for _, r := range results {
		line := fmt.Sprintf("[%s] key=%q live=%v chart=%v",
			strings.ToUpper(string(r.Severity)),
			r.Key,
			r.LiveValue,
			r.ChartValue,
		)
		if _, err := fmt.Fprintln(f.w, line); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) writeJSON(results []drift.DriftResult) error {
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	if results == nil {
		results = []drift.DriftResult{}
	}
	return enc.Encode(results)
}
