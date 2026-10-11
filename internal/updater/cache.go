package updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultCacheFile = ".goslide/update_check.json"
	DefaultCacheTTL  = 24 * time.Hour
)

// CacheData stores cached update check results on local disk.
type CacheData struct {
	LastCheckedAt time.Time `json:"last_checked_at"`
	LatestVersion string    `json:"latest_version"`
	ReleaseURL    string    `json:"release_url"`
}

// CacheManager manages reading and writing the update check cache.
type CacheManager struct {
	filePath string
}

// NewCacheManager creates a new CacheManager. If customPath is empty, it defaults to ~/.goslide/update_check.json.
func NewCacheManager(customPath string) *CacheManager {
	if customPath != "" {
		return &CacheManager{filePath: customPath}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to local working directory if home cannot be determined
		return &CacheManager{filePath: filepath.Join(".", ".goslide", "update_check.json")}
	}
	return &CacheManager{filePath: filepath.Join(home, defaultCacheFile)}
}

// Load reads the cache file from disk. Returns nil if file does not exist or is corrupted.
func (cm *CacheManager) Load() (*CacheData, error) {
	data, err := os.ReadFile(cm.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read update cache: %w", err)
	}

	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update cache: %w", err)
	}
	return &cache, nil
}

// Save writes the update cache to disk, creating parent directories if necessary.
func (cm *CacheManager) Save(cache *CacheData) error {
	if cache == nil {
		return errors.New("cannot save nil cache")
	}

	dir := filepath.Dir(cm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize update cache: %w", err)
	}

	if err := os.WriteFile(cm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write update cache file: %w", err)
	}

	return nil
}

// IsExpired checks if the given cache data has exceeded the specified TTL.
func (cm *CacheManager) IsExpired(cache *CacheData, ttl time.Duration) bool {
	if cache == nil || cache.LastCheckedAt.IsZero() {
		return true
	}
	return time.Since(cache.LastCheckedAt) > ttl
}
