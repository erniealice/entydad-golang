// Package home is the dashboard stub for the client-delegate portal (/portal/client-delegate/).
package home

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
