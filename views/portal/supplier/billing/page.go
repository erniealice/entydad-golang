package billing

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the supplier billing page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the supplier billing view (GET /portal/supplier/billing).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.Supplier.Billing.PageTitle == "" {
		labels.Supplier.Billing.PageTitle = "Billing"
	}
	if labels.Supplier.Name == "" {
		labels.Supplier.Name = "Supplier Portal"
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
					Title:           labels.Supplier.Billing.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-billing-content",
					ActiveNav:       "billing",
				},
				PrincipalKind: "supplier",
				PortalName:    labels.Supplier.Name,
				ProfileURL:    "/portal/supplier/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-supplier-billing", pageData)
	})
}
