package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAssertGolden(t *testing.T) {
	tempDir := t.TempDir()
	goldenPath := filepath.Join(tempDir, "test.golden")
	content := []byte("hello golden test\nsecond line\n")

	// 1. Manually write golden file
	if err := os.WriteFile(goldenPath, content, 0644); err != nil {
		t.Fatalf("failed to write test golden: %v", err)
	}

	// 2. Exact match should pass
	AssertGolden(t, content, goldenPath)

	// 3. Test -update mode
	*updateGolden = true
	defer func() { *updateGolden = false }()

	newContent := []byte("updated content\n")
	AssertGolden(t, newContent, goldenPath)

	readBack, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read back golden: %v", err)
	}
	if string(readBack) != string(newContent) {
		t.Errorf("expected updated content %q, got %q", string(newContent), string(readBack))
	}
}
