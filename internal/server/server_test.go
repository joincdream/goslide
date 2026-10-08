package server

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yundream/goslide/internal/model"
)

func createSampleDeck(t *testing.T, dir string) string {
	t.Helper()
	mdPath := filepath.Join(dir, "presentation.md")
	content := `---
title: "Live Preview Test"
theme: "clean"
---

# Slide 1: Welcome
Live preview test content.

---

# Slide 2: Details
More details here.
`
	if err := os.WriteFile(mdPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write sample deck: %v", err)
	}
	return mdPath
}

func TestServer_ServeSlideHTML_And_SSEInjection(t *testing.T) {
	tempDir := t.TempDir()
	mdPath := createSampleDeck(t, tempDir)

	s, err := NewServer(Config{
		MarkdownPath: mdPath,
		Port:         8080,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer func() { _ = s.Close() }()

	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	// 1. Request root
	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("failed to GET root: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	html := string(bodyBytes)

	if !strings.Contains(html, "Slide 1: Welcome") {
		t.Errorf("expected HTML to contain slide title, got:\n%s", html)
	}

	if !strings.Contains(html, `id="goslide-live-reload"`) {
		t.Errorf("expected HTML to contain injected SSE script")
	}
}

func TestServer_ServeStaticAsset(t *testing.T) {
	tempDir := t.TempDir()
	mdPath := createSampleDeck(t, tempDir)

	// Create a dummy image asset
	imgAsset := filepath.Join(tempDir, "diagram.svg")
	svgContent := `<svg width="100" height="100"><circle cx="50" cy="50" r="40" fill="red"/></svg>`
	if err := os.WriteFile(imgAsset, []byte(svgContent), 0644); err != nil {
		t.Fatalf("failed to write asset: %v", err)
	}

	s, err := NewServer(Config{
		MarkdownPath: mdPath,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer func() { _ = s.Close() }()

	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	// 1. Valid asset
	res, err := http.Get(ts.URL + "/diagram.svg")
	if err != nil {
		t.Fatalf("failed to GET static asset: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for static asset, got %d", res.StatusCode)
	}

	// 2. Non-existent asset
	res404, err := http.Get(ts.URL + "/non-existent.png")
	if err != nil {
		t.Fatalf("failed to GET non-existent asset: %v", err)
	}
	defer func() { _ = res404.Body.Close() }()

	if res404.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 NotFound, got %d", res404.StatusCode)
	}
}

func TestServer_SSE_BroadcastReload(t *testing.T) {
	tempDir := t.TempDir()
	mdPath := createSampleDeck(t, tempDir)

	s, err := NewServer(Config{
		MarkdownPath: mdPath,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer func() { _ = s.Close() }()

	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL+"/events", nil)
	if err != nil {
		t.Fatalf("failed to create SSE request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to connect to SSE stream: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for SSE, got %d", res.StatusCode)
	}

	reader := bufio.NewReader(res.Body)

	// Read handshake comment
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	if !strings.HasPrefix(line, ": connected") {
		t.Errorf("expected connected handshake, got %q", line)
	}

	// Trigger broadcast
	go func() {
		time.Sleep(50 * time.Millisecond)
		s.BroadcastReload()
	}()

	// Wait for reload message
	var receivedReload bool
	for i := 0; i < 5; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("error reading SSE stream: %v", err)
		}
		if strings.Contains(line, `"type":"reload"`) || strings.TrimSpace(line) == "data: reload" {
			receivedReload = true
			break
		}
	}

	if !receivedReload {
		t.Fatal("expected reload message from SSE stream")
	}
}

func TestServer_SSE_BroadcastError(t *testing.T) {
	tempDir := t.TempDir()
	mdPath := createSampleDeck(t, tempDir)

	s, err := NewServer(Config{
		MarkdownPath: mdPath,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer func() { _ = s.Close() }()

	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL+"/events", nil)
	if err != nil {
		t.Fatalf("failed to create SSE request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to connect to SSE stream: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	reader := bufio.NewReader(res.Body)

	// Read handshake comment
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	if !strings.HasPrefix(line, ": connected") {
		t.Errorf("expected connected handshake, got %q", line)
	}

	// Trigger error broadcast
	go func() {
		time.Sleep(50 * time.Millisecond)
		s.BroadcastError("Markdown Syntax Error", "yaml: line 4: mapping values not allowed", "test.md")
	}()

	// Wait for error message
	var receivedError bool
	for i := 0; i < 5; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("error reading SSE stream: %v", err)
		}
		if strings.Contains(line, `"type":"error"`) && strings.Contains(line, "Markdown Syntax Error") && strings.Contains(line, "test.md") {
			receivedError = true
			break
		}
	}

	if !receivedError {
		t.Fatal("expected structured error message from SSE stream")
	}
}

func TestServer_FallbackOnSyntaxError(t *testing.T) {
	tempDir := t.TempDir()
	mdPath := createSampleDeck(t, tempDir)

	s, err := NewServer(Config{
		MarkdownPath: mdPath,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer func() { _ = s.Close() }()

	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	// 1. Initial valid render warms cache
	res1, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("failed to get initial page: %v", err)
	}
	_ = res1.Body.Close()

	// 2. Corrupt markdown with invalid YAML frontmatter
	brokenContent := `---
title: [unclosed yaml bracket
---
# Corrupted slide`
	if writeErr := os.WriteFile(mdPath, []byte(brokenContent), 0644); writeErr != nil {
		t.Fatalf("failed to overwrite with broken markdown: %v", writeErr)
	}

	// 3. Second request should fallback to last valid HTML
	res2, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("failed to get page after corruption: %v", err)
	}
	defer func() { _ = res2.Body.Close() }()

	if res2.StatusCode != http.StatusOK {
		t.Errorf("expected fallback 200 OK, got %d", res2.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res2.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	html := string(bodyBytes)

	if !strings.Contains(html, "Slide 1: Welcome") {
		t.Errorf("expected cached presentation to be served during syntax error")
	}
}

func TestServer_SSE_BroadcastPatches(t *testing.T) {
	tempDir := t.TempDir()
	mdPath := createSampleDeck(t, tempDir)

	s, err := NewServer(Config{
		MarkdownPath: mdPath,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer func() { _ = s.Close() }()

	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL+"/events", nil)
	if err != nil {
		t.Fatalf("failed to create SSE request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to connect to SSE stream: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	reader := bufio.NewReader(res.Body)

	// Read handshake comment
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	if !strings.HasPrefix(line, ": connected") {
		t.Errorf("expected connected handshake, got %q", line)
	}

	// Trigger patch broadcast
	go func() {
		time.Sleep(50 * time.Millisecond)
		s.BroadcastPatches([]SlidePatch{
			{
				Index: 2,
				HTML:  `<section class="slide-card" data-slide="2">Updated Slide 2</section>`,
			},
		})
	}()

	// Wait for patch message
	var receivedPatch bool
	for i := 0; i < 5; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("error reading SSE stream: %v", err)
		}
		if strings.Contains(line, `"type":"patch"`) && strings.Contains(line, `"index":2`) && strings.Contains(line, "Updated Slide 2") {
			receivedPatch = true
			break
		}
	}

	if !receivedPatch {
		t.Fatal("expected structured patch message from SSE stream")
	}
}

func TestServer_SSE_BroadcastWarnings(t *testing.T) {
	tempDir := t.TempDir()
	mdPath := createSampleDeck(t, tempDir)

	s, err := NewServer(Config{
		MarkdownPath: mdPath,
	})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer func() { _ = s.Close() }()

	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL+"/events", nil)
	if err != nil {
		t.Fatalf("failed to create SSE request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to connect to SSE stream: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	reader := bufio.NewReader(res.Body)

	// Read handshake comment
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	if !strings.HasPrefix(line, ": connected") {
		t.Errorf("expected connected handshake, got %q", line)
	}

	// Trigger warning broadcast
	go func() {
		time.Sleep(50 * time.Millisecond)
		s.BroadcastWarnings([]model.Diagnostic{
			{
				Severity:   model.SeverityWarning,
				SlideIndex: 3,
				Line:       2,
				Rule:       "directive.missing_underscore",
				Message:    "Slide-local directive missing underscore prefix",
				RawSnippet: "<!-- class: lead -->",
				Candidates: []string{"_class"},
			},
		})
	}()

	var receivedWarning bool
	for i := 0; i < 5; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("error reading SSE stream: %v", err)
		}
		if strings.Contains(line, `"type":"warning"`) && strings.Contains(line, "directive.missing_underscore") && strings.Contains(line, "_class") {
			receivedWarning = true
			break
		}
	}

	if !receivedWarning {
		t.Fatal("expected structured warning message with candidates from SSE stream")
	}
}


