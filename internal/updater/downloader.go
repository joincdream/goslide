package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

// FindMatchingAsset finds the binary archive asset for the current OS and architecture.
func FindMatchingAsset(assets []Asset, targetOS, targetArch string) (*Asset, error) {
	if targetOS == "" {
		targetOS = runtime.GOOS
	}
	if targetArch == "" {
		targetArch = runtime.GOARCH
	}

	for _, a := range assets {
		lowerName := strings.ToLower(a.Name)
		if strings.HasSuffix(lowerName, ".txt") || strings.HasSuffix(lowerName, ".sha256") {
			continue
		}

		// Match OS and Arch in filename (e.g., goslide_linux_amd64.tar.gz or goslide_v0.2.0_linux_amd64.tar.gz)
		if strings.Contains(lowerName, targetOS) && strings.Contains(lowerName, targetArch) {
			return &a, nil
		}
	}

	return nil, fmt.Errorf("no matching release asset found for %s/%s", targetOS, targetArch)
}

// FindChecksumAsset finds the checksums.txt asset.
func FindChecksumAsset(assets []Asset) (*Asset, error) {
	for _, a := range assets {
		lowerName := strings.ToLower(a.Name)
		if strings.Contains(lowerName, "checksum") || strings.HasSuffix(lowerName, ".sha256") {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("no checksums file found in release assets")
}

// DownloadBytes downloads arbitrary bytes from a given URL using context.
func DownloadBytes(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with HTTP status %d: %s", resp.StatusCode, resp.Status)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read response stream: %w", err)
	}

	return buf.Bytes(), nil
}

// ExtractBinaryFromArchive extracts the executable binary file from a .tar.gz or .zip archive.
func ExtractBinaryFromArchive(archiveData []byte, assetName, binaryName string) ([]byte, error) {
	if binaryName == "" {
		binaryName = "goslide"
		if runtime.GOOS == "windows" {
			binaryName = "goslide.exe"
		}
	}

	lowerAsset := strings.ToLower(assetName)
	if strings.HasSuffix(lowerAsset, ".tar.gz") || strings.HasSuffix(lowerAsset, ".tgz") {
		return extractFromTarGz(archiveData, binaryName)
	} else if strings.HasSuffix(lowerAsset, ".zip") {
		return extractFromZip(archiveData, binaryName)
	}

	return nil, fmt.Errorf("unsupported archive format: %s", assetName)
}

func extractFromTarGz(data []byte, binaryName string) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() { _ = gzReader.Close() }()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed reading tar stream: %w", err)
		}

		baseName := filepath.Base(header.Name)
		if baseName == binaryName && header.Typeflag == tar.TypeReg {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, tarReader); err != nil {
				return nil, fmt.Errorf("failed extracting binary from tar: %w", err)
			}
			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("binary %q not found in tar.gz archive", binaryName)
}

func extractFromZip(data []byte, binaryName string) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	for _, file := range zipReader.File {
		baseName := filepath.Base(file.Name)
		if baseName == binaryName && !file.FileInfo().IsDir() {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file inside zip: %w", err)
			}
			defer func() { _ = rc.Close() }()

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, rc); err != nil {
				return nil, fmt.Errorf("failed extracting binary from zip: %w", err)
			}
			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("binary %q not found in zip archive", binaryName)
}
