// Package home is the dashboard stub for the client portal (/portal/client/).
package home

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
