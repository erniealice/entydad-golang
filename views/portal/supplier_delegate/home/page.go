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

// PageData holds data for the supplier-delegate home page.
type PageData struct {
	portal.PageData
	Labels entydad.PortalLabels
	Cards  []NavCard
}

// Deps holds view dependencies.
type Deps struct {
	Labels entydad.PortalLabels
}

var defaultSupplierDelegateNavCards = []NavCard{
	{Key: "billing", URL: "/portal/supplier-delegate/billing", Icon: "icon-credit-card"},
	{Key: "profile", URL: "/portal/supplier-delegate/profile", Icon: "icon-user"},
	{Key: "account", URL: "/portal/supplier-delegate/account", Icon: "icon-shield"},
	{Key: "preferences", URL: "/portal/supplier-delegate/preferences", Icon: "icon-settings"},
}

// NewView creates the supplier-delegate home view (GET /portal/supplier-delegate/).
func NewView(deps *Deps) view.View {
	var labels entydad.PortalLabels
	if deps != nil {
		labels = deps.Labels
	}
	if labels.SupplierDelegate.Home.PageTitle == "" {
		labels.SupplierDelegate.Home.PageTitle = "Delegate Portal — Home"
	}
	if labels.SupplierDelegate.Home.Heading == "" {
		labels.SupplierDelegate.Home.Heading = "Welcome to your delegate portal"
	}
	if labels.SupplierDelegate.Home.Subheading == "" {
		labels.SupplierDelegate.Home.Subheading = "You are acting on behalf of a supplier. Choose a section to continue."
	}
	if labels.SupplierDelegate.Name == "" {
		labels.SupplierDelegate.Name = "Delegate Portal"
	}
	if labels.Home.RecentActivityTitle == "" {
		labels.Home.RecentActivityTitle = "Recent activity"
	}
	if labels.Home.NoRecentActivity == "" {
		labels.Home.NoRecentActivity = "No recent activity to show."
	}
	navLabels := []string{"Billing", "Profile", "Account", "Preferences"}
	if labels.Sidebar.Nav.Billing != "" {
		navLabels[0] = labels.Sidebar.Nav.Billing
	}
	if labels.Sidebar.Nav.Profile != "" {
		navLabels[1] = labels.Sidebar.Nav.Profile
	}
	if labels.Sidebar.Nav.Account != "" {
		navLabels[2] = labels.Sidebar.Nav.Account
	}
	if labels.Sidebar.Nav.Preferences != "" {
		navLabels[3] = labels.Sidebar.Nav.Preferences
	}
	navCards := make([]NavCard, len(defaultSupplierDelegateNavCards))
	copy(navCards, defaultSupplierDelegateNavCards)
	for i, lbl := range navLabels {
		navCards[i].Label = lbl
	}

	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: portal.PageData{
				PageData: types.PageData{
					CacheVersion:    viewCtx.CacheVersion,
					Title:           labels.SupplierDelegate.Home.PageTitle,
					CurrentPath:     viewCtx.CurrentPath,
					ContentTemplate: "portal-supplier-delegate-home-content",
					ActiveNav:       "home",
				},
				PrincipalKind: "supplier-delegate",
				PortalName:    labels.SupplierDelegate.Name,
				ProfileURL:    "/portal/supplier-delegate/profile",
			},
			Labels: labels,
			Cards:  navCards,
		}
		return view.OK("portal-supplier-delegate-home", pageData)
	})
}
