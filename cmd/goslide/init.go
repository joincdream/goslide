package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/i18n"
	"github.com/yundream/goslide/internal/theme"
)

var (
	initThemeFlag string
	initForceFlag bool
)

var initCmd = &cobra.Command{
	Use:   "init [filename]",
	Short: "Create a new starter presentation Markdown file",
	Long:  "Initializes a rich sample presentation Markdown file (.md) showcasing Goslide features (layouts, code highlighting, math, mermaid diagrams, incremental fragments, and presenter notes).",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetFile := "presentation.md"
		if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
			targetFile = strings.TrimSpace(args[0])
		}

		if !strings.HasSuffix(strings.ToLower(targetFile), ".md") {
			targetFile += ".md"
		}

		// Check if file already exists
		if _, err := os.Stat(targetFile); err == nil && !initForceFlag {
			return fmt.Errorf("file %q already exists; use --force to overwrite", targetFile)
		}

		themeMgr := theme.NewManager(nil)
		content, err := themeMgr.GetStarterTemplate(initThemeFlag)
		if err != nil {
			return fmt.Errorf("failed to load starter template: %w", err)
		}

		// Ensure parent directory exists if specified
		dir := filepath.Dir(targetFile)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %q: %w", dir, err)
			}
		}

		if err := os.WriteFile(targetFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write starter presentation: %w", err)
		}

		fmt.Printf(i18n.T("cli.init.success")+"\n", targetFile, targetFile, targetFile)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&initThemeFlag, "theme", "t", "clean", "Starter theme name (default, clean, dark)")
	initCmd.Flags().BoolVarP(&initForceFlag, "force", "f", false, "Overwrite existing file without confirmation")
}
