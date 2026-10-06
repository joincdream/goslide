package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	imgTagRegex      = regexp.MustCompile(`(?i)<img\s+([^>]*?)>`)
	widthDirective   = regexp.MustCompile(`(?i)\b(?:width|w):([0-9]+%|[0-9]+(?:px|em|rem|vw|vh)?\b)`)
	heightDirective  = regexp.MustCompile(`(?i)\b(?:height|h):([0-9]+%|[0-9]+(?:px|em|rem|vw|vh)?\b)`)
	centerDirective  = regexp.MustCompile(`(?i)\bcenter\b`)
	blockquoteRegex  = regexp.MustCompile(`(?s)<blockquote>(.*?)</blockquote>`)
	alertHeaderRegex = regexp.MustCompile(`(?s)^\s*<p>\s*\[!([a-zA-Z0-9_-]+)\](?:\s+([^\n<]+?))?(?:<br\s*/?>|\n|$)(.*?)(?:</p>)?$`)
)

type alertDefinition struct {
	Icon         string
	DefaultTitle string
}

var defaultAlertRegistry = map[string]alertDefinition{
	"note":      {Icon: "ℹ️", DefaultTitle: "Note"},
	"tip":       {Icon: "💡", DefaultTitle: "Tip"},
	"important": {Icon: "📌", DefaultTitle: "Important"},
	"warning":   {Icon: "⚠️", DefaultTitle: "Warning"},
	"caution":   {Icon: "🚨", DefaultTitle: "Caution"},
}

func postProcessHTML(html string) string {
	html = transformMediaElements(html)
	html = transformImages(html)
	html = transformAlerts(html)
	html = transformFragments(html)
	return html
}

func transformImages(html string) string {
	return imgTagRegex.ReplaceAllStringFunc(html, processSingleImageTag)
}

func processSingleImageTag(fullImgTag string) string {
	rawAlt, hasAlt := getAttribute(fullImgTag, "alt")
	if !hasAlt {
		return fullImgTag
	}

	styles, classes, cleanAlt, changed := parseImageDirectives(rawAlt)
	if !changed {
		return fullImgTag
	}

	resTag := setAttribute(fullImgTag, "alt", cleanAlt)
	return applyStylesAndClasses(resTag, styles, classes)
}

func parseImageDirectives(rawAlt string) (styles []string, classes []string, cleanAlt string, changed bool) {
	// Check width
	if wMatch := widthDirective.FindStringSubmatch(rawAlt); len(wMatch) >= 2 {
		val := wMatch[1]
		if isPureDigits(val) {
			val += "px"
		}
		styles = append(styles, fmt.Sprintf("width: %s;", val))
	}

	// Check height
	if hMatch := heightDirective.FindStringSubmatch(rawAlt); len(hMatch) >= 2 {
		val := hMatch[1]
		if isPureDigits(val) {
			val += "px"
		}
		styles = append(styles, fmt.Sprintf("height: %s;", val))
	}

	// Check center
	if centerDirective.MatchString(rawAlt) {
		classes = append(classes, "img-center")
	}

	if len(styles) == 0 && len(classes) == 0 {
		return nil, nil, rawAlt, false
	}

	clean := widthDirective.ReplaceAllString(rawAlt, "")
	clean = heightDirective.ReplaceAllString(clean, "")
	clean = centerDirective.ReplaceAllString(clean, "")
	cleanAlt = strings.Join(strings.Fields(clean), " ")
	return styles, classes, cleanAlt, true
}

func applyStylesAndClasses(tag string, styles, classes []string) string {
	if len(styles) > 0 {
		newStyleStr := strings.Join(styles, " ")
		existingStyle, hasStyle := getAttribute(tag, "style")
		if hasStyle && existingStyle != "" {
			tag = setAttribute(tag, "style", strings.TrimRight(existingStyle, "; ")+"; "+newStyleStr)
		} else {
			tag = setAttribute(tag, "style", newStyleStr)
		}
	}

	if len(classes) > 0 {
		newClassStr := strings.Join(classes, " ")
		existingClass, hasClass := getAttribute(tag, "class")
		if hasClass && existingClass != "" {
			tag = setAttribute(tag, "class", existingClass+" "+newClassStr)
		} else {
			tag = setAttribute(tag, "class", newClassStr)
		}
	}

	return tag
}

func getAttribute(tag, attrName string) (string, bool) {
	re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(attrName) + `=(?:"([^"]*)"|'([^']*)')`)
	match := re.FindStringSubmatch(tag)
	if len(match) == 3 {
		if match[1] != "" {
			return match[1], true
		}
		return match[2], true
	}
	return "", false
}

func setAttribute(tag, attrName, value string) string {
	re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(attrName) + `=(?:"[^"]*"|'[^']*')`)
	newAttr := fmt.Sprintf(`%s="%s"`, attrName, value)
	if re.MatchString(tag) {
		return re.ReplaceAllString(tag, newAttr)
	}

	// Insert before closing >
	trimmed := strings.TrimSuffix(tag, ">")
	if strings.HasSuffix(trimmed, "/") {
		trimmed = strings.TrimSpace(strings.TrimSuffix(trimmed, "/"))
		return fmt.Sprintf("%s %s />", trimmed, newAttr)
	}
	return fmt.Sprintf("%s %s>", trimmed, newAttr)
}

func isPureDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func transformAlerts(html string) string {
	return blockquoteRegex.ReplaceAllStringFunc(html, renderSingleAlert)
}

func renderSingleAlert(fullBQ string) string {
	sub := blockquoteRegex.FindStringSubmatch(fullBQ)
	if len(sub) < 2 {
		return fullBQ
	}
	bqContent := strings.TrimSpace(sub[1])

	// Look for [!NOTE], etc. in the first paragraph
	pStart := strings.Index(bqContent, "<p>")
	if pStart == -1 {
		return fullBQ
	}
	pEnd := strings.Index(bqContent, "</p>")
	var firstP string
	var restBQ string
	if pEnd != -1 {
		firstP = bqContent[pStart : pEnd+4]
		restBQ = strings.TrimSpace(bqContent[pEnd+4:])
	} else {
		firstP = bqContent[pStart:]
		restBQ = ""
	}

	headerMatch := alertHeaderRegex.FindStringSubmatch(firstP)
	if len(headerMatch) < 4 {
		return fullBQ
	}

	alertType := strings.ToLower(headerMatch[1])
	alertDef, exists := defaultAlertRegistry[alertType]
	if !exists {
		return fullBQ
	}

	customTitle := strings.TrimSpace(headerMatch[2])
	firstPRest := strings.TrimSpace(headerMatch[3])

	title := customTitle
	if title == "" {
		title = alertDef.DefaultTitle
	}

	var bodyBuilder strings.Builder
	if firstPRest != "" {
		bodyBuilder.WriteString(fmt.Sprintf("<p>%s</p>\n", firstPRest))
	}
	if restBQ != "" {
		bodyBuilder.WriteString(restBQ)
	}

	return fmt.Sprintf(
		`<div class="markdown-alert markdown-alert-%s"><div class="markdown-alert-title"><span class="markdown-alert-icon">%s</span> %s</div>%s</div>`,
		alertType,
		alertDef.Icon,
		title,
		bodyBuilder.String(),
	)
}
