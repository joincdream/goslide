package browser

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

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

	return "", fmt.Errorf("%w: could not locate Google Chrome or Chromium in PATH or default locations. Please install Chrome or set GOSLIDE_CHROME_BIN", model.ErrChromeNotFound)
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
