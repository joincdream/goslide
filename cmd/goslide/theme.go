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

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage built-in themes and export custom stylesheets (CSS)",
	Long:  "Commands to list built-in themes and export their stylesheets for local customization.",
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available built-in themes",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		themes := []struct {
			name string
			desc string
		}{
			{"clean", "Modern, clean light theme for business and general talks (default)"},
			{"dark", "Deep dark theme with high contrast, ideal for developer seminars"},
			{"academic", "Formal serif typography optimized for research and academic conferences"},
			{"cyber-dark", "Neon-accented dark theme for tech, hacking, and gaming topics"},
			{"default", "Standard neutral theme"},
		}

		fmt.Println("Available Built-in Themes in Goslide:")
		fmt.Println("--------------------------------------------------")
		for _, th := range themes {
			fmt.Printf("  • %-12s - %s\n", th.name, th.desc)
		}
		fmt.Println("--------------------------------------------------")
		fmt.Println("Usage: goslide theme export <theme-name> [output.css]")
		return nil
	},
}

var themeExportCmd = &cobra.Command{
	Use:   "export <theme_name> [output.css]",
	Short: "Export built-in theme CSS to a local file for customization",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		themeName := strings.ToLower(strings.TrimSpace(args[0]))
		outputFile := themeName + ".css"
		if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
			outputFile = strings.TrimSpace(args[1])
		}

		themeMgr := theme.NewManager(nil)
		if !themeMgr.HasBuiltinTheme(themeName) {
			return fmt.Errorf("unknown built-in theme %q; run 'goslide theme list' to see available themes", themeName)
		}

		cssData, err := themeMgr.GetThemeCSS(themeName)
		if err != nil {
			return fmt.Errorf("failed to retrieve theme CSS: %w", err)
		}

		// Ensure parent directory exists
		dir := filepath.Dir(outputFile)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %q: %w", dir, err)
			}
		}

		if err := os.WriteFile(outputFile, []byte(cssData), 0644); err != nil {
			return fmt.Errorf("failed to write stylesheet %q: %w", outputFile, err)
		}

		fmt.Printf(i18n.T("cli.theme.export.success")+"\n", themeName, outputFile)
		return nil
	},
}

func init() {
	themeCmd.Short = i18n.Lookup("cli.theme.desc")
	themeListCmd.Short = i18n.Lookup("cli.theme.list.desc")
	themeExportCmd.Short = i18n.Lookup("cli.theme.export.desc")

	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeExportCmd)
}
