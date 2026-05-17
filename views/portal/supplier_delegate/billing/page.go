package billing

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the supplier-delegate billing page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the supplier-delegate billing view (GET /portal/supplier-delegate/billing).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.SupplierDelegate.Billing.PageTitle == "" {
		labels.SupplierDelegate.Billing.PageTitle = "Billing"
	}
	if labels.SupplierDelegate.Name == "" {
		labels.SupplierDelegate.Name = "Delegate Portal"
	}
	if labels.Page.Billing.Title == "" {
		labels.Page.Billing.Title = "Billing"
	}
	if labels.Page.Billing.ComingSoon == "" {
		labels.Page.Billing.ComingSoon = "Acting-as supplier billing context — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.SupplierDelegate.Billing.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-delegate-billing-content",
					ActiveNav:       "billing",
				},
				PrincipalKind: "supplier-delegate",
				PortalName:    labels.SupplierDelegate.Name,
				ProfileURL:    "/portal/supplier-delegate/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-supplier-delegate-billing", pageData)
	})
}
