package server

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/yundream/goslide/internal/i18n"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/parser"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
)

const sseScript = `
  <!-- Goslide Live Reload & Error Overlay Engine -->
  <script id="goslide-live-reload">
  (function() {
    if (!window.EventSource) return;
    const es = new EventSource('/events');

    function removeOverlay() {
      var el = document.getElementById('goslide-error-overlay');
      if (el) el.remove();
    }

    function removeToast() {
      var el = document.getElementById('goslide-warning-toast');
      if (el) el.remove();
    }

    function showWarningToast(diagnostics) {
      removeToast();
      if (!diagnostics || !diagnostics.length) return;

      var toast = document.createElement('div');
      toast.id = 'goslide-warning-toast';
      toast.style.cssText = 'position:fixed;top:20px;right:20px;max-width:440px;width:calc(100vw - 40px);z-index:999998;background:rgba(24, 24, 27, 0.96);backdrop-filter:blur(10px);-webkit-backdrop-filter:blur(10px);border:1.5px solid #f59e0b;border-radius:8px;box-shadow:0 15px 30px rgba(0,0,0,0.5);padding:14px 16px;color:#f3f4f6;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;font-size:13px;line-height:1.4;box-sizing:border-box;transition:opacity 0.3s ease, transform 0.3s ease;';

      var items = diagnostics.map(function(d) {
        var avail = (d.candidates && d.candidates.length) ? '<div style="margin-top:4px;color:#fbbf24;font-size:12px;">💡 Available: <code>' + d.candidates.join(', ') + '</code></div>' : '';
        return '<div style="margin-bottom:8px;padding-bottom:8px;border-bottom:1px solid rgba(245,158,11,0.2);">' +
          '<div style="font-weight:600;color:#fde68a;">Slide ' + d.slideIndex + ' (Line ' + d.line + ')</div>' +
          '<div style="color:#d1d5db;margin-top:2px;">' + d.message + '</div>' +
          avail +
          '</div>';
      }).join('');

      var content = '<div style="display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:8px;">' +
        '<div style="display:flex;align-items:center;gap:6px;font-weight:700;color:#fbbf24;font-size:14px;">' +
        '<span>⚠️</span> <span>DSL Warning</span>' +
        '</div>' +
        '<button id="goslide-toast-close" style="background:transparent;border:none;color:#9ca3af;font-size:18px;cursor:pointer;line-height:1;" title="Dismiss">&times;</button>' +
        '</div>' +
        '<div style="max-height:220px;overflow-y:auto;">' + items + '</div>';

      toast.innerHTML = content;
      document.body.appendChild(toast);

      var closeBtn = document.getElementById('goslide-toast-close');
      if (closeBtn) {
        closeBtn.onclick = function() { removeToast(); };
      }

      setTimeout(function() {
        if (toast.parentNode) {
          toast.style.opacity = '0';
          toast.style.transform = 'translateY(-10px)';
          setTimeout(removeToast, 350);
        }
      }, 6000);
    }

    function showOverlay(title, message, file) {
      removeOverlay();
      var overlay = document.createElement('div');
      overlay.id = 'goslide-error-overlay';
      overlay.style.cssText = 'position:fixed;top:20px;left:50%;transform:translateX(-50%);width:min(900px, 92vw);z-index:999999;background:rgba(24, 24, 27, 0.95);backdrop-filter:blur(12px);-webkit-backdrop-filter:blur(12px);border:1.5px solid #ef4444;border-radius:10px;box-shadow:0 25px 50px -12px rgba(0,0,0,0.7);padding:20px;color:#f3f4f6;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;box-sizing:border-box;';

      var header = '<div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:12px;border-bottom:1px solid rgba(239,68,68,0.3);padding-bottom:10px;">' +
        '<div style="display:flex;align-items:center;gap:10px;">' +
        '<span style="font-size:20px;">⚠️</span>' +
        '<span style="font-weight:700;font-size:16px;color:#fca5a5;">' + (title || 'Build Error') + '</span>' +
        (file ? '<span style="font-size:12px;background:rgba(239,68,68,0.2);color:#fca5a5;padding:2px 8px;border-radius:4px;font-family:monospace;">' + file + '</span>' : '') +
        '</div>' +
        '<button id="goslide-err-close" style="background:transparent;border:none;color:#9ca3af;font-size:22px;cursor:pointer;line-height:1;padding:2px 6px;border-radius:4px;" title="Dismiss">&times;</button>' +
        '</div>';

      var body = '<pre style="margin:0;padding:14px;background:rgba(0,0,0,0.55);border-radius:6px;color:#f87171;font-family:ui-monospace,Menlo,Monaco,Consolas,monospace;font-size:13px;line-height:1.5;overflow-x:auto;max-height:45vh;white-space:pre-wrap;word-break:break-word;">' + (message || 'Unknown error occurred') + '</pre>';

      var footer = '<div style="margin-top:12px;font-size:12px;color:#9ca3af;display:flex;align-items:center;justify-content:space-between;">' +
        '<span>💡 Fix the syntax error in your editor and save to resume live preview.</span>' +
        '<span style="font-size:11px;color:#6b7280;">Goslide Live Resilience Engine</span>' +
        '</div>';

      overlay.innerHTML = header + body + footer;
      document.body.appendChild(overlay);

      var closeBtn = document.getElementById('goslide-err-close');
      if (closeBtn) {
        closeBtn.onclick = function() { removeOverlay(); };
      }
    }

    es.onmessage = function(e) {
      var data = e.data;
      if (!data) return;
      var payload;
      try {
        payload = JSON.parse(data);
      } catch (err) {
        if (data === 'reload') {
          payload = { type: 'reload' };
        }
      }

      if (!payload) return;

      if (payload.type === 'reload') {
        removeOverlay();
        console.log('[goslide] Live reload triggered');
        location.reload();
      } else if (payload.type === 'patch') {
        removeOverlay();
        var patches = payload.patches || [];
        patches.forEach(function(patch) {
          var oldEl = document.querySelector('.slide-card[data-slide="' + patch.index + '"]');
          if (!oldEl) return;

          var temp = document.createElement('div');
          temp.innerHTML = patch.html.trim();
          var newEl = temp.firstElementChild;
          if (!newEl) return;

          if (oldEl.classList.contains('active')) {
            newEl.classList.add('active');
          }

          newEl.style.transition = 'opacity 0.15s ease-in-out';
          newEl.style.opacity = '0.7';

          oldEl.replaceWith(newEl);
          requestAnimationFrame(function() {
            newEl.style.opacity = '1';
          });

          // Notify Svelte DeckStore to synchronize slide DOM element reference
          window.dispatchEvent(new CustomEvent('goslide:slide-patched', {
            detail: { index: patch.index - 1, element: newEl }
          }));

          if (typeof renderMathInElement === 'function') {
            renderMathInElement(newEl, {
              delimiters: [
                {left: '$$', right: '$$', display: true},
                {left: '$', right: '$', display: false}
              ],
              throwOnError: false
            });
          }

          if (typeof mermaid !== 'undefined' && newEl.querySelector('.mermaid')) {
            mermaid.run({ nodes: newEl.querySelectorAll('.mermaid') });
          }
        });
      } else if (payload.type === 'error') {
        console.error('[goslide] Build error:', payload.message);
        showOverlay(payload.title, payload.message, payload.file);
      } else if (payload.type === 'warning') {
        console.warn('[goslide] DSL warnings:', payload.diagnostics);
        showWarningToast(payload.diagnostics);
      }
    };
    es.onerror = function() {
      // EventSource automatically reconnects upon disconnection
    };
  })();
  </script>
`

