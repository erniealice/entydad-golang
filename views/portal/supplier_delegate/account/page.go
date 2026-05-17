package account

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the supplier-delegate account page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the supplier-delegate account view (GET /portal/supplier-delegate/account).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.SupplierDelegate.Account.PageTitle == "" {
		labels.SupplierDelegate.Account.PageTitle = "Account"
	}
	if labels.SupplierDelegate.Name == "" {
		labels.SupplierDelegate.Name = "Delegate Portal"
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
					Title:           labels.SupplierDelegate.Account.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-delegate-account-content",
					ActiveNav:       "account",
				},
				PrincipalKind: "supplier-delegate",
				PortalName:    labels.SupplierDelegate.Name,
				ProfileURL:    "/portal/supplier-delegate/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-supplier-delegate-account", pageData)
	})
}
