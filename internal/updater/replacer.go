package updater

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrPermissionDenied is returned when writing to the target executable path fails due to permissions.
var ErrPermissionDenied = errors.New("permission denied")

// ReplaceBinary atomically replaces the executable at targetPath with newBinaryData.
// If targetPath is empty, it resolves os.Executable().
func ReplaceBinary(targetPath string, newBinaryData []byte) error {
	resolvedPath, err := resolveTargetExecutable(targetPath)
	if err != nil {
		return err
	}

	fileMode, err := getTargetFileMode(resolvedPath)
	if err != nil {
		return err
	}

	targetDir := filepath.Dir(resolvedPath)
	tmpPath, err := writeTempBinary(targetDir, newBinaryData, fileMode)
	if err != nil {
		return err
	}
	defer func() {
		if tmpPath != "" {
			_ = os.Remove(tmpPath)
		}
	}()

	backupPath, hasBackup, err := backupExisting(resolvedPath, targetDir)
	if err != nil {
		return err
	}
	defer func() {
		if hasBackup && backupPath != "" {
			_ = os.Remove(backupPath)
		}
	}()

	if err := commitReplacement(tmpPath, resolvedPath, backupPath, hasBackup); err != nil {
		return err
	}

	tmpPath = ""
	return nil
}

func resolveTargetExecutable(targetPath string) (string, error) {
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

func writeTempBinary(dir string, data []byte, mode os.FileMode) (string, error) {
	tmpFile, err := os.CreateTemp(dir, "goslide-new-*.tmp")
	if err != nil {
		if os.IsPermission(err) {
			return "", fmt.Errorf("%w: cannot create temporary file in %s", ErrPermissionDenied, dir)
		}
		return "", fmt.Errorf("failed to create temporary file in %s: %w", dir, err)
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("failed to write new binary data: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("failed to close temporary binary file: %w", err)
	}

	if err := os.Chmod(tmpPath, mode); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("failed to set permissions on %s: %w", tmpPath, err)
	}

	return tmpPath, nil
}

func backupExisting(targetPath, dir string) (string, bool, error) {
	if _, err := os.Stat(targetPath); err != nil {
		return "", false, nil
	}

	backupPath := filepath.Join(dir, fmt.Sprintf("goslide-old-%d.tmp", os.Getpid()))
	if err := os.Rename(targetPath, backupPath); err != nil {
		if os.IsPermission(err) {
			return "", false, fmt.Errorf("%w: cannot rename existing binary %s", ErrPermissionDenied, targetPath)
		}
		return "", false, fmt.Errorf("failed to move existing binary to backup %s: %w", backupPath, err)
	}
	return backupPath, true, nil
}

func commitReplacement(tmpPath, targetPath, backupPath string, hasBackup bool) error {
	if err := os.Rename(tmpPath, targetPath); err != nil {
		if hasBackup {
			_ = os.Rename(backupPath, targetPath)
		}
		return fmt.Errorf("failed to install new binary to %s: %w", targetPath, err)
	}
	return nil
}
