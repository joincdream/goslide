package html_test

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
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
)

func sampleDeck() *model.Deck {
	return &model.Deck{
		Title:     "Sample Tech Presentation",
		Author:    "Presenter",
		CreatedAt: time.Now(),
		GlobalAttrs: model.GlobalDirectives{
			Theme:    "default",
			Layout:   model.LayoutDefault,
			Size:     model.Ratio16x9,
			Paginate: true,
			Header:   "Tech Conference 2026",
			Footer:   "Confidential",
		},
		Slides: []*model.Slide{
			{
				Index:  1,
				Layout: model.LayoutCover,
				Directives: model.SlideDirectives{
					Layout: model.LayoutCover,
					Class:  []string{"lead"},
				},
				HTMLContent: "<h1>Welcome to Goslide</h1><p>High Performance Presentations</p>",
			},
			{
				Index:  2,
				Layout: model.LayoutDefault,
				Directives: model.SlideDirectives{
					Layout:   model.LayoutDefault,
					Paginate: true,
				},
				HTMLContent: "<h2>Architecture Overview</h2><p>Here is how it works.</p>",
			},
			{
				Index:  3,
				Layout: model.LayoutTwoCols,
				Directives: model.SlideDirectives{
					Layout: model.LayoutTwoCols,
				},
				LeftHTML:  "<h3>Left Column</h3><p>Data stream</p>",
				RightHTML: "<h3>Right Column</h3><p>Graph node</p>",
			},
		},
	}
}

func renderSampleHTML(t *testing.T, opts ...htmlrenderer.Option) string {
	t.Helper()
	renderer := htmlrenderer.NewRenderer(opts...)
	deck := sampleDeck()
	var buf bytes.Buffer
	if err := renderer.Render(context.Background(), deck, &buf); err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}
	return buf.String()
}

func TestHTMLRenderer_Render_DocumentStructure(t *testing.T) {
	t.Parallel()

	html := renderSampleHTML(t)

	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Errorf("missing <!DOCTYPE html>")
	}
	if !strings.Contains(html, "<title>Sample Tech Presentation</title>") {
		t.Errorf("missing or incorrect title tag")
	}
	if !strings.Contains(html, `id="goslide-stage"`) {
		t.Errorf("missing goslide-stage container")
	}
	if !strings.Contains(html, `id="goslide-deck"`) {
		t.Errorf("missing goslide-deck container")
	}
}

func TestHTMLRenderer_Render_SlidesAndLayouts(t *testing.T) {
	t.Parallel()

	html := renderSampleHTML(t)

	if !strings.Contains(html, `data-slide="1"`) || !strings.Contains(html, `data-slide="2"`) {
		t.Errorf("missing slide data attributes")
	}
	if !strings.Contains(html, `active`) {
		t.Errorf("first slide should have active class")
	}
	if !strings.Contains(html, `class="two-cols"`) {
		t.Errorf("missing two-cols container for slide 3")
	}
	if !strings.Contains(html, "Left Column") || !strings.Contains(html, "Right Column") {
		t.Errorf("missing two-cols left/right content")
	}
}

func TestHTMLRenderer_Render_CanvasAndRuntime(t *testing.T) {
	t.Parallel()

	html := renderSampleHTML(t)

	if !strings.Contains(html, `id="goslide-stage"`) {
		t.Errorf("missing goslide-stage container")
	}
	if !strings.Contains(html, `id="goslide-app"`) {
		t.Errorf("missing goslide-app mount point for svelte runtime")
	}
	if !strings.Contains(html, `id="goslide-runtime-script"`) {
		t.Errorf("missing goslide-runtime-script tag")
	}
	if !strings.Contains(html, "goslide-deck") {
		t.Errorf("master template missing goslide-deck container")
	}
}

func TestHTMLRenderer_Render_Themes(t *testing.T) {
	t.Parallel()

	deck := sampleDeck()

	themes := []string{"default", "clean", "dark"}
	for _, th := range themes {
		t.Run("theme_"+th, func(t *testing.T) {
			t.Parallel()
			renderer := htmlrenderer.NewRenderer(htmlrenderer.WithTheme(th))
			var buf bytes.Buffer
			err := renderer.Render(context.Background(), deck, &buf)
			if err != nil {
				t.Fatalf("failed to render theme %s: %v", th, err)
			}
			if !strings.Contains(buf.String(), "Theme: "+th) {
				t.Errorf("rendered output missing header comment for theme %s", th)
			}
		})
	}
}

func TestHTMLRenderer_Render_ContextCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel before call

	renderer := htmlrenderer.NewRenderer()
	deck := sampleDeck()

	var buf bytes.Buffer
	err := renderer.Render(ctx, deck, &buf)
	if err == nil {
		t.Fatalf("expected error for canceled context, got nil")
	}
	if !errors.Is(err, model.ErrCanceled) {
		t.Errorf("expected ErrCanceled, got %v", err)
	}
}

