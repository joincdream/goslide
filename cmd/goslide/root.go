package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/i18n"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/pkg/goslide"
)

var (
	langFlag    string
	versionFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "goslide",
	Short: "Goslide - Pure Go Markdown Presentation Builder",
	RunE: func(cmd *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Printf(i18n.T("cli.version")+"\n", goslide.Version)
			return nil
		}
		return cmd.Help()
	},
}

var buildCmd = &cobra.Command{
	Use:   "build <input.md>",
	Short: "Build presentation slides from Markdown source",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Build command will be implemented in upcoming milestone (M1-6).")
		return nil
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve <input.md>",
	Short: "Start local development server with live reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Serve command will be implemented in upcoming milestone (M2-3).")
		return nil
	},
}

func initCLI() {
	parseEarlyLang()

	rootCmd.Short = i18n.T("cli.description")
	buildCmd.Short = i18n.T("cli.build.desc")
	serveCmd.Short = i18n.T("cli.serve.desc")

	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "auto", i18n.T("cli.lang.flag"))
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serveCmd)
}

func parseEarlyLang() {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--lang=") {
			i18n.SetLocale(strings.TrimPrefix(arg, "--lang="))
			return
		}
		if arg == "--lang" && i+1 < len(os.Args) {
			i18n.SetLocale(os.Args[i+1])
			return
		}
	}
}

// determineExitCode resolves the appropriate exit code from a given error.
func determineExitCode(err error) int {
	if err == nil {
		return model.ExitSuccess
	}
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr.Code
	}
	if errors.Is(err, os.ErrNotExist) {
		return model.ExitFileNotFound
	}
	if errors.Is(err, model.ErrInvalidFrontmatter) {
		return model.ExitParseError
	}
	errStr := err.Error()
	if strings.Contains(errStr, "unknown flag") ||
		strings.Contains(errStr, "unknown shorthand flag") ||
		strings.Contains(errStr, "flag needs an argument") ||
		strings.Contains(errStr, "invalid argument") ||
		strings.Contains(errStr, "accepts ") {
		return model.ExitInvalidUsage
	}
	return model.ExitGeneralError
}

// Execute runs the root CLI command and terminates with standard exit codes.
func Execute() {
	initCLI()
	if err := rootCmd.Execute(); err != nil {
		code := determineExitCode(err)
		os.Exit(code)
	}
}
