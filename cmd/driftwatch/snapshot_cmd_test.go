package main

import (
	"bytes"
	"testing"
)

func TestSnapshotCaptureCmd_RequiresRelease(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot", "capture", "--namespace", "default"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --release is missing")
	}
}

func TestSnapshotDiffCmd_RequiresRelease(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot", "diff", "--namespace", "staging"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --release is missing")
	}
}

func TestSnapshotCaptureCmd_FlagsParsed(t *testing.T) {
	snapshotCmd := rootCmd.Commands()
	var found bool
	for _, c := range snapshotCmd {
		if c.Use == "snapshot" {
			for _, sub := range c.Commands() {
				if sub.Use == "capture" {
					found = true
					if sub.Flag("file") == nil {
						t.Error("expected --file flag on snapshot capture")
					}
					if sub.Flag("namespace") == nil {
						t.Error("expected --namespace flag on snapshot capture")
					}
				}
			}
		}
	}
	if !found {
		t.Error("snapshot capture subcommand not registered")
	}
}

func TestSnapshotDiffCmd_FlagsParsed(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "snapshot" {
			for _, sub := range c.Commands() {
				if sub.Use == "diff" {
					if sub.Flag("format") == nil {
						t.Error("expected --format flag on snapshot diff")
					}
					if sub.Flag("file") == nil {
						t.Error("expected --file flag on snapshot diff")
					}
					return
				}
			}
		}
	}
	t.Error("snapshot diff subcommand not registered")
}
