package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCommand(t *testing.T) {
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

	// Reset flags
	themeFlag = ""
	themePathFlag = ""
	outputPathFlag = ""
}

func TestBuildCommand_Standalone(t *testing.T) {
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

	// Reset flags
	standaloneFlag = false
	outputPathFlag = ""
}
