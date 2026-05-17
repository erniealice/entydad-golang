// Package home is the dashboard stub for the supplier-delegate portal (/portal/supplier-delegate/).
package home

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
