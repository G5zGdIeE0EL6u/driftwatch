package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/helm"
)

var (
	baselineFile      string
	baselineNamespace string
	baselineRelease   string
)

func init() {
	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Manage drift baselines",
	}

	captureCmd := &cobra.Command{
		Use:   "capture",
		Short: "Capture current drift state as a baseline",
		RunE:  runCaptureBaseline,
	}
	captureCmd.Flags().StringVarP(&baselineRelease, "release", "r", "", "Helm release name (required)")
	captureCmd.Flags().StringVarP(&baselineNamespace, "namespace", "n", "default", "Kubernetes namespace")
	captureCmd.Flags().StringVarP(&baselineFile, "file", "f", "baseline.json", "Output baseline file path")
	_ = captureCmd.MarkFlagRequired("release")

	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Show drifts not present in the baseline",
		RunE:  runDiffBaseline,
	}
	diffCmd.Flags().StringVarP(&baselineRelease, "release", "r", "", "Helm release name (required)")
	diffCmd.Flags().StringVarP(&baselineNamespace, "namespace", "n", "default", "Kubernetes namespace")
	diffCmd.Flags().StringVarP(&baselineFile, "file", "f", "baseline.json", "Baseline file to compare against")
	_ = diffCmd.MarkFlagRequired("release")

	baselineCmd.AddCommand(captureCmd, diffCmd)
	rootCmd.AddCommand(baselineCmd)
}

func runCaptureBaseline(cmd *cobra.Command, _ []string) error {
	hc, err := helm.NewClient(baselineNamespace)
	if err != nil {
		return fmt.Errorf("helm client: %w", err)
	}
	rel, err := hc.GetRelease(baselineRelease)
	if err != nil {
		return fmt.Errorf("get release: %w", err)
	}
	det := drift.NewDetector(rel)
	results := det.Detect()
	key := baselineNamespace + "/" + baselineRelease
	b := drift.NewBaseline(map[string][]drift.DriftResult{key: results})
	if err := drift.SaveBaseline(baselineFile, b); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Baseline saved to %s (%d drift(s) recorded)\n", baselineFile, len(results))
	return nil
}

func runDiffBaseline(cmd *cobra.Command, _ []string) error {
	b, err := drift.LoadBaseline(baselineFile)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}
	hc, err := helm.NewClient(baselineNamespace)
	if err != nil {
		return fmt.Errorf("helm client: %w", err)
	}
	rel, err := hc.GetRelease(baselineRelease)
	if err != nil {
		return fmt.Errorf("get release: %w", err)
	}
	det := drift.NewDetector(rel)
	current := det.Detect()
	key := baselineNamespace + "/" + baselineRelease
	novel := b.NewDriftsSince(key, current)
	if len(novel) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No new drift since baseline.")
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%d new drift(s) since baseline captured at %s:\n", len(novel), b.CapturedAt.Format("2006-01-02T15:04:05Z"))
	for _, r := range novel {
		fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s: live=%v chart=%v\n", r.Severity, r.Path, r.LiveVal, r.ChartVal)
	}
	os.Exit(1)
	return nil
}
