package browser

import (
	"errors"
	"strings"
	"testing"

	"github.com/yundream/goslide/internal/model"
)

func TestFormatChromeInstallGuide(t *testing.T) {
	tests := []struct {
		name         string
		lang         string
		expectSubstr []string
	}{
		{
			name: "Korean guide",
			lang: "ko",
			expectSubstr: []string{
				"PDF 및 PPTX 생성을 위해 Google Chrome 또는 Chromium 브라우저가 필요합니다.",
				"GOSLIDE_CHROME_BIN",
			},
		},
		{
			name: "English guide",
			lang: "en",
			expectSubstr: []string{
				"Google Chrome or Chromium is required for PDF and PPTX export.",
				"GOSLIDE_CHROME_BIN",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guide := FormatChromeInstallGuide(tt.lang)
			for _, substr := range tt.expectSubstr {
				if !strings.Contains(guide, substr) {
					t.Errorf("FormatChromeInstallGuide(%q) expected to contain %q, but got:\n%s", tt.lang, substr, guide)
				}
			}
		})
	}
}

func TestFindChromeCustomPathInvalid(t *testing.T) {
	_, err := FindChrome("/path/that/definitely/does/not/exist/chrome-binary")
	if err == nil {
		t.Fatal("expected error for non-existent custom path, got nil")
	}
	if !errors.Is(err, model.ErrChromeNotFound) {
		t.Errorf("expected ErrChromeNotFound, got: %v", err)
	}
}
