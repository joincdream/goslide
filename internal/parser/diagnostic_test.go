package parser

import (
	"context"
	"strings"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

func TestInspectSlideComments_UnknownLayout(t *testing.T) {
	content := "# Title\n<!-- _layout: two-col -->\nBody text"
	diags := inspectSlideComments(1, content)

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}

	d := diags[0]
	if d.Rule != "layout.unknown" {
		t.Errorf("expected rule 'layout.unknown', got %q", d.Rule)
	}
	if d.Line != 2 {
		t.Errorf("expected line 2, got %d", d.Line)
	}
	if len(d.Candidates) == 0 {
		t.Errorf("expected non-empty candidates for unknown layout")
	}

	// Verify that canonical layouts are present in candidates
	hasTwoCols := false
	for _, c := range d.Candidates {
		if c == string(model.LayoutTwoCols) {
			hasTwoCols = true
			break
		}
	}
	if !hasTwoCols {
		t.Errorf("expected candidate list to include %q, got %v", model.LayoutTwoCols, d.Candidates)
	}
}

func TestInspectSlideComments_MissingUnderscore(t *testing.T) {
	content := "<!-- layout: two-cols -->\n# Slide Content"
	diags := inspectSlideComments(2, content)

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}

	d := diags[0]
	if d.Rule != "directive.missing_underscore" {
		t.Errorf("expected rule 'directive.missing_underscore', got %q", d.Rule)
	}
	if d.SlideIndex != 2 {
		t.Errorf("expected slideIndex 2, got %d", d.SlideIndex)
	}
	if d.Line != 1 {
		t.Errorf("expected line 1, got %d", d.Line)
	}
	if len(d.Candidates) != 1 || d.Candidates[0] != "_layout" {
		t.Errorf("expected candidate ['_layout'], got %v", d.Candidates)
	}
}

func TestInspectSlideComments_CodeBlockIsolation(t *testing.T) {
	content := "# Code Example\n\n```html\n<!-- layout: two-cols -->\n<!-- _layout: two-col -->\n```\n\nSome description"
	diags := inspectSlideComments(1, content)

	if len(diags) != 0 {
		t.Errorf("expected 0 diagnostics inside code block, got %d: %+v", len(diags), diags)
	}
}

func TestInspectSlideComments_MalformedDelimiter(t *testing.T) {
	content := "# Slide\n<!-- _layout: two-cols --->\nContent"
	diags := inspectSlideComments(1, content)

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic for malformed delimiter, got %d", len(diags))
	}
	if diags[0].Rule != "syntax.malformed_delimiter" {
		t.Errorf("expected rule 'syntax.malformed_delimiter', got %q", diags[0].Rule)
	}
}

func TestInspectSlideComments_GeneralUserComments(t *testing.T) {
	content := "# Meeting Notes\n<!-- memo: remember to send follow-up email -->\n<!-- todo: add roadmap diagram -->\n<!-- section note -->"
	diags := inspectSlideComments(1, content)

	if len(diags) != 0 {
		t.Errorf("expected 0 diagnostics for general comments, got %d: %+v", len(diags), diags)
	}
}

func TestInspectSlideComments_MarkerTypo(t *testing.T) {
	content := "# Column Left\n<!-- splti -->\n# Column Right"
	diags := inspectSlideComments(1, content)

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic for marker typo, got %d", len(diags))
	}
	if diags[0].Rule != "marker.unknown" {
		t.Errorf("expected rule 'marker.unknown', got %q", diags[0].Rule)
	}
	if len(diags[0].Candidates) != 1 || diags[0].Candidates[0] != "split" {
		t.Errorf("expected candidate ['split'], got %v", diags[0].Candidates)
	}
}

func TestParser_DeckDiagnosticsIntegration(t *testing.T) {
	input := `---
title: Sample Deck
---
# Slide 1
<!-- layout: two-cols -->

---
# Slide 2
<!-- _layout: two-col -->
`
	p := NewParser()
	deck, err := p.Parse(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if len(deck.Diagnostics) != 2 {
		t.Fatalf("expected 2 deck diagnostics, got %d: %+v", len(deck.Diagnostics), deck.Diagnostics)
	}

	// Verify Slide 1 diagnostic
	if len(deck.Slides[0].Diagnostics) != 1 {
		t.Errorf("expected 1 diagnostic on slide 1, got %d", len(deck.Slides[0].Diagnostics))
	}
	if deck.Slides[0].Diagnostics[0].Rule != "directive.missing_underscore" {
		t.Errorf("expected 'directive.missing_underscore' on slide 1, got %q", deck.Slides[0].Diagnostics[0].Rule)
	}

	// Verify Slide 2 diagnostic
	if len(deck.Slides[1].Diagnostics) != 1 {
		t.Errorf("expected 1 diagnostic on slide 2, got %d", len(deck.Slides[1].Diagnostics))
	}
	if deck.Slides[1].Diagnostics[0].Rule != "layout.unknown" {
		t.Errorf("expected 'layout.unknown' on slide 2, got %q", deck.Slides[1].Diagnostics[0].Rule)
	}
}
