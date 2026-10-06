package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCmd(t *testing.T) {
	tempDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir to tempDir: %v", err)
	}

	initCLI()

	t.Run("default presentation.md creation", func(t *testing.T) {
		initThemeFlag = "clean"
		initForceFlag = false

		rootCmd.SetArgs([]string{"init"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error executing init: %v", err)
		}

		target := filepath.Join(tempDir, "presentation.md")
		data, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("expected presentation.md to exist: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "theme: \"clean\"") {
			t.Errorf("expected clean theme in generated template")
		}
		if !strings.Contains(content, "<!-- _layout: cover -->") {
			t.Errorf("expected cover layout in generated template")
		}
	})

	t.Run("existing file fails without force", func(t *testing.T) {
		initThemeFlag = "clean"
		initForceFlag = false

		rootCmd.SetArgs([]string{"init", "presentation.md"})
		err := rootCmd.Execute()
		if err == nil {
			t.Fatalf("expected error when file already exists without --force")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("expected 'already exists' error message, got: %v", err)
		}
	})

	t.Run("existing file overwrites with force", func(t *testing.T) {
		initThemeFlag = "dark"
		initForceFlag = true

		rootCmd.SetArgs([]string{"init", "presentation.md", "--force", "--theme", "dark"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error executing init with force: %v", err)
		}

		target := filepath.Join(tempDir, "presentation.md")
		data, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("failed to read overwritten file: %v", err)
		}

		if !strings.Contains(string(data), "theme: \"dark\"") {
			t.Errorf("expected dark theme after overwrite")
		}
	})

	t.Run("auto-append .md extension", func(t *testing.T) {
		initThemeFlag = "clean"
		initForceFlag = false

		rootCmd.SetArgs([]string{"init", "talk"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error executing init with extension-less name: %v", err)
		}

		target := filepath.Join(tempDir, "talk.md")
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected talk.md to be created: %v", err)
		}
	})
}
