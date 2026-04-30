package main

import (
	"fmt"
	"os"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/helm"
	"github.com/example/driftwatch/internal/output"
	"github.com/example/driftwatch/internal/watch"
	"github.com/spf13/cobra"
)

var (
	watchInterval  time.Duration
	watchNamespace string
	watchFormat    string
)

var watchCmd = &cobra.Command{
	Use:   "watch <release>",
	Short: "Continuously watch a Helm release for configuration drift",
	Args:  cobra.ExactArgs(1),
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().DurationVarP(&watchInterval, "interval", "i", 60*time.Second, "polling interval")
	watchCmd.Flags().StringVarP(&watchNamespace, "namespace", "n", "default", "Kubernetes namespace")
	watchCmd.Flags().StringVarP(&watchFormat, "output", "o", "text", "output format: text or json")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	releaseName := args[0]

	helmClient, err := helm.NewClient(watchNamespace)
	if err != nil {
		return fmt.Errorf("helm client: %w", err)
	}

	rel, err := helmClient.GetRelease(releaseName, watchNamespace)
	if err != nil {
		return fmt.Errorf("get release: %w", err)
	}

	chartDefaults := rel.Chart.Values
	detector := drift.NewDetector(chartDefaults)
	fmt := output.NewFormatter(watchFormat, os.Stdout)

	cfg := watch.Config{
		Release:   releaseName,
		Namespace: watchNamespace,
		Interval:  watchInterval,
		Formatter: fmt,
	}

	w := watch.New(cfg, helmClient, detector)
	return w.Run(cmd.Context())
}
