package pdf

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/parser"
)

func TestFindChrome(t *testing.T) {
	t.Run("Finds system chrome if available", func(t *testing.T) {
		path, err := FindChrome("")
		if err != nil {
			t.Skipf("Chrome/Chromium not installed on this system: %v", err)
		}
		if path == "" {
			t.Errorf("expected non-empty browser path")
		}
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			t.Errorf("browser path %q is not a valid file", path)
		}
	})

	t.Run("Fails gracefully with invalid custom path", func(t *testing.T) {
		_, err := FindChrome("/path/to/nonexistent/chrome_binary")
		if err == nil {
			t.Fatalf("expected error for non-existent chrome binary, got nil")
		}
		if !errors.Is(err, model.ErrChromeNotFound) {
			t.Errorf("expected ErrChromeNotFound, got: %v", err)
		}
	})
}

func TestPDFExporter_NilDeck(t *testing.T) {
	e := NewExporter()
	err := e.Export(context.Background(), nil, "output.pdf")
	if err == nil {
		t.Fatalf("expected error for nil deck, got nil")
	}
	if !errors.Is(err, model.ErrExportFailed) {
		t.Errorf("expected ErrExportFailed, got: %v", err)
	}
}

func TestPDFExporter_Export(t *testing.T) {
	// Skip if Chrome is not installed in the test environment
	browserPath, err := FindChrome("")
	if err != nil {
		t.Skipf("skipping PDF export integration test: chrome not available: %v", err)
	}

	p := parser.NewParser()
	sampleMD := `---
title: "Vector PDF Test Deck"
theme: "clean"
---

<!-- _layout: cover -->
# Vector PDF Slide 1

Welcome to Goslide PDF Export

---

<!-- _layout: two-cols -->
## Slide 2: Two Columns

Left Column Content

<!-- split -->

Right Column Content

---

## Slide 3: Code Block

` + "```go\npackage main\n\nfunc main() {\n\tprintln(\"Hello PDF\")\n}\n```"

	deck, err := p.Parse(context.Background(), strings.NewReader(sampleMD))
	if err != nil {
		t.Fatalf("failed to parse sample markdown: %v", err)
	}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test_output.pdf")

	exporter := NewExporter(
		WithTheme("clean"),
		WithBrowserPath(browserPath),
		WithTimeout(30*time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	if err := exporter.Export(ctx, deck, outputPath); err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	// 1. Verify file exists and has non-zero size
	pdfData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read generated PDF: %v", err)
	}
	if len(pdfData) < 1000 {
		t.Errorf("generated PDF is unusually small: %d bytes", len(pdfData))
	}

	// 2. Verify PDF Magic Number (%PDF-)
	if !bytes.HasPrefix(pdfData, []byte("%PDF-")) {
		t.Errorf("generated file does not start with %%PDF- magic header, prefix: %q", string(pdfData[:min(len(pdfData), 10)]))
	}

	// 3. Verify PDF EOF Marker
	if !bytes.Contains(pdfData, []byte("%%EOF")) {
		t.Errorf("generated file does not contain %%%%EOF marker")
	}
}
