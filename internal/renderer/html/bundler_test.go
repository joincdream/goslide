package html_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
)

func TestAssetBundler_BundleImages(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testImagePath := filepath.Join(tmpDir, "test.png")
	// Minimal 1x1 PNG bytes
	pngBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
		0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
	if err := os.WriteFile(testImagePath, pngBytes, 0600); err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	bundler := htmlrenderer.NewAssetBundler(tmpDir)

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "local png image is converted to base64 data uri",
			input:    `<p><img src="test.png" alt="sample"></p>`,
			contains: `src="data:image/png;base64,`,
		},
		{
			name:     "remote http url is preserved",
			input:    `<img src="https://example.com/logo.png">`,
			contains: `src="https://example.com/logo.png"`,
		},
		{
			name:     "existing data uri is preserved",
			input:    `<img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7">`,
			contains: `src="data:image/gif;base64,`,
		},
		{
			name:     "missing local file is preserved gracefully",
			input:    `<img src="not-found.png">`,
			contains: `src="not-found.png"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := bundler.BundleImages(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(result, tt.contains) {
				t.Errorf("expected result to contain %q, got %q", tt.contains, result)
			}
		})
	}
}
