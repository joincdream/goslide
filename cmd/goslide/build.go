package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/parser"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
	"github.com/yundream/goslide/internal/theme"
)

var (
	outputPathFlag string
	themeFlag      string
	themePathFlag  string
	standaloneFlag bool
)

func init() {
	buildCmd.Flags().StringVarP(&outputPathFlag, "output", "o", "", "Output HTML file path (default: <input>.html)")
	buildCmd.Flags().StringVarP(&themeFlag, "theme", "t", "", "Theme name (default, clean, dark; overrides frontmatter)")
	buildCmd.Flags().StringVar(&themePathFlag, "theme-path", "", "Path to custom external CSS stylesheet")
	buildCmd.Flags().BoolVar(&standaloneFlag, "standalone", false, "Inline local image assets as Base64 Data URIs")
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
	baseDir := filepath.Dir(inputPath)

	renderer := htmlrenderer.NewRenderer(
		htmlrenderer.WithTheme(chosenTheme),
		htmlrenderer.WithCustomCSS(themePathFlag),
		htmlrenderer.WithStandalone(standaloneFlag),
		htmlrenderer.WithBaseDir(baseDir),
	)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", outputPath, err)
	}
	defer func() { _ = outFile.Close() }()

	if err := renderer.Render(cmd.Context(), deck, outFile); err != nil {
		return fmt.Errorf("failed to render HTML slides: %w", err)
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
