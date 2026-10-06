package parser

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
)

func TestParseCodeInfo(t *testing.T) {
	tests := []struct {
		info       string
		wantLang   string
		wantRanges [][2]int
	}{
		{"go", "go", nil},
		{"go {2,4-6}", "go", [][2]int{{2, 2}, {4, 6}}},
		{"python { 1-3, 5 }", "python", [][2]int{{1, 3}, {5, 5}}},
		{"rust {10}", "rust", [][2]int{{10, 10}}},
		{"go {invalid}", "go", nil},
		{"", "", nil},
	}

	for _, tt := range tests {
		lang, ranges := parseCodeInfo(tt.info)
		if lang != tt.wantLang {
			t.Errorf("parseCodeInfo(%q) lang = %q, want %q", tt.info, lang, tt.wantLang)
		}
		if len(ranges) != len(tt.wantRanges) {
			t.Errorf("parseCodeInfo(%q) ranges len = %d, want %d", tt.info, len(ranges), len(tt.wantRanges))
			continue
		}
		for i := range ranges {
			if ranges[i] != tt.wantRanges[i] {
				t.Errorf("parseCodeInfo(%q) range[%d] = %v, want %v", tt.info, i, ranges[i], tt.wantRanges[i])
			}
		}
	}
}

func TestFencedCodeBlock_LineHighlighting(t *testing.T) {
	gm := goldmark.New(
		goldmark.WithExtensions(newChromaHighlightExtension("dracula")),
	)

	input := "```go {2,4}\npackage main\nimport \"fmt\"\nfunc main() {\n\tfmt.Println(\"hi\")\n}\n```"
	var buf bytes.Buffer
	if err := gm.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	html := buf.String()

	// 1. Should have highlight-container has-highlights in <pre>
	if !strings.Contains(html, "highlight-container has-highlights") {
		t.Errorf("expected <pre> to contain 'highlight-container has-highlights', got:\n%s", html)
	}

	// 2. Should contain line hl class for highlighted lines
	if !strings.Contains(html, "class=\"line hl\"") {
		t.Errorf("expected output to contain 'class=\"line hl\"', got:\n%s", html)
	}

	// 3. Should contain line class for standard lines
	if !strings.Contains(html, "class=\"line\"") {
		t.Errorf("expected output to contain 'class=\"line\"', got:\n%s", html)
	}
}

func TestFencedCodeBlock_NoHighlights(t *testing.T) {
	gm := goldmark.New(
		goldmark.WithExtensions(newChromaHighlightExtension("dracula")),
	)

	input := "```go\npackage main\n```"
	var buf bytes.Buffer
	if err := gm.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	html := buf.String()

	if !strings.Contains(html, "highlight-container") {
		t.Errorf("expected <pre> to contain 'highlight-container', got:\n%s", html)
	}
	if strings.Contains(html, "has-highlights") {
		t.Errorf("did not expect 'has-highlights', got:\n%s", html)
	}
}
