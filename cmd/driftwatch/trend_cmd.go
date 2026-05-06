package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
)

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Record and display drift trend for a release",
}

var trendRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Append current drift results to a trend file",
	RunE:  runTrendRecord,
}

var trendShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show recorded drift trend for a release",
	RunE:  runTrendShow,
}

var trendFile string
var trendRelease string
var trendNamespace string

func init() {
	trendRecordCmd.Flags().StringVar(&trendFile, "file", "drift-trend.json", "path to trend file")
	trendRecordCmd.Flags().StringVar(&trendRelease, "release", "", "release name (required)")
	trendRecordCmd.Flags().StringVar(&trendNamespace, "namespace", "default", "Kubernetes namespace")
	_ = trendRecordCmd.MarkFlagRequired("release")

	trendShowCmd.Flags().StringVar(&trendFile, "file", "drift-trend.json", "path to trend file")

	trendCmd.AddCommand(trendRecordCmd, trendShowCmd)
	rootCmd.AddCommand(trendCmd)
}

func runTrendRecord(cmd *cobra.Command, _ []string) error {
	tr := drift.NewTrend(trendRelease, trendNamespace)

	// Load existing trend if file exists.
	if data, err := os.ReadFile(trendFile); err == nil {
		if jsonErr := json.Unmarshal(data, tr); jsonErr != nil {
			return fmt.Errorf("parsing trend file: %w", jsonErr)
		}
	}

	// Record an empty snapshot (real usage would pass actual drift results).
	tr.Record([]drift.DriftResult{})

	out, err := json.MarshalIndent(tr, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling trend: %w", err)
	}
	if err := os.WriteFile(trendFile, out, 0o644); err != nil {
		return fmt.Errorf("writing trend file: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Trend recorded to %s (%d entries)\n", trendFile, len(tr.Entries))
	return nil
}

func runTrendShow(cmd *cobra.Command, _ []string) error {
	data, err := os.ReadFile(trendFile)
	if err != nil {
		return fmt.Errorf("reading trend file: %w", err)
	}
	var tr drift.Trend
	if err := json.Unmarshal(data, &tr); err != nil {
		return fmt.Errorf("parsing trend file: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Trend for %s/%s (%d entries)\n", tr.Namespace, tr.Release, len(tr.Entries))
	for _, e := range tr.Sorted() {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s  drift_count=%d\n", e.Timestamp.Format("2006-01-02T15:04:05Z"), e.DriftCount)
	}
	if tr.Growing() {
		fmt.Fprintln(cmd.OutOrStdout(), "WARNING: drift is growing")
	}
	return nil
}
