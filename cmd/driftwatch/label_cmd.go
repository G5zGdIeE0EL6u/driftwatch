package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Apply metadata labels to drift results based on key-prefix rules",
}

func init() {
	rootCmd.AddCommand(labelCmd)

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply label rules to a drift JSON file and print labelled results",
		RunE:  runLabelApply,
	}
	applyCmd.Flags().StringP("input", "i", "", "Path to drift results JSON (required)")
	applyCmd.Flags().StringArrayP("rule", "r", nil, "Label rule in format 'key_prefix:label=value' (repeatable)")
	_ = applyCmd.MarkFlagRequired("input")
	labelCmd.AddCommand(applyCmd)
}

func runLabelApply(cmd *cobra.Command, _ []string) error {
	inputPath, _ := cmd.Flags().GetString("input")
	ruleStrs, _ := cmd.Flags().GetStringArray("rule")

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var results []drift.DriftResult
	if err := json.Unmarshal(data, &results); err != nil {
		return fmt.Errorf("parsing drift results: %w", err)
	}

	rules, err := parseLabelRules(ruleStrs)
	if err != nil {
		return err
	}

	labeler := drift.NewLabeler(rules)
	labelled := labeler.Label(results)

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(labelled)
}

// parseLabelRules converts CLI strings like "image.:team=platform" into LabelRule values.
func parseLabelRules(raw []string) ([]drift.LabelRule, error) {
	var rules []drift.LabelRule
	for _, s := range raw {
		colon := strings.Index(s, ":")
		if colon < 0 {
			return nil, fmt.Errorf("invalid rule %q: expected 'key_prefix:label=value'", s)
		}
		prefix := s[:colon]
		kv := s[colon+1:]
		eq := strings.Index(kv, "=")
		if eq < 0 {
			return nil, fmt.Errorf("invalid rule %q: label part must be 'label=value'", s)
		}
		rules = append(rules, drift.LabelRule{
			KeyPrefix: prefix,
			Labels:    map[string]string{kv[:eq]: kv[eq+1:]},
		})
	}
	return rules, nil
}
