package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yundream/goslide/internal/model"
)

var (
	malformedClosingDelimRegex = regexp.MustCompile(`<!--(.*?)--{2,}>`)
	commentPattern             = regexp.MustCompile(`<!--(?s:(.*?))-->`)
)

// inspectSlideComments checks a slide chunk for directive syntax errors,
// unknown layouts, and missing underscores, returning diagnostics and sanitized content.
func inspectSlideComments(slideIndex int, rawContent string) ([]model.Diagnostic, string) {
	var diagnostics []model.Diagnostic

	// 1. Sanitize malformed delimiters (e.g. ---> to -->)
	sanitizedContent := rawContent
	if malformedClosingDelimRegex.MatchString(rawContent) {
		lines := strings.Split(rawContent, "\n")
		for lineIdx, line := range lines {
			if matches := malformedClosingDelimRegex.FindAllString(line, -1); len(matches) > 0 {
				for _, match := range matches {
					diagnostics = append(diagnostics, model.Diagnostic{
						Severity:   model.SeverityWarning,
						SlideIndex: slideIndex,
						Line:       lineIdx + 1,
						Rule:       "syntax.malformed_delimiter",
						Message:    "Malformed HTML comment closing delimiter",
						RawSnippet: match,
						Candidates: []string{"-->"},
					})
				}
			}
		}
		sanitizedContent = malformedClosingDelimRegex.ReplaceAllString(rawContent, "<!--$1-->")
	}

	// 2. Identify code block line ranges to isolate code block comments
	inCodeBlock := false
	codeBlockLines := make(map[int]bool)
	lines := strings.Split(sanitizedContent, "\n")
	for lineIdx, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inCodeBlock = !inCodeBlock
			codeBlockLines[lineIdx+1] = true
			continue
		}
		if inCodeBlock {
			codeBlockLines[lineIdx+1] = true
		}
	}

	// 3. Inspect each HTML comment outside code blocks
	commentIndices := commentPattern.FindAllStringIndex(sanitizedContent, -1)
	for _, loc := range commentIndices {
		startIdx := loc[0]
		rawSnippet := sanitizedContent[loc[0]:loc[1]]

		// Calculate 1-based relative line number in slide
		lineNum := strings.Count(sanitizedContent[:startIdx], "\n") + 1

		// Skip comments inside code blocks (False-Positive Zero)
		if codeBlockLines[lineNum] {
			continue
		}

		commentBody := strings.TrimSpace(sanitizedContent[startIdx+4 : loc[1]-3])

		// Case A: Speaker notes - ignore
		if isNoteComment(commentBody) {
			continue
		}

		// Case B: Known single markers
		lowerBody := strings.ToLower(commentBody)
		if lowerBody == "split" || lowerBody == "pause" {
			continue
		}

		// Check for marker typos (e.g. splti, puase)
		if lowerBody == "splti" || lowerBody == "spilt" {
			diagnostics = append(diagnostics, model.Diagnostic{
				Severity:   model.SeverityWarning,
				SlideIndex: slideIndex,
				Line:       lineNum,
				Rule:       "marker.unknown",
				Message:    fmt.Sprintf("Unknown marker '%s'", commentBody),
				RawSnippet: rawSnippet,
				Candidates: []string{"split"},
			})
			continue
		}
		if lowerBody == "puase" || lowerBody == "paues" {
			diagnostics = append(diagnostics, model.Diagnostic{
				Severity:   model.SeverityWarning,
				SlideIndex: slideIndex,
				Line:       lineNum,
				Rule:       "marker.unknown",
				Message:    fmt.Sprintf("Unknown marker '%s'", commentBody),
				RawSnippet: rawSnippet,
				Candidates: []string{"pause"},
			})
			continue
		}

		// Case C: Check key-value directives
		commentLines := strings.Split(commentBody, "\n")
		for offset, cLine := range commentLines {
			cTrimmed := strings.TrimSpace(cLine)
			if cTrimmed == "" {
				continue
			}

			colonIdx := strings.Index(cTrimmed, ":")
			if colonIdx == -1 {
				continue
			}

			key := strings.TrimSpace(cTrimmed[:colonIdx])
			val := unquoteValue(cTrimmed[colonIdx+1:])

			isLocal := strings.HasPrefix(key, "_")
			cleanKey := strings.TrimPrefix(key, "_")
			currentLineNum := lineNum + offset

			// 1) Underscore missing on slide-local directive
			if !isLocal {
				if canonical, ok := MatchCanonicalDirective(cleanKey); ok {
					diagnostics = append(diagnostics, model.Diagnostic{
						Severity:   model.SeverityWarning,
						SlideIndex: slideIndex,
						Line:       currentLineNum,
						Rule:       "directive.missing_underscore",
						Message:    "Slide-local directive missing underscore prefix",
						RawSnippet: rawSnippet,
						Candidates: []string{canonical},
					})

					if strings.EqualFold(cleanKey, "layout") && !model.IsKnownLayout(val) {
						diagnostics = append(diagnostics, model.Diagnostic{
							Severity:   model.SeverityWarning,
							SlideIndex: slideIndex,
							Line:       currentLineNum,
							Rule:       "layout.unknown",
							Message:    fmt.Sprintf("Unknown layout '%s'", val),
							RawSnippet: rawSnippet,
							Candidates: model.AvailableLayouts(),
						})
					}
					continue
				}
				// General user comment (memo:, todo:, etc.) - ignore!
				continue
			}

			// 2) Directive key is prefixed with '_'
			if strings.EqualFold(cleanKey, "layout") {
				if !model.IsKnownLayout(val) {
					diagnostics = append(diagnostics, model.Diagnostic{
						Severity:   model.SeverityWarning,
						SlideIndex: slideIndex,
						Line:       currentLineNum,
						Rule:       "layout.unknown",
						Message:    fmt.Sprintf("Unknown layout '%s'", val),
						RawSnippet: rawSnippet,
						Candidates: model.AvailableLayouts(),
					})
				}
				continue
			}

			// 3) Check if unknown '_' directive key
			if _, ok := MatchCanonicalDirective(cleanKey); !ok {
				diagnostics = append(diagnostics, model.Diagnostic{
					Severity:   model.SeverityWarning,
					SlideIndex: slideIndex,
					Line:       currentLineNum,
					Rule:       "directive.unknown",
					Message:    fmt.Sprintf("Unknown directive '%s'", key),
					RawSnippet: rawSnippet,
					Candidates: AvailableDirectives(),
				})
			}
		}
	}

	return diagnostics, sanitizedContent
}
