// Package me embeds the shared /me/* shell templates.
package me

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
