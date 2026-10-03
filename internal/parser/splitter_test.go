package parser

import (
	"reflect"
	"testing"
)

type splitterTestCase struct {
	name     string
	content  string
	expected []string
}

func runSplitterTests(t *testing.T, tests []splitterTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitSlides(tt.content)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("splitSlides mismatch:\ngot:  %#v\nwant: %#v", got, tt.expected)
			}
		})
	}
}

func TestSplitSlides_Basic(t *testing.T) {
	tests := []splitterTestCase{
		{
			name:     "single slide without delimiters",
			content:  "# Title\nThis is a single slide.",
			expected: []string{"# Title\nThis is a single slide."},
		},
		{
			name:     "two slides with standard delimiter",
			content:  "# Slide 1\n---\n# Slide 2",
			expected: []string{"# Slide 1", "# Slide 2"},
		},
		{
			name:     "three slides with trailing whitespace on delimiter",
			content:  "Slide 1\n---   \nSlide 2\n---\t\nSlide 3",
			expected: []string{"Slide 1", "Slide 2", "Slide 3"},
		},
		{
			name:     "four dashes is not a delimiter",
			content:  "Slide 1\n----\nStill Slide 1 with hr",
			expected: []string{"Slide 1\n----\nStill Slide 1 with hr"},
		},
		{
			name:     "asterisks and underscores are not delimiters",
			content:  "Slide 1\n***\nStill Slide 1\n___\nEnd of Slide 1",
			expected: []string{"Slide 1\n***\nStill Slide 1\n___\nEnd of Slide 1"},
		},
	}
	runSplitterTests(t, tests)
}

func TestSplitSlides_FencesAndEmpty(t *testing.T) {
	tests := []splitterTestCase{
		{
			name:    "fenced code block with three backticks containing separator",
			content: "Slide 1\n```yaml\n---\nkey: value\n---\n```\nStill Slide 1\n---\nSlide 2",
			expected: []string{
				"Slide 1\n```yaml\n---\nkey: value\n---\n```\nStill Slide 1",
				"Slide 2",
			},
		},
		{
			name:    "nested fenced code block with four backticks",
			content: "Slide 1\n````markdown\n```yaml\n---\n```\n````\nStill Slide 1\n---\nSlide 2",
			expected: []string{
				"Slide 1\n````markdown\n```yaml\n---\n```\n````\nStill Slide 1",
				"Slide 2",
			},
		},
		{
			name:    "tilde fenced code block containing separator",
			content: "Slide 1\n~~~bash\n---\n~~~\n---\nSlide 2",
			expected: []string{
				"Slide 1\n~~~bash\n---\n~~~",
				"Slide 2",
			},
		},
		{
			name:     "consecutive delimiters skip empty slides",
			content:  "Slide 1\n---\n\n---\nSlide 2",
			expected: []string{"Slide 1", "Slide 2"},
		},
		{
			name:     "leading and trailing delimiters skip empty slides",
			content:  "---\nSlide 1\n---\nSlide 2\n---",
			expected: []string{"Slide 1", "Slide 2"},
		},
		{
			name:     "empty content yields one empty slide",
			content:  "",
			expected: []string{""},
		},
		{
			name:     "only whitespace and delimiters yields one empty slide",
			content:  "   \n---\n   \n---\n\t\n",
			expected: []string{""},
		},
	}
	runSplitterTests(t, tests)
}
