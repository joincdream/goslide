package parser

import (
	"strings"
	"testing"
)

func TestTransformImages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain image without directives",
			input:    `<img src="foo.png" alt="Just a photo">`,
			expected: `<img src="foo.png" alt="Just a photo">`,
		},
		{
			name:     "width directive in px",
			input:    `<img src="foo.png" alt="width:400px Photo">`,
			expected: `<img src="foo.png" alt="Photo" style="width: 400px;">`,
		},
		{
			name:     "width directive bare digits",
			input:    `<img src="foo.png" alt="w:300 Photo">`,
			expected: `<img src="foo.png" alt="Photo" style="width: 300px;">`,
		},
		{
			name:     "width and height and center",
			input:    `<img src="foo.png" alt="width:500px height:300px center Diagram">`,
			expected: `<img src="foo.png" alt="Diagram" style="width: 500px; height: 300px;" class="img-center">`,
		},
		{
			name:     "percentage width and center only",
			input:    `<img src="foo.png" alt="center width:80%">`,
			expected: `<img src="foo.png" alt="" style="width: 80%;" class="img-center">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformImages(tt.input)
			if got != tt.expected {
				t.Errorf("transformImages() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTransformAlerts(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    string
		wantTitle   string
		shouldMatch bool
	}{
		{
			name:        "plain blockquote without alert",
			input:       "<blockquote>\n<p>This is a quote</p>\n</blockquote>",
			shouldMatch: false,
		},
		{
			name:        "unregistered alert tag fallbacks to plain blockquote",
			input:       "<blockquote>\n<p>[!UNKNOWN] Not supported<br>Content</p>\n</blockquote>",
			shouldMatch: false,
		},
		{
			name:        "note alert default title",
			input:       "<blockquote>\n<p>[!NOTE]\nThis is a note content.</p>\n</blockquote>",
			wantType:    "note",
			wantTitle:   "Note",
			shouldMatch: true,
		},
		{
			name:        "warning alert custom title",
			input:       "<blockquote>\n<p>[!WARNING] Danger Ahead<br />\nPlease proceed with caution.</p>\n</blockquote>",
			wantType:    "warning",
			wantTitle:   "Danger Ahead",
			shouldMatch: true,
		},
		{
			name:        "tip alert",
			input:       "<blockquote>\n<p>[!TIP] Pro Tip<br>Helpful advice.</p>\n</blockquote>",
			wantType:    "tip",
			wantTitle:   "Pro Tip",
			shouldMatch: true,
		},
		{
			name:        "important alert",
			input:       "<blockquote>\n<p>[!IMPORTANT] Key Principle<br>Pure Go.</p>\n</blockquote>",
			wantType:    "important",
			wantTitle:   "Key Principle",
			shouldMatch: true,
		},
		{
			name:        "caution alert",
			input:       "<blockquote>\n<p>[!CAUTION] Critical Alert<br>Do not delete.</p>\n</blockquote>",
			wantType:    "caution",
			wantTitle:   "Critical Alert",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformAlerts(tt.input)
			if !tt.shouldMatch {
				if got != tt.input {
					t.Errorf("expected unchanged, got %q", got)
				}
				return
			}

			if !strings.Contains(got, "markdown-alert-"+tt.wantType) {
				t.Errorf("expected alert class %q in %q", tt.wantType, got)
			}
			if !strings.Contains(got, tt.wantTitle) {
				t.Errorf("expected title %q in %q", tt.wantTitle, got)
			}
		})
	}
}
