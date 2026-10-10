package browser

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/yundream/goslide/internal/model"
)

// defaultExecNames contains common executable names to search in PATH.
var defaultExecNames = []string{
	"google-chrome",
	"google-chrome-stable",
	"chromium",
	"chromium-browser",
	"chrome",
}

// defaultKnownPaths contains OS-specific standard binary installation locations.
var defaultKnownPaths = map[string][]string{
	"linux": {
		"/usr/bin/google-chrome",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/snap/bin/chromium",
	},
	"darwin": {
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
	},
	"windows": {
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
	},
}

// FindChrome discovers an installed Chrome or Chromium binary executable.
// The search order is:
// 1. Explicit customPath (if provided)
// 2. GOSLIDE_CHROME_BIN and CHROME_BIN environment variables
// 3. System PATH via exec.LookPath
// 4. OS-specific known installation paths
func FindChrome(customPath string) (string, error) {
	// 1. Explicit custom path
	if customPath != "" {
		if path, err := verifyExecutable(customPath); err == nil {
			return path, nil
		}
		return "", fmt.Errorf("%w: custom browser path %q is not executable", model.ErrChromeNotFound, customPath)
	}

	// 2. Environment variables
	for _, envKey := range []string{"GOSLIDE_CHROME_BIN", "CHROME_BIN"} {
		if envVal := os.Getenv(envKey); envVal != "" {
			if path, err := verifyExecutable(envVal); err == nil {
				return path, nil
			}
		}
	}

	// 3. Search in system PATH
	for _, name := range defaultExecNames {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	// 4. Check known OS installation paths
	if paths, ok := defaultKnownPaths[runtime.GOOS]; ok {
		for _, p := range paths {
			if path, err := verifyExecutable(p); err == nil {
				return path, nil
			}
		}
	}

	guide := FormatChromeInstallGuide(detectSystemLocale())
	return "", fmt.Errorf("%w: %s", model.ErrChromeNotFound, guide)
}

// FormatChromeInstallGuide returns platform-specific, copy-pasteable installation instructions.
func FormatChromeInstallGuide(lang string) string {
	isKorean := strings.HasPrefix(strings.ToLower(lang), "ko")

	if isKorean {
		var installCmd string
		switch runtime.GOOS {
		case "darwin":
			installCmd = "  • macOS: brew install --cask google-chrome"
		case "windows":
			installCmd = "  • Windows: Google Chrome을 설치하거나 Microsoft Edge(msedge.exe) 경로를 확인하세요."
		default: // linux / others
			installCmd = "  • Ubuntu / Debian: sudo apt update && sudo apt install -y chromium-browser\n" +
				"  • Fedora / RHEL:   sudo dnf install -y chromium\n" +
				"  • Arch Linux:      sudo pacman -S chromium"
		}

		return fmt.Sprintf("PDF 및 PPTX 생성을 위해 Google Chrome 또는 Chromium 브라우저가 필요합니다.\n\n"+
			"설치 방법:\n%s\n\n"+
			"임의 경로 지정:\n"+
			"  환경 변수 GOSLIDE_CHROME_BIN=<경로> 또는 --chrome-path 플래그를 설정하세요.", installCmd)
	}

	var installCmd string
	switch runtime.GOOS {
	case "darwin":
		installCmd = "  • macOS: brew install --cask google-chrome"
	case "windows":
		installCmd = "  • Windows: Install Google Chrome or verify Microsoft Edge (msedge.exe) path."
	default: // linux / others
		installCmd = "  • Ubuntu / Debian: sudo apt update && sudo apt install -y chromium-browser\n" +
			"  • Fedora / RHEL:   sudo dnf install -y chromium\n" +
			"  • Arch Linux:      sudo pacman -S chromium"
	}

	return fmt.Sprintf("Google Chrome or Chromium is required for PDF and PPTX export.\n\n"+
		"Installation:\n%s\n\n"+
		"Custom Path:\n"+
		"  Set environment variable GOSLIDE_CHROME_BIN=<path> or use --chrome-path flag.", installCmd)
}

func detectSystemLocale() string {
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		val := strings.ToLower(os.Getenv(env))
		if strings.HasPrefix(val, "ko") {
			return "ko"
		}
		if strings.HasPrefix(val, "en") {
			return "en"
		}
	}
	return "en"
}

func verifyExecutable(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s is a directory", path)
	}
	return path, nil
}
