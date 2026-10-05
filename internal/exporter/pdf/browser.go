package pdf

import (
	"github.com/yundream/goslide/internal/browser"
)

// FindChrome discovers an installed Chrome or Chromium binary executable.
func FindChrome(customPath string) (string, error) {
	return browser.FindChrome(customPath)
}
