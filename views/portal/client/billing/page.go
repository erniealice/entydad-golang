package billing

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the client billing page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the client billing view (GET /portal/client/billing).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.Client.Billing.PageTitle == "" {
		labels.Client.Billing.PageTitle = "Billing"
	}
	if labels.Client.Name == "" {
		labels.Client.Name = "My Account"
	}
	if labels.Page.Billing.Title == "" {
		labels.Page.Billing.Title = "Billing"
	}
	if labels.Page.Billing.ComingSoon == "" {
		labels.Page.Billing.ComingSoon = "Your billing context — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.Client.Billing.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-billing-content",
					ActiveNav:       "billing",
				},
				PrincipalKind: "client",
				PortalName:    labels.Client.Name,
				ProfileURL:    "/portal/client/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-client-billing", pageData)
	})
}
