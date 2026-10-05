package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	pdfexporter "github.com/yundream/goslide/internal/exporter/pdf"
	pptxexporter "github.com/yundream/goslide/internal/exporter/pptx"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/parser"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
	"github.com/yundream/goslide/internal/theme"
)

// CLIError wraps an error with an associated standard POSIX exit code.
type CLIError struct {
	Code int
	Err  error
}

func (e *CLIError) Error() string {
	return e.Err.Error()
}

func (e *CLIError) Unwrap() error {
	return e.Err
}

func newCLIError(code int, err error) error {
	return &CLIError{
		Code: code,
		Err:  err,
	}
}

var (
	outputPathFlag string
	formatFlag     string
	themeFlag      string
	themePathFlag  string
	standaloneFlag bool
	quietFlag      bool
	verboseFlag    bool
)

func init() {
	buildCmd.Flags().StringVarP(&outputPathFlag, "output", "o", "", "Output file path (default: <input>.<ext>)")
	buildCmd.Flags().StringVarP(&formatFlag, "format", "f", "html", "Output format: html, pdf, pptx (default: html)")
	buildCmd.Flags().StringVarP(&themeFlag, "theme", "t", "", "Theme name (default, clean, dark; overrides frontmatter)")
	buildCmd.Flags().StringVar(&themePathFlag, "theme-path", "", "Path to custom external CSS stylesheet")
	buildCmd.Flags().BoolVar(&standaloneFlag, "standalone", false, "Inline local image assets as Base64 Data URIs")
	buildCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress informational build messages")
	buildCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable detailed diagnostic build output")
	buildCmd.SilenceUsage = true
	buildCmd.RunE = runBuild
}

func validateBuildOptions(args []string) error {
	if len(args) < 1 {
		return newCLIError(model.ExitInvalidUsage, errors.New("input markdown file is required: goslide build <input.md>"))
	}
	if quietFlag && verboseFlag {
		return newCLIError(model.ExitInvalidUsage, errors.New("cannot specify both --quiet and --verbose"))
	}
	fmtLower := strings.ToLower(formatFlag)
	if fmtLower != "html" && fmtLower != "pdf" && fmtLower != "pptx" {
		return newCLIError(model.ExitInvalidUsage, fmt.Errorf("unsupported output format %q: supported formats are 'html', 'pdf', 'pptx'", formatFlag))
	}
	if themePathFlag != "" {
		if _, err := os.Stat(themePathFlag); err != nil {
			if os.IsNotExist(err) {
				return newCLIError(model.ExitFileNotFound, fmt.Errorf("custom CSS file %q not found", themePathFlag))
			}
			return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to access custom CSS file %q: %w", themePathFlag, err))
		}
	}
	return nil
}

func openInputDeck(cmd *cobra.Command, inputPath string) (*model.Deck, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, newCLIError(model.ExitFileNotFound, fmt.Errorf("input file %q not found", inputPath))
		}
		return nil, newCLIError(model.ExitGeneralError, fmt.Errorf("failed to open input file %q: %w", inputPath, err))
	}
	defer func() { _ = file.Close() }()

	p := parser.NewParser()
	deck, err := p.Parse(cmd.Context(), file)
	if err != nil {
		return nil, newCLIError(model.ExitParseError, fmt.Errorf("failed to parse markdown: %w", err))
	}
	return deck, nil
}

func renderAndWriteHTML(cmd *cobra.Command, deck *model.Deck, inputPath, outputPath, chosenTheme string) error {
	baseDir := filepath.Dir(inputPath)
	renderer := htmlrenderer.NewRenderer(
		htmlrenderer.WithTheme(chosenTheme),
		htmlrenderer.WithCustomCSS(themePathFlag),
		htmlrenderer.WithStandalone(standaloneFlag),
		htmlrenderer.WithBaseDir(baseDir),
	)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to create output file %q: %w", outputPath, err))
	}
	defer func() { _ = outFile.Close() }()

	if err := renderer.Render(cmd.Context(), deck, outFile); err != nil {
		return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to render HTML slides: %w", err))
	}
	return nil
}

func exportPDF(cmd *cobra.Command, deck *model.Deck, inputPath, outputPath, chosenTheme string) error {
	baseDir := filepath.Dir(inputPath)
	exporter := pdfexporter.NewExporter(
		pdfexporter.WithTheme(chosenTheme),
		pdfexporter.WithCustomCSS(themePathFlag),
		pdfexporter.WithStandalone(standaloneFlag),
		pdfexporter.WithBaseDir(baseDir),
	)

	if err := exporter.Export(cmd.Context(), deck, outputPath); err != nil {
		return newCLIError(model.ExitExportFailed, fmt.Errorf("failed to export PDF slides: %w", err))
	}
	return nil
}

func exportPPTX(cmd *cobra.Command, deck *model.Deck, inputPath, outputPath, chosenTheme string) error {
	baseDir := filepath.Dir(inputPath)
	exporter := pptxexporter.NewExporter(
		pptxexporter.WithTheme(chosenTheme),
		pptxexporter.WithCustomCSS(themePathFlag),
		pptxexporter.WithStandalone(standaloneFlag),
		pptxexporter.WithBaseDir(baseDir),
	)

	if err := exporter.Export(cmd.Context(), deck, outputPath); err != nil {
		return newCLIError(model.ExitExportFailed, fmt.Errorf("failed to export PPTX slides: %w", err))
	}
	return nil
}

func logBuildStatus(slidesCount int, inputPath, outputPath, chosenTheme string) {
	if verboseFlag {
		fmt.Printf("[verbose] Successfully parsed %d slides from %s\n", slidesCount, inputPath)
		fmt.Printf("[verbose] Theme: %s\n", chosenTheme)
		if themePathFlag != "" {
			fmt.Printf("[verbose] External stylesheet: %s\n", themePathFlag)
		}
		fmt.Printf("[verbose] Standalone mode: %t\n", standaloneFlag)
		fmt.Printf("[verbose] Output written to: %s\n", outputPath)
	}
	if !quietFlag {
		fmt.Printf("✓ Successfully built %d slides to %s [theme: %s]\n", slidesCount, outputPath, chosenTheme)
	}
}

func runBuild(cmd *cobra.Command, args []string) error {
	if err := validateBuildOptions(args); err != nil {
		return err
	}

	inputPath := args[0]
	format := strings.ToLower(formatFlag)
	outputPath := resolveOutputPath(inputPath, outputPathFlag, format)

	deck, err := openInputDeck(cmd, inputPath)
	if err != nil {
		return err
	}

	chosenTheme := resolveThemeName(themeFlag, deck.GlobalAttrs.Theme)
	switch format {
	case "pdf":
		if err := exportPDF(cmd, deck, inputPath, outputPath, chosenTheme); err != nil {
			return err
		}
	case "pptx":
		if err := exportPPTX(cmd, deck, inputPath, outputPath, chosenTheme); err != nil {
			return err
		}
	default:
		if err := renderAndWriteHTML(cmd, deck, inputPath, outputPath, chosenTheme); err != nil {
			return err
		}
	}

	logBuildStatus(len(deck.Slides), inputPath, outputPath, chosenTheme)
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

func resolveOutputPath(inputPath, outputFlag, format string) string {
	if outputFlag != "" {
		return outputFlag
	}
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	switch format {
	case "pdf":
		return base + ".pdf"
	case "pptx":
		return base + ".pptx"
	default:
		return base + ".html"
	}
}
