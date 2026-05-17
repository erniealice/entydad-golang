// Package profile is the profile stub for the client portal (/portal/client/profile).
package profile

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
