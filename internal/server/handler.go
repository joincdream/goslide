package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler returns the HTTP handler used by the server.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", s.handleEvents)
	mux.HandleFunc("/", s.handleRoot)
	return mux
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	// Root or index.html requests are served with rendered slide HTML
	path := r.URL.Path
	if path == "/" || path == "/index.html" {
		s.serveSlideHTML(w, r)
		return
	}

	// All other requests serve relative assets (e.g. ./images/diagram.png)
	s.serveStaticAsset(w, r)
}

func (s *Server) serveSlideHTML(w http.ResponseWriter, r *http.Request) {
	htmlBytes, err := s.renderMarkdown(r.Context())
	if err != nil {
		s.lastHTMLMu.RLock()
		cached := s.lastHTML
		s.lastHTMLMu.RUnlock()

		if len(cached) > 0 {
			log.Printf("[goslide] Markdown parse error (%v), serving last valid presentation", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(cached)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to render markdown slide:\n\n%v", err), http.StatusInternalServerError)
		return
	}

	s.lastHTMLMu.Lock()
	s.lastHTML = htmlBytes
	s.lastHTMLMu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(htmlBytes)
}

func (s *Server) serveStaticAsset(w http.ResponseWriter, r *http.Request) {
	cleanPath := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	targetPath := filepath.Join(s.baseDir, cleanPath)

	// Prevent directory traversal outside of baseDir
	rel, err := filepath.Rel(s.baseDir, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
		http.NotFound(w, r)
		return
	}

	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, targetPath)
}

const sseLiveReloadMarker = `<script id="goslide-live-reload">/* Goslide Live Reload & HMR is managed by Svelte core */</script>`

// injectSSEScript injects the live reload marker into HTML before </body>.
// Note: Detailed HMR and overlay UI logic is handled natively in the Svelte frontend bundle.
func injectSSEScript(rendered []byte) []byte {
	targetTag := []byte("</body>")
	if idx := bytes.LastIndex(rendered, targetTag); idx != -1 {
		var injected bytes.Buffer
		injected.Write(rendered[:idx])
		injected.WriteString(sseLiveReloadMarker)
		injected.Write(rendered[idx:])
		return injected.Bytes()
	}
	return append(rendered, []byte(sseLiveReloadMarker)...)
}

func (s *Server) renderMarkdown(ctx context.Context) ([]byte, error) {
	raw, err := os.ReadFile(s.absFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read markdown file: %w", err)
	}

	deck, err := s.parser.Parse(ctx, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}

	s.lastSnapshotMu.Lock()
	if s.lastSnapshot == nil {
		s.lastSnapshot = NewDeckSnapshot(deck)
	}
	s.lastSnapshotMu.Unlock()

	var buf bytes.Buffer
	if err := s.renderer.Render(ctx, deck, &buf); err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return injectSSEScript(buf.Bytes()), nil
}
