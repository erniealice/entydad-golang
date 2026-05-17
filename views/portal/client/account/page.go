package account

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the client account page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the client account view (GET /portal/client/account).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.Client.Account.PageTitle == "" {
		labels.Client.Account.PageTitle = "Account"
	}
	if labels.Client.Name == "" {
		labels.Client.Name = "My Account"
	}
	if labels.Page.Account.Title == "" {
		labels.Page.Account.Title = "Account & security"
	}
	if labels.Page.Account.ComingSoon == "" {
		labels.Page.Account.ComingSoon = "Email address and session management — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.Client.Account.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-account-content",
					ActiveNav:       "account",
				},
				PrincipalKind: "client",
				PortalName:    labels.Client.Name,
				ProfileURL:    "/portal/client/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-client-account", pageData)
	})
}
