package testutil

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden files with actual test outputs")

// AssertGolden compares actual bytes against the golden file at goldenPath.
// If the -update flag is passed, it writes actual bytes to goldenPath instead.
func AssertGolden(t *testing.T, actual []byte, goldenPath string) {
	t.Helper()

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
			t.Fatalf("failed to create golden file directory: %v", err)
		}
		if err := os.WriteFile(goldenPath, actual, 0644); err != nil {
			t.Fatalf("failed to write golden file %q: %v", goldenPath, err)
		}
		t.Logf("updated golden file: %s (%d bytes)", goldenPath, len(actual))
		return
	}

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("golden file %q does not exist; run tests with -update to create it", goldenPath)
		}
		t.Fatalf("failed to read golden file %q: %v", goldenPath, err)
	}

	if !bytes.Equal(actual, expected) {
		t.Errorf("actual output does not match golden file %q:\nexpected %d bytes, got %d bytes\nrun with -update to review and accept changes",
			goldenPath, len(expected), len(actual))
		printFirstDiff(t, expected, actual)
	}
}

func printFirstDiff(t *testing.T, expected, actual []byte) {
	t.Helper()
	expLines := bytes.Split(expected, []byte("\n"))
	actLines := bytes.Split(actual, []byte("\n"))

	maxLines := len(expLines)
	if len(actLines) < maxLines {
		maxLines = len(actLines)
	}

	for i := 0; i < maxLines; i++ {
		if !bytes.Equal(expLines[i], actLines[i]) {
			t.Errorf("first diff at line %d:\n  exp: %s\n  got: %s", i+1, string(expLines[i]), string(actLines[i]))
			return
		}
	}
	if len(expLines) != len(actLines) {
		t.Errorf("line count mismatch: exp %d lines, got %d lines", len(expLines), len(actLines))
	}
}
