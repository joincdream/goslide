package server

import (
	"context"
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

func TestWatcher_FileModificationAndDebounce(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "presentation.md")
	if err := os.WriteFile(testFile, []byte("# Slide 1"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	debounce := 40 * time.Millisecond
	w, err := NewWatcher(testFile, debounce)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer func() { _ = w.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Start(ctx)

	// Rapid modifications within debounce window
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Millisecond)
		if err := os.WriteFile(testFile, []byte("# Updated Slide "+string(rune('A'+i))), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}
	}

	select {
	case <-w.Events():
		// Received event as expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for reload event")
	}

	// Ensure no extra duplicate events immediately queued
	select {
	case <-w.Events():
		t.Fatal("unexpected duplicate event received after debounce window")
	case <-time.After(100 * time.Millisecond):
		// Success: debounced into a single event
	}
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
