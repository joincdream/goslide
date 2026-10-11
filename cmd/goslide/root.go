package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/i18n"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/updater"
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
}

func initCLI() {
	parseEarlyLang()

	refreshCLITexts()

	if rootCmd.PersistentFlags().Lookup("lang") == nil {
		rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "auto", i18n.T("cli.lang.flag"))
	}
	if rootCmd.Flags().Lookup("version") == nil {
		rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, i18n.T("cli.version.flag"))
	}

	initUpdateFlags()

	if !hasCommand(rootCmd, buildCmd.Name()) {
		rootCmd.AddCommand(buildCmd)
	}
	if !hasCommand(rootCmd, serveCmd.Name()) {
		rootCmd.AddCommand(serveCmd)
	}
	if !hasCommand(rootCmd, initCmd.Name()) {
		rootCmd.AddCommand(initCmd)
	}
	if !hasCommand(rootCmd, demoCmd.Name()) {
		rootCmd.AddCommand(demoCmd)
	}
	if !hasCommand(rootCmd, promptCmd.Name()) {
		rootCmd.AddCommand(promptCmd)
	}
	if !hasCommand(rootCmd, themeCmd.Name()) {
		rootCmd.AddCommand(themeCmd)
	}
	if !hasCommand(rootCmd, updateCmd.Name()) {
		rootCmd.AddCommand(updateCmd)
	}
}

func hasCommand(root *cobra.Command, name string) bool {
	for _, c := range root.Commands() {
		if c.Name() == name {
			return true
		}
	}
	return false
}

func refreshCLITexts() {
	rootCmd.Short = i18n.T("cli.description")
	rootCmd.Long = i18n.T("cli.long")
	rootCmd.Example = i18n.T("cli.example")

	buildCmd.Short = i18n.T("cli.build.desc")
	buildCmd.Long = i18n.T("cli.build.long")
	buildCmd.Example = i18n.T("cli.build.example")

	serveCmd.Short = i18n.T("cli.serve.desc")
	serveCmd.Long = i18n.T("cli.serve.long")
	serveCmd.Example = i18n.T("cli.serve.example")

	initCmd.Short = i18n.T("cli.init.desc")
	initCmd.Long = i18n.T("cli.init.long")
	initCmd.Example = i18n.T("cli.init.example")

	demoCmd.Short = i18n.T("cli.demo.desc")
	promptCmd.Short = i18n.T("cli.prompt.desc")
	themeCmd.Short = i18n.T("cli.theme.desc")
	themeListCmd.Short = i18n.T("cli.theme.list.desc")
	themeExportCmd.Short = i18n.T("cli.theme.export.desc")

	updateCmd.Short = i18n.T("cli.update.desc")
	updateCmd.Long = i18n.T("cli.update.long")
	updateCmd.Example = i18n.T("cli.update.example")

	updateFlagDescriptions()
}

func updateFlagDescriptions() {
	type flagSpec struct {
		cmd    *cobra.Command
		name   string
		msgKey string
	}

	specs := []flagSpec{
		{rootCmd, "lang", "cli.lang.flag"},
		{rootCmd, "version", "cli.version.flag"},
		{buildCmd, "output", "cli.build.flag.output"},
		{buildCmd, "format", "cli.build.flag.format"},
		{buildCmd, "theme", "cli.build.flag.theme"},
		{buildCmd, "theme-path", "cli.build.flag.theme_path"},
		{buildCmd, "standalone", "cli.build.flag.standalone"},
		{buildCmd, "quiet", "cli.build.flag.quiet"},
		{buildCmd, "verbose", "cli.build.flag.verbose"},
		{serveCmd, "port", "cli.serve.flag.port"},
		{serveCmd, "bind", "cli.serve.flag.bind"},
		{serveCmd, "open", "cli.serve.flag.open"},
		{serveCmd, "theme", "cli.serve.flag.theme"},
		{serveCmd, "theme-path", "cli.serve.flag.theme_path"},
		{serveCmd, "debounce", "cli.serve.flag.debounce"},
		{initCmd, "theme", "cli.init.flag.theme"},
		{initCmd, "force", "cli.init.flag.force"},
		{demoCmd, "force", "cli.demo.flag.force"},
		{updateCmd, "check", "cli.update.flag.check"},
		{updateCmd, "dry-run", "cli.update.flag.dry_run"},
		{updateCmd, "force", "cli.update.flag.force"},
	}

	for _, s := range specs {
		if f := s.cmd.Flags().Lookup(s.name); f != nil {
			f.Usage = i18n.Lookup(s.msgKey)
		} else if f := s.cmd.PersistentFlags().Lookup(s.name); f != nil {
			f.Usage = i18n.Lookup(s.msgKey)
		}
	}
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

// checkAndNotifyUpdate prints a non-intrusive tip if a newer version of Goslide is available.
func checkAndNotifyUpdate() {
	if os.Getenv("CI") != "" || os.Getenv("GOSLIDE_NO_UPDATE_NOTIFIER") != "" {
		return
	}

	for _, arg := range os.Args[1:] {
		if arg == "update" || arg == "-v" || arg == "--version" || arg == "-h" || arg == "--help" || arg == "--quiet" || arg == "-q" {
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	checker := updater.NewChecker(nil, nil)
	result, err := checker.CheckUpdate(ctx, goslide.Version)
	if err != nil || result == nil || !result.HasUpdate {
		return
	}

	fmt.Println()
	fmt.Printf(i18n.T("update.tip")+"\n", result.LatestVersion, result.CurrentVersion, result.ReleaseURL)
}

// Execute runs the root CLI command and terminates with standard exit codes.
func Execute() {
	initCLI()
	if err := rootCmd.Execute(); err != nil {
		code := determineExitCode(err)
		os.Exit(code)
	}
	checkAndNotifyUpdate()
}
