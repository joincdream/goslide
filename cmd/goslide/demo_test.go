package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDemoCmd_Success(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}

	demoForceFlag = false
	if err := demoCmd.RunE(demoCmd, []string{}); err != nil {
		t.Fatalf("demoCmd failed unexpectedly: %v", err)
	}

	// 1. Verify demo/demo.md created
	demoPath := filepath.Join(tempDir, "demo", "demo.md")
	if info, err := os.Stat(demoPath); err != nil || info.IsDir() {
		t.Fatalf("expected demo/demo.md to exist: %v", err)
	}

	// 2. Verify root themes/ folder and stylesheets created
	for _, th := range []string{"clean.css", "dark.css", "academic.css"} {
		thPath := filepath.Join(tempDir, "themes", th)
		if info, err := os.Stat(thPath); err != nil || info.IsDir() {
			t.Errorf("expected %s to exist in root themes/: %v", th, err)
		}
	}

	// 3. Verify demo images extracted into demo/
	for _, img := range []string{"hands-on-bg.jpeg", "sample-img-01.png", "sample-img-02.jpg", "llm-architecture.jpg"} {
		imgPath := filepath.Join(tempDir, "demo", img)
		if info, err := os.Stat(imgPath); err != nil || info.IsDir() {
			t.Errorf("expected demo image %s to exist in demo/: %v", img, err)
		}
	}
}

func TestDemoCmd_ExistingFile_SafetyAndForce(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}

	// Create dummy demo/demo.md
	if err := os.MkdirAll("demo", 0755); err != nil {
		t.Fatalf("failed to create demo dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join("demo", "demo.md"), []byte("# Existing"), 0644); err != nil {
		t.Fatalf("failed to create existing demo.md: %v", err)
	}

	// Run without force: must error
	demoForceFlag = false
	if err := demoCmd.RunE(demoCmd, []string{}); err == nil {
		t.Fatal("expected error when demo/demo.md already exists without --force, got nil")
	}

	// Run with force: must succeed
	demoForceFlag = true
	if err := demoCmd.RunE(demoCmd, []string{}); err != nil {
		t.Fatalf("expected success with --force, got error: %v", err)
	}
}

func TestPromptCmd(t *testing.T) {
	// 1. Test slide prompt
	var buf bytes.Buffer
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	err = promptCmd.RunE(promptCmd, []string{})
	_ = w.Close()
	os.Stdout = origStdout

	if err != nil {
		t.Fatalf("promptCmd failed: %v", err)
	}

	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "dsl-guide.md") {
		t.Errorf("prompt output missing dsl-guide URL: %s", output)
	}

	// 2. Test theme prompt
	var themeBuf bytes.Buffer
	r2, w2, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w2

	err = promptCmd.RunE(promptCmd, []string{"theme"})
	_ = w2.Close()
	os.Stdout = origStdout

	if err != nil {
		t.Fatalf("promptCmd theme failed: %v", err)
	}

	_, _ = themeBuf.ReadFrom(r2)
	themeOutput := themeBuf.String()

	if !strings.Contains(themeOutput, "theme-guide.md") {
		t.Errorf("prompt theme output missing theme-guide URL: %s", themeOutput)
	}
}

func TestThemeCmd_ListAndExport(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}

	// 1. Test theme list
	if err := themeListCmd.RunE(themeListCmd, []string{}); err != nil {
		t.Fatalf("themeListCmd failed: %v", err)
	}

	// 2. Test theme export valid theme
	exportOut := filepath.Join(tempDir, "my-theme.css")
	if err := themeExportCmd.RunE(themeExportCmd, []string{"clean", exportOut}); err != nil {
		t.Fatalf("themeExportCmd failed: %v", err)
	}
	if info, err := os.Stat(exportOut); err != nil || info.IsDir() {
		t.Fatalf("expected exported css file %s to exist", exportOut)
	}

	// 3. Test theme export invalid theme returns error
	if err := themeExportCmd.RunE(themeExportCmd, []string{"nonexistent-theme"}); err == nil {
		t.Fatal("expected error for nonexistent theme, got nil")
	}
}
