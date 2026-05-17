package profile

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the client-delegate profile page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the client-delegate profile view (GET /portal/client-delegate).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.ClientDelegate.Profile.PageTitle == "" {
		labels.ClientDelegate.Profile.PageTitle = "Profile"
	}
	if labels.ClientDelegate.Name == "" {
		labels.ClientDelegate.Name = "Delegate Portal"
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
					Title:           labels.ClientDelegate.Profile.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-delegate-profile-content",
					ActiveNav:       "profile",
				},
				PrincipalKind: "client-delegate",
				PortalName:    labels.ClientDelegate.Name,
				ProfileURL:    "/portal/client-delegate/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-client-delegate-profile", pageData)
	})
}
