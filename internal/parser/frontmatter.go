package parser

import (
	"fmt"
	"strings"

	"github.com/yundream/goslide/internal/model"
	"gopkg.in/yaml.v3"
)

type frontmatterData struct {
	Title     string `yaml:"title"`
	Author    string `yaml:"author"`
	Theme     string `yaml:"theme"`
	Layout    string `yaml:"layout"`
	Size      string `yaml:"size"`
	Paginate  *bool  `yaml:"paginate"`
	Header    string `yaml:"header"`
	Footer    string `yaml:"footer"`
	Autofit   *bool  `yaml:"autofit"`
	CustomCSS string `yaml:"custom_css"`
	Style     string `yaml:"style"` // Marp style compatibility
	Marp      any    `yaml:"marp"`  // Marp compatibility flag
}

// FrontmatterResult holds parsed metadata and remaining markdown body.
type FrontmatterResult struct {
	Title       string
	Author      string
	CustomCSS   string
	GlobalAttrs model.GlobalDirectives
	Body        string
}

// extractFrontmatter extracts YAML frontmatter from normalized markdown content.
func extractFrontmatter(content string) (*FrontmatterResult, error) {
	defaultResult := &FrontmatterResult{
		GlobalAttrs: model.GlobalDirectives{
			Theme:    "default",
			Layout:   model.LayoutDefault,
			Size:     model.Ratio16x9,
			Paginate: false,
		},
		Body: content,
	}

	trimmed := strings.TrimLeft(content, " \t")
	if !strings.HasPrefix(trimmed, "---") {
		return defaultResult, nil
	}

	// Must be on its own line: "---" followed by optional whitespace and newline
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return defaultResult, nil
	}

	firstLine := strings.TrimRight(lines[0], " \t")
	if firstLine != "---" {
		return defaultResult, nil
	}

	// Find the closing "---" line
	closingIdx := -1
	for i := 1; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], " \t")
		if line == "---" {
			closingIdx = i
			break
		}
	}

	if closingIdx == -1 {
		// No closing delimiter found; treat entire content as body
		return defaultResult, nil
	}

	yamlBlock := strings.Join(lines[1:closingIdx], "\n")
	remainingBody := strings.Join(lines[closingIdx+1:], "\n")

	var data frontmatterData
	if strings.TrimSpace(yamlBlock) != "" {
		if err := yaml.Unmarshal([]byte(yamlBlock), &data); err != nil {
			return nil, fmt.Errorf("%w: %v", model.ErrInvalidFrontmatter, err)
		}
	}

	res := &FrontmatterResult{
		Title:  data.Title,
		Author: data.Author,
		GlobalAttrs: model.GlobalDirectives{
			Theme:    resolveTheme(data.Theme),
			Layout:   resolveLayout(data.Layout),
			Size:     resolveSize(data.Size),
			Paginate: data.Paginate != nil && *data.Paginate,
			Header:   data.Header,
			Footer:   data.Footer,
			Autofit:  data.Autofit != nil && *data.Autofit,
		},
		CustomCSS: resolveCustomCSS(data.CustomCSS, data.Style),
		Body:      remainingBody,
	}

	return res, nil
}

func resolveTheme(theme string) string {
	t := strings.TrimSpace(theme)
	if t == "" {
		return "default"
	}
	// If the theme specifies a file path (contains path separators or ends with .css), preserve it
	if strings.Contains(t, "/") || strings.Contains(t, "\\") || strings.HasSuffix(strings.ToLower(t), ".css") {
		return t
	}
	switch t {
	case "default", "clean", "dark", "corporate", "academic", "cyber-dark":
		return t
	default:
		// Unsupported themes fallback to "default" as per Decision 3
		return "default"
	}
}

func resolveSize(size string) model.SizeRatio {
	switch strings.TrimSpace(size) {
	case string(model.Ratio16x9):
		return model.Ratio16x9
	case string(model.Ratio4x3):
		return model.Ratio4x3
	default:
		return model.Ratio16x9
	}
}

func resolveLayout(layout string) model.LayoutType {
	return model.NormalizeLayout(layout)
}

func resolveCustomCSS(customCSS, style string) string {
	if customCSS != "" {
		return customCSS
	}
	return style
}
