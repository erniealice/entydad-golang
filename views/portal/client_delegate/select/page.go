// Package selectclient is the acting-as picker for CLIENT_DELEGATE principals.
//
// Shown when a CLIENT_DELEGATE principal represents 2+ clients. The user picks
// which client to act as; picking POSTs to /action/auth/switch-principal with:
//
//	principal_id=<current delegate grant id>
//	principal_kind=CLIENT_DELEGATE
//	acting_as_client_id=<chosen client id>
//
// The switch-principal handler (domain_auth.go) rotates the session and
// redirects to /portal/client-delegate/.
package selectclient

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ClientOption is one selectable client card on the picker page.
type ClientOption struct {
	// ClientID is the client entity id — sent as acting_as_client_id.
	ClientID string
	// DisplayName is the client's human-readable name shown on the card.
	DisplayName string
}

// PageData holds data for the acting-as picker page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
	// PrincipalID is the delegate grant id — sent as principal_id on POST.
	PrincipalID   string
	SwitchPostURL string
	Clients       []ClientOption
}

// Deps holds view dependencies for the select view.
type Deps struct {
	Labels entydad.PortalLabels
	// ResolveClients returns the selectable client options for the current
	// request. Called once per request in NewView; the host wires the
	// correct resolver at startup.
	ResolveClients func(ctx context.Context) []ClientOption
	// ResolvePrincipalID returns the delegate grant id from the session
	// (used as principal_id in the POST form).
	ResolvePrincipalID func(ctx context.Context) string
}

// NewView creates the acting-as client picker view (GET /portal/client-delegate/select).
func NewView(deps *Deps) view.View {
	switchPostURL := "/action/auth/switch-principal"
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.ClientDelegate.Select.PageTitle == "" {
		labels.ClientDelegate.Select.PageTitle = "Choose a client"
	}
	if labels.ClientDelegate.Name == "" {
		labels.ClientDelegate.Name = "Delegate Portal"
	}
	if labels.ClientDelegate.Select.EmptyState == "" {
		labels.ClientDelegate.Select.EmptyState = "No client accounts found. Contact your administrator."
	}
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		var clients []ClientOption
		if deps != nil && deps.ResolveClients != nil {
			clients = deps.ResolveClients(ctx)
		}
		principalID := ""
		if deps != nil && deps.ResolvePrincipalID != nil {
			principalID = deps.ResolvePrincipalID(ctx)
		}

		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.ClientDelegate.Select.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-delegate-select-content",
					ActiveNav:       "home",
				},
				PrincipalKind: "client-delegate",
				PortalName:    labels.ClientDelegate.Name,
				ProfileURL:    "/portal/client-delegate/profile",
			},
			Labels:        labels,
			PrincipalID:   principalID,
			SwitchPostURL: switchPostURL,
			Clients:       clients,
		}
		return view.OK("portal-client-delegate-select", pageData)
	})
}
