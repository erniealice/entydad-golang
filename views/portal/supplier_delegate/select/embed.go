// Package selectsupplier is the acting-as picker for SUPPLIER_DELEGATE principals.
package selectsupplier

import "embed"

//go:embed templates/*.html
var TemplatesFS embed.FS
