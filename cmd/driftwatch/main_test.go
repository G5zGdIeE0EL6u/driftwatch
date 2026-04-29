package main

import (
	"bytes"
	"testing"
)

func TestRootCmdHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	// Execute should not return an error for --help
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected help output, got empty string")
	}
}

func TestDetectCmdRequiresRelease(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"detect"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --release flag is missing, got nil")
	}
}

func TestDetectCmdFlagsParsed(t *testing.T) {
	rootCmd.SetArgs([]string{"detect", "--release", "my-app", "--namespace", "staging"})

	// Override RunE to avoid actual execution
	detectCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if release != "my-app" {
			t.Errorf("expected release=my-app, got %q", release)
		}
		if namespace != "staging" {
			t.Errorf("expected namespace=staging, got %q", namespace)
		}
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
