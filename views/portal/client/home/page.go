package home

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	portal "github.com/erniealice/entydad-golang/views/portal"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// NavCard is one quick-link card on the home dashboard grid.
type NavCard struct {
	Key   string // used in data-testid="portal-client-card-{key}"
	Label string
	URL   string
	Icon  string
}

// PageData holds data for the client home page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
	Cards  []NavCard
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

var defaultClientNavCards = []NavCard{
	{Key: "profile", URL: "/portal/client/profile", Icon: "icon-user"},
	{Key: "account", URL: "/portal/client/account", Icon: "icon-shield"},
	{Key: "billing", URL: "/portal/client/billing", Icon: "icon-credit-card"},
	{Key: "preferences", URL: "/portal/client/preferences", Icon: "icon-settings"},
}

// NewView creates the client home view (GET /portal/client/).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	// Apply English defaults when labels are not wired yet.
	if labels.Client.Home.PageTitle == "" {
		labels.Client.Home.PageTitle = "My Account — Home"
	}
	if labels.Client.Home.Heading == "" {
		labels.Client.Home.Heading = "Welcome to your account"
	}
	if labels.Client.Home.Subheading == "" {
		labels.Client.Home.Subheading = "Choose a section to get started."
	}
	if labels.Client.Name == "" {
		labels.Client.Name = "My Account"
	}
	if labels.Home.RecentActivityTitle == "" {
		labels.Home.RecentActivityTitle = "Recent activity"
	}
	if labels.Home.NoRecentActivity == "" {
		labels.Home.NoRecentActivity = "No recent activity to show."
	}
	// Nav card labels
	navLabels := []string{"Profile", "Account", "Billing", "Preferences"}
	if labels.Sidebar.Nav.Profile != "" {
		navLabels[0] = labels.Sidebar.Nav.Profile
	}
	if labels.Sidebar.Nav.Account != "" {
		navLabels[1] = labels.Sidebar.Nav.Account
	}
	if labels.Sidebar.Nav.Billing != "" {
		navLabels[2] = labels.Sidebar.Nav.Billing
	}
	if labels.Sidebar.Nav.Preferences != "" {
		navLabels[3] = labels.Sidebar.Nav.Preferences
	}
	navCards := make([]NavCard, len(defaultClientNavCards))
	copy(navCards, defaultClientNavCards)
	for i, lbl := range navLabels {
		navCards[i].Label = lbl
	}

	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.Client.Home.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-client-home-content",
					ActiveNav:       "home",
				},
				PrincipalKind: "client",
				PortalName:    labels.Client.Name,
				ProfileURL:    "/portal/client/profile",
			},
			Labels: labels,
			Cards:  navCards,
		}
		return view.OK("portal-client-home", pageData)
	})
}
