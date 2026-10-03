package theme

import "embed"

// embeddedAssets holds all built-in CSS stylesheets packaged into the binary.
//
//go:embed assets/css/*.css
var embeddedAssets embed.FS
