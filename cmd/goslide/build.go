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
	"github.com/yundream/goslide/internal/theme"
)

var (
	outputPathFlag string
	themeFlag      string
	themePathFlag  string
)

func init() {
	buildCmd.Flags().StringVarP(&outputPathFlag, "output", "o", "", "Output HTML file path (default: <input>.html)")
	buildCmd.Flags().StringVarP(&themeFlag, "theme", "t", "", "Theme name (default, clean, dark; overrides frontmatter)")
	buildCmd.Flags().StringVar(&themePathFlag, "theme-path", "", "Path to custom external CSS stylesheet")
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

	chosenTheme := resolveThemeName(themeFlag, deck.GlobalAttrs.Theme)
	if err := renderPreviewHTML(cmd.Context(), deck, outputPath, chosenTheme, themePathFlag); err != nil {
		return fmt.Errorf("failed to write output HTML: %w", err)
	}

	fmt.Printf("✓ Successfully built %d slides to %s [theme: %s]\n", len(deck.Slides), outputPath, chosenTheme)
	return nil
}

func resolveThemeName(cliTheme, frontmatterTheme string) string {
	if cliTheme != "" {
		return cliTheme
	}
	if frontmatterTheme != "" {
		return frontmatterTheme
	}
	return theme.DefaultTheme
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
	Title       string
	ComposedCSS template.CSS
	Slides      []slideViewData
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

func renderPreviewHTML(_ context.Context, deck *model.Deck, outputPath, themeName, themePath string) error {
	mgr := theme.NewManager(nil)
	composedCSS, err := mgr.ComposeFullCSS(themeName, themePath, deck.CustomCSS)
	if err != nil {
		return fmt.Errorf("failed to compose presentation styles: %w", err)
	}

	data := buildPreviewData(deck, composedCSS)

	tmpl, err := template.New("preview").Parse(previewHTMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse preview template: %w", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() { _ = out.Close() }()

	return tmpl.Execute(out, data)
}

func buildPreviewData(deck *model.Deck, composedCSS string) previewTemplateData {
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
		Title:       title,
		ComposedCSS: template.CSS(composedCSS), // nolint:gosec
		Slides:      views,
	}
}

const previewHTMLTemplate = `<!DOCTYPE html>
<html lang="ko">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .Title }}</title>
  <style>
{{ .ComposedCSS }}
  </style>
</head>
<body>
  {{ range .Slides }}
  <div class="slide-wrapper">
    <div class="slide-badge">Slide {{ .Index }} [{{ .Layout }}]</div>
    <div class="slide-card {{ .Layout }} layout-{{ .Layout }} {{ .Classes }}"
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
