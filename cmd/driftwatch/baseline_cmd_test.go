package main

import (
	"bytes"
	"testing"
)

func TestBaselineCaptureCmd_RequiresRelease(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"baseline", "capture", "--namespace", "default"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --release is missing")
	}
}

func TestBaselineDiffCmd_RequiresRelease(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"baseline", "diff", "--namespace", "default"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --release is missing")
	}
}

func TestBaselineCaptureCmd_FlagsParsed(t *testing.T) {
	// Verify flags are registered and parse without error (no real cluster needed).
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	// --help exits cleanly and should not return an error.
	rootCmd.SetArgs([]string{"baseline", "capture", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBaselineDiffCmd_FlagsParsed(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"baseline", "diff", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
