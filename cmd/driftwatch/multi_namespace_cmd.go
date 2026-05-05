package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/helm"
	"github.com/yourorg/driftwatch/internal/output"
)

var (
	mnNamespaces []string
	mnOutputFmt  string
)

func init() {
	mnCmd := &cobra.Command{
		Use:   "scan-namespaces",
		Short: "Detect drift across all releases in one or more namespaces",
		RunE:  runScanNamespaces,
	}
	mnCmd.Flags().StringSliceVarP(&mnNamespaces, "namespaces", "n", []string{"default"}, "Namespaces to scan")
	mnCmd.Flags().StringVarP(&mnOutputFmt, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(mnCmd)
}

func runScanNamespaces(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	helmClient, err := helm.NewClient(kubeconfig, "")
	if err != nil {
		return fmt.Errorf("creating helm client: %w", err)
	}

	lister, err := helm.NewNamespaceLister(kubeconfig)
	if err != nil {
		return fmt.Errorf("creating namespace lister: %w", err)
	}

	targets, err := helm.FetchAllNamespaceReleases(ctx, helmClient, lister, mnNamespaces)
	if err != nil {
		return fmt.Errorf("fetching releases: %w", err)
	}

	detector := drift.NewDetector()
	var results []drift.Result

	for _, t := range targets {
		if t.Err != nil {
			fmt.Fprintf(os.Stderr, "warn: skipping %s: %v\n", t.Ref.Key(), t.Err)
			continue
		}
		res := detector.Detect(t.Release)
		results = append(results, res)
	}

	formatter := output.NewFormatter(os.Stdout)
	if err := formatter.Write(mnOutputFmt, results); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	summary := drift.Summarise(results)
	os.Exit(summary.ExitCode())
	return nil
}
