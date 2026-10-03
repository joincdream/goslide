package main

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/parser"
)

var outputPathFlag string

func init() {
	buildCmd.Flags().StringVarP(&outputPathFlag, "output", "o", "", "Output HTML file path (default: <input>.html)")
	buildCmd.RunE = runBuild
}

func runBuild(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("input markdown file is required: goslide build <input.md>")
	}

	inputPath := args[0]
	outputPath := resolveOutputPath(inputPath, outputPathFlag)

	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file %q: %w", inputPath, err)
	}
	defer func() { _ = file.Close() }()

	p := parser.NewParser()
	deck, err := p.Parse(cmd.Context(), file)
	if err != nil {
		return fmt.Errorf("failed to parse markdown: %w", err)
	}

	if err := renderPreviewHTML(cmd.Context(), deck, outputPath); err != nil {
		return fmt.Errorf("failed to write output HTML: %w", err)
	}

	fmt.Printf("✓ Successfully built %d slides to %s\n", len(deck.Slides), outputPath)
	return nil
}

func resolveOutputPath(inputPath, outputFlag string) string {
	if outputFlag != "" {
		return outputFlag
	}
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	return base + ".html"
}

type previewTemplateData struct {
	Title     string
	CustomCSS template.CSS
	Slides    []slideViewData
}

type slideViewData struct {
	Index       int
	Layout      string
	Classes     string
	BgColor     string
	Color       string
	Header      string
	Footer      string
	Paginate    bool
	HTMLContent template.HTML
}

func renderPreviewHTML(_ context.Context, deck *model.Deck, outputPath string) error {
	data := buildPreviewData(deck)

	tmpl, err := template.New("preview").Parse(previewHTMLTemplate)
	if err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	return tmpl.Execute(out, data)
}

func buildPreviewData(deck *model.Deck) previewTemplateData {
	title := deck.Title
	if title == "" {
		title = "Goslide Preview"
	}

	views := make([]slideViewData, 0, len(deck.Slides))
	for _, s := range deck.Slides {
		views = append(views, slideViewData{
			Index:       s.Index,
			Layout:      string(s.Layout),
			Classes:     strings.Join(s.Directives.Class, " "),
			BgColor:     s.Directives.BackgroundColor,
			Color:       s.Directives.Color,
			Header:      s.Directives.Header,
			Footer:      s.Directives.Footer,
			Paginate:    s.Directives.Paginate,
			HTMLContent: template.HTML(s.HTMLContent), // nolint:gosec
		})
	}

	return previewTemplateData{
		Title:     title,
		CustomCSS: template.CSS(deck.CustomCSS), // nolint:gosec
		Slides:    views,
	}
}

const previewHTMLTemplate = `<!DOCTYPE html>
<html lang="ko">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .Title }}</title>
  <style>
    * { box-sizing: border-box; }
    body {
      margin: 0;
      padding: 2rem;
      background: #181825;
      color: #cdd6f4;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Pretendard", sans-serif;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 2.5rem;
    }
    .slide-wrapper {
      position: relative;
    }
    .slide-badge {
      position: absolute;
      top: -1.5rem;
      left: 0.5rem;
      font-size: 0.75rem;
      color: #a6adc8;
      font-family: monospace;
    }
    .slide-card {
      width: 960px;
      height: 540px;
      background: #ffffff;
      color: #1e1e2e;
      box-shadow: 0 12px 36px rgba(0, 0, 0, 0.45);
      border-radius: 8px;
      overflow: hidden;
      position: relative;
      padding: 3rem 4rem;
      display: flex;
      flex-direction: column;
    }
    .slide-card.cover {
      justify-content: center;
      align-items: center;
      text-align: center;
    }
    .slide-header {
      font-size: 0.85rem;
      color: #888;
      margin-bottom: 1rem;
    }
    .slide-body {
      flex: 1;
      overflow: hidden;
    }
    .slide-footer {
      display: flex;
      justify-content: space-between;
      font-size: 0.8rem;
      color: #888;
      margin-top: 1rem;
    }
    .two-cols {
      display: flex;
      gap: 2rem;
      height: 100%;
    }
    .col-left, .col-right {
      flex: 1;
    }
    pre {
      border-radius: 6px;
      padding: 1rem;
      overflow-x: auto;
      font-size: 0.9rem;
    }
    table {
      border-collapse: collapse;
      width: 100%;
      margin: 1rem 0;
    }
    th, td {
      border: 1px solid #ddd;
      padding: 8px 12px;
    }
    th {
      background-color: #f2f2f2;
    }
    {{ .CustomCSS }}
  </style>
</head>
<body>
  {{ range .Slides }}
  <div class="slide-wrapper">
    <div class="slide-badge">Slide {{ .Index }} [{{ .Layout }}]</div>
    <div class="slide-card {{ .Layout }} {{ .Classes }}"
         style="{{ if .BgColor }}background-color: {{ .BgColor }};{{ end }}{{ if .Color }}color: {{ .Color }};{{ end }}">
      {{ if .Header }}<div class="slide-header">{{ .Header }}</div>{{ end }}
      <div class="slide-body">
        {{ .HTMLContent }}
      </div>
      <div class="slide-footer">
        <span>{{ .Footer }}</span>
        {{ if .Paginate }}<span>{{ .Index }}</span>{{ end }}
      </div>
    </div>
  </div>
  {{ end }}
</body>
</html>
`
