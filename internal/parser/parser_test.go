package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

func TestParser_Parse_Success(t *testing.T) {
	p := NewParser()
	input := `---
title: "Goslide Presentation"
author: "Goslide Team"
theme: "clean"
size: "16:9"
paginate: true
---
# First Slide

Welcome to Goslide!

---

# Second Slide

- Item 1
- Item 2
`
	deck, err := p.Parse(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertDeckMetadata(t, deck)
	assertDeckSlides(t, deck.Slides)
}

func assertDeckMetadata(t *testing.T, deck *model.Deck) {
	t.Helper()
	if deck.Title != "Goslide Presentation" {
		t.Errorf("expected title 'Goslide Presentation', got %q", deck.Title)
	}
	if deck.Author != "Goslide Team" {
		t.Errorf("expected author 'Goslide Team', got %q", deck.Author)
	}
	if deck.GlobalAttrs.Theme != "clean" {
		t.Errorf("expected theme 'clean', got %q", deck.GlobalAttrs.Theme)
	}
	if deck.GlobalAttrs.Size != model.Ratio16x9 {
		t.Errorf("expected size '16:9', got %q", deck.GlobalAttrs.Size)
	}
	if !deck.GlobalAttrs.Paginate {
		t.Errorf("expected paginate true, got false")
	}
}

func assertDeckSlides(t *testing.T, slides []*model.Slide) {
	t.Helper()
	if len(slides) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(slides))
	}

	s1 := slides[0]
	if s1.Index != 1 || s1.Layout != model.LayoutDefault || s1.Notes != "" {
		t.Errorf("slide 1 fields mismatch: index=%d, layout=%v, notes=%q", s1.Index, s1.Layout, s1.Notes)
	}
	if !strings.Contains(s1.HTMLContent, "First Slide</h1>") {
		t.Errorf("slide 1 html missing h1: %s", s1.HTMLContent)
	}

	s2 := slides[1]
	if s2.Index != 2 || !strings.Contains(s2.HTMLContent, "<li>Item 1</li>") {
		t.Errorf("slide 2 html missing list items: %s", s2.HTMLContent)
	}
}

func TestParser_Parse_GFM(t *testing.T) {
	p := NewParser()

	t.Run("GFM table rendering", func(t *testing.T) {
		input := "| Header 1 | Header 2 |\n| :--- | :--- |\n| Val 1 | Val 2 |"
		deck, err := p.Parse(context.Background(), strings.NewReader(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		html := deck.Slides[0].HTMLContent
		if !strings.Contains(html, "<table>") || !strings.Contains(html, "Header 1</th>") {
			t.Errorf("expected GFM table tag in HTML, got:\n%s", html)
		}
	})

	t.Run("GFM strikethrough and task list", func(t *testing.T) {
		input := "- [x] Done task\n- [ ] Open task\n- ~~Deleted text~~"
		deck, err := p.Parse(context.Background(), strings.NewReader(input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		html := deck.Slides[0].HTMLContent
		if !strings.Contains(html, "<del>Deleted text</del>") || !strings.Contains(html, `type="checkbox"`) {
			t.Errorf("expected strikethrough and checkbox in HTML, got:\n%s", html)
		}
	})
}

func TestParser_Parse_CRLF(t *testing.T) {
	p := NewParser()
	input := "---\r\ntitle: Windows CRLF\r\n---\r\n# Slide 1\r\n---\r\n# Slide 2\r\n"
	deck, err := p.Parse(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error with CRLF: %v", err)
	}
	if deck.Title != "Windows CRLF" || len(deck.Slides) != 2 {
		t.Errorf("CRLF deck mismatch: title=%q, slides=%d", deck.Title, len(deck.Slides))
	}
}

func TestParser_Parse_Errors(t *testing.T) {
	p := NewParser()

	t.Run("context cancellation returns ErrCanceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := p.Parse(ctx, strings.NewReader("# Slide"))
		if err == nil || !errors.Is(err, model.ErrCanceled) {
			t.Fatalf("expected error wrapping ErrCanceled, got: %v", err)
		}
	})

	t.Run("nil reader returns error", func(t *testing.T) {
		if _, err := p.Parse(context.Background(), nil); err == nil {
			t.Fatal("expected error for nil reader, got nil")
		}
	})

	t.Run("invalid frontmatter returns ErrInvalidFrontmatter", func(t *testing.T) {
		input := "---\ntitle: \"broken\n  [bad yaml]\n---\n# Slide"
		_, err := p.Parse(context.Background(), strings.NewReader(input))
		if err == nil || !errors.Is(err, model.ErrInvalidFrontmatter) {
			t.Fatalf("expected ErrInvalidFrontmatter, got: %v", err)
		}
	})
}

func TestParser_Parse_DirectivesAndHighlighting(t *testing.T) {
	p := NewParser()
	input := `---
title: "Directives Demo"
---
<!-- _layout: cover -->
<!-- _backgroundColor: #0f172a -->
<!-- _color: #ffffff -->
<!-- note: Cover slide speech -->
# Title

---
<!-- _layout: two-cols -->
Col Left
<!-- split -->
Col Right

---
` + "```go\nfunc add(a, b int) int { return a + b }\n```"

	deck, err := p.Parse(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(deck.Slides) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(deck.Slides))
	}

	// Slide 1 checks
	s1 := deck.Slides[0]
	if s1.Layout != model.LayoutCover {
		t.Errorf("s1 layout mismatch: got %v, want cover", s1.Layout)
	}
	if s1.Directives.BackgroundColor != "#0f172a" || s1.Directives.Color != "#ffffff" {
		t.Errorf("s1 directives mismatch: bg=%q, color=%q", s1.Directives.BackgroundColor, s1.Directives.Color)
	}
	if s1.Notes != "Cover slide speech" {
		t.Errorf("s1 notes mismatch: got %q, want 'Cover slide speech'", s1.Notes)
	}

	// Slide 2 checks
	s2 := deck.Slides[1]
	if s2.Layout != model.LayoutTwoCols {
		t.Errorf("s2 layout mismatch: got %v, want two-cols", s2.Layout)
	}
	if !strings.Contains(s2.HTMLContent, "<div class=\"two-cols\">") {
		t.Errorf("s2 missing two-cols container: %s", s2.HTMLContent)
	}

	// Slide 3 checks (Chroma)
	s3 := deck.Slides[2]
	if !strings.Contains(s3.HTMLContent, "<span style=") || !strings.Contains(s3.HTMLContent, "func") {
		t.Errorf("s3 missing Chroma highlighted span: %s", s3.HTMLContent)
	}
}
