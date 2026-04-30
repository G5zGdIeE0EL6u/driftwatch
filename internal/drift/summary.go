package drift

import (
	"fmt"
	"strings"
)

// DriftResult holds the outcome of a drift detection run.
type DriftResult struct {
	Release   string
	Namespace string
	Drifted   bool
	Changes   []Change
}

// Change represents a single detected drift between live and chart values.
type Change struct {
	Key      string
	LiveVal  interface{}
	ChartVal interface{}
}

// Summary produces a human-readable summary of the drift result.
func (r *DriftResult) Summary() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Release : %s\n", r.Release)
	fmt.Fprintf(&sb, "Namespace: %s\n", r.Namespace)

	if !r.Drifted {
		sb.WriteString("Status  : ✓ No drift detected\n")
		return sb.String()
	}

	fmt.Fprintf(&sb, "Status  : ✗ Drift detected (%d change(s))\n", len(r.Changes))
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	for _, c := range r.Changes {
		fmt.Fprintf(&sb, "  key   : %s\n", c.Key)
		fmt.Fprintf(&sb, "  live  : %v\n", c.LiveVal)
		fmt.Fprintf(&sb, "  chart : %v\n", c.ChartVal)
		sb.WriteString(strings.Repeat("-", 40) + "\n")
	}

	return sb.String()
}

// ExitCode returns 1 when drift is detected, 0 otherwise.
// Useful for CI pipelines that rely on process exit codes.
func (r *DriftResult) ExitCode() int {
	if r.Drifted {
		return 1
	}
	return 0
}
