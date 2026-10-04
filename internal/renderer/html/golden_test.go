package html_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yundream/goslide/internal/parser"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
	"github.com/yundream/goslide/internal/testutil"
)

func TestGoldenRenderer(t *testing.T) {
	testCases := []struct {
		name       string
		sourceFile string
		goldenFile string
	}{
		{
			name:       "basic slides",
			sourceFile: "basic.md",
			goldenFile: "basic.html",
		},
		{
			name:       "two columns layout",
			sourceFile: "two_cols.md",
			goldenFile: "two_cols.html",
		},
		{
			name:       "syntax highlighting",
			sourceFile: "highlight.md",
			goldenFile: "highlight.html",
		},
		{
			name:       "cover layout",
			sourceFile: "cover.md",
			goldenFile: "cover.html",
		},
		{
			name:       "comprehensive elements",
			sourceFile: "comprehensive_elements.md",
			goldenFile: "comprehensive_elements.html",
		},
	}

	p := parser.NewParser()
	r := htmlrenderer.NewRenderer()
	ctx := context.Background()

	slidesDir := filepath.Join("..", "..", "..", "testdata", "slides")
	goldenDir := filepath.Join("..", "..", "..", "testdata", "golden")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srcPath := filepath.Join(slidesDir, tc.sourceFile)
			srcBytes, err := os.ReadFile(srcPath)
			if err != nil {
				t.Fatalf("failed to read test fixture %s: %v", srcPath, err)
			}

			deck, err := p.Parse(ctx, bytes.NewReader(srcBytes))
			if err != nil {
				t.Fatalf("failed to parse markdown: %v", err)
			}

			var buf bytes.Buffer
			if err := r.Render(ctx, deck, &buf); err != nil {
				t.Fatalf("failed to render HTML: %v", err)
			}

			goldenPath := filepath.Join(goldenDir, tc.goldenFile)
			testutil.AssertGolden(t, buf.Bytes(), goldenPath)
		})
	}
}
