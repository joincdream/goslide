package theme

import "embed"

// embeddedAssets holds all built-in CSS stylesheets, runtime JS, and starter templates packaged into the binary.
//
//go:embed assets/css/*.css assets/js/*.js assets/templates/*.md
var embeddedAssets embed.FS
