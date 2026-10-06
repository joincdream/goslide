package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var pauseCommentRegex = regexp.MustCompile(`(?is)<!--\s*pause\s*-->\s*(<([a-zA-Z0-9]+)([^>]*)>)`)
var strayPauseRegex = regexp.MustCompile(`(?i)<!--\s*pause\s*-->`)

// transformFragments detects <!-- pause --> comments and injects class="fragment"
// and data-fragment-index attributes into the subsequent HTML block elements.
func transformFragments(html string) string {
	if !strings.Contains(html, "pause") {
		return html
	}

	fragmentIndex := 1
	res := pauseCommentRegex.ReplaceAllStringFunc(html, func(fullMatch string) string {
		sub := pauseCommentRegex.FindStringSubmatch(fullMatch)
		if len(sub) < 4 {
			return fullMatch
		}

		rawTag := sub[1] // e.g. <li class="item"> or <p>

		// Inject fragment class
		var newTag string
		existingClass, hasClass := getAttribute(rawTag, "class")
		if hasClass && existingClass != "" {
			newTag = setAttribute(rawTag, "class", existingClass+" fragment")
		} else {
			newTag = setAttribute(rawTag, "class", "fragment")
		}

		// Inject data-fragment-index
		newTag = setAttribute(newTag, "data-fragment-index", fmt.Sprintf("%d", fragmentIndex))
		fragmentIndex++

		return newTag
	})

	// Clean up any stray <!-- pause --> comments if there were no following tags
	return strayPauseRegex.ReplaceAllString(res, "")
}
