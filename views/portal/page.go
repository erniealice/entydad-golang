// Package portal holds the shared types used by all portal view sub-packages.
//
// Every portal view (client, supplier, client_delegate, supplier_delegate) embeds
// PortalPageData into its own PageData struct so it gets PrincipalKind, PortalName,
// and ProfileURL — the three fields read by portal-shell.html.
package portal

import (
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the shared base embedded by every portal view's concrete PageData.
// portal-shell.html reads PrincipalKind, PortalName, and ProfileURL directly.
type PageData struct {
	types.PageData
	// PrincipalKind is the kebab-case portal kind token
	// ("client" | "supplier" | "client-delegate" | "supplier-delegate").
	// Drives data-principal-kind on <body> and data-testid on <header>.
	PrincipalKind string
	// PortalName is the human-readable portal name shown in the header.
	PortalName string
	// ProfileURL is the portal-relative profile link in the header.
	ProfileURL string
}
