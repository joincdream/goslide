package updater

import (
	"context"
	"time"
)

// UpdateCheckResult holds the outcome of checking for available updates.
type UpdateCheckResult struct {
	CurrentVersion string
	LatestVersion  string
	HasUpdate      bool
	ReleaseURL     string
	ReleaseBody    string
}

// Checker coordinates version checks between cache and remote GitHub API.
type Checker struct {
	client       *Client
	cacheManager *CacheManager
	cacheTTL     time.Duration
}

// NewChecker creates a Checker with the given client and cache manager.
func NewChecker(client *Client, cacheMgr *CacheManager) *Checker {
	if client == nil {
		client = NewClient()
	}
	if cacheMgr == nil {
		cacheMgr = NewCacheManager("")
	}
	return &Checker{
		client:       client,
		cacheManager: cacheMgr,
		cacheTTL:     DefaultCacheTTL,
	}
}

// CheckUpdate checks whether a newer version of Goslide is available.
// It leverages local cache to avoid hitting GitHub API on every CLI run.
func (c *Checker) CheckUpdate(ctx context.Context, currentVersion string) (*UpdateCheckResult, error) {
	// 1. Try reading valid cache first
	if cached, _ := c.cacheManager.Load(); cached != nil {
		if !c.cacheManager.IsExpired(cached, c.cacheTTL) && cached.LatestVersion != "" {
			hasUpdate := IsNewer(cached.LatestVersion, currentVersion)
			return &UpdateCheckResult{
				CurrentVersion: currentVersion,
				LatestVersion:  cached.LatestVersion,
				HasUpdate:      hasUpdate,
				ReleaseURL:     cached.ReleaseURL,
			}, nil
		}
	}

	// 2. Cache is missing or expired, fetch from GitHub API
	info, err := c.client.GetLatestRelease(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Save to cache
	_ = c.cacheManager.Save(&CacheData{
		LastCheckedAt: time.Now(),
		LatestVersion: info.TagName,
		ReleaseURL:    info.HTMLURL,
	})

	hasUpdate := IsNewer(info.TagName, currentVersion)
	return &UpdateCheckResult{
		CurrentVersion: currentVersion,
		LatestVersion:  info.TagName,
		HasUpdate:      hasUpdate,
		ReleaseURL:     info.HTMLURL,
		ReleaseBody:    info.Body,
	}, nil
}
