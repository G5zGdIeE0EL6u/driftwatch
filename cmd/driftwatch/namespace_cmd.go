package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/helm"
)

var namespaceCmdFlags struct {
	namespaces []string
	kubeconfig string
	kubeContext string
}

func init() {
	nsCmd := &cobra.Command{
		Use:   "list-releases",
		Short: "List all Helm releases across specified namespaces",
		RunE:  runListReleases,
	}

	nsCmd.Flags().StringSliceVarP(
		&namespaceCmdFlags.namespaces, "namespace", "n", nil,
		"Namespaces to query (repeatable); omit to query all",
	)
	nsCmd.Flags().StringVar(
		&namespaceCmdFlags.kubeconfig, "kubeconfig", "",
		"Path to kubeconfig file (defaults to in-cluster or KUBECONFIG env)",
	)
	nsCmd.Flags().StringVar(
		&namespaceCmdFlags.kubeContext, "kube-context", "",
		"Kubernetes context to use",
	)

	rootCmd.AddCommand(nsCmd)
}

func runListReleases(cmd *cobra.Command, _ []string) error {
	lister, err := helm.NewNamespaceLister(
		namespaceCmdFlags.kubeconfig,
		namespaceCmdFlags.kubeContext,
	)
	if err != nil {
		return fmt.Errorf("create namespace lister: %w", err)
	}

	refs, err := lister.ListReleaseNames(cmd.Context(), namespaceCmdFlags.namespaces)
	if err != nil {
		return fmt.Errorf("list releases: %w", err)
	}

	if len(refs) == 0 {
		fmt.Fprintln(os.Stdout, "No releases found.")
		return nil
	}

	for _, ref := range refs {
		fmt.Fprintln(os.Stdout, ref.String())
	}
	return nil
}
