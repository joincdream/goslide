package theme

import "embed"

// embeddedAssets holds all built-in CSS stylesheets packaged into the binary.
//
//go:embed assets/css/*.css assets/js/*.js
var embeddedAssets embed.FS
