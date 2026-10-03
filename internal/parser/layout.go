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

// renderSlideContent renders slide markdown into HTMLContent, splitting two-cols if applicable.
func renderSlideContent(
	gm goldmark.Markdown,
	layout model.LayoutType,
	content string,
) (string, string, string, error) {
	if layout == model.LayoutTwoCols && splitCommentRegex.MatchString(content) {
		parts := splitCommentRegex.Split(content, 2)
		leftMD := strings.TrimSpace(parts[0])
		rightMD := strings.TrimSpace(parts[1])

		leftBuf, err := convertToHTML(gm, leftMD)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to render left column: %w", err)
		}

		rightBuf, err := convertToHTML(gm, rightMD)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to render right column: %w", err)
		}

		combinedHTML := fmt.Sprintf(
			"<div class=\"two-cols\"><div class=\"col-left\">\n%s\n</div><div class=\"col-right\">\n%s\n</div></div>",
			leftBuf,
			rightBuf,
		)

		return combinedHTML, leftBuf, rightBuf, nil
	}

	fullHTML, err := convertToHTML(gm, content)
	if err != nil {
		return "", "", "", err
	}

	return fullHTML, "", "", nil
}

func convertToHTML(gm goldmark.Markdown, markdown string) (string, error) {
	var buf bytes.Buffer
	if err := gm.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
