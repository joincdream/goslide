package pptx

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/yundream/goslide/internal/browser"
	"github.com/yundream/goslide/internal/model"
)

// CapturedSlide holds the rendered slide PNG image buffer and clickable link geometries.
type CapturedSlide struct {
	Image []byte
	Links []SlideLink
}

// CaptureOptions defines configuration parameters for headless slide capture.
type CaptureOptions struct {
	Timeout     time.Duration
	BrowserPath string
}

// DefaultCaptureOptions returns default options for slide image capture.
func DefaultCaptureOptions() CaptureOptions {
	return CaptureOptions{
		Timeout: 60 * time.Second,
	}
}

const pptxCaptureCSS = `
html, body {
	margin: 0 !important;
	padding: 0 !important;
	width: 1920px !important;
	height: 1080px !important;
	overflow: hidden !important;
}
.goslide-stage, .goslide-deck {
	width: 1920px !important;
	height: 1080px !important;
	margin: 0 !important;
	padding: 0 !important;
	overflow: hidden !important;
	transform: none !important;
}
.slide-card {
	box-shadow: none !important;
	border-radius: 0 !important;
	border: none !important;
}
#goslide-app, .goslide-canvas, .goslide-indicator, .slide-badge, .slide-controls {
	display: none !important;
}
`

// CaptureSlides launches headless Chrome, iterates through each slide in the intermediate HTML,
// and returns the 1920x1080 PNG screenshot buffers and clickable link geometries.
func CaptureSlides(ctx context.Context, htmlFilePath string, slideCount int, opts CaptureOptions) ([]CapturedSlide, error) {
	if slideCount <= 0 {
		return nil, nil
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	browserPath := opts.BrowserPath
	if browserPath == "" {
		var err error
		browserPath, err = browser.FindChrome("")
		if err != nil {
			return nil, err
		}
	}

	absHTMLPath, err := filepath.Abs(htmlFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute HTML path: %w", err)
	}
	fileURL := "file://" + absHTMLPath

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
		chromedp.Flag("headless", "new"),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("allow-file-access-from-files", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, allocOpts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	timeoutCtx, cancelTimeout := context.WithTimeout(browserCtx, timeout)
	defer cancelTimeout()

	// Initial DOM load, styling injection, and layout stabilization
	prepJS := fmt.Sprintf(`(function() {
		var style = document.createElement('style');
		style.id = 'goslide-pptx-capture-style';
		style.textContent = %q;
		document.head.appendChild(style);
	})();`, pptxCaptureCSS)

	err = chromedp.Do(timeoutCtx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.Evaluate[chromedp.Void](prepJS),
		chromedp.Sleep(200*time.Millisecond),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to initialize headless browser for PPTX capture: %w", model.ErrExportFailed, err)
	}

	slides := make([]CapturedSlide, 0, slideCount)
	for i := 0; i < slideCount; i++ {
		s, err := captureSingleSlide(timeoutCtx, i)
		if err != nil {
			return nil, fmt.Errorf("%w: slide %d capture failed: %w", model.ErrExportFailed, i+1, err)
		}
		slides = append(slides, s)
	}

	return slides, nil
}

const extractLinksJS = `(function() {
	var links = [];
	var els = document.querySelectorAll('.slide-card.active .goslide-youtube-link');
	els.forEach(function(el) {
		var rect = el.getBoundingClientRect();
		if (rect.width > 0 && rect.height > 0 && el.href) {
			links.push({
				url: el.href,
				x: Math.round(rect.x * 6350),
				y: Math.round(rect.y * 6350),
				w: Math.round(rect.width * 6350),
				h: Math.round(rect.height * 6350)
			});
		}
	});
	return JSON.stringify(links);
})()`

func captureSingleSlide(ctx context.Context, slideIdx int) (CapturedSlide, error) {
	activateJS := fmt.Sprintf(`(function(targetIdx) {
		var cards = document.querySelectorAll('.slide-card');
		cards.forEach(function(card, idx) {
			if (idx === targetIdx) {
				card.classList.add('active');
				card.style.visibility = 'visible';
				card.style.opacity = '1';
				card.style.zIndex = '100';
			} else {
				card.classList.remove('active');
				card.style.visibility = 'hidden';
				card.style.opacity = '0';
				card.style.zIndex = '0';
			}
		});
	})(%d);`, slideIdx)

	err := chromedp.Do(ctx,
		chromedp.Evaluate[chromedp.Void](activateJS),
		chromedp.Sleep(50*time.Millisecond),
	)
	if err != nil {
		return CapturedSlide{}, fmt.Errorf("activation script error: %w", err)
	}

	var rawJSON string
	rawJSON, _ = chromedp.Run(ctx, chromedp.Evaluate[string](extractLinksJS))
	var links []SlideLink
	if rawJSON != "" {
		_ = json.Unmarshal([]byte(rawJSON), &links)
	}

	imgBuf, err := chromedp.Run(ctx, chromedp.CaptureScreenshot())
	if err != nil {
		return CapturedSlide{}, fmt.Errorf("screenshot capture failed: %w", err)
	}
	if len(imgBuf) == 0 {
		return CapturedSlide{}, fmt.Errorf("captured empty image buffer")
	}

	return CapturedSlide{Image: imgBuf, Links: links}, nil
}
