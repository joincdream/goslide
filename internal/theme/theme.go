package theme

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// ErrThemeNotFound indicates that the requested theme name does not exist in built-in assets.
var ErrThemeNotFound = errors.New("theme not found")

const (
	// DefaultTheme is the fallback theme when none is specified.
	DefaultTheme = "default"

	baseCSSPath = "assets/css/base.css"
	coreJSPath  = "assets/js/goslide-core.js"
)

// Manager handles retrieving built-in themes and composing presentation stylesheets.
type Manager struct {
	fsys fs.FS
}

// NewManager creates a new theme Manager. If fsys is nil, the embedded filesystem is used.
func NewManager(fsys fs.FS) *Manager {
	if fsys == nil {
		fsys = embeddedAssets
	}
	return &Manager{fsys: fsys}
}

// GetBaseCSS returns the raw contents of the base stylesheet.
func (m *Manager) GetBaseCSS() (string, error) {
	data, err := fs.ReadFile(m.fsys, baseCSSPath)
	if err != nil {
		return "", fmt.Errorf("failed to read base css: %w", err)
	}
	return string(data), nil
}

// GetCoreJS returns the raw contents of the core presentation runtime script.
func (m *Manager) GetCoreJS() (string, error) {
	data, err := fs.ReadFile(m.fsys, coreJSPath)
	if err != nil {
		return "", fmt.Errorf("failed to read core js: %w", err)
	}
	return string(data), nil
}

// GetThemeCSS returns the CSS string for the specified built-in theme.
// If themeName is empty, it falls back to DefaultTheme ("default").
func (m *Manager) GetThemeCSS(themeName string) (string, error) {
	name := strings.ToLower(strings.TrimSpace(themeName))
	if name == "" {
		name = DefaultTheme
	}

	targetPath := fmt.Sprintf("assets/css/%s.css", name)
	data, err := fs.ReadFile(m.fsys, targetPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("%w: %q", ErrThemeNotFound, name)
		}
		return "", fmt.Errorf("failed to read theme css %q: %w", name, err)
	}

	return string(data), nil
}

// ComposeFullCSS combines base, theme, external custom file, and inline CSS in cascading order.
func (m *Manager) ComposeFullCSS(themeName, customCSSPath, inlineCSS string) (string, error) {
	var b strings.Builder

	baseCSS, err := m.GetBaseCSS()
	if err != nil {
		return "", err
	}
	b.WriteString("/* --- Base Styles --- */\n")
	b.WriteString(baseCSS)
	b.WriteString("\n\n")

	themeCSS, err := m.GetThemeCSS(themeName)
	if err != nil {
		return "", err
	}
	b.WriteString(fmt.Sprintf("/* --- Theme: %s --- */\n", themeName))
	b.WriteString(themeCSS)
	b.WriteString("\n\n")

	if customCSSPath != "" {
		customData, err := os.ReadFile(customCSSPath)
		if err != nil {
			return "", fmt.Errorf("failed to read custom css file %q: %w", customCSSPath, err)
		}
		b.WriteString(fmt.Sprintf("/* --- Custom CSS File: %s --- */\n", customCSSPath))
		b.Write(customData)
		b.WriteString("\n\n")
	}

	if strings.TrimSpace(inlineCSS) != "" {
		b.WriteString("/* --- Inline Styles --- */\n")
		b.WriteString(inlineCSS)
		b.WriteString("\n")
	}

	return b.String(), nil
}
