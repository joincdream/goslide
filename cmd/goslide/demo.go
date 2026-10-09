package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/i18n"
	"github.com/yundream/goslide/internal/theme"
)

var demoForceFlag bool

var demoCmd = &cobra.Command{
	Use:   "demo [directory]",
	Short: "Unpack demo presentation into dedicated directory (./demo) and themes into current directory (themes/)",
	Long:  "Unpacks a rich showcase presentation Markdown file (demo.md) and image assets into a dedicated demo directory, and built-in themes into themes/ in the current directory, providing an interactive playground to experience Goslide.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		demoDir := "demo"
		if len(args) > 0 && args[0] != "" {
			demoDir = args[0]
		}

		targetFile := filepath.Join(demoDir, "demo.md")
		themesDir := "themes"

		// 1. Safety check: ensure targetFile doesn't already exist without --force
		if !demoForceFlag {
			if _, err := os.Stat(targetFile); err == nil {
				return errors.New(fmt.Sprintf(i18n.Lookup("cli.demo.already_exists"), targetFile))
			}
		}

		// 2. Ensure directories exist
		if err := os.MkdirAll(demoDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", demoDir, err)
		}
		if err := os.MkdirAll(themesDir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", themesDir, err)
		}

		themeMgr := theme.NewManager(nil)

		// 3. Load and write demo.md
		demoContent, err := themeMgr.GetStarterTemplate("clean")
		if err != nil {
			return fmt.Errorf("failed to load demo template: %w", err)
		}

		if err := os.WriteFile(targetFile, []byte(demoContent), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetFile, err)
		}

		// 4. Extract 3 showcase theme stylesheets into current directory's themes/
		showcaseThemes := []string{"clean", "dark", "academic"}
		for _, thName := range showcaseThemes {
			cssContent, err := themeMgr.GetThemeCSS(thName)
			if err != nil {
				continue
			}
			destPath := filepath.Join(themesDir, thName+".css")
			if !demoForceFlag {
				if _, err := os.Stat(destPath); err == nil {
					continue // Preserve existing file unless --force
				}
			}
			_ = os.WriteFile(destPath, []byte(cssContent), 0644)
		}

		// 5. Extract demo image assets into demo directory
		images, err := themeMgr.GetDemoImages()
		if err == nil {
			for imgName, imgData := range images {
				destImgPath := filepath.Join(demoDir, imgName)
				if !demoForceFlag {
					if _, err := os.Stat(destImgPath); err == nil {
						continue // Preserve existing file unless --force
					}
				}
				_ = os.WriteFile(destImgPath, imgData, 0644)
			}
		}

		// 6. Print actionable, localized step-by-step guide
		fmt.Printf(i18n.T("cli.demo.success")+"\n", demoDir, themesDir, demoDir, demoDir, demoDir, themesDir, demoDir)
		return nil
	},
}

func init() {
	demoCmd.Short = i18n.Lookup("cli.demo.desc")
	demoCmd.Flags().BoolVarP(&demoForceFlag, "force", "f", false, i18n.Lookup("cli.demo.flag.force"))
}
