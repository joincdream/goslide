package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

func TestValidateServeOptions(t *testing.T) {
	// 1. Missing argument
	err := validateServeOptions([]string{})
	if err == nil {
		t.Fatal("expected error when no arguments provided, got nil")
	}
	var cliErr *CLIError
	if !errors.As(err, &cliErr) || cliErr.Code != model.ExitInvalidUsage {
		t.Errorf("expected ExitInvalidUsage (2), got %v", err)
	}

	// 2. Non-existent file
	err = validateServeOptions([]string{"non_existent_file.md"})
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	if !errors.As(err, &cliErr) || cliErr.Code != model.ExitFileNotFound {
		t.Errorf("expected ExitFileNotFound (3), got %v", err)
	}

	// 3. Valid file
	tempDir := t.TempDir()
	validFile := filepath.Join(tempDir, "talk.md")
	if writeErr := os.WriteFile(validFile, []byte("# Valid Slide"), 0644); writeErr != nil {
		t.Fatalf("failed to create valid test file: %v", writeErr)
	}

	err = validateServeOptions([]string{validFile})
	if err != nil {
		t.Errorf("unexpected error for valid file: %v", err)
	}
}
