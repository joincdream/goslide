package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

func TestIsPortInUseError(t *testing.T) {
	tests := []struct {
		errText  string
		expected bool
	}{
		{"bind: address already in use", true},
		{"listen tcp 127.0.0.1:8080: bind: address already in use", true},
		{"connect: connection refused", false},
		{"permission denied", false},
	}

	for _, tt := range tests {
		res := isPortInUseError(errors.New(tt.errText))
		if res != tt.expected {
			t.Errorf("isPortInUseError(%q) = %v, expected %v", tt.errText, res, tt.expected)
		}
	}

	if isPortInUseError(nil) {
		t.Errorf("isPortInUseError(nil) should be false")
	}
}

func TestFormatPortInUseError(t *testing.T) {
	err := formatPortInUseError(8080, "demo.md")
	if err == nil {
		t.Fatal("expected formatted error, got nil")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "8080") {
		t.Errorf("expected error to mention port 8080, got: %s", errStr)
	}
	if !strings.Contains(errStr, "--port") {
		t.Errorf("expected error to suggest --port flag, got: %s", errStr)
	}
}

func TestTruncateBannerText(t *testing.T) {
	short := "short text"
	if truncateBannerText(short, 20) != short {
		t.Errorf("expected %q, got %q", short, truncateBannerText(short, 20))
	}

	long := "this is a very long text that should be truncated to fit"
	truncated := truncateBannerText(long, 20)
	if len(truncated) > 20 {
		t.Errorf("expected length <= 20, got %d", len(truncated))
	}
	if !strings.HasSuffix(truncated, "...") {
		t.Errorf("expected suffix ..., got: %s", truncated)
	}
}