// EventType defines the category of SSE message sent to clients.
type EventType string

const (
	EventReload  EventType = "reload"
	EventError   EventType = "error"
	EventPatch   EventType = "patch"
	EventWarning EventType = "warning"
)

// SSEMessage represents a structured payload broadcasted via Server-Sent Events.
type SSEMessage struct {
	Type        EventType          `json:"type"`
	Title       string             `json:"title,omitempty"`
	Message     string             `json:"message,omitempty"`
	File        string             `json:"file,omitempty"`
	Patches     []SlidePatch       `json:"patches,omitempty"`
	Diagnostics []model.Diagnostic `json:"diagnostics,omitempty"`
}

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

	lastSnapshotMu sync.RWMutex
	lastSnapshot   *DeckSnapshot

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

func injectSSEScript(rendered []byte) []byte {
	targetTag := []byte("</body>")
	if idx := bytes.LastIndex(rendered, targetTag); idx != -1 {
		var injected bytes.Buffer
		injected.Write(rendered[:idx])
		injected.WriteString(sseScript)
		injected.Write(rendered[idx:])
		return injected.Bytes()
	}
	return append(rendered, []byte(sseScript)...)
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

// BroadcastMessage sends a structured SSEMessage (marshaled to JSON) to all connected SSE clients.
func (s *Server) BroadcastMessage(msg SSEMessage) {
	dataBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[goslide] Failed to marshal SSE message: %v", err)
		return
	}
	payload := string(dataBytes)

	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	for ch := range s.clients {
		select {
		case ch <- payload:
		default:
			// Client buffer full or slow; skip without blocking
		}
	}
}

// BroadcastReload sends a reload command to all connected SSE clients.
func (s *Server) BroadcastReload() {
	s.BroadcastMessage(SSEMessage{Type: EventReload})
}

// BroadcastError sends an error payload to all connected SSE clients to display an error overlay.
func (s *Server) BroadcastError(title, msg, file string) {
	s.BroadcastMessage(SSEMessage{
		Type:    EventError,
		Title:   title,
		Message: msg,
		File:    file,
	})
}

// BroadcastPatches sends incremental slide patches to all connected SSE clients.
func (s *Server) BroadcastPatches(patches []SlidePatch) {
	s.BroadcastMessage(SSEMessage{
		Type:    EventPatch,
		Patches: patches,
	})
}

