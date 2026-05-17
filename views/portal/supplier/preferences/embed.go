// Package preferences is the preferences stub for the supplier portal (/portal/supplier/preferences).
package preferences

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
