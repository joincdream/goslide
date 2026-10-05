package parser

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
)

func TestChromaHighlightExtension(t *testing.T) {
	gm := goldmark.New(
		goldmark.WithExtensions(newChromaHighlightExtension("dracula")),
	)

	t.Run("highlight Go fenced code block", func(t *testing.T) {
		input := "```go\npackage main\nfunc main() {}\n```"
		var buf bytes.Buffer
		if err := gm.Convert([]byte(input), &buf); err != nil {
			t.Fatalf("conversion failed: %v", err)
		}

		html := buf.String()
		// Chroma wraps code with pre/code and spans with inline color styles
		if !strings.Contains(html, "<pre") || !strings.Contains(html, "<span style=") {
			t.Errorf("expected syntax highlighted spans with styles in output: %s", html)
		}
		if !strings.Contains(html, "package") || !strings.Contains(html, "func") {
			t.Errorf("missing code tokens in output: %s", html)
		}
	})

	t.Run("highlight Python fenced code block", func(t *testing.T) {
		input := "```python\ndef hello():\n    return 'world'\n```"
		var buf bytes.Buffer
		if err := gm.Convert([]byte(input), &buf); err != nil {
			t.Fatalf("conversion failed: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "<span style=") || !strings.Contains(html, "def") {
			t.Errorf("missing python highlighting in output: %s", html)
		}
	})

	t.Run("code block with unknown language falls back safely", func(t *testing.T) {
		input := "```unknownlang\nsome raw code\n```"
		var buf bytes.Buffer
		if err := gm.Convert([]byte(input), &buf); err != nil {
			t.Fatalf("conversion failed: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "some raw code") {
			t.Errorf("missing code content in fallback output: %s", html)
		}
	})
}

func TestChromaHighlightExtension_Mermaid(t *testing.T) {
	gm := goldmark.New(
		goldmark.WithExtensions(newChromaHighlightExtension("dracula")),
	)

	input := "```mermaid\nflowchart LR\nA --> B\n```"
	var buf bytes.Buffer
	if err := gm.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "<pre class=\"mermaid\">") {
		t.Errorf("expected <pre class=\"mermaid\"> in output: %s", html)
	}
	if !strings.Contains(html, "flowchart LR") || !strings.Contains(html, "A --&gt; B") {
		t.Errorf("missing escaped mermaid content in output: %s", html)
	}
}
