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
// unknown layouts, and missing underscores, returning diagnostics without mutating the content.
func inspectSlideComments(slideIndex int, rawContent string) []model.Diagnostic {
	var diagnostics []model.Diagnostic

	// 1. Inspect malformed delimiters (e.g. ---> instead of -->)
	diagnostics = append(diagnostics, inspectMalformedDelimiters(slideIndex, rawContent)...)

	// 2. Identify code block line ranges to isolate code block comments
	codeBlockLines := findCodeBlockLines(rawContent)

	// 3. Inspect each HTML comment outside code blocks
	commentIndices := commentPattern.FindAllStringIndex(rawContent, -1)
	for _, loc := range commentIndices {
		startIdx := loc[0]
		rawSnippet := rawContent[loc[0]:loc[1]]
		lineNum := strings.Count(rawContent[:startIdx], "\n") + 1

		if codeBlockLines[lineNum] {
			continue
		}

		commentBody := strings.TrimSpace(rawContent[startIdx+4 : loc[1]-3])
		commentBody = strings.TrimSpace(strings.TrimRight(commentBody, "-"))
		if isNoteComment(commentBody) {
			continue
		}

		lowerBody := strings.ToLower(commentBody)
		if lowerBody == "split" || lowerBody == "pause" {
			continue
		}

		if diag, ok := inspectMarkerTypo(lowerBody, rawSnippet, slideIndex, lineNum); ok {
			diagnostics = append(diagnostics, diag)
			continue
		}

		commentLines := strings.Split(commentBody, "\n")
		for offset, cLine := range commentLines {
			cTrimmed := strings.TrimSpace(cLine)
			if cTrimmed != "" {
				diagnostics = append(diagnostics, inspectDirectiveLine(slideIndex, lineNum+offset, cTrimmed, rawSnippet)...)
			}
		}
	}

	return diagnostics
}

func inspectMalformedDelimiters(slideIndex int, rawContent string) []model.Diagnostic {
	if !malformedClosingDelimRegex.MatchString(rawContent) {
		return nil
	}

	var diagnostics []model.Diagnostic
	lines := strings.Split(rawContent, "\n")
	for lineIdx, line := range lines {
		matches := malformedClosingDelimRegex.FindAllString(line, -1)
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
	return diagnostics
}

func findCodeBlockLines(content string) map[int]bool {
	inCodeBlock := false
	codeBlockLines := make(map[int]bool)
	lines := strings.Split(content, "\n")
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
	return codeBlockLines
}

func inspectMarkerTypo(lowerBody, rawSnippet string, slideIndex, lineNum int) (model.Diagnostic, bool) {
	var candidate string
	switch lowerBody {
	case "splti", "spilt":
		candidate = "split"
	case "puase", "paues":
		candidate = "pause"
	default:
		return model.Diagnostic{}, false
	}

	return model.Diagnostic{
		Severity:   model.SeverityWarning,
		SlideIndex: slideIndex,
		Line:       lineNum,
		Rule:       "marker.unknown",
		Message:    fmt.Sprintf("Unknown marker '%s'", lowerBody),
		RawSnippet: rawSnippet,
		Candidates: []string{candidate},
	}, true
}

func inspectDirectiveLine(slideIndex, lineNum int, cTrimmed, rawSnippet string) []model.Diagnostic {
	colonIdx := strings.Index(cTrimmed, ":")
	if colonIdx == -1 {
		return nil
	}

	key := strings.TrimSpace(cTrimmed[:colonIdx])
	val := unquoteValue(cTrimmed[colonIdx+1:])
	cleanKey := strings.TrimPrefix(key, "_")

	// 1) Underscore missing on slide-local directive
	if !strings.HasPrefix(key, "_") {
		canonical, ok := MatchCanonicalDirective(cleanKey)
		if !ok {
			return nil // General user comment (memo:, todo:, etc.) - ignore
		}

		diags := []model.Diagnostic{{
			Severity:   model.SeverityWarning,
			SlideIndex: slideIndex,
			Line:       lineNum,
			Rule:       "directive.missing_underscore",
			Message:    "Slide-local directive missing underscore prefix",
			RawSnippet: rawSnippet,
			Candidates: []string{canonical},
		}}

		if strings.EqualFold(cleanKey, "layout") && !model.IsKnownLayout(val) {
			diags = append(diags, model.Diagnostic{
				Severity:   model.SeverityWarning,
				SlideIndex: slideIndex,
				Line:       lineNum,
				Rule:       "layout.unknown",
				Message:    fmt.Sprintf("Unknown layout '%s'", val),
				RawSnippet: rawSnippet,
				Candidates: model.AvailableLayouts(),
			})
		}
		return diags
	}

	// 2) Directive key is prefixed with '_'
	if strings.EqualFold(cleanKey, "layout") {
		if !model.IsKnownLayout(val) {
			return []model.Diagnostic{{
				Severity:   model.SeverityWarning,
				SlideIndex: slideIndex,
				Line:       lineNum,
				Rule:       "layout.unknown",
				Message:    fmt.Sprintf("Unknown layout '%s'", val),
				RawSnippet: rawSnippet,
				Candidates: model.AvailableLayouts(),
			}}
		}
		return nil
	}

	// 3) Check if unknown '_' directive key
	if _, ok := MatchCanonicalDirective(cleanKey); !ok {
		return []model.Diagnostic{{
			Severity:   model.SeverityWarning,
			SlideIndex: slideIndex,
			Line:       lineNum,
			Rule:       "directive.unknown",
			Message:    fmt.Sprintf("Unknown directive '%s'", key),
			RawSnippet: rawSnippet,
			Candidates: AvailableDirectives(),
		}}
	}

	return nil
}
