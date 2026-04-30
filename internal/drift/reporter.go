package drift

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Reporter formats and writes DriftResult output.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter that writes to the given writer.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Report writes a human-readable summary of the drift result.
func (r *Reporter) Report(result *DriftResult) {
	if result == nil {
		fmt.Fprintln(r.w, "no result to report")
		return
	}

	fmt.Fprintf(r.w, "Release : %s\n", result.ReleaseName)
	fmt.Fprintf(r.w, "Namespace: %s\n", result.Namespace)

	if !result.HasDrift {
		fmt.Fprintln(r.w, "Status  : ✓ No drift detected")
		return
	}

	fmt.Fprintf(r.w, "Status  : ✗ Drift detected (%d change(s))\n", len(result.Changes))
	fmt.Fprintln(r.w)

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tLIVE VALUE\tDESIRED VALUE")
	fmt.Fprintln(tw, "---\t----------\t-------------")
	for _, c := range result.Changes {
		fmt.Fprintf(tw, "%s\t%v\t%v\n", c.Key, formatVal(c.OldValue), formatVal(c.NewValue))
	}
	tw.Flush()
}

func formatVal(v interface{}) string {
	if v == nil {
		return "<unset>"
	}
	return fmt.Sprintf("%v", v)
}
