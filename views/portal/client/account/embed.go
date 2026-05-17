// Package account is the account stub for the client portal (/portal/client/account).
package account

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
