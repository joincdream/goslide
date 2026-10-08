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

var (
	adjacentUlFragmentRegex = regexp.MustCompile(`(?is)</ul>\s*<ul([^>]*(?:class=["'][^"']*fragment[^"']*|data-fragment-index)[^>]*)>(.*?)</ul>`)
	adjacentOlFragmentRegex = regexp.MustCompile(`(?is)</ol>\s*<ol([^>]*(?:class=["'][^"']*fragment[^"']*|data-fragment-index)[^>]*)>(.*?)</ol>`)
)

// normalizeListFragments detects adjacent lists created by <!-- pause --> comments between items
// and merges them into a single list container, transferring fragment attributes to the direct <li> children.
func normalizeListFragments(html string) string {
	if !strings.Contains(html, "fragment") {
		return html
	}

	html = mergeAdjacentList(html, adjacentUlFragmentRegex, "ul")
	html = mergeAdjacentList(html, adjacentOlFragmentRegex, "ol")
	return html
}

func mergeAdjacentList(html string, re *regexp.Regexp, tag string) string {
	for {
		matched := false
		newHTML := re.ReplaceAllStringFunc(html, func(fullMatch string) string {
			sub := re.FindStringSubmatch(fullMatch)
			if len(sub) < 3 {
				return fullMatch
			}

			rawAttributes := "<dummy " + sub[1] + ">"
			fragClass, _ := getAttribute(rawAttributes, "class")
			fragIndex, _ := getAttribute(rawAttributes, "data-fragment-index")
			innerContent := sub[2]

			transformedInner := injectFragmentToListItems(innerContent, fragClass, fragIndex)
			matched = true
			return transformedInner + "</" + tag + ">"
		})

		if !matched || newHTML == html {
			break
		}
		html = newHTML
	}
	return html
}

// injectFragmentToListItems traverses top-level <li> tags within a list chunk and injects fragment attributes.
// It avoids modifying nested sub-lists by tracking list element depth.
func injectFragmentToListItems(content string, fragmentClass string, fragmentIndex string) string {
	var sb strings.Builder
	idx := 0
	listDepth := 0

	for idx < len(content) {
		if content[idx] != '<' {
			sb.WriteByte(content[idx])
			idx++
			continue
		}

		end := strings.IndexByte(content[idx:], '>')
		if end == -1 {
			sb.WriteString(content[idx:])
			break
		}

		tag := content[idx : idx+end+1]
		lowerTag := strings.ToLower(tag)

		switch {
		case strings.HasPrefix(lowerTag, "<ul") || strings.HasPrefix(lowerTag, "<ol"):
			listDepth++
			sb.WriteString(tag)
		case strings.HasPrefix(lowerTag, "</ul") || strings.HasPrefix(lowerTag, "</ol"):
			if listDepth > 0 {
				listDepth--
			}
			sb.WriteString(tag)
		case listDepth == 0 && (strings.HasPrefix(lowerTag, "<li ") || strings.HasPrefix(lowerTag, "<li>")):
			sb.WriteString(applyFragmentToLiTag(tag, fragmentClass, fragmentIndex))
		default:
			sb.WriteString(tag)
		}

		idx += end + 1
	}

	return sb.String()
}

func applyFragmentToLiTag(tag, fragmentClass, fragmentIndex string) string {
	modifiedTag := tag
	if fragmentClass != "" {
		existingClass, hasClass := getAttribute(modifiedTag, "class")
		switch {
		case !hasClass || existingClass == "":
			modifiedTag = setAttribute(modifiedTag, "class", "fragment")
		case !strings.Contains(existingClass, "fragment"):
			modifiedTag = setAttribute(modifiedTag, "class", existingClass+" fragment")
		}
	}
	if fragmentIndex != "" {
		modifiedTag = setAttribute(modifiedTag, "data-fragment-index", fragmentIndex)
	}
	return modifiedTag
}
