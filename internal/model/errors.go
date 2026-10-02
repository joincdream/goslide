package model

import "errors"

// Domain sentinel errors.
var (
	ErrSlideNotFound      = errors.New("slide not found")
	ErrInvalidFrontmatter = errors.New("invalid frontmatter syntax")
	ErrThemeNotFound      = errors.New("theme not found")
	ErrBrowserNotFound    = errors.New("headless browser executable not found")
	ErrExportFailed       = errors.New("export operation failed")
	ErrCanceled           = errors.New("operation canceled")
)

// CLI standard exit codes.
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitInvalidUsage = 2
	ExitFileNotFound = 3
	ExitParseError   = 4
	ExitExportFailed = 5
)
