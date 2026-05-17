package profile

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the client profile page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the client profile view (GET /portal/client/profile).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.Client.Profile.PageTitle == "" {
		labels.Client.Profile.PageTitle = "Profile"
	}
	if labels.Client.Name == "" {
		labels.Client.Name = "My Account"
	}
	if labels.Page.Profile.Title == "" {
		labels.Page.Profile.Title = "Your profile"
	}
	if labels.Page.Profile.ComingSoon == "" {
		labels.Page.Profile.ComingSoon = "Profile information — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.Client.Profile.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-profile-content",
					ActiveNav:       "profile",
				},
				PrincipalKind: "client",
				PortalName:    labels.Client.Name,
				ProfileURL:    "/portal/client/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-client-profile", pageData)
	})
}
