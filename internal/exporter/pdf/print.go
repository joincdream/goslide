package pdf

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// PrintOptions defines options for printing slides to PDF via Chrome DevTools Protocol.
type PrintOptions struct {
	Timeout     time.Duration
	BrowserPath string
}

// DefaultPrintOptions returns recommended defaults for PDF export.
func DefaultPrintOptions() PrintOptions {
	return PrintOptions{
		Timeout: 30 * time.Second,
	}
}

// PrintHTMLToPDF launches a headless Chrome instance and prints the given HTML file to a vector PDF buffer.
func PrintHTMLToPDF(ctx context.Context, htmlFilePath string, opts PrintOptions) ([]byte, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}

	browserPath := opts.BrowserPath
	if browserPath == "" {
		var err error
		browserPath, err = FindChrome("")
		if err != nil {
			return nil, err
		}
	}

	absHTMLPath, err := filepath.Abs(htmlFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute HTML path: %w", err)
	}
	slashPath := filepath.ToSlash(absHTMLPath)
	if !strings.HasPrefix(slashPath, "/") {
		slashPath = "/" + slashPath
	}
	fileURL := "file://" + slashPath

	// 1. Configure Allocator with headless Chrome flags
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

	// 2. Create Browser Context
	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	// 3. Apply Timeout
	timeoutCtx, cancelTimeout := context.WithTimeout(browserCtx, opts.Timeout)
	defer cancelTimeout()

	// 4. Navigate & Wait for Render stabilization
	err = chromedp.Do(timeoutCtx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(200*time.Millisecond),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to render HTML in headless browser: %w", err)
	}

	// 5. Execute 16:9 zero-margin vector print
	pdfBuf, err := chromedp.Run(timeoutCtx, chromedp.PrintToPDF(
		chromedp.PDFLandscape(),
		chromedp.PDFPrintBackground(),
		chromedp.PDFPreferCSSPageSize(),
		chromedp.PDFMargin(0),
	))
	if err != nil {
		return nil, fmt.Errorf("chromedp PrintToPDF failed: %w", err)
	}

	if len(pdfBuf) == 0 {
		return nil, fmt.Errorf("generated PDF buffer is empty")
	}

	return pdfBuf, nil
}
