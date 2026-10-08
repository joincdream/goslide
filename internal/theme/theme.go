package theme

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ErrThemeNotFound indicates that the requested theme name does not exist in built-in assets.
var ErrThemeNotFound = errors.New("theme not found")

var themeFrontmatterRegex = regexp.MustCompile(`(?m)^theme:\s*["']?[^"'\n]+["']?`)

const (
	// DefaultTheme is the fallback theme when none is specified.
	DefaultTheme = "default"

	coreJSPath          = "assets/js/goslide-core.js"
	starterTemplatePath = "assets/templates/demo.md"
)

// baseCSSFiles defines the modular stylesheets that compose the foundational presentation styling.
var baseCSSFiles = []string{
	"assets/css/deck-canvas.css",
	"assets/css/deck-content.css",
	"assets/css/presenter.css",
}

// Manager handles retrieving built-in themes and composing presentation stylesheets.
type Manager struct {
	fsys    fs.FS
	baseDir string
}

// NewManager creates a new theme Manager. If fsys is nil, the embedded filesystem is used.
func NewManager(fsys fs.FS) *Manager {
	if fsys == nil {
		fsys = embeddedAssets
	}
	return &Manager{fsys: fsys, baseDir: "."}
}

// SetBaseDir configures the base directory for resolving relative external themes.
func (m *Manager) SetBaseDir(dir string) {
	if dir == "" {
		dir = "."
	}
	m.baseDir = dir
}

// GetBaseCSS returns the composed raw contents of the base stylesheets:
// deck-canvas.css, deck-content.css, and presenter.css.
func (m *Manager) GetBaseCSS() (string, error) {
	var b strings.Builder
	for i, path := range baseCSSFiles {
		data, err := fs.ReadFile(m.fsys, path)
		if err != nil {
			return "", fmt.Errorf("failed to read base css file %q: %w", path, err)
		}
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.Write(data)
	}
	return b.String(), nil
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

// HasBuiltinTheme checks if the specified theme name exists in built-in assets.
func (m *Manager) HasBuiltinTheme(name string) bool {
	clean := strings.ToLower(strings.TrimSpace(name))
	if clean == "" {
		return false
	}
	targetPath := fmt.Sprintf("assets/css/%s.css", clean)
	info, err := fs.Stat(m.fsys, targetPath)
	return err == nil && !info.IsDir()
}

// ResolveTheme determines the effective built-in theme and any external custom CSS file path.
// It supports:
//  1. Explicit customCSSPath (CLI --theme-path flag) -> highest priority.
//  2. Frontmatter theme as external file path (e.g. "themes/corporate.css", "themes/corporate", "./my.css").
//  3. Frontmatter theme as simple name (e.g. "corporate" -> checks themes/corporate.css).
//  4. Built-in themes (e.g. "default", "clean", "dark").
//  5. Fallback to raw name if unmatched (for built-in error handling).
func (m *Manager) ResolveTheme(themeName, customCSSPath string) (baseTheme string, resolvedCustomPath string) {
	if strings.TrimSpace(customCSSPath) != "" {
		name := strings.TrimSpace(themeName)
		if m.HasBuiltinTheme(name) {
			return strings.ToLower(name), customCSSPath
		}
		return DefaultTheme, customCSSPath
	}

	name := strings.TrimSpace(themeName)
	if name == "" {
		return DefaultTheme, ""
	}

	// 1. If themeName is a path (contains "/" or "\" or ends with ".css")
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.HasSuffix(strings.ToLower(name), ".css") {
		if path := m.findExistingCSSFile(name); path != "" {
			return DefaultTheme, path
		}
	}

	// 2. If themeName is a simple identifier, check local project themes directory (e.g. themes/<name>.css)
	localThemeCandidate := filepath.Join("themes", name+".css")
	if path := m.findExistingCSSFile(localThemeCandidate); path != "" {
		return DefaultTheme, path
	}

	// Also check directly in working directory or baseDir (e.g. <name>.css)
	directCandidate := name + ".css"
	if path := m.findExistingCSSFile(directCandidate); path != "" {
		return DefaultTheme, path
	}

	// 3. Built-in theme check
	if m.HasBuiltinTheme(name) {
		return strings.ToLower(name), ""
	}

	return name, ""
}

func (m *Manager) findExistingCSSFile(path string) string {
	candidates := []string{path}
	if !strings.HasSuffix(strings.ToLower(path), ".css") {
		candidates = append(candidates, path+".css")
	}

	baseDir := m.baseDir
	if baseDir == "" {
		baseDir = "."
	}

	for _, cand := range candidates {
		// Check relative to current working directory
		if info, err := os.Stat(cand); err == nil && !info.IsDir() {
			return cand
		}
		// Check relative to baseDir if baseDir != "."
		if baseDir != "." && !filepath.IsAbs(cand) {
			joined := filepath.Join(baseDir, cand)
			if info, err := os.Stat(joined); err == nil && !info.IsDir() {
				return joined
			}
		}
	}
	return ""
}

// ComposeFullCSS combines base, theme, external custom file, and inline CSS in cascading order.
func (m *Manager) ComposeFullCSS(themeName, customCSSPath, inlineCSS string) (string, error) {
	effectiveTheme, effectiveCustomPath := m.ResolveTheme(themeName, customCSSPath)

	var b strings.Builder

	baseCSS, err := m.GetBaseCSS()
	if err != nil {
		return "", err
	}
	b.WriteString("/* --- Base Styles --- */\n")
	b.WriteString(baseCSS)
	b.WriteString("\n\n")

	themeCSS, err := m.GetThemeCSS(effectiveTheme)
	if err != nil {
		return "", err
	}
	b.WriteString(fmt.Sprintf("/* --- Theme: %s --- */\n", effectiveTheme))
	b.WriteString(themeCSS)
	b.WriteString("\n\n")

	if effectiveCustomPath != "" {
		customData, err := os.ReadFile(effectiveCustomPath)
		if err != nil {
			return "", fmt.Errorf("failed to read custom css file %q: %w", effectiveCustomPath, err)
		}
		b.WriteString(fmt.Sprintf("/* --- Custom CSS File: %s --- */\n", effectiveCustomPath))
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

// GetStarterTemplate returns the embedded starter presentation Markdown template.
// If themeName is provided, it replaces the theme property in the frontmatter.
func (m *Manager) GetStarterTemplate(themeName string) (string, error) {
	data, err := fs.ReadFile(m.fsys, starterTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read starter template %q: %w", starterTemplatePath, err)
	}

	content := string(data)
	themeName = strings.TrimSpace(themeName)
	if themeName != "" {
		content = themeFrontmatterRegex.ReplaceAllString(content, fmt.Sprintf("theme: %q", themeName))
	}

	return content, nil
}
