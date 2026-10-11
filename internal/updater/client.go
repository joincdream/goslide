package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	defaultGitHubAPI = "https://api.github.com/repos/joincdream/goslide/releases/latest"
	defaultUserAgent = "Goslide-Updater"
	defaultTimeout   = 10 * time.Second
)

// Asset represents a single release asset artifact from GitHub Releases.
type Asset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// ReleaseInfo represents metadata for a GitHub release.
type ReleaseInfo struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	HTMLURL     string    `json:"html_url"`
	Body        string    `json:"body"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

// Client is an HTTP client for querying GitHub Releases.
type Client struct {
	apiURL     string
	httpClient *http.Client
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithAPIURL overrides the default GitHub API endpoint (useful for testing).
func WithAPIURL(url string) ClientOption {
	return func(c *Client) {
		if url != "" {
			c.apiURL = url
		}
	}
}

// WithHTTPClient injects a custom *http.Client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// NewClient creates a new GitHub Releases client.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		apiURL: defaultGitHubAPI,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// GetLatestRelease fetches the latest published release metadata.
func (c *Client) GetLatestRelease(ctx context.Context) (*ReleaseInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("no release found on repository")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, resp.Status)
	}

	var info ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to parse release json: %w", err)
	}

	info.TagName = strings.TrimSpace(info.TagName)
	return &info, nil
}
