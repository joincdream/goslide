package server

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestNewWatcher_Validation(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Non-existent file
	_, err := NewWatcher(filepath.Join(tempDir, "non_existent.md"), 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}

	// 2. Directory instead of file
	_, err = NewWatcher(tempDir, 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error when target is directory, got nil")
	}

	// 3. Valid file
	testFile := filepath.Join(tempDir, "slides.md")
	if writeErr := os.WriteFile(testFile, []byte("# Test Slide"), 0644); writeErr != nil {
		t.Fatalf("failed to create test file: %v", writeErr)
	}

	w, err := NewWatcher(testFile, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error for valid file: %v", err)
	}
	defer func() { _ = w.Close() }()
}

func TestWatcher_FilterIrrelevantFiles(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "talk.md")
	if err := os.WriteFile(testFile, []byte("# Slide"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	w, err := NewWatcher(testFile, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer func() { _ = w.Close() }()

	// Irrelevant: vim swp file
	swpEvent := fsnotify.Event{
		Name: filepath.Join(tempDir, ".talk.md.swp"),
		Op:   fsnotify.Create,
	}
	if w.isRelevantEvent(swpEvent) {
		t.Errorf("expected swp file to be filtered out")
	}

	// Irrelevant: backup file
	backupEvent := fsnotify.Event{
		Name: filepath.Join(tempDir, "talk.md~"),
		Op:   fsnotify.Write,
	}
	if w.isRelevantEvent(backupEvent) {
		t.Errorf("expected backup file to be filtered out")
	}

	// Relevant: image asset update
	imgEvent := fsnotify.Event{
		Name: filepath.Join(tempDir, "diagram.png"),
		Op:   fsnotify.Write,
	}
	if !w.isRelevantEvent(imgEvent) {
		t.Errorf("expected image asset event to be relevant")
	}

	// Relevant: target markdown update
	mdEvent := fsnotify.Event{
		Name: testFile,
		Op:   fsnotify.Write,
	}
	if !w.isRelevantEvent(mdEvent) {
		t.Errorf("expected target markdown event to be relevant")
	}
}
