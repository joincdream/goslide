package pdf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yundream/goslide/internal/model"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
)

// Option configures the PDFExporter.
type Option func(*PDFExporter)

// WithTheme sets the theme name for the HTML renderer.
func WithTheme(theme string) Option {
	return func(e *PDFExporter) {
		e.theme = theme
	}
}

// WithCustomCSS sets the path to an external CSS file.
func WithCustomCSS(cssPath string) Option {
	return func(e *PDFExporter) {
		e.customCSS = cssPath
	}
}

// WithStandalone controls Base64 embedding of local images.
func WithStandalone(standalone bool) Option {
	return func(e *PDFExporter) {
		e.standalone = standalone
	}
}

// WithBaseDir sets the base working directory for asset resolution.
func WithBaseDir(baseDir string) Option {
	return func(e *PDFExporter) {
		e.baseDir = baseDir
	}
}

// WithTimeout sets the maximum time allowed for headless Chrome export.
func WithTimeout(timeout time.Duration) Option {
	return func(e *PDFExporter) {
		e.timeout = timeout
	}
}

// WithBrowserPath explicitly sets the path to the Chrome/Chromium binary.
func WithBrowserPath(path string) Option {
	return func(e *PDFExporter) {
		e.browserPath = path
	}
}

// PDFExporter implements the model.Exporter interface using headless Chrome.
type PDFExporter struct {
	theme       string
	customCSS   string
	standalone  bool
	baseDir     string
	timeout     time.Duration
	browserPath string
}

// NewExporter creates a new PDFExporter with the given options.
func NewExporter(opts ...Option) *PDFExporter {
	e := &PDFExporter{
		standalone: true,
		timeout:    30 * time.Second,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Export renders the deck to an intermediate HTML representation and prints it to a vector PDF file.
func (e *PDFExporter) Export(ctx context.Context, deck *model.Deck, outputPath string) error {
	if deck == nil {
		return fmt.Errorf("%w: cannot export nil deck", model.ErrExportFailed)
	}
	if outputPath == "" {
		return fmt.Errorf("%w: output path cannot be empty", model.ErrExportFailed)
	}

	// 1. Prepare HTML Renderer
	renderer := htmlrenderer.NewRenderer(
		htmlrenderer.WithTheme(e.theme),
		htmlrenderer.WithCustomCSS(e.customCSS),
		htmlrenderer.WithStandalone(e.standalone),
		htmlrenderer.WithBaseDir(e.baseDir),
	)

	// 2. Create temporary HTML file
	tempDir := e.baseDir
	if tempDir == "" || tempDir == "." {
		tempDir = ""
	}
	tempFile, err := os.CreateTemp(tempDir, ".goslide-pdf-*.html")
	if err != nil {
		tempFile, err = os.CreateTemp("", "goslide-pdf-*.html")
		if err != nil {
			return fmt.Errorf("%w: failed to create temporary HTML file: %w", model.ErrExportFailed, err)
		}
	}
	tempPath := tempFile.Name()
	defer func() { _ = os.Remove(tempPath) }()

	// 3. Render deck into temporary HTML file
	if err := renderer.Render(ctx, deck, tempFile); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("%w: failed to render intermediate HTML: %w", model.ErrExportFailed, err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("%w: failed to flush temporary HTML file: %w", model.ErrExportFailed, err)
	}

	// 4. Print HTML to PDF using headless Chrome
	printOpts := PrintOptions{
		Timeout:     e.timeout,
		BrowserPath: e.browserPath,
	}

	pdfData, err := PrintHTMLToPDF(ctx, tempPath, printOpts)
	if err != nil {
		return fmt.Errorf("%w: %w", model.ErrExportFailed, err)
	}

	// 5. Ensure output directory exists
	if dir := filepath.Dir(outputPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("%w: failed to create output directory %q: %w", model.ErrExportFailed, dir, err)
		}
	}

	// 6. Write final vector PDF
	if err := os.WriteFile(outputPath, pdfData, 0644); err != nil {
		return fmt.Errorf("%w: failed to write PDF to %q: %w", model.ErrExportFailed, outputPath, err)
	}

	return nil
}
