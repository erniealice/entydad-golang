// Package billing is the billing stub for the client portal (/portal/client/billing).
package billing

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
