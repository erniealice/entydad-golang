package profile

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the supplier profile page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the supplier profile view (GET /portal/supplier).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.Supplier.Profile.PageTitle == "" {
		labels.Supplier.Profile.PageTitle = "Profile"
	}
	if labels.Supplier.Name == "" {
		labels.Supplier.Name = "Supplier Portal"
	}
	if labels.Page.Profile.Title == "" {
		labels.Page.Profile.Title = "Your profile"
	}
	if labels.Page.Profile.ComingSoon == "" {
		labels.Page.Profile.ComingSoon = "Supplier profile information — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.Supplier.Profile.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-profile-content",
					ActiveNav:       "profile",
				},
				PrincipalKind: "supplier",
				PortalName:    labels.Supplier.Name,
				ProfileURL:    "/portal/supplier/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-supplier-profile", pageData)
	})
}
