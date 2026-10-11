package updater

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReleaseNoteHighlights(t *testing.T) {
	body := `## What's Changed
* feat(update): add self-update command by @yundream in [#37](https://github.com/joincdream/goslide/pull/37)
- fix(test): resolve locale issue with **golden test**
* docs: update readme
+ chore: bump dependencies
`
	highlights := ExtractHighlights(body, 2)
	if len(highlights) != 2 {
		t.Fatalf("expected 2 highlights, got %d", len(highlights))
	}

	if !strings.Contains(highlights[0], "add self-update command") || strings.Contains(highlights[0], "https://") {
		t.Errorf("unexpected highlight 0: %s", highlights[0])
	}
	if !strings.Contains(highlights[1], "golden test") || strings.Contains(highlights[1], "**") {
		t.Errorf("unexpected highlight 1: %s", highlights[1])
	}
}

func TestVerifier(t *testing.T) {
	data := []byte("hello goslide binary content")
	hashBytes := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hashBytes[:])

	checksumsContent := fmt.Sprintf("%s  goslide_linux_amd64.tar.gz\n", hashStr)
	parsed := ParseChecksums(checksumsContent)

	if parsed["goslide_linux_amd64.tar.gz"] != hashStr {
		t.Errorf("expected hash %s, got %s", hashStr, parsed["goslide_linux_amd64.tar.gz"])
	}

	if err := VerifySHA256(data, hashStr); err != nil {
		t.Errorf("VerifySHA256 failed: %v", err)
	}

	if err := VerifySHA256(data, "invalidhash"); err == nil {
		t.Errorf("expected error on invalid hash, got nil")
	}
}

func TestCacheManager(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "update_cache.json")
	cm := NewCacheManager(cachePath)

	cached, err := cm.Load()
	if err != nil || cached != nil {
		t.Fatalf("expected nil cache on missing file, got %v, err: %v", cached, err)
	}

	now := time.Now()
	err = cm.Save(&CacheData{
		LastCheckedAt: now,
		LatestVersion: "v0.2.0",
		ReleaseURL:    "https://github.com/joincdream/goslide/releases/tag/v0.2.0",
	})
	if err != nil {
		t.Fatalf("failed to save cache: %v", err)
	}

	loaded, err := cm.Load()
	if err != nil {
		t.Fatalf("failed to load cache: %v", err)
	}
	if loaded.LatestVersion != "v0.2.0" {
		t.Errorf("expected v0.2.0, got %s", loaded.LatestVersion)
	}

	if cm.IsExpired(loaded, 1*time.Hour) {
		t.Errorf("expected cache not to be expired")
	}
	if !cm.IsExpired(loaded, -1*time.Hour) {
		t.Errorf("expected cache to be expired with negative ttl")
	}
}

func createTestTarGz(t *testing.T, binaryName string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	hdr := &tar.Header{
		Name:     binaryName,
		Mode:     0755,
		Size:     int64(len(content)),
		Typeflag: tar.TypeReg,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("failed to write tar body: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("failed to close tar: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("failed to close gzip: %v", err)
	}

	return buf.Bytes()
}

func TestUpdater_FullFlow_MockServer(t *testing.T) {
	fakeBinary := []byte("#!/bin/sh\necho 'goslide v0.2.0'\n")
	archiveBytes := createTestTarGz(t, "goslide", fakeBinary)

	hashBytes := sha256.Sum256(archiveBytes)
	hashStr := hex.EncodeToString(hashBytes[:])

	checksumsBytes := []byte(fmt.Sprintf("%s  goslide_linux_amd64.tar.gz\n", hashStr))

	// Start Mock HTTP Server
	var serverURL string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/releases/latest":
			rel := ReleaseInfo{
				TagName: "v0.2.0",
				Name:    "Goslide v0.2.0",
				HTMLURL: "https://github.com/joincdream/goslide/releases/tag/v0.2.0",
				Body:    "## Highlights\n* feat: in-app update support (#37)\n* fix: golden test",
				Assets: []Asset{
					{
						Name:               "goslide_linux_amd64.tar.gz",
						Size:               int64(len(archiveBytes)),
						BrowserDownloadURL: serverURL + "/download/goslide_linux_amd64.tar.gz",
					},
					{
						Name:               "checksums.txt",
						Size:               int64(len(checksumsBytes)),
						BrowserDownloadURL: serverURL + "/download/checksums.txt",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(rel)
		case "/download/goslide_linux_amd64.tar.gz":
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(archiveBytes)
		case "/download/checksums.txt":
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write(checksumsBytes)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()
	serverURL = ts.URL

	tmpDir := t.TempDir()
	targetBinaryPath := filepath.Join(tmpDir, "goslide")
	if err := os.WriteFile(targetBinaryPath, []byte("old-binary-content"), 0755); err != nil {
		t.Fatalf("failed to write target binary: %v", err)
	}

	cacheFile := filepath.Join(tmpDir, "cache.json")
	cacheMgr := NewCacheManager(cacheFile)

	client := NewClient(
		WithAPIURL(serverURL+"/releases/latest"),
		WithHTTPClient(ts.Client()),
	)

	updater := NewUpdater(client, cacheMgr)

	// 1. Check when already latest
	resLatest, err := updater.Execute(context.Background(), UpdateOptions{
		CurrentVersion: "v0.2.0",
		TargetOS:       "linux",
		TargetArch:     "amd64",
		TargetPath:     targetBinaryPath,
		DryRun:         false,
	})
	if err != nil {
		t.Fatalf("updater.Execute error: %v", err)
	}
	if !resLatest.AlreadyLatest {
		t.Errorf("expected AlreadyLatest to be true")
	}

	// 2. Perform actual update from v0.1.0 to v0.2.0
	resUpdate, err := updater.Execute(context.Background(), UpdateOptions{
		CurrentVersion: "v0.1.0",
		TargetOS:       "linux",
		TargetArch:     "amd64",
		TargetPath:     targetBinaryPath,
		DryRun:         false,
	})
	if err != nil {
		t.Fatalf("updater.Execute failed: %v", err)
	}

	if resUpdate.AlreadyLatest {
		t.Errorf("expected AlreadyLatest to be false")
	}
	if resUpdate.UpdatedVersion != "v0.2.0" {
		t.Errorf("expected UpdatedVersion v0.2.0, got %s", resUpdate.UpdatedVersion)
	}
	if len(resUpdate.Highlights) == 0 {
		t.Errorf("expected non-empty highlights")
	}

	// Verify target file has been updated
	installedBytes, err := os.ReadFile(targetBinaryPath)
	if err != nil {
		t.Fatalf("failed to read installed binary: %v", err)
	}
	if !bytes.Equal(installedBytes, fakeBinary) {
		t.Errorf("installed binary content mismatch: got %q, want %q", string(installedBytes), string(fakeBinary))
	}
}
