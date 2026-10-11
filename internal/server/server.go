package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/yundream/goslide/internal/parser"
	htmlrenderer "github.com/yundream/goslide/internal/renderer/html"
)

// Config holds options for running the development server.
type Config struct {
	Port         int
	Bind         string
	MarkdownPath string
	Theme        string
	ThemePath    string
	Lang         string
	OpenBrowser  bool
	Debounce     time.Duration
	OnReady      func()
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
	if cfg.Lang != "" {
		opts = append(opts, htmlrenderer.WithLang(cfg.Lang))
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
	if s.cfg.OnReady != nil {
		s.cfg.OnReady()
	}

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
