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
