package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/helm"
	"github.com/yourorg/driftwatch/internal/output"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Evaluate drift results against a policy file",
}

func init() {
	policyCmd.AddCommand(policyEvalCmd)
	rootCmd.AddCommand(policyCmd)

	policyEvalCmd.Flags().StringP("release", "r", "", "Helm release name (required)")
	policyEvalCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace")
	policyEvalCmd.Flags().StringP("policy", "p", "policy.json", "Path to policy file")
	policyEvalCmd.Flags().StringP("format", "f", "text", "Output format: text|json")
	policyEvalCmd.Flags().StringP("kubeconfig", "k", "", "Path to kubeconfig")
	_ = policyEvalCmd.MarkFlagRequired("release")
}

var policyEvalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Evaluate a release's drift against policy rules",
	RunE:  runPolicyEval,
}

func runPolicyEval(cmd *cobra.Command, _ []string) error {
	releaseName, _ := cmd.Flags().GetString("release")
	namespace, _ := cmd.Flags().GetString("namespace")
	policyPath, _ := cmd.Flags().GetString("policy")
	format, _ := cmd.Flags().GetString("format")
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")

	hc, err := helm.NewClient(namespace, kubeconfig)
	if err != nil {
		return fmt.Errorf("helm client: %w", err)
	}

	rel, err := hc.GetRelease(releaseName)
	if err != nil {
		return fmt.Errorf("get release: %w", err)
	}

	values, err := hc.GetValues(releaseName)
	if err != nil {
		return fmt.Errorf("get values: %w", err)
	}

	detector := drift.NewDetector(rel, values)
	driftResults, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("detect drift: %w", err)
	}

	policy, err := drift.LoadPolicy(policyPath)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}

	policyResults := policy.Evaluate(driftResults)

	hasViolation := false
	for _, pr := range policyResults {
		if pr.Violation {
			hasViolation = true
			break
		}
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(policyResults)
	}

	formatter := output.NewFormatter(os.Stdout)
	_ = formatter.WriteText(driftResults)

	if hasViolation {
		fmt.Fprintln(os.Stderr, "\n[POLICY] One or more drift violations detected.")
		for _, pr := range policyResults {
			if pr.Violation {
				fmt.Fprintf(os.Stderr, "  VIOLATION key=%s reason=%s\n", pr.Key, pr.Reason)
			}
		}
		os.Exit(2)
	}
	return nil
}
