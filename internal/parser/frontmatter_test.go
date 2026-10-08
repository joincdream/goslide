package parser

import (
	"errors"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

type frontmatterTestCase struct {
	name            string
	content         string
	expected        *FrontmatterResult
	wantErrSentinel error
}

func runFrontmatterTests(t *testing.T, tests []frontmatterTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := extractFrontmatter(tt.content)
			if tt.wantErrSentinel != nil {
				if !errors.Is(err, tt.wantErrSentinel) {
					t.Fatalf("expected error wrapping %v, got: %v", tt.wantErrSentinel, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertFrontmatterFields(t, res, tt.expected)
		})
	}
}

func TestExtractFrontmatter_Valid(t *testing.T) {
	tests := []frontmatterTestCase{
		{
			name: "full frontmatter with all valid fields",
			content: `---
title: "Goslide Arch"
author: "Alice"
theme: "clean"
layout: "cover"
size: "4:3"
paginate: true
header: "Header text"
footer: "Footer text"
custom_css: "h1 { color: blue; }"
---
# Main Content`,
			expected: &FrontmatterResult{
				Title:     "Goslide Arch",
				Author:    "Alice",
				CustomCSS: "h1 { color: blue; }",
				Body:      "# Main Content",
				GlobalAttrs: model.GlobalDirectives{
					Theme:    "clean",
					Layout:   model.LayoutCover,
					Size:     model.Ratio4x3,
					Paginate: true,
					Header:   "Header text",
					Footer:   "Footer text",
				},
			},
		},
		{
			name: "marp compatibility with style and marp flag",
			content: `---
marp: true
theme: dark
size: 16:9
style: "section { padding: 20px; }"
---
# Marp Slide`,
			expected: &FrontmatterResult{
				CustomCSS: "section { padding: 20px; }",
				Body:      "# Marp Slide",
				GlobalAttrs: model.GlobalDirectives{
					Theme:  "dark",
					Layout: model.LayoutDefault,
					Size:   model.Ratio16x9,
				},
			},
		},
		{
			name: "external theme path preserved",
			content: `---
theme: "themes/corporate.css"
---
# External Theme Slide`,
			expected: &FrontmatterResult{
				Body: "# External Theme Slide",
				GlobalAttrs: model.GlobalDirectives{
					Theme:  "themes/corporate.css",
					Layout: model.LayoutDefault,
					Size:   model.Ratio16x9,
				},
			},
		},
		{
			name: "external theme relative path without css extension preserved",
			content: `---
theme: "themes/corporate"
---
# External Theme Relative Slide`,
			expected: &FrontmatterResult{
				Body: "# External Theme Relative Slide",
				GlobalAttrs: model.GlobalDirectives{
					Theme:  "themes/corporate",
					Layout: model.LayoutDefault,
					Size:   model.Ratio16x9,
				},
			},
		},
	}
	runFrontmatterTests(t, tests)
}

func TestExtractFrontmatter_Fallbacks(t *testing.T) {
	tests := []frontmatterTestCase{
		{
			name: "unsupported theme and size fallback to default",
			content: `---
theme: "super-fancy-unknown-theme"
size: "21:9"
layout: "unknown-layout"
---
# Content`,
			expected: &FrontmatterResult{
				Body: "# Content",
				GlobalAttrs: model.GlobalDirectives{
					Theme:  "default",
					Layout: model.LayoutDefault,
					Size:   model.Ratio16x9,
				},
			},
		},
		{
			name:    "no frontmatter present",
			content: "# Plain Markdown\nJust some text.",
			expected: &FrontmatterResult{
				Body: "# Plain Markdown\nJust some text.",
				GlobalAttrs: model.GlobalDirectives{
					Theme:  "default",
					Layout: model.LayoutDefault,
					Size:   model.Ratio16x9,
				},
			},
		},
		{
			name: "unclosed frontmatter treated as body",
			content: `---
title: Unclosed
# Body starts here without closing delimiter`,
			expected: &FrontmatterResult{
				Body: "---\ntitle: Unclosed\n# Body starts here without closing delimiter",
				GlobalAttrs: model.GlobalDirectives{
					Theme:  "default",
					Layout: model.LayoutDefault,
					Size:   model.Ratio16x9,
				},
			},
		},
		{
			name: "empty frontmatter",
			content: `---
---
# Body`,
			expected: &FrontmatterResult{
				Body: "# Body",
				GlobalAttrs: model.GlobalDirectives{
					Theme:  "default",
					Layout: model.LayoutDefault,
					Size:   model.Ratio16x9,
				},
			},
		},
		{
			name: "invalid yaml syntax returns ErrInvalidFrontmatter",
			content: `---
title: "Unclosed string
  bad: indentation: [
---
# Body`,
			wantErrSentinel: model.ErrInvalidFrontmatter,
		},
	}
	runFrontmatterTests(t, tests)
}

func assertFrontmatterFields(t *testing.T, res, exp *FrontmatterResult) {
	t.Helper()
	if res.Title != exp.Title {
		t.Errorf("title mismatch: got %q, want %q", res.Title, exp.Title)
	}
	if res.Author != exp.Author {
		t.Errorf("author mismatch: got %q, want %q", res.Author, exp.Author)
	}
	if res.GlobalAttrs.Theme != exp.GlobalAttrs.Theme {
		t.Errorf("theme mismatch: got %q, want %q", res.GlobalAttrs.Theme, exp.GlobalAttrs.Theme)
	}
	if res.GlobalAttrs.Size != exp.GlobalAttrs.Size {
		t.Errorf("size mismatch: got %q, want %q", res.GlobalAttrs.Size, exp.GlobalAttrs.Size)
	}
	if res.GlobalAttrs.Layout != exp.GlobalAttrs.Layout {
		t.Errorf("layout mismatch: got %q, want %q", res.GlobalAttrs.Layout, exp.GlobalAttrs.Layout)
	}
	if res.GlobalAttrs.Paginate != exp.GlobalAttrs.Paginate {
		t.Errorf("paginate mismatch: got %v, want %v", res.GlobalAttrs.Paginate, exp.GlobalAttrs.Paginate)
	}
	if res.CustomCSS != exp.CustomCSS {
		t.Errorf("custom_css mismatch: got %q, want %q", res.CustomCSS, exp.CustomCSS)
	}
	if res.Body != exp.Body {
		t.Errorf("body mismatch:\ngot:  %q\nwant: %q", res.Body, exp.Body)
	}
}
