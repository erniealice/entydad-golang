// Package selectclient is the acting-as picker for CLIENT_DELEGATE principals.
// Shown when the delegate represents 2+ clients so they can choose which to act as.
package selectclient

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
