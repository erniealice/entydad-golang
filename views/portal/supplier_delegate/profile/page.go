package profile

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the supplier-delegate profile page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the supplier-delegate profile view (GET /portal/supplier-delegate).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.SupplierDelegate.Profile.PageTitle == "" {
		labels.SupplierDelegate.Profile.PageTitle = "Profile"
	}
	if labels.SupplierDelegate.Name == "" {
		labels.SupplierDelegate.Name = "Delegate Portal"
	}
	if labels.Page.Profile.Title == "" {
		labels.Page.Profile.Title = "Your profile"
	}
	if labels.Page.Profile.ComingSoon == "" {
		labels.Page.Profile.ComingSoon = "Delegate profile information — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.SupplierDelegate.Profile.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-delegate-profile-content",
					ActiveNav:       "profile",
				},
				PrincipalKind: "supplier-delegate",
				PortalName:    labels.SupplierDelegate.Name,
				ProfileURL:    "/portal/supplier-delegate/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-supplier-delegate-profile", pageData)
	})
}
