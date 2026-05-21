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
	// User holds the logged-in user's basics (name, email, mobile).
	// Populated by portal pages that load via ReadUser from espyna
	// (e.g. profile, account). Empty struct otherwise — templates
	// guard with {{if .User.FirstName}}.
	//
	// Added 2026-05-22 as the first Pre-B wave per
	// docs/plan/20260516-self-domain/ §P6 and
	// docs/plan/20260521-workspace-keyed-routing/phases.md Pre-B.
	User ProfileUser
}

// ProfileUser holds the User basics rendered by portal profile/account pages.
// Mirror of userpb.User but flattened to template-friendly shape (the proto's
// optional/pointer fields confuse Go templates, so we resolve them here).
type ProfileUser struct {
	ID           string
	FirstName    string
	LastName     string
	EmailAddress string
	MobileNumber string
}
