// Package billing is the billing stub for the supplier portal (/portal/supplier/billing).
package billing

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
