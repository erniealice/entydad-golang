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
	Key   string
	Label string
	URL   string
	Icon  string
}

// PageData holds data for the supplier home page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
	Cards  []NavCard
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

var defaultSupplierNavCards = []NavCard{
	{Key: "profile", URL: "/portal/supplier/profile", Icon: "icon-user"},
	{Key: "account", URL: "/portal/supplier/account", Icon: "icon-shield"},
	{Key: "billing", URL: "/portal/supplier/billing", Icon: "icon-credit-card"},
	{Key: "preferences", URL: "/portal/supplier/preferences", Icon: "icon-settings"},
}

// NewView creates the supplier home view (GET /portal/supplier/).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.Supplier.Home.PageTitle == "" {
		labels.Supplier.Home.PageTitle = "Supplier Portal — Home"
	}
	if labels.Supplier.Home.Heading == "" {
		labels.Supplier.Home.Heading = "Welcome to your supplier portal"
	}
	if labels.Supplier.Home.Subheading == "" {
		labels.Supplier.Home.Subheading = "Choose a section to get started."
	}
	if labels.Supplier.Name == "" {
		labels.Supplier.Name = "Supplier Portal"
	}
	if labels.Home.RecentActivityTitle == "" {
		labels.Home.RecentActivityTitle = "Recent activity"
	}
	if labels.Home.NoRecentActivity == "" {
		labels.Home.NoRecentActivity = "No recent activity to show."
	}
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
	navCards := make([]NavCard, len(defaultSupplierNavCards))
	copy(navCards, defaultSupplierNavCards)
	for i, lbl := range navLabels {
		navCards[i].Label = lbl
	}

	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.Supplier.Home.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-home-content",
					ActiveNav:       "home",
				},
				PrincipalKind: "supplier",
				PortalName:    labels.Supplier.Name,
				ProfileURL:    "/portal/supplier/profile",
			},
			Labels: labels,
			Cards:  navCards,
		}
		return view.OK("portal-supplier-home", pageData)
	})
}
