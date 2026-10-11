package updater

import (
	"regexp"
	"strings"
)

var (
	mdLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	boldRegex   = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

// ExtractHighlights extracts up to maxItems bullet points from a markdown release body.
func ExtractHighlights(body string, maxItems int) []string {
	if maxItems <= 0 {
		maxItems = 3
	}

	lines := strings.Split(body, "\n")
	highlights := make([]string, 0, maxItems)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Look for bullet points
		if strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "+ ") {
			item := strings.TrimSpace(trimmed[2:])

			// Clean markdown links: [text](url) -> text
			item = mdLinkRegex.ReplaceAllString(item, "$1")
			// Clean markdown bold: **text** -> text
			item = boldRegex.ReplaceAllString(item, "$1")

			if item != "" {
				highlights = append(highlights, item)
				if len(highlights) >= maxItems {
					break
				}
			}
		}
	}

	return highlights
}
