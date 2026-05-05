package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/driftwatch/internal/drift"
	"github.com/yourusername/driftwatch/internal/helm"
	"github.com/yourusername/driftwatch/internal/output"
)

var (
	aggNamespace string
	aggOutput    string
)

func init() {
	aggregateCmd := &cobra.Command{
		Use:   "aggregate",
		Short: "Detect drift across all releases in a namespace",
		RunE:  runAggregate,
	}
	aggregateCmd.Flags().StringVarP(&aggNamespace, "namespace", "n", "default", "Kubernetes namespace to scan")
	aggregateCmd.Flags().StringVarP(&aggOutput, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(aggregateCmd)
}

func runAggregate(cmd *cobra.Command, args []string) error {
	client, err := helm.NewClient(aggNamespace, kubeconfig)
	if err != nil {
		return fmt.Errorf("helm client: %w", err)
	}

	lister := helm.NewNamespaceLister(client)
	refs, err := lister.ListReleaseNames(cmd.Context(), []string{aggNamespace})
	if err != nil {
		return fmt.Errorf("listing releases: %w", err)
	}

	fetched, fetchErrs := helm.FetchMultipleReleases(cmd.Context(), client, refs)
	for ref, e := range fetchErrs {
		fmt.Fprintf(os.Stderr, "warn: skipping %s: %v\n", ref, e)
	}

	releases := make(map[string]*helm.Release, len(fetched))
	overrides := make(map[string]map[string]interface{}, len(fetched))
	for _, ref := range refs {
		k := ref.Key()
		if rel, ok := fetched[k]; ok {
			releases[k] = rel
			vals, _ := client.GetValues(cmd.Context(), ref.Name)
			overrides[k] = vals
		}
	}

	detector := drift.NewDetector()
	agg := drift.NewAggregator(detector)
	aggResult := agg.Run(releases, overrides)

	var driftResults []drift.Result
	for _, k := range aggResult.SortedKeys() {
		driftResults = append(driftResults, *aggResult.Results[k])
	}

	formatter := output.NewFormatter(os.Stdout)
	if aggOutput == "json" {
		return formatter.WriteJSON(driftResults)
	}
	return formatter.WriteText(driftResults)
}
