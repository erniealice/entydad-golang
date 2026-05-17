// Package home is the dashboard stub for the supplier portal (/portal/supplier/).
package home

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
