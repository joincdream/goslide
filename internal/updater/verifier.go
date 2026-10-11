package updater

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ParseChecksums parses the content of a standard checksums.txt file.
// Format per line: "<sha256_hash>  <filename>" or "<sha256_hash> *<filename>"
func ParseChecksums(content string) map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			hash := strings.ToLower(strings.TrimSpace(parts[0]))
			filename := strings.TrimSpace(parts[1])
			filename = strings.TrimPrefix(filename, "*") // Strip binary mode prefix if present
			result[filename] = hash
		}
	}

	return result
}

// VerifySHA256 compares the SHA-256 hash of the provided byte slice against expectedHash.
func VerifySHA256(data []byte, expectedHash string) error {
	hasher := sha256.New()
	hasher.Write(data)
	actualHash := hex.EncodeToString(hasher.Sum(nil))

	expected := strings.ToLower(strings.TrimSpace(expectedHash))
	if actualHash != expected {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actualHash)
	}

	return nil
}
