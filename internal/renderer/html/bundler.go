package html

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	imgSrcRegex = regexp.MustCompile(`(?i)(<img\s+[^>]*src=["'])([^"']+)(["'][^>]*>)`)
)

// AssetBundler provides utilities for inlining local image resources into HTML.
type AssetBundler struct {
	baseDir string
}

// NewAssetBundler creates a bundler with the given base directory for relative asset resolution.
func NewAssetBundler(baseDir string) *AssetBundler {
	if baseDir == "" {
		baseDir = "."
	}
	return &AssetBundler{baseDir: baseDir}
}

// BundleImages converts local image src attributes in the given HTML string into Base64 Data URIs.
// Remote URLs (http://, https://) and existing Data URIs are left untouched.
func (b *AssetBundler) BundleImages(htmlContent string) (string, error) {
	if htmlContent == "" {
		return "", nil
	}

	var replaceErr error
	bundled := imgSrcRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		if replaceErr != nil {
			return match
		}

		submatches := imgSrcRegex.FindStringSubmatch(match)
		if len(submatches) < 4 {
			return match
		}

		prefix := submatches[1]
		src := submatches[2]
		suffix := submatches[3]

		// Skip remote URLs and existing data URIs
		lowerSrc := strings.ToLower(strings.TrimSpace(src))
		if strings.HasPrefix(lowerSrc, "http://") ||
			strings.HasPrefix(lowerSrc, "https://") ||
			strings.HasPrefix(lowerSrc, "data:") {
			return match
		}

		dataURI, err := b.encodeFileToDataURI(src)
		if err != nil {
			// If file does not exist or cannot be read, preserve original src
			return match
		}

		return prefix + dataURI + suffix
	})

	if replaceErr != nil {
		return "", replaceErr
	}

	return bundled, nil
}

func (b *AssetBundler) encodeFileToDataURI(relPath string) (string, error) {
	cleanPath := relPath
	if !filepath.IsAbs(relPath) {
		cleanPath = filepath.Join(b.baseDir, relPath)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to read local image %q: %w", cleanPath, err)
	}

	mimeType := detectImageMimeType(cleanPath, data)
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func detectImageMimeType(filePath string, data []byte) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	default:
		contentType := http.DetectContentType(data)
		if contentType != "" && contentType != "application/octet-stream" {
			return contentType
		}
		return "application/octet-stream"
	}
}
