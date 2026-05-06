package main

import (
	"bytes"
	"testing"
)

func TestPolicyCmdHasSubcommands(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"policy", "--help"})
	_ = rootCmd.Execute()
	if !bytes.Contains(buf.Bytes(), []byte("eval")) {
		t.Errorf("expected 'eval' subcommand in help output, got: %s", buf.String())
	}
}

func TestPolicyEvalCmd_RequiresRelease(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"policy", "eval", "--policy", "policy.json"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when --release is missing")
	}
}

func TestPolicyEvalCmd_FlagsParsed(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"policy", "eval"})
	if err != nil || cmd == nil {
		t.Fatalf("could not find policy eval command: %v", err)
	}
	flags := []string{"release", "namespace", "policy", "format", "kubeconfig"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("expected flag --%s to be registered", f)
		}
	}
}
