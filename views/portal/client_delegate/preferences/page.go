package preferences

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PageData holds data for the client-delegate preferences page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

// NewView creates the client-delegate preferences view (GET /portal/client-delegate/preferences).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.ClientDelegate.Preferences.PageTitle == "" {
		labels.ClientDelegate.Preferences.PageTitle = "Preferences"
	}
	if labels.ClientDelegate.Name == "" {
		labels.ClientDelegate.Name = "Delegate Portal"
	}
	if labels.Page.Preferences.Title == "" {
		labels.Page.Preferences.Title = "Preferences"
	}
	if labels.Page.Preferences.ComingSoon == "" {
		labels.Page.Preferences.ComingSoon = "User preference settings — coming soon."
	}
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.ClientDelegate.Preferences.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-delegate-preferences-content",
					ActiveNav:       "preferences",
				},
				PrincipalKind: "client-delegate",
				PortalName:    labels.ClientDelegate.Name,
				ProfileURL:    "/portal/client-delegate/profile",
			},
			Labels: labels,
		}
		return view.OK("portal-client-delegate-preferences", pageData)
	})
}
