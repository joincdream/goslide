package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors markdown files and asset directories for changes with debouncing.
type Watcher struct {
	targetFile string
	targetDir  string
	debounce   time.Duration

	fswatcher *fsnotify.Watcher
	events    chan struct{}
	errors    chan error
	closeOnce sync.Once
	done      chan struct{}
}

// NewWatcher creates a new Watcher for the given target markdown file path.
func NewWatcher(targetFile string, debounce time.Duration) (*Watcher, error) {
	absFile, err := filepath.Abs(targetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute file path: %w", err)
	}

	info, err := os.Stat(absFile)
	if err != nil {
		return nil, fmt.Errorf("failed to stat target file: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("target path %q is a directory, expected a file", targetFile)
	}

	absDir := filepath.Dir(absFile)

	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	// Watch the target directory to handle atomic saves (write-and-rename) and asset modifications
	if err := fsw.Add(absDir); err != nil {
		_ = fsw.Close()
		return nil, fmt.Errorf("failed to watch directory %q: %w", absDir, err)
	}

	w := &Watcher{
		targetFile: absFile,
		targetDir:  absDir,
		debounce:   debounce,
		fswatcher:  fsw,
		events:     make(chan struct{}, 1),
		errors:     make(chan error, 1),
		done:       make(chan struct{}),
	}

	return w, nil
}

// Events returns the receive-only channel that emits notifications when relevant changes occur.
func (w *Watcher) Events() <-chan struct{} {
	return w.events
}

// Errors returns the receive-only channel for watching errors.
func (w *Watcher) Errors() <-chan error {
	return w.errors
}

// Start runs the event listening loop in background until context is canceled or Close is called.
func (w *Watcher) Start(ctx context.Context) {
	go w.loop(ctx)
}

type debouncer struct {
	mu       sync.Mutex
	timer    *time.Timer
	duration time.Duration
}

func newDebouncer(d time.Duration) *debouncer {
	return &debouncer{duration: d}
}

func (d *debouncer) trigger(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.duration, fn)
}

func (d *debouncer) stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
}

func (w *Watcher) loop(ctx context.Context) {
	db := newDebouncer(w.debounce)
	defer db.stop()

	onTimeout := func() {
		select {
		case w.events <- struct{}{}:
		case <-w.done:
		case <-ctx.Done():
		default:
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.done:
			return
		case event, ok := <-w.fswatcher.Events:
			if !ok {
				return
			}
			if w.isRelevantEvent(event) {
				db.trigger(onTimeout)
			}
		case err, ok := <-w.fswatcher.Errors:
			if !ok || !w.dispatchError(ctx, err) {
				return
			}
		}
	}
}

func (w *Watcher) dispatchError(ctx context.Context, err error) bool {
	select {
	case w.errors <- err:
		return true
	case <-w.done:
		return false
	case <-ctx.Done():
		return false
	default:
		return true
	}
}

// isRelevantEvent determines whether a filesystem event should trigger a reload.
func (w *Watcher) isRelevantEvent(event fsnotify.Event) bool {
	// Only respond to Write, Create, Rename, Remove
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Remove) == 0 {
		return false
	}

	baseName := filepath.Base(event.Name)

	// Ignore hidden files and vim/emacs temporary backup files
	if strings.HasPrefix(baseName, ".") || strings.HasSuffix(baseName, "~") || strings.HasSuffix(baseName, ".swp") {
		return false
	}

	absPath, err := filepath.Abs(event.Name)
	if err != nil {
		return false
	}

	// Direct match with the target markdown file
	if absPath == w.targetFile {
		return true
	}

	// Or any asset file inside target directory with common web extensions
	ext := strings.ToLower(filepath.Ext(baseName))
	switch ext {
	case ".md", ".markdown", ".css", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".json":
		return true
	}

	return false
}

// Close terminates the filesystem watcher and cleans up resources.
func (w *Watcher) Close() error {
	var err error
	w.closeOnce.Do(func() {
		close(w.done)
		err = w.fswatcher.Close()
	})
	return err
}
