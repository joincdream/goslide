package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/server"
)

var (
	servePortFlag      int
	serveBindFlag      string
	serveOpenFlag      bool
	serveThemeFlag     string
	serveThemePathFlag string
	serveDebounceFlag  time.Duration
)

func init() {
	serveCmd.Flags().IntVarP(&servePortFlag, "port", "p", 8080, "HTTP server listening port")
	serveCmd.Flags().StringVar(&serveBindFlag, "bind", "localhost", "Network interface address to bind to")
	serveCmd.Flags().BoolVar(&serveOpenFlag, "open", false, "Open presentation in default browser on start")
	serveCmd.Flags().StringVarP(&serveThemeFlag, "theme", "t", "", "Theme name (default, clean, dark; overrides frontmatter)")
	serveCmd.Flags().StringVar(&serveThemePathFlag, "theme-path", "", "Path to custom external CSS stylesheet")
	serveCmd.Flags().DurationVar(&serveDebounceFlag, "debounce", 50*time.Millisecond, "File watch debounce interval")
	serveCmd.SilenceUsage = true
	serveCmd.RunE = runServe
}

func validateServeOptions(args []string) error {
	if len(args) < 1 {
		return newCLIError(model.ExitInvalidUsage, errors.New("input markdown file is required: goslide serve <input.md>"))
	}
	inputFile := args[0]
	if _, err := os.Stat(inputFile); err != nil {
		if os.IsNotExist(err) {
			return newCLIError(model.ExitFileNotFound, fmt.Errorf("input file %q not found", inputFile))
		}
		return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to access input file %q: %w", inputFile, err))
	}
	return nil
}

func runServe(cmd *cobra.Command, args []string) error {
	if err := validateServeOptions(args); err != nil {
		return err
	}

	inputFile := args[0]
	absFile, err := filepath.Abs(inputFile)
	if err != nil {
		return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to resolve absolute path for %q: %w", inputFile, err))
	}

	cfg := server.Config{
		Port:         servePortFlag,
		Bind:         serveBindFlag,
		MarkdownPath: absFile,
		Theme:        serveThemeFlag,
		ThemePath:    serveThemePathFlag,
		Lang:         langFlag,
		OpenBrowser:  serveOpenFlag,
		Debounce:     serveDebounceFlag,
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to initialize live server: %w", err))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	themeName := "default"
	if serveThemeFlag != "" {
		themeName = serveThemeFlag
	}

	fmt.Printf("🚀 Goslide Live Server running at http://%s:%d\n", serveBindFlag, servePortFlag)
	fmt.Printf("📂 Watching: %s (Theme: %s)\n", filepath.Base(inputFile), themeName)
	fmt.Printf("💡 Press Ctrl+C to stop\n\n")

	if err := srv.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return newCLIError(model.ExitGeneralError, fmt.Errorf("server error: %w", err))
	}

	fmt.Println("\n👋 Goslide Live Server stopped gracefully.")
	return nil
}
