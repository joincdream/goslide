package updater

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/minio/selfupdate"
)

// ErrPermissionDenied is returned when writing to the target executable path fails due to permissions.
var ErrPermissionDenied = errors.New("permission denied")

// ReplaceBinary atomically replaces the executable at targetPath with newBinaryData using minio/selfupdate.
// If targetPath is empty, it resolves os.Executable().
func ReplaceBinary(targetPath string, newBinaryData []byte) error {
	resolvedPath, err := resolveTargetPath(targetPath)
	if err != nil {
		return err
	}

	fileMode, err := getTargetFileMode(resolvedPath)
	if err != nil {
		return err
	}

	opts := selfupdate.Options{
		TargetPath: resolvedPath,
		TargetMode: fileMode,
	}

	if err := selfupdate.Apply(bytes.NewReader(newBinaryData), opts); err != nil {
		if os.IsPermission(err) || errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("%w: %v", ErrPermissionDenied, err)
		}
		return fmt.Errorf("failed to apply binary update to %s: %w", resolvedPath, err)
	}

	return nil
}

func resolveTargetPath(targetPath string) (string, error) {
	if targetPath != "" {
		return targetPath, nil
	}
	exec, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to determine executable path: %w", err)
	}
	if realPath, err := filepath.EvalSymlinks(exec); err == nil && realPath != "" {
		return realPath, nil
	}
	return exec, nil
}

func getTargetFileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.Mode().Perm(), nil
	}
	if os.IsPermission(err) {
		return 0, fmt.Errorf("%w: cannot access %s", ErrPermissionDenied, path)
	}
	return 0755, nil
}