// BroadcastWarnings sends non-blocking DSL diagnostics as floating warning toasts to connected SSE clients.
func (s *Server) BroadcastWarnings(diagnostics []model.Diagnostic) {
	if len(diagnostics) == 0 {
		return
	}
	s.BroadcastMessage(SSEMessage{
		Type:        EventWarning,
		Diagnostics: diagnostics,
	})
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
	go s.watchLoop(ctx)

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

func (s *Server) watchLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.watcher.Events():
			s.handleFileChange(ctx)
		case err, ok := <-s.watcher.Errors():
			if !ok {
				return
			}
			log.Printf("[goslide] Watcher error: %v", err)
		}
	}
}

func (s *Server) handleFileChange(ctx context.Context) {
	start := time.Now()

	deck, err := s.rebuildAndWarmCache(ctx)
	if err != nil {
		return
	}

	diff := s.reconcileSnapshot(deck)
	s.dispatchDiff(ctx, deck, diff, start)
}

func (s *Server) rebuildAndWarmCache(ctx context.Context) (*model.Deck, error) {
	raw, readErr := os.ReadFile(s.absFile)
	if readErr != nil {
		log.Printf("[goslide] Warning: failed to read file on change: %v", readErr)
		return nil, readErr
	}

	deck, parseErr := s.parser.Parse(ctx, bytes.NewReader(raw))
	if parseErr != nil {
		log.Printf("[goslide] Warning: parse failed on file change: %v", parseErr)
		s.BroadcastError("Markdown Syntax Error", parseErr.Error(), filepath.Base(s.absFile))
		return nil, parseErr
	}

	s.logDiagnostics(deck.Diagnostics)
	s.BroadcastWarnings(deck.Diagnostics)

	// Pre-render full HTML to warm cache and update lastHTML
	var fullBuf bytes.Buffer
	if renderErr := s.renderer.Render(ctx, deck, &fullBuf); renderErr != nil {
		log.Printf("[goslide] Warning: render failed on file change: %v", renderErr)
		s.BroadcastError("Render Error", renderErr.Error(), filepath.Base(s.absFile))
		return nil, renderErr
	}
	renderedBytes := injectSSEScript(fullBuf.Bytes())
	s.lastHTMLMu.Lock()
	s.lastHTML = renderedBytes
	s.lastHTMLMu.Unlock()

	return deck, nil
}

func (s *Server) reconcileSnapshot(deck *model.Deck) DiffResult {
	s.lastSnapshotMu.RLock()
	oldSnap := s.lastSnapshot
	s.lastSnapshotMu.RUnlock()

	diff := CompareSnapshots(oldSnap, deck)

	s.lastSnapshotMu.Lock()
	s.lastSnapshot = NewDeckSnapshot(deck)
	s.lastSnapshotMu.Unlock()

	return diff
}

func (s *Server) dispatchDiff(ctx context.Context, deck *model.Deck, diff DiffResult, start time.Time) {
	bundle, _ := i18n.GetDefaultBundle()

	if diff.NeedsReload {
		reason := diff.ReloadReasonKey
		if bundle != nil {
			reason = bundle.T(diff.ReloadReasonKey)
		}
		log.Printf("[goslide] Rebuild completed in %v. Full reload triggered (%s)...", time.Since(start), reason)
		s.BroadcastReload()
		return
	}

	if len(diff.ChangedIndices) > 0 {
		s.dispatchSlidePatches(ctx, deck, diff.ChangedIndices, start, bundle)
		return
	}

	log.Printf("[goslide] File saved with no slide content changes (%v)", time.Since(start))
}

func (s *Server) dispatchSlidePatches(ctx context.Context, deck *model.Deck, indices []int, start time.Time, bundle *i18n.Bundle) {
	var patches []SlidePatch
	for _, idx := range indices {
		var patchBuf bytes.Buffer
		if err := s.renderer.RenderSlide(ctx, deck, idx, &patchBuf); err == nil {
			patches = append(patches, SlidePatch{
				Index: idx + 1,
				HTML:  patchBuf.String(),
			})
		} else {
			log.Printf("[goslide] Warning: failed to render slide patch %d: %v", idx+1, err)
		}
	}
	if len(patches) > 0 {
		if bundle != nil {
			log.Print(bundle.T("server.hmr.patch.log", patches[0].Index, time.Since(start)))
		} else {
			log.Printf("[goslide] Slide %d incrementally patched in %v", patches[0].Index, time.Since(start))
		}
		s.BroadcastPatches(patches)
	}
}

func (s *Server) logDiagnostics(diagnostics []model.Diagnostic) {
	for _, diag := range diagnostics {
		if len(diag.Candidates) > 0 {
			log.Printf("[goslide] ⚠️  Slide %d (Line %d): %s (Raw: %q, Available: %v)",
				diag.SlideIndex, diag.Line, diag.Message, diag.RawSnippet, diag.Candidates)
		} else {
			log.Printf("[goslide] ⚠️  Slide %d (Line %d): %s (Raw: %q)",
				diag.SlideIndex, diag.Line, diag.Message, diag.RawSnippet)
		}
	}
}


