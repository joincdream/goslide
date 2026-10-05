package parser

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// MediaHandler processes a media URL and caption, returning the responsive HTML embed and whether it matched.
type MediaHandler func(u *url.URL, caption string) (htmlContent string, matched bool)

// mediaProviders maps lowercase hostnames to their corresponding MediaHandler.
var mediaProviders = map[string]MediaHandler{
	"youtube.com":     handleYouTube,
	"www.youtube.com": handleYouTube,
	"m.youtube.com":   handleYouTube,
	"youtu.be":        handleYouTube,
}

var (
	timeUnitRegex = regexp.MustCompile(`^(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s?)?$`)
	pImgRegex     = regexp.MustCompile(`(?is)<p>\s*(<img\s+[^>]*?>)\s*</p>`)
)

// TransformMedia checks if the rawURL belongs to a registered media provider and returns the embed HTML.
func TransformMedia(rawURL, caption string) (string, bool) {
	if rawURL == "" {
		return "", false
	}

	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return "", false
	}

	host := strings.ToLower(u.Host)
	handler, exists := mediaProviders[host]
	if !exists {
		return "", false
	}

	return handler(u, caption)
}

func handleYouTube(u *url.URL, caption string) (string, bool) {
	videoID := extractYouTubeVideoID(u)
	if videoID == "" {
		return "", false
	}

	watchURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	q := u.Query()

	// Parse start time (t or start)
	startVal := q.Get("start")
	if startVal == "" {
		startVal = q.Get("t")
	}
	if startVal != "" {
		parsedStart := parseTimeString(startVal)
		if parsedStart != "" {
			watchURL = fmt.Sprintf("%s&t=%ss", watchURL, parsedStart)
		}
	}

	safeCaption := html.EscapeString(caption)
	captionHTML := ""
	ariaLabel := ""
	if safeCaption != "" {
		captionHTML = fmt.Sprintf("\n    <figcaption class=\"goslide-youtube-title\">%s</figcaption>", safeCaption)
		ariaLabel = fmt.Sprintf(": %s", safeCaption)
	}

	embedHTML := fmt.Sprintf(`<figure class="goslide-youtube-wrapper" data-video-id="%s">
    <a class="goslide-youtube-link" href="%s" target="_blank" rel="noopener noreferrer" aria-label="Play YouTube video%s">
      <div class="goslide-youtube-poster">
        <img src="https://img.youtube.com/vi/%s/maxresdefault.jpg" onerror="this.onerror=null;this.src='https://img.youtube.com/vi/%s/hqdefault.jpg';" alt="%s" class="goslide-youtube-thumb">
        <div class="goslide-youtube-play-btn" aria-hidden="true">
          <svg viewBox="0 0 68 48" class="goslide-youtube-play-icon">
            <path class="play-bg" d="M66.52,7.74c-0.78-2.93-2.49-5.41-5.42-6.19C55.79,.13,34,0,34,0S12.21,.13,6.9,1.55 C3.97,2.33,2.27,4.81,1.48,7.74C0.06,13.05,0,24,0,24s0.06,10.95,1.48,16.26c0.78,2.93,2.49,5.41,5.42,6.19 C12.21,47.87,34,48,34,48s21.79-0.13,27.1-1.55c2.93-0.78,4.64-3.26,5.42-6.19C67.94,34.95,68,24,68,24S67.94,13.05,66.52,7.74z"></path>
            <path class="play-arrow" d="M 45,24 27,14 27,34" fill="#ffffff"></path>
          </svg>
        </div>
      </div>
    </a>%s
  </figure>`, videoID, watchURL, ariaLabel, videoID, videoID, safeCaption, captionHTML)

	return embedHTML, true
}

func extractYouTubeVideoID(u *url.URL) string {
	host := strings.ToLower(u.Host)
	path := strings.TrimPrefix(u.Path, "/")

	if host == "youtu.be" {
		// Short URL: https://youtu.be/VIDEO_ID
		parts := strings.Split(path, "/")
		if len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
		return ""
	}

	// Standard URL: https://www.youtube.com/watch?v=VIDEO_ID
	if strings.HasPrefix(path, "watch") {
		return u.Query().Get("v")
	}

	// Embed URL: https://www.youtube.com/embed/VIDEO_ID
	if strings.HasPrefix(path, "embed/") {
		return strings.TrimPrefix(path, "embed/")
	}

	// Shorts URL: https://www.youtube.com/shorts/VIDEO_ID
	if strings.HasPrefix(path, "shorts/") {
		return strings.TrimPrefix(path, "shorts/")
	}

	return ""
}

func parseTimeString(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// Pure integer seconds (e.g. "828" or "828s")
	trimmedS := strings.TrimSuffix(s, "s")
	if sec, err := strconv.Atoi(trimmedS); err == nil && sec >= 0 {
		return strconv.Itoa(sec)
	}

	// Formats like "13m48s", "1h2m3s", "48s"
	matches := timeUnitRegex.FindStringSubmatch(s)
	if len(matches) == 4 {
		hours, _ := strconv.Atoi(matches[1])
		mins, _ := strconv.Atoi(matches[2])
		secs, _ := strconv.Atoi(matches[3])
		total := hours*3600 + mins*60 + secs
		if total > 0 {
			return strconv.Itoa(total)
		}
	}

	return trimmedS
}

// transformMediaElements replaces <img> tags with media embeds if the src is a supported media provider.
func transformMediaElements(htmlContent string) string {
	// 1. First replace <p><img ...></p> blocks to avoid invalid <p><figure></p> nesting
	htmlContent = pImgRegex.ReplaceAllStringFunc(htmlContent, func(pBlock string) string {
		sub := pImgRegex.FindStringSubmatch(pBlock)
		if len(sub) < 2 {
			return pBlock
		}
		imgTag := sub[1]
		src, hasSrc := getAttribute(imgTag, "src")
		if !hasSrc {
			return pBlock
		}
		alt, _ := getAttribute(imgTag, "alt")
		if embedHTML, ok := TransformMedia(src, alt); ok {
			return embedHTML
		}
		return pBlock
	})

	// 2. Replace any remaining standalone <img> tags targeting media
	htmlContent = imgTagRegex.ReplaceAllStringFunc(htmlContent, func(imgTag string) string {
		src, hasSrc := getAttribute(imgTag, "src")
		if !hasSrc {
			return imgTag
		}
		alt, _ := getAttribute(imgTag, "alt")
		if embedHTML, ok := TransformMedia(src, alt); ok {
			return embedHTML
		}
		return imgTag
	})

	return htmlContent
}
