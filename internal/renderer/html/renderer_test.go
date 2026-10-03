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

	if !strings.Contains(html, `id="goslide-canvas"`) {
		t.Errorf("missing goslide-canvas overlay")
	}
	if !strings.Contains(html, `id="goslide-indicator"`) {
		t.Errorf("missing goslide-indicator")
	}
	if !strings.Contains(html, `1 / 3`) {
		t.Errorf("indicator should display initial page '1 / 3'")
	}
	if !strings.Contains(html, `id="goslide-runtime-script"`) {
		t.Errorf("missing goslide-runtime-script tag")
	}
	if !strings.Contains(html, "goslide-deck") || !strings.Contains(html, "toggleDrawMode") {
		t.Errorf("master template missing core js runtime logic")
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
