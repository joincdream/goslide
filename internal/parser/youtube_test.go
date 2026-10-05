package parser

import (
	"context"
	"strings"
	"testing"
)

func TestTransformMedia_YouTube(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		caption      string
		wantMatched  bool
		wantVideoID  string
		wantContains []string
	}{
		{
			name:        "Standard watch URL with start and end",
			url:         "https://www.youtube.com/watch?v=oDjESsByBmM&start=828&end=863",
			caption:     "Demo Video",
			wantMatched: true,
			wantVideoID: "oDjESsByBmM",
			wantContains: []string{
				`data-video-id="oDjESsByBmM"`,
				`href="https://www.youtube.com/watch?v=oDjESsByBmM&t=828s"`,
				`<figcaption class="goslide-youtube-title">Demo Video</figcaption>`,
				`img.youtube.com/vi/oDjESsByBmM/maxresdefault.jpg`,
				`class="goslide-youtube-play-btn"`,
			},
		},
		{
			name:        "Short youtu.be URL with timestamp format",
			url:         "https://youtu.be/dQw4w9WgXcQ?t=1m30s",
			caption:     "Music",
			wantMatched: true,
			wantVideoID: "dQw4w9WgXcQ",
			wantContains: []string{
				`data-video-id="dQw4w9WgXcQ"`,
				`href="https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=90s"`,
			},
		},
		{
			name:        "Standard watch URL without timestamp",
			url:         "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			caption:     "",
			wantMatched: true,
			wantVideoID: "dQw4w9WgXcQ",
			wantContains: []string{
				`data-video-id="dQw4w9WgXcQ"`,
				`href="https://www.youtube.com/watch?v=dQw4w9WgXcQ"`,
				`img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg`,
			},
		},
		{
			name:        "Regular image URL should not match",
			url:         "https://example.com/image.png",
			caption:     "Example",
			wantMatched: false,
		},
		{
			name:        "Local image path should not match",
			url:         "./assets/architecture.png",
			caption:     "Local",
			wantMatched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHTML, matched := TransformMedia(tt.url, tt.caption)
			if matched != tt.wantMatched {
				t.Fatalf("TransformMedia() matched = %v, want %v", matched, tt.wantMatched)
			}
			if !tt.wantMatched {
				return
			}
			for _, exp := range tt.wantContains {
				if !strings.Contains(gotHTML, exp) {
					t.Errorf("TransformMedia() output missing expected substring %q\nGot:\n%s", exp, gotHTML)
				}
			}
		})
	}
}

func TestParser_YouTubeMarkdownEmbed(t *testing.T) {
	p := NewParser()

	md := `---
title: "Test Deck"
---

# Video Slide

![YouTube Presentation](https://www.youtube.com/watch?v=oDjESsByBmM&start=828&end=863)

Check out [normal link](https://www.youtube.com/watch?v=oDjESsByBmM) as well.
`

	deck, err := p.Parse(context.Background(), strings.NewReader(md))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(deck.Slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(deck.Slides))
	}

	slideHTML := deck.Slides[0].HTMLContent

	// 1. Must contain responsive youtube wrapper
	if !strings.Contains(slideHTML, `class="goslide-youtube-wrapper"`) {
		t.Errorf("slide HTML does not contain goslide-youtube-wrapper\nGot:\n%s", slideHTML)
	}

	// 2. Must contain link with target="_blank" and youtube poster
	if !strings.Contains(slideHTML, `class="goslide-youtube-poster"`) || !strings.Contains(slideHTML, `target="_blank"`) {
		t.Errorf("slide HTML does not contain valid poster card link\nGot:\n%s", slideHTML)
	}

	// 3. Normal link must remain <a> tag, not replaced by youtube card
	if !strings.Contains(slideHTML, `<a href="https://www.youtube.com/watch?v=oDjESsByBmM">normal link</a>`) {
		t.Errorf("slide HTML did not preserve standard text link\nGot:\n%s", slideHTML)
	}
}
