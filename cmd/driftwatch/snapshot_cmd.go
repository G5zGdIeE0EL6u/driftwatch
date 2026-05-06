package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/helm"
	"github.com/yourorg/driftwatch/internal/output"
)

var (
	snapshotFile      string
	snapshotDiffFile  string
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Capture or diff drift snapshots",
	}

	captureCmd := &cobra.Command{
		Use:   "capture",
		Short: "Capture current drift state to a snapshot file",
		RunE:  runCaptureSnapshot,
	}
	captureCmd.Flags().StringVarP(&releaseFlag, "release", "r", "", "Helm release name (required)")
	captureCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "default", "Kubernetes namespace")
	captureCmd.Flags().StringVar(&snapshotFile, "file", "drift-snapshot.json", "Path to write snapshot")
	_ = captureCmd.MarkFlagRequired("release")

	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Diff current drift against a saved snapshot",
		RunE:  runDiffSnapshot,
	}
	diffCmd.Flags().StringVarP(&releaseFlag, "release", "r", "", "Helm release name (required)")
	diffCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "default", "Kubernetes namespace")
	diffCmd.Flags().StringVar(&snapshotDiffFile, "file", "drift-snapshot.json", "Path to existing snapshot")
	diffCmd.Flags().StringVarP(&formatFlag, "format", "f", "text", "Output format: text|json")
	_ = diffCmd.MarkFlagRequired("release")

	snapshotCmd.AddCommand(captureCmd, diffCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runCaptureSnapshot(cmd *cobra.Command, _ []string) error {
	client, err := helm.NewCachedClient(namespaceFlag)
	if err != nil {
		return fmt.Errorf("helm client: %w", err)
	}
	rel, err := client.GetRelease(releaseFlag)
	if err != nil {
		return fmt.Errorf("get release: %w", err)
	}
	det := drift.NewDetector(rel)
	results := det.Detect()

	path := snapshotFile
	if !filepath.IsAbs(path) {
		path = filepath.Join(".", path)
	}
	s := drift.NewSnapshot(releaseFlag, namespaceFlag, results)
	if err := drift.SaveSnapshot(path, s); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved to %s (%d drifts)\n", path, len(results))
	return nil
}

func runDiffSnapshot(cmd *cobra.Command, _ []string) error {
	prev, err := drift.LoadSnapshot(snapshotDiffFile)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	client, err := helm.NewCachedClient(namespaceFlag)
	if err != nil {
		return fmt.Errorf("helm client: %w", err)
	}
	rel, err := client.GetRelease(releaseFlag)
	if err != nil {
		return fmt.Errorf("get release: %w", err)
	}
	det := drift.NewDetector(rel)
	current := drift.NewSnapshot(releaseFlag, namespaceFlag, det.Detect())

	newDrifts := drift.DiffSnapshots(prev, current)
	if len(newDrifts) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No new drifts since snapshot.")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%d new drift(s) since %s:\n", len(newDrifts), prev.CapturedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Fprintln(cmd.OutOrStdout())
	formatter := output.NewFormatter(cmd.OutOrStdout())
	if err := formatter.Write(formatFlag, newDrifts); err != nil {
		return err
	}
	os.Exit(1)
	return nil
}
