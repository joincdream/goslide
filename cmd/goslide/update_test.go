package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestUpdateCommand_Help(t *testing.T) {
	initCLI()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"update", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error on 'update --help': %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "--check") {
		t.Errorf("expected --check flag in help output, got: %s", output)
	}
	if !strings.Contains(output, "--dry-run") {
		t.Errorf("expected --dry-run flag in help output, got: %s", output)
	}
	if !strings.Contains(output, "--force") {
		t.Errorf("expected --force flag in help output, got: %s", output)
	}
}
