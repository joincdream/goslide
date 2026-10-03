package parser

import (
	"strings"
	"unicode"
)

type fenceState struct {
	inFence bool
	char    rune
	length  int
}

func (fs *fenceState) update(line string) {
	trimmedLeft := strings.TrimLeftFunc(line, unicode.IsSpace)
	leadingSpaces := len(line) - len(trimmedLeft)
	if leadingSpaces > 3 {
		return
	}

	if !fs.inFence {
		char, count := detectOpeningFence(trimmedLeft)
		if count >= 3 {
			fs.inFence = true
			fs.char = char
			fs.length = count
		}
		return
	}

	if isClosingFence(trimmedLeft, fs.char, fs.length) {
		fs.inFence = false
		fs.char = 0
		fs.length = 0
	}
}

// splitSlides splits markdown content into raw slide chunks using strict "---" delimiter.
func splitSlides(content string) []string {
	lines := strings.Split(content, "\n")

	var slides []string
	var currentLines []string
	var fs fenceState

	for _, line := range lines {
		fs.update(line)

		// When outside of a code block, check for slide separator
		if !fs.inFence && isSlideSeparator(line) {
			chunk := strings.Join(currentLines, "\n")
			if strings.TrimSpace(chunk) != "" {
				slides = append(slides, chunk)
			}
			currentLines = nil
			continue
		}

		currentLines = append(currentLines, line)
	}

	// Flush trailing slide
	if len(currentLines) > 0 {
		chunk := strings.Join(currentLines, "\n")
		if strings.TrimSpace(chunk) != "" {
			slides = append(slides, chunk)
		}
	}

	// If no slides exist (empty input or only empty delimiters), provide 1 default slide
	if len(slides) == 0 {
		slides = append(slides, "")
	}

	return slides
}

// isSlideSeparator checks if line is strictly 3 dashes with optional trailing whitespace.
func isSlideSeparator(line string) bool {
	trimmed := strings.TrimRight(line, " \t")
	return trimmed == "---"
}

// detectOpeningFence checks if trimmed line starts with at least 3 backticks or tildes.
func detectOpeningFence(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}
	first := rune(s[0])
	if first != '`' && first != '~' {
		return 0, 0
	}

	count := 0
	for _, r := range s {
		if r == first {
			count++
		} else {
			break
		}
	}

	if count < 3 {
		return 0, 0
	}
	return first, count
}

// isClosingFence checks if line is a valid closing fence (same char, >= opening len, trailing spaces only).
func isClosingFence(s string, char rune, minLen int) bool {
	if len(s) == 0 {
		return false
	}
	count := 0
	rest := ""
	for i, r := range s {
		if r == char {
			count++
		} else {
			rest = s[i:]
			break
		}
	}

	if count < minLen {
		return false
	}

	// Closing fence cannot have content other than trailing spaces
	return strings.TrimSpace(rest) == ""
}
