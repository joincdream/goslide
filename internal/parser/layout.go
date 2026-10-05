package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yundream/goslide/internal/model"
)

var splitCommentRegex = regexp.MustCompile(`(?i)<!--\s*split\s*-->`)

// renderSlideContent renders slide markdown into TitleHTML, HTMLContent, and optional two-cols HTML.
func renderSlideContent(
	gm goldmark.Markdown,
	layout model.LayoutType,
	content string,
) (titleHTML string, htmlContent string, leftHTML string, rightHTML string, err error) {
	// For cover, section, and blank layouts, retain all content centered as a whole (no separated title box)
	if layout == model.LayoutCover || layout == model.LayoutSection || layout == model.LayoutBlank {
		fullHTML, err := convertToHTML(gm, content)
		if err != nil {
			return "", "", "", "", err
		}
		return "", fullHTML, "", "", nil
	}

	// For standard layouts, extract the first leading heading as Slide Title
	titleMD, bodyMD := extractFirstHeading(content)
	if titleMD != "" {
		tBuf, err := convertToHTML(gm, titleMD)
		if err != nil {
			return "", "", "", "", fmt.Errorf("failed to render slide title: %w", err)
		}
		titleHTML = tBuf
	}

	if layout == model.LayoutTwoCols && splitCommentRegex.MatchString(bodyMD) {
		parts := splitCommentRegex.Split(bodyMD, 2)
		leftMD := strings.TrimSpace(parts[0])
		rightMD := strings.TrimSpace(parts[1])

		leftBuf, err := convertToHTML(gm, leftMD)
		if err != nil {
			return "", "", "", "", fmt.Errorf("failed to render left column: %w", err)
		}

		rightBuf, err := convertToHTML(gm, rightMD)
		if err != nil {
			return "", "", "", "", fmt.Errorf("failed to render right column: %w", err)
		}

		combinedHTML := fmt.Sprintf(
			"<div class=\"two-cols\"><div class=\"col-left\">\n%s\n</div><div class=\"col-right\">\n%s\n</div></div>",
			leftBuf,
			rightBuf,
		)

		return titleHTML, combinedHTML, leftBuf, rightBuf, nil
	}

	fullHTML, err := convertToHTML(gm, bodyMD)
	if err != nil {
		return "", "", "", "", err
	}

	return titleHTML, fullHTML, "", "", nil
}

// extractFirstHeading finds the first leading # or ## heading in the markdown text.
// It returns the heading markdown string and the remaining markdown body.
func extractFirstHeading(content string) (string, string) {
	lines := strings.Split(content, "\n")
	firstNonEmpty := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			firstNonEmpty = i
			break
		}
	}
	if firstNonEmpty == -1 {
		return "", content
	}

	firstLine := strings.TrimSpace(lines[firstNonEmpty])
	if strings.HasPrefix(firstLine, "# ") || strings.HasPrefix(firstLine, "## ") {
		titleMD := firstLine
		remaining := append([]string{}, lines[:firstNonEmpty]...)
		remaining = append(remaining, lines[firstNonEmpty+1:]...)
		return titleMD, strings.TrimSpace(strings.Join(remaining, "\n"))
	}

	return "", content
}

func convertToHTML(gm goldmark.Markdown, markdown string) (string, error) {
	var buf bytes.Buffer
	if err := gm.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return postProcessHTML(buf.String()), nil
}
