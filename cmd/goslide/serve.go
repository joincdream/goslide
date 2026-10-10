package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/i18n"
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

	printServeBanner(serveBindFlag, servePortFlag, inputFile, themeName)

	if err := srv.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		if isPortInUseError(err) {
			return newCLIError(model.ExitGeneralError, formatPortInUseError(servePortFlag, inputFile))
		}
		return newCLIError(model.ExitGeneralError, fmt.Errorf("server error: %w", err))
	}

	fmt.Println("\nGoslide Live Server stopped gracefully.")
	return nil
}

func printServeBanner(bind string, port int, inputFile, theme string) {
	lang := i18n.DetectLocale()
	isKorean := strings.HasPrefix(strings.ToLower(lang), "ko")

	url := fmt.Sprintf("http://%s:%d", bind, port)
	fileInfo := fmt.Sprintf("%s (Theme: %s)", filepath.Base(inputFile), theme)

	if isKorean {
		fmt.Println("┌─────────────────────────────────────────────────────────────┐")
		fmt.Printf("│  Goslide 실시간 개발 서버                                    │\n")
		fmt.Printf("│  주소:     %-48s │\n", url)
		fmt.Printf("│  감시 중:  %-48s │\n", truncateBannerText(fileInfo, 48))
		fmt.Printf("│  단축키:   [?] 도움말  [P] 발표자 콘솔  [D] 판서  [F] 전체화면  │\n")
		fmt.Println("└─────────────────────────────────────────────────────────────┘")
		fmt.Printf("종료하려면 Ctrl+C를 누르세요\n\n")
		return
	}

	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Printf("│  Goslide Live Server                                        │\n")
	fmt.Printf("│  URL:      %-48s │\n", url)
	fmt.Printf("│  Watching: %-48s │\n", truncateBannerText(fileInfo, 48))
	fmt.Printf("│  Hotkeys:  [?] Help  [P] Presenter  [D] Pen  [F] Fullscreen │\n")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Printf("Press Ctrl+C to stop\n\n")
}

func truncateBannerText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func isPortInUseError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "address already in use") ||
		strings.Contains(msg, "eaddrinuse") ||
		strings.Contains(msg, "bind: address already in use")
}

func formatPortInUseError(port int, inputFile string) error {
	lang := i18n.DetectLocale()
	isKorean := strings.HasPrefix(strings.ToLower(lang), "ko")

	altPort := port + 1
	if port == 8080 {
		altPort = 8081
	}

	if isKorean {
		return fmt.Errorf("포트 %d이(가) 이미 다른 프로세스에서 사용 중입니다.\n해결 방법: --port 플래그로 다른 포트를 지정해 보세요 (예: goslide serve %s --port %d)",
			port, filepath.Base(inputFile), altPort)
	}

	return fmt.Errorf("port %d is already in use by another process.\nSolution: specify a different port using --port (e.g. goslide serve %s --port %d)",
		port, filepath.Base(inputFile), altPort)
}
