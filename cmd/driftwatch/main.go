package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeconfig string
	namespace  string
	release    string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "driftwatch",
	Short: "Detect configuration drift between running Kubernetes workloads and their source Helm charts",
	Long: `driftwatch compares live Kubernetes resources against the manifests
defined in Helm charts to identify configuration drift.`,
}

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect drift for a Helm release",
	RunE:  runDetect,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file (defaults to in-cluster config)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")

	detectCmd.Flags().StringVarP(&release, "release", "r", "", "Helm release name (required)")
	_ = detectCmd.MarkFlagRequired("release")

	rootCmd.AddCommand(detectCmd)
}

func runDetect(cmd *cobra.Command, args []string) error {
	fmt.Printf("Detecting drift for release %q in namespace %q\n", release, namespace)
	// TODO: wire up drift detection logic
	return nil
}