func TestHTMLRenderer_Render_NilDeck(t *testing.T) {
	t.Parallel()

	renderer := htmlrenderer.NewRenderer()
	var buf bytes.Buffer
	err := renderer.Render(context.Background(), nil, &buf)
	if err == nil {
		t.Fatalf("expected error for nil deck, got nil")
	}
}

func TestHTMLRenderer_Render_CustomCSS(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	customCSSFile := filepath.Join(tmpDir, "custom.css")
	if err := os.WriteFile(customCSSFile, []byte(".custom-rule { color: #f00; }"), 0600); err != nil {
		t.Fatalf("failed to write custom css: %v", err)
	}

	renderer := htmlrenderer.NewRenderer(htmlrenderer.WithCustomCSS(customCSSFile))
	deck := sampleDeck()

	var buf bytes.Buffer
	err := renderer.Render(context.Background(), deck, &buf)
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}

	if !strings.Contains(buf.String(), ".custom-rule { color: #f00; }") {
		t.Errorf("custom css rule not found in rendered output")
	}
}

func TestHTMLRenderer_Render_BackgroundDim(t *testing.T) {
	t.Parallel()

	deck := &model.Deck{
		Title: "Dim Test",
		Slides: []*model.Slide{
			{
				Index:  1,
				Layout: model.LayoutDefault,
				Directives: model.SlideDirectives{
					BackgroundImage: "bg.jpg",
					BackgroundDim:   "0.6",
				},
				HTMLContent: "<p>Dimmed Slide</p>",
			},
			{
				Index:  2,
				Layout: model.LayoutDefault,
				Directives: model.SlideDirectives{
					Class: []string{"dim"},
				},
				HTMLContent: "<p>Dim Class Slide</p>",
			},
		},
	}

	renderer := htmlrenderer.NewRenderer()
	var buf bytes.Buffer
	if err := renderer.Render(context.Background(), deck, &buf); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `has-bg-dim`) {
		t.Errorf("expected has-bg-dim class in output")
	}
	if !strings.Contains(html, `background-color: rgba(0, 0, 0, 0.6);`) {
		t.Errorf("expected rgba(0, 0, 0, 0.6) in output")
	}
	if !strings.Contains(html, `background-color: rgba(0, 0, 0, 0.5);`) {
		t.Errorf("expected default dim rgba(0, 0, 0, 0.5) for dim class in output")
	}
}

func TestHTMLRenderer_Render_BackgroundImage_Standalone(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testImagePath := filepath.Join(tmpDir, "test-bg.png")
	pngBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
		0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
	if err := os.WriteFile(testImagePath, pngBytes, 0600); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}

	deck := &model.Deck{
		Title: "Standalone Background Deck",
		Slides: []*model.Slide{
			{
				Index:  1,
				Layout: model.LayoutDefault,
				Directives: model.SlideDirectives{
					BackgroundImage: "test-bg.png",
				},
				HTMLContent: "<p>Slide 1</p>",
			},
			{
				Index:  2,
				Layout: model.LayoutDefault,
				Directives: model.SlideDirectives{
					BackgroundImage: "url('test-bg.png')",
				},
				HTMLContent: "<p>Slide 2</p>",
			},
			{
				Index:  3,
				Layout: model.LayoutDefault,
				Directives: model.SlideDirectives{
					BackgroundImage: "linear-gradient(to right, #111, #222)",
				},
				HTMLContent: "<p>Slide 3</p>",
			},
			{
				Index:  4,
				Layout: model.LayoutDefault,
				Directives: model.SlideDirectives{
					BackgroundImage: "my'pic.png",
				},
				HTMLContent: "<p>Slide 4</p>",
			},
		},
	}

	renderer := htmlrenderer.NewRenderer(
		htmlrenderer.WithStandalone(true),
		htmlrenderer.WithBaseDir(tmpDir),
	)
	var buf bytes.Buffer
	if err := renderer.Render(context.Background(), deck, &buf); err != nil {
		t.Fatalf("standalone render failed: %v", err)
	}

	output := buf.String()

	// Must NOT contain #ZgotmplZ
	if strings.Contains(output, "#ZgotmplZ") {
		t.Errorf("rendered output contains #ZgotmplZ sanitization failure")
	}

	// Slide 1 & 2 must contain base64 data uri
	if !strings.Contains(output, "data:image/png;base64,") {
		t.Errorf("expected base64 data URI in background-image, got:\n%s", output)
	}

	// Slide 3 must contain gradient without url()
	if !strings.Contains(output, "background-image: linear-gradient(to right, #111, #222);") {
		t.Errorf("expected linear gradient in background-image")
	}

	// Slide 4 must escape single quotes as %27
	if !strings.Contains(output, "my%27pic.png") {
		t.Errorf("expected single quote to be escaped as %%27 in background-image")
	}
}

