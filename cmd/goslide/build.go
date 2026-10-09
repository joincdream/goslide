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
	buildCmd.Flags().StringVarP(&outputPathFlag, "output", "o", "", "Output file path or directory (default: <input>.<ext>)")
	buildCmd.Flags().StringVarP(&formatFlag, "format", "f", "html", "Output format: html, pdf, pptx, or all (comma-separated; default: html)")
	buildCmd.Flags().StringVarP(&themeFlag, "theme", "t", "", "Theme name (default, clean, dark; overrides frontmatter)")
	buildCmd.Flags().StringVar(&themePathFlag, "theme-path", "", "Path to custom external CSS stylesheet")
	buildCmd.Flags().BoolVar(&standaloneFlag, "standalone", false, "Inline local image assets as Base64 Data URIs")
	buildCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress informational build messages")
	buildCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable detailed diagnostic build output")
	buildCmd.SilenceUsage = true
	buildCmd.RunE = runBuild
}

func parseAndValidateFormats(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{"html"}, nil
	}

	parts := strings.Split(raw, ",")
	seen := make(map[string]bool)
	var formats []string
	hasAll := false

	for _, p := range parts {
		item := strings.TrimSpace(strings.ToLower(p))
		if item == "" {
			continue
		}
		if item == "all" {
			hasAll = true
			continue
		}
		if item != "html" && item != "pdf" && item != "pptx" {
			return nil, newCLIError(model.ExitInvalidUsage, fmt.Errorf("unsupported output format %q: supported formats are 'html', 'pdf', 'pptx', or 'all'", p))
		}
		if !seen[item] {
			seen[item] = true
			formats = append(formats, item)
		}
	}

	if hasAll {
		return []string{"html", "pdf", "pptx"}, nil
	}

	if len(formats) == 0 {
		return []string{"html"}, nil
	}

	return formats, nil
}

func validateBuildOptions(args []string) ([]string, error) {
	if len(args) < 1 {
		return nil, newCLIError(model.ExitInvalidUsage, errors.New("input markdown file is required: goslide build <input.md>"))
	}
	if quietFlag && verboseFlag {
		return nil, newCLIError(model.ExitInvalidUsage, errors.New("cannot specify both --quiet and --verbose"))
	}
	formats, err := parseAndValidateFormats(formatFlag)
	if err != nil {
		return nil, err
	}
	if themePathFlag != "" {
		if _, err := os.Stat(themePathFlag); err != nil {
			if os.IsNotExist(err) {
				return nil, newCLIError(model.ExitFileNotFound, fmt.Errorf("custom CSS file %q not found", themePathFlag))
			}
			return nil, newCLIError(model.ExitGeneralError, fmt.Errorf("failed to access custom CSS file %q: %w", themePathFlag, err))
		}
	}
	return formats, nil
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

	if !quietFlag && len(deck.Diagnostics) > 0 {
		for _, diag := range deck.Diagnostics {
			if len(diag.Candidates) > 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "[goslide] ⚠️  Slide %d (Line %d): %s\n          Raw: %q\n          Available: %v\n",
					diag.SlideIndex, diag.Line, diag.Message, diag.RawSnippet, diag.Candidates)
			} else {
				fmt.Fprintf(cmd.ErrOrStderr(), "[goslide] ⚠️  Slide %d (Line %d): %s\n          Raw: %q\n",
					diag.SlideIndex, diag.Line, diag.Message, diag.RawSnippet)
			}
		}
	}

	return deck, nil
}

func renderAndWriteHTML(cmd *cobra.Command, deck *model.Deck, inputPath, outputPath, chosenTheme string) error {
	baseDir := filepath.Dir(inputPath)
	renderer := htmlrenderer.NewRenderer(
		htmlrenderer.WithTheme(chosenTheme),
		htmlrenderer.WithCustomCSS(themePathFlag),
		htmlrenderer.WithLang(langFlag),
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

type formatBuilder func(cmd *cobra.Command, deck *model.Deck, inputPath, outputPath, chosenTheme string) error

var formatBuilders = map[string]formatBuilder{
	"html": renderAndWriteHTML,
	"pdf":  exportPDF,
	"pptx": exportPPTX,
}

func ensureOutputDir(outputPath string) error {
	if outDir := filepath.Dir(outputPath); outDir != "" && outDir != "." {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to create output directory %q: %w", outDir, err))
		}
	}
	return nil
}

func runBuild(cmd *cobra.Command, args []string) error {
	formats, err := validateBuildOptions(args)
	if err != nil {
		return err
	}

	inputPath := args[0]
	outputPaths := resolveOutputPaths(inputPath, outputPathFlag, formats)

	deck, err := openInputDeck(cmd, inputPath)
	if err != nil {
		return err
	}

	chosenTheme := resolveThemeName(themeFlag, deck.GlobalAttrs.Theme)

	for _, format := range formats {
		outputPath := outputPaths[format]
		if err := ensureOutputDir(outputPath); err != nil {
			return err
		}

		builder, ok := formatBuilders[format]
		if !ok {
			return newCLIError(model.ExitInvalidUsage, fmt.Errorf("unsupported output format %q", format))
		}

		if err := builder(cmd, deck, inputPath, outputPath, chosenTheme); err != nil {
			return err
		}

		logBuildStatus(len(deck.Slides), inputPath, outputPath, chosenTheme)
	}

	if len(formats) > 1 && !quietFlag {
		fmt.Printf("🎉 Successfully built %d formats [%s] from %s\n", len(formats), strings.Join(formats, ", "), inputPath)
	}

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

func isDirectoryPath(path string) bool {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return true
	}
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}

func resolveSingleOutputPath(inputPath, outputFlag, format string, isMulti bool) string {
	ext := filepath.Ext(inputPath)
	inputBase := strings.TrimSuffix(filepath.Base(inputPath), ext)
	inputDir := filepath.Dir(inputPath)

	if outputFlag == "" {
		return filepath.Join(inputDir, fmt.Sprintf("%s.%s", inputBase, format))
	}

	if isDirectoryPath(outputFlag) {
		return filepath.Join(outputFlag, fmt.Sprintf("%s.%s", inputBase, format))
	}

	if isMulti {
		dir := filepath.Dir(outputFlag)
		outExt := filepath.Ext(outputFlag)
		outBase := strings.TrimSuffix(filepath.Base(outputFlag), outExt)
		return filepath.Join(dir, fmt.Sprintf("%s.%s", outBase, format))
	}

	return outputFlag
}

func resolveOutputPaths(inputPath, outputFlag string, formats []string) map[string]string {
	result := make(map[string]string, len(formats))

	if len(formats) == 1 && outputFlag != "" && !isDirectoryPath(outputFlag) {
		result[formats[0]] = outputFlag
		return result
	}

	for _, f := range formats {
		result[f] = resolveSingleOutputPath(inputPath, outputFlag, f, len(formats) > 1)
	}
	return result
}

func resolveOutputPath(inputPath, outputFlag, format string) string {
	return resolveSingleOutputPath(inputPath, outputFlag, format, false)
}
