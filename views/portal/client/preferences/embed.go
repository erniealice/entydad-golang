// Package preferences is the preferences stub for the client portal (/portal/client/preferences).
package preferences

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
