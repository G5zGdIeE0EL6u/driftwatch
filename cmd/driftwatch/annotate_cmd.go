package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Annotate drift results loaded from a snapshot with rule-based messages",
	RunE:  runAnnotate,
}

func init() {
	annotateCmd.Flags().String("snapshot", "", "path to snapshot file (required)")
	annotateCmd.Flags().StringArray("rule", nil, "annotation rule as PREFIX:MESSAGE (repeatable)")
	annotateCmd.Flags().String("output", "text", "output format: text|json")
	_ = annotateCmd.MarkFlagRequired("snapshot")
	rootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	rules, _ := cmd.Flags().GetStringArray("rule")
	format, _ := cmd.Flags().GetString("output")

	snap, err := drift.LoadSnapshot(snapshotPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	annotator := drift.NewAnnotator()
	for _, rule := range rules {
		prefix, message, ok := splitRule(rule)
		if !ok {
			return fmt.Errorf("invalid rule %q: expected PREFIX:MESSAGE", rule)
		}
		annotator.AddRule(prefix, message)
	}

	annotations := annotator.Annotate(snap.Results)
	merged := drift.Merge(snap.Results, annotations)

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(merged)
	default:
		for _, m := range merged {
			line := fmt.Sprintf("[%s] %s: %q -> %q", m.Severity, m.Key, m.ChartVal, m.LiveVal)
			if m.HasAnnotation {
				line += fmt.Sprintf(" // %s", m.Annotation.Message)
			}
			fmt.Fprintln(os.Stdout, line)
		}
	}
	return nil
}

// splitRule splits "PREFIX:MESSAGE" into its two parts.
func splitRule(rule string) (prefix, message string, ok bool) {
	for i := 0; i < len(rule); i++ {
		if rule[i] == ':' {
			return rule[:i], rule[i+1:], true
		}
	}
	return "", "", false
}
