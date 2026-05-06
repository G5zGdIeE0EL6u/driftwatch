package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
)

var suppressCmd = &cobra.Command{
	Use:   "suppress",
	Short: "Manage drift suppression rules",
}

var suppressAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Append a suppression rule to a rules file",
	RunE:  runSuppressAdd,
}

var (
	suppressFile      string
	suppressRelease   string
	suppressNamespace string
	suppressKeyPrefix string
	suppressReason    string
)

func init() {
	suppressAddCmd.Flags().StringVar(&suppressFile, "file", "suppression-rules.json", "Path to suppression rules JSON file")
	suppressAddCmd.Flags().StringVar(&suppressRelease, "release", "", "Release name to match (empty = any)")
	suppressAddCmd.Flags().StringVar(&suppressNamespace, "namespace", "", "Namespace to match (empty = any)")
	suppressAddCmd.Flags().StringVar(&suppressKeyPrefix, "key-prefix", "", "Key prefix to suppress")
	suppressAddCmd.Flags().StringVar(&suppressReason, "reason", "", "Human-readable reason for suppression")
	_ = suppressAddCmd.MarkFlagRequired("key-prefix")

	suppressCmd.AddCommand(suppressAddCmd)
	rootCmd.AddCommand(suppressCmd)
}

func runSuppressAdd(cmd *cobra.Command, _ []string) error {
	var existing []drift.SuppressionRule

	data, err := os.ReadFile(suppressFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading rules file: %w", err)
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &existing); err != nil {
			return fmt.Errorf("parsing rules file: %w", err)
		}
	}

	newRule := drift.SuppressionRule{
		Release:   suppressRelease,
		Namespace: suppressNamespace,
		KeyPrefix: suppressKeyPrefix,
		Reason:    suppressReason,
	}
	existing = append(existing, newRule)

	out, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling rules: %w", err)
	}
	if err := os.WriteFile(suppressFile, out, 0o644); err != nil {
		return fmt.Errorf("writing rules file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Added suppression rule for key prefix %q to %s\n", suppressKeyPrefix, suppressFile)
	return nil
}
