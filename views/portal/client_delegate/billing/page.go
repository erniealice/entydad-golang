package billing

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the client-delegate billing page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the client-delegate billing view (GET /portal/client-delegate/billing).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.ClientDelegate.Billing.PageTitle == "" {
		labels.ClientDelegate.Billing.PageTitle = "Billing"
	}
	if labels.ClientDelegate.Name == "" {
		labels.ClientDelegate.Name = "Delegate Portal"
	}
	if labels.Page.Billing.Title == "" {
		labels.Page.Billing.Title = "Billing"
	}
	if labels.Page.Billing.ComingSoon == "" {
		labels.Page.Billing.ComingSoon = "Acting-as client billing context — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.ClientDelegate.Billing.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-delegate-billing-content",
					ActiveNav:       "billing",
				},
				PrincipalKind: "client-delegate",
				PortalName:    labels.ClientDelegate.Name,
				ProfileURL:    "/portal/client-delegate/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-client-delegate-billing", pageData)
	})
}
