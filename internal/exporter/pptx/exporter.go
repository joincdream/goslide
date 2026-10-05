package pptx

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yundream/goslide/internal/model"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
)

// Option configures the PPTXExporter.
type Option func(*PPTXExporter)

// WithTheme sets the theme name for the HTML renderer.
func WithTheme(theme string) Option {
	return func(e *PPTXExporter) {
		e.theme = theme
	}
}

// WithCustomCSS sets the path to an external CSS file.
func WithCustomCSS(cssPath string) Option {
	return func(e *PPTXExporter) {
		e.customCSS = cssPath
	}
}

// WithStandalone controls Base64 embedding of local images.
func WithStandalone(standalone bool) Option {
	return func(e *PPTXExporter) {
		e.standalone = standalone
	}
}

// WithBaseDir sets the base working directory for asset resolution.
func WithBaseDir(baseDir string) Option {
	return func(e *PPTXExporter) {
		e.baseDir = baseDir
	}
}

// WithTimeout sets the maximum duration for headless Chrome capture.
func WithTimeout(timeout time.Duration) Option {
	return func(e *PPTXExporter) {
		e.timeout = timeout
	}
}

// WithBrowserPath explicitly sets the path to the Chrome/Chromium binary.
func WithBrowserPath(path string) Option {
	return func(e *PPTXExporter) {
		e.browserPath = path
	}
}

// PPTXExporter implements the model.Exporter interface for PowerPoint 16:9 presentation generation.
type PPTXExporter struct {
	theme       string
	customCSS   string
	standalone  bool
	baseDir     string
	timeout     time.Duration
	browserPath string
}

// NewExporter creates a new PPTXExporter instance with the given options.
func NewExporter(opts ...Option) *PPTXExporter {
	e := &PPTXExporter{
		standalone: true,
		timeout:    60 * time.Second,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Export renders the deck to an intermediate HTML representation, captures 1920x1080 slide images,
// and packages them into a native OpenXML PPTX presentation file with speaker notes preserved.
func (e *PPTXExporter) Export(ctx context.Context, deck *model.Deck, outputPath string) error {
	if deck == nil {
		return fmt.Errorf("%w: cannot export nil deck", model.ErrExportFailed)
	}
	if len(deck.Slides) == 0 {
		return fmt.Errorf("%w: deck contains no slides", model.ErrExportFailed)
	}
	if outputPath == "" {
		return fmt.Errorf("%w: output path cannot be empty", model.ErrExportFailed)
	}

	// 1. Render intermediate HTML
	tempPath, err := e.renderIntermediateHTML(ctx, deck)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tempPath) }()

	// 2. Capture slide screenshot buffers via headless Chrome
	slides, err := CaptureSlides(ctx, tempPath, len(deck.Slides), CaptureOptions{
		Timeout:     e.timeout,
		BrowserPath: e.browserPath,
	})
	if err != nil {
		return err
	}

	// 3. Assemble and package OpenXML zip archive
	if err := writePPTXPackage(outputPath, deck, slides); err != nil {
		return fmt.Errorf("%w: failed to assemble PPTX package: %w", model.ErrExportFailed, err)
	}

	return nil
}

func (e *PPTXExporter) renderIntermediateHTML(ctx context.Context, deck *model.Deck) (string, error) {
	renderer := htmlrenderer.NewRenderer(
		htmlrenderer.WithTheme(e.theme),
		htmlrenderer.WithCustomCSS(e.customCSS),
		htmlrenderer.WithStandalone(e.standalone),
		htmlrenderer.WithBaseDir(e.baseDir),
	)

	tempDir := e.baseDir
	if tempDir == "" || tempDir == "." {
		tempDir = ""
	}
	tempFile, err := os.CreateTemp(tempDir, ".goslide-pptx-*.html")
	if err != nil {
		tempFile, err = os.CreateTemp("", "goslide-pptx-*.html")
		if err != nil {
			return "", fmt.Errorf("%w: failed to create temporary HTML file: %w", model.ErrExportFailed, err)
		}
	}
	tempPath := tempFile.Name()

	if err := renderer.Render(ctx, deck, tempFile); err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
		return "", fmt.Errorf("%w: failed to render intermediate HTML: %w", model.ErrExportFailed, err)
	}
	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return "", fmt.Errorf("%w: failed to flush temporary HTML file: %w", model.ErrExportFailed, err)
	}

	return tempPath, nil
}

func writePPTXPackage(outputPath string, deck *model.Deck, slides []CapturedSlide) error {
	if dir := filepath.Dir(outputPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %q: %w", dir, err)
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", outputPath, err)
	}
	defer func() { _ = outFile.Close() }()

	zw := zip.NewWriter(outFile)
	defer func() { _ = zw.Close() }()

	slideCount := len(deck.Slides)
	hasNotes := make([]bool, slideCount)
	for i, s := range deck.Slides {
		hasNotes[i] = strings.TrimSpace(s.Notes) != ""
	}

	// Root entries
	if err := writeZipEntry(zw, "[Content_Types].xml", buildContentTypesXML(slideCount, hasNotes)); err != nil {
		return err
	}
	if err := writeZipEntry(zw, "_rels/.rels", buildRootRelsXML()); err != nil {
		return err
	}
	if err := writeZipEntry(zw, "ppt/presentation.xml", buildPresentationXML(slideCount)); err != nil {
		return err
	}
	if err := writeZipEntry(zw, "ppt/_rels/presentation.xml.rels", buildPresentationRelsXML(slideCount)); err != nil {
		return err
	}

	// Slide entries
	return writeSlideEntries(zw, deck, slides, hasNotes)
}

func writeSlideEntries(zw *zip.Writer, deck *model.Deck, slides []CapturedSlide, hasNotes []bool) error {
	for i := 1; i <= len(deck.Slides); i++ {
		idx := i - 1
		var slide CapturedSlide
		if idx < len(slides) {
			slide = slides[idx]
		}
		if err := writeSingleSlide(zw, i, slide, hasNotes[idx], deck.Slides[idx].Notes); err != nil {
			return err
		}
	}
	return nil
}

func writeSingleSlide(zw *zip.Writer, i int, slide CapturedSlide, hasNotes bool, notes string) error {
	if len(slide.Image) > 0 {
		if err := writeZipEntry(zw, fmt.Sprintf("ppt/media/slide%d.png", i), slide.Image); err != nil {
			return err
		}
	}

	if err := writeZipEntry(zw, fmt.Sprintf("ppt/slides/slide%d.xml", i), buildSlideXML(slide.Links)); err != nil {
		return err
	}

	if err := writeZipEntry(zw, fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", i), buildSlideRelsXML(i, hasNotes, slide.Links)); err != nil {
		return err
	}

	if !hasNotes {
		return nil
	}

	if err := writeZipEntry(zw, fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", i), buildNotesSlideXML(notes)); err != nil {
		return err
	}
	return writeZipEntry(zw, fmt.Sprintf("ppt/notesSlides/_rels/notesSlide%d.xml.rels", i), buildNotesSlideRelsXML(i))
}

func writeZipEntry(zw *zip.Writer, path string, content []byte) error {
	w, err := zw.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create zip entry %q: %w", path, err)
	}
	if _, err := w.Write(content); err != nil {
		return fmt.Errorf("failed to write zip entry %q: %w", path, err)
	}
	return nil
}
