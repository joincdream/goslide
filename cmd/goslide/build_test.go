package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

func resetFlags() {
	outputPathFlag = ""
	themeFlag = ""
	themePathFlag = ""
	standaloneFlag = false
	quietFlag = false
	verboseFlag = false
}

func TestBuildCommand(t *testing.T) {
	defer resetFlags()
	resetFlags()

	tempDir := t.TempDir()
	inputMD := filepath.Join(tempDir, "sample.md")
	outputHTML := filepath.Join(tempDir, "sample.html")

	content := `---
title: "CLI Preview Test"
theme: "clean"
---
<!-- _layout: cover -->
<!-- _backgroundColor: #112233 -->
<!-- _color: #ffffff -->
# CLI Test Slide

---
<!-- _layout: two-cols -->
Left Column
<!-- split -->
Right Column

---
` + "```go\nfunc hello() {}\n```"

	if err := os.WriteFile(inputMD, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	buildCmd.SetArgs([]string{inputMD, "-o", outputHTML})
	if err := buildCmd.Execute(); err != nil {
		t.Fatalf("build command failed: %v", err)
	}

	outBytes, err := os.ReadFile(outputHTML)
	if err != nil {
		t.Fatalf("failed to read output HTML: %v", err)
	}

	outStr := string(outBytes)
	if !strings.Contains(outStr, "<title>CLI Preview Test</title>") {
		t.Errorf("missing title in output HTML")
	}
	if !strings.Contains(outStr, "Goslide Clean Theme") {
		t.Errorf("expected clean theme css in output HTML")
	}
	if !strings.Contains(outStr, "background-color: #112233;") {
		t.Errorf("missing background-color in output HTML")
	}
	if !strings.Contains(outStr, "class=\"two-cols\"") {
		t.Errorf("missing two-cols in output HTML")
	}
	if !strings.Contains(outStr, "<span style=") {
		t.Errorf("missing Chroma highlighted code in output HTML")
	}
}

func TestBuildCommand_ThemeFlags(t *testing.T) {
	defer resetFlags()
	resetFlags()

	tempDir := t.TempDir()
	inputMD := filepath.Join(tempDir, "theme_override.md")
	customCSS := filepath.Join(tempDir, "brand.css")
	outputHTML := filepath.Join(tempDir, "theme_override.html")

	mdContent := "# Test Theme Override\n"
	if err := os.WriteFile(inputMD, []byte(mdContent), 0600); err != nil {
		t.Fatalf("failed to write input md: %v", err)
	}

	cssContent := ".custom-brand { color: #abcdef; }"
	if err := os.WriteFile(customCSS, []byte(cssContent), 0600); err != nil {
		t.Fatalf("failed to write custom css: %v", err)
	}

	// Test CLI theme flag overriding frontmatter and custom theme path
	buildCmd.SetArgs([]string{inputMD, "-o", outputHTML, "--theme", "dark", "--theme-path", customCSS})
	if err := buildCmd.Execute(); err != nil {
		t.Fatalf("build command failed with theme flags: %v", err)
	}

	outBytes, err := os.ReadFile(outputHTML)
	if err != nil {
		t.Fatalf("failed to read output HTML: %v", err)
	}

	outStr := string(outBytes)
	if !strings.Contains(outStr, "Goslide Dark Theme") {
		t.Errorf("expected Dark theme css in output HTML")
	}
	if !strings.Contains(outStr, ".custom-brand { color: #abcdef; }") {
		t.Errorf("expected custom css content in output HTML")
	}
}

func TestBuildCommand_Standalone(t *testing.T) {
	defer resetFlags()
	resetFlags()

	tempDir := t.TempDir()
	inputMD := filepath.Join(tempDir, "standalone.md")
	outputHTML := filepath.Join(tempDir, "standalone.html")
	imagePath := filepath.Join(tempDir, "icon.png")

	// 1x1 png
	pngBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
		0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
	if err := os.WriteFile(imagePath, pngBytes, 0600); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}

	mdContent := "# Standalone Slide\n\n![Icon](icon.png)\n"
	if err := os.WriteFile(inputMD, []byte(mdContent), 0600); err != nil {
		t.Fatalf("failed to write input md: %v", err)
	}

	buildCmd.SetArgs([]string{inputMD, "-o", outputHTML, "--standalone"})
	if err := buildCmd.Execute(); err != nil {
		t.Fatalf("build command failed with --standalone: %v", err)
	}

	outBytes, err := os.ReadFile(outputHTML)
	if err != nil {
		t.Fatalf("failed to read output HTML: %v", err)
	}

	outStr := string(outBytes)
	if !strings.Contains(outStr, "data:image/png;base64,") {
		t.Errorf("expected base64 data uri in standalone output HTML")
	}
}

