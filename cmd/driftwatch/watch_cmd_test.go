package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWatchCmdRequiresRelease(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"watch"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no release argument provided")
	}
}

func TestWatchCmdFlagsParsed(t *testing.T) {
	tests := []struct {
		args    []string
		wantErr bool
	}{
		{[]string{"watch", "--help"}, false},
		{[]string{"watch", "--interval", "30s", "--help"}, false},
		{[]string{"watch", "--output", "json", "--help"}, false},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, "_"), func(t *testing.T) {
			buf := &bytes.Buffer{}
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("args %v: wantErr=%v, got %v", tt.args, tt.wantErr, err)
			}
		})
	}
}
