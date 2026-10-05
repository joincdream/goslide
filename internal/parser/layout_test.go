package parser

import (
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yundream/goslide/internal/model"
)

func TestRenderSlideContent_TwoCols_WithSplit(t *testing.T) {
	gm := goldmark.New()
	input := "## Two Cols Slide\n\n### Left Side\n- Item L\n\n<!-- split -->\n\n### Right Side\n- Item R"
	title, html, left, right, err := renderSlideContent(gm, model.LayoutTwoCols, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(title, "Two Cols Slide</h2>") {
		t.Errorf("expected title in title return, got: %s", title)
	}
	if !strings.Contains(html, "<div class=\"two-cols\">") {
		t.Errorf("expected two-cols wrapper in html: %s", html)
	}
	if !strings.Contains(left, "Left Side</h3>") || !strings.Contains(left, "<li>Item L</li>") {
		t.Errorf("left column missing content: %s", left)
	}
	if !strings.Contains(right, "Right Side</h3>") || !strings.Contains(right, "<li>Item R</li>") {
		t.Errorf("right column missing content: %s", right)
	}
}

func TestRenderSlideContent_TwoCols_WithoutSplit(t *testing.T) {
	gm := goldmark.New()
	input := "### Single Column Content"
	title, html, left, right, err := renderSlideContent(gm, model.LayoutTwoCols, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if title != "" {
		t.Errorf("expected empty title for h3, got %q", title)
	}
	if strings.Contains(html, "<div class=\"two-cols\">") {
		t.Errorf("did not expect two-cols wrapper when split marker is missing: %s", html)
	}
	if left != "" || right != "" {
		t.Errorf("expected empty left and right, got left=%q, right=%q", left, right)
	}
}

func TestRenderSlideContent_DefaultLayout(t *testing.T) {
	gm := goldmark.New()
	input := "## Slide Header\n\nLeft\n<!-- split -->\nRight"
	title, html, left, right, err := renderSlideContent(gm, model.LayoutDefault, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(title, "Slide Header</h2>") {
		t.Errorf("expected title to be extracted: %s", title)
	}
	if strings.Contains(html, "<div class=\"two-cols\">") {
		t.Errorf("default layout should not create two-cols container: %s", html)
	}
	if left != "" || right != "" {
		t.Errorf("expected empty left and right for default layout")
	}
}