func TestBuildCommand_QuietAndVerbose(t *testing.T) {
	defer resetFlags()
	resetFlags()

	tempDir := t.TempDir()
	inputMD := filepath.Join(tempDir, "sample.md")
	outputHTML := filepath.Join(tempDir, "sample.html")

	if err := os.WriteFile(inputMD, []byte("# Simple Slide\n"), 0600); err != nil {
		t.Fatalf("failed to write input md: %v", err)
	}

	// 1. Mutually exclusive test: --quiet and --verbose together
	buildCmd.SetArgs([]string{inputMD, "-o", outputHTML, "-q", "-v"})
	err := buildCmd.Execute()
	if err == nil {
		t.Fatalf("expected error when specifying both -q and -v")
	}
	if code := determineExitCode(err); code != model.ExitInvalidUsage {
		t.Errorf("expected ExitInvalidUsage (%d), got %d", model.ExitInvalidUsage, code)
	}

	// 2. Quiet mode: capture stdout to verify no success message is printed
	resetFlags()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	buildCmd.SetArgs([]string{inputMD, "-o", outputHTML, "--quiet"})
	execErr := buildCmd.Execute()

	_ = w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	if execErr != nil {
		t.Fatalf("buildCmd failed in quiet mode: %v", execErr)
	}
	if strings.Contains(buf.String(), "Successfully built") {
		t.Errorf("quiet mode should suppress success message, got: %s", buf.String())
	}

	// 3. Verbose mode: verify [verbose] diagnostics are printed
	resetFlags()
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	buildCmd.SetArgs([]string{inputMD, "-o", outputHTML, "--verbose"})
	execErr2 := buildCmd.Execute()

	_ = w2.Close()
	os.Stdout = oldStdout
	var buf2 bytes.Buffer
	_, _ = io.Copy(&buf2, r2)

	if execErr2 != nil {
		t.Fatalf("buildCmd failed in verbose mode: %v", execErr2)
	}
	if !strings.Contains(buf2.String(), "[verbose]") {
		t.Errorf("verbose mode should print [verbose] lines, got: %s", buf2.String())
	}
}

func TestBuildCommand_ExitCodes(t *testing.T) {
	defer resetFlags()

	tempDir := t.TempDir()

	// 1. Missing input file: ExitFileNotFound (3)
	resetFlags()
	buildCmd.SetArgs([]string{filepath.Join(tempDir, "non_existent.md")})
	err := buildCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for non-existent file")
	}
	if code := determineExitCode(err); code != model.ExitFileNotFound {
		t.Errorf("expected ExitFileNotFound (%d), got %d (err: %v)", model.ExitFileNotFound, code, err)
	}

	// 2. Missing custom CSS file: ExitFileNotFound (3)
	resetFlags()
	validMD := filepath.Join(tempDir, "valid.md")
	if writeErr := os.WriteFile(validMD, []byte("# Valid Slide\n"), 0600); writeErr != nil {
		t.Fatalf("failed to write valid md: %v", writeErr)
	}
	buildCmd.SetArgs([]string{validMD, "--theme-path", filepath.Join(tempDir, "missing.css")})
	err = buildCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for non-existent theme-path CSS")
	}
	if code := determineExitCode(err); code != model.ExitFileNotFound {
		t.Errorf("expected ExitFileNotFound (%d) for missing CSS, got %d (err: %v)", model.ExitFileNotFound, code, err)
	}

	// 3. Invalid YAML frontmatter: ExitParseError (4)
	resetFlags()
	invalidFmMD := filepath.Join(tempDir, "bad_frontmatter.md")
	badFM := "---\ntitle: [invalid: yaml: syntax:\n---\n# Bad FM\n"
	if writeErr := os.WriteFile(invalidFmMD, []byte(badFM), 0600); writeErr != nil {
		t.Fatalf("failed to write bad fm md: %v", writeErr)
	}
	buildCmd.SetArgs([]string{invalidFmMD})
	err = buildCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for invalid frontmatter")
	}
	if code := determineExitCode(err); code != model.ExitParseError {
		t.Errorf("expected ExitParseError (%d), got %d (err: %v)", model.ExitParseError, code, err)
	}

	// 4. Missing required input argument: ExitInvalidUsage (2)
	resetFlags()
	buildCmd.SetArgs([]string{})
	err = buildCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for missing input argument")
	}
	if code := determineExitCode(err); code != model.ExitInvalidUsage {
		t.Errorf("expected ExitInvalidUsage (%d), got %d", model.ExitInvalidUsage, code)
	}
}
