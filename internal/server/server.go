package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/yundream/goslide/internal/parser"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
)

const sseScript = `
  <!-- Goslide Live Reload Engine -->
  <script id="goslide-live-reload">
  (function() {
    if (!window.EventSource) return;
    const es = new EventSource('/events');
    es.onmessage = function(e) {
      if (e.data === 'reload') {
        console.log('[goslide] Live reload triggered');
        location.reload();
      }
    };
    es.onerror = function() {
      // EventSource automatically reconnects upon disconnection
    };
  })();
  </script>
`

// Config holds options for running the development server.
type Config struct {
	Port         int
	Bind         string
	MarkdownPath string
	Theme        string
	ThemePath    string
	OpenBrowser  bool
	Debounce     time.Duration
}

// Server provides local HTTP hosting, asset serving, and SSE live reload broadcasting.
type Server struct {
	cfg      Config
	absFile  string
	baseDir  string
	parser   *parser.Parser
	renderer *htmlrenderer.HTMLRenderer
	watcher  *Watcher

	clientsMu sync.RWMutex
	clients   map[chan string]struct{}

	lastHTMLMu sync.RWMutex
	lastHTML   []byte

	httpServer *http.Server
	listener   net.Listener
}

// NewServer initializes a new development Server with the provided configuration.
func NewServer(cfg Config) (*Server, error) {
	if cfg.Port <= 0 {
		cfg.Port = 8080
	}
	if cfg.Bind == "" {
		cfg.Bind = "localhost"
	}
	if cfg.Debounce <= 0 {
		cfg.Debounce = 50 * time.Millisecond
	}

	absFile, err := filepath.Abs(cfg.MarkdownPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve markdown file path: %w", err)
	}

	if _, statErr := os.Stat(absFile); statErr != nil {
		return nil, fmt.Errorf("markdown file not found: %w", statErr)
	}

	baseDir := filepath.Dir(absFile)

	w, err := NewWatcher(absFile, cfg.Debounce)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file watcher: %w", err)
	}

	var opts []htmlrenderer.Option
	if cfg.Theme != "" {
		opts = append(opts, htmlrenderer.WithTheme(cfg.Theme))
	}
	if cfg.ThemePath != "" {
		opts = append(opts, htmlrenderer.WithCustomCSS(cfg.ThemePath))
	}
	opts = append(opts, htmlrenderer.WithBaseDir(baseDir))

	s := &Server{
		cfg:      cfg,
		absFile:  absFile,
		baseDir:  baseDir,
		parser:   parser.NewParser(),
		renderer: htmlrenderer.NewRenderer(opts...),
		watcher:  w,
		clients:  make(map[chan string]struct{}),
	}

	return s, nil
}

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

	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, targetPath)
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

	var buf bytes.Buffer
	if err := s.renderer.Render(ctx, deck, &buf); err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	rendered := buf.Bytes()
	// Inject SSE live reload script right before </body>
	targetTag := []byte("</body>")
	if idx := bytes.LastIndex(rendered, targetTag); idx != -1 {
		var injected bytes.Buffer
		injected.Write(rendered[:idx])
		injected.WriteString(sseScript)
		injected.Write(rendered[idx:])
		return injected.Bytes(), nil
	}

	// Fallback if </body> is missing
	return append(rendered, []byte(sseScript)...), nil
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	clientChan := make(chan string, 4)

	s.clientsMu.Lock()
	s.clients[clientChan] = struct{}{}
	s.clientsMu.Unlock()

	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, clientChan)
		close(clientChan)
		s.clientsMu.Unlock()
	}()

	// Send initial handshake comment
	_, _ = fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, ok := <-clientChan:
			if !ok {
				return
			}
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// BroadcastReload sends a reload command to all connected SSE clients.
func (s *Server) BroadcastReload() {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	for ch := range s.clients {
		select {
		case ch <- "reload":
		default:
			// Client buffer full or slow; skip without blocking
		}
	}
}

// URL returns the addressable base URL for this server (e.g. http://localhost:8080).
func (s *Server) URL() string {
	if s.listener != nil {
		addr := s.listener.Addr().String()
		return fmt.Sprintf("http://%s", addr)
	}
	return fmt.Sprintf("http://%s:%d", s.cfg.Bind, s.cfg.Port)
}

// Start begins listening and serving HTTP traffic. It blocks until context is canceled or an error occurs.
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Bind, s.cfg.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind address %s: %w", addr, err)
	}
	s.listener = ln

	s.httpServer = &http.Server{
		Handler: s.Handler(),
	}

	s.watcher.Start(ctx)

	// Watcher loop: notify SSE clients upon changes
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.watcher.Events():
				start := time.Now()
				// Pre-render in background to warm cache and log latency
				_, renderErr := s.renderMarkdown(ctx)
				if renderErr != nil {
					log.Printf("[goslide] Warning: rebuild failed on file change: %v", renderErr)
				} else {
					log.Printf("[goslide] Rebuild completed in %v. Broadcasting reload...", time.Since(start))
				}
				s.BroadcastReload()
			case err, ok := <-s.watcher.Errors():
				if !ok {
					return
				}
				log.Printf("[goslide] Watcher error: %v", err)
			}
		}
	}()

	if s.cfg.OpenBrowser {
		go func() {
			time.Sleep(100 * time.Millisecond)
			s.openBrowser(s.URL())
		}()
	}

	errChan := make(chan error, 1)
	go func() {
		if serveErr := s.httpServer.Serve(s.listener); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errChan <- serveErr
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.watcher.Close()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errChan:
		_ = s.watcher.Close()
		return err
	}
}

// Close terminates the server immediately.
func (s *Server) Close() error {
	_ = s.watcher.Close()
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) openBrowser(targetURL string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", targetURL)
	case "darwin":
		cmd = exec.Command("open", targetURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
	default:
		return
	}
	_ = cmd.Start()
}
