package updater

import (
	"context"
	"fmt"
	"time"
)

// UpdateOptions configures the self-update execution.
type UpdateOptions struct {
	CurrentVersion string
	TargetOS       string
	TargetArch     string
	TargetPath     string
	DryRun         bool
	Force          bool
	Client         *Client
	CacheManager   *CacheManager
	OnProgress     func(stage string, detail string)
}

// UpdateResult contains the result of a self-update operation.
type UpdateResult struct {
	CurrentVersion string
	UpdatedVersion string
	AlreadyLatest  bool
	ReleaseURL     string
	ReleaseBody    string
	Highlights     []string
}

// Updater orchestrates checking, downloading, verifying, and atomically replacing the binary.
type Updater struct {
	client   *Client
	cacheMgr *CacheManager
}

// NewUpdater creates a new Updater.
func NewUpdater(client *Client, cacheMgr *CacheManager) *Updater {
	if client == nil {
		client = NewClient()
	}
	if cacheMgr == nil {
		cacheMgr = NewCacheManager("")
	}
	return &Updater{
		client:   client,
		cacheMgr: cacheMgr,
	}
}

// Execute performs the full self-update sequence.
func (u *Updater) Execute(ctx context.Context, opts UpdateOptions) (*UpdateResult, error) {
	report := func(stage, detail string) {
		if opts.OnProgress != nil {
			opts.OnProgress(stage, detail)
		}
	}

	report("checking", "")
	release, err := u.client.GetLatestRelease(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check latest release: %w", err)
	}

	highlights := ExtractHighlights(release.Body, 5)

	isNewer := IsNewer(release.TagName, opts.CurrentVersion)
	if !isNewer && !opts.Force {
		return &UpdateResult{
			CurrentVersion: opts.CurrentVersion,
			UpdatedVersion: release.TagName,
			AlreadyLatest:  true,
			ReleaseURL:     release.HTMLURL,
			ReleaseBody:    release.Body,
			Highlights:     highlights,
		}, nil
	}

	// If dry run, stop before download
	if opts.DryRun {
		return &UpdateResult{
			CurrentVersion: opts.CurrentVersion,
			UpdatedVersion: release.TagName,
			AlreadyLatest:  false,
			ReleaseURL:     release.HTMLURL,
			ReleaseBody:    release.Body,
			Highlights:     highlights,
		}, nil
	}

	// 1. Find matching asset
	binaryAsset, err := FindMatchingAsset(release.Assets, opts.TargetOS, opts.TargetArch)
	if err != nil {
		return nil, err
	}

	// 2. Download binary asset
	report("downloading", binaryAsset.Name)
	archiveBytes, err := DownloadBytes(ctx, u.client.httpClient, binaryAsset.BrowserDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download release asset: %w", err)
	}

	// 3. Verify Checksum if checksums.txt is available
	checksumAsset, err := FindChecksumAsset(release.Assets)
	if err == nil && checksumAsset != nil {
		report("verifying", checksumAsset.Name)
		checksumBytes, err := DownloadBytes(ctx, u.client.httpClient, checksumAsset.BrowserDownloadURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download checksums: %w", err)
		}

		checksumMap := ParseChecksums(string(checksumBytes))
		expectedHash, found := checksumMap[binaryAsset.Name]
		if !found {
			return nil, fmt.Errorf("checksum for %s not found in %s", binaryAsset.Name, checksumAsset.Name)
		}

		if err := VerifySHA256(archiveBytes, expectedHash); err != nil {
			return nil, fmt.Errorf("integrity verification failed: %w", err)
		}
	}

	// 4. Extract executable from archive
	newBinary, err := ExtractBinaryFromArchive(archiveBytes, binaryAsset.Name, "goslide")
	if err != nil {
		return nil, fmt.Errorf("failed to extract binary: %w", err)
	}

	// 5. Replace binary atomically
	report("replacing", opts.TargetPath)
	if err := ReplaceBinary(opts.TargetPath, newBinary); err != nil {
		return nil, err
	}

	// 6. Update cache
	_ = u.cacheMgr.Save(&CacheData{
		LastCheckedAt: time.Now(),
		LatestVersion: release.TagName,
		ReleaseURL:    release.HTMLURL,
	})

	return &UpdateResult{
		CurrentVersion: opts.CurrentVersion,
		UpdatedVersion: release.TagName,
		AlreadyLatest:  false,
		ReleaseURL:     release.HTMLURL,
		ReleaseBody:    release.Body,
		Highlights:     highlights,
	}, nil
}
