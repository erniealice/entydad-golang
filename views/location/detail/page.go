package detail

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       entydad.LocationRoutes
	ReadLocation func(ctx context.Context, req *locationpb.ReadLocationRequest) (*locationpb.ReadLocationResponse, error)
	Labels       entydad.LocationLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the location detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.LocationLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	ID              string
	LocationName    string
	LocationAddress string
	LocationDesc    string
	LocationStatus  string
	StatusVariant   string
	EditDetailURL   string
}

// NewView creates the location detail view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("location-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial -- returns only the tab content).
// Handles GET /action/locations/{id}/tab/{tab}
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		// Return only the tab partial template
		templateName := "location-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// buildPageData loads location data and builds the PageData for the given active tab.
func buildPageData(ctx context.Context, deps *Deps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadLocation(ctx, &locationpb.ReadLocationRequest{
		Data: &locationpb.Location{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read location %s: %v", id, err)
		return nil, fmt.Errorf("failed to load location: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("location not found")
	}
	loc := data[0]

	name := loc.GetName()
	address := loc.GetAddress()
	description := loc.GetDescription()

	locationStatus := "active"
	if !loc.GetActive() {
		locationStatus = "inactive"
	}
	statusVariant := "success"
	if locationStatus == "inactive" {
		statusVariant = "warning"
	}

	tabItems := buildTabItems(id, deps.Labels, deps.Routes)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "locations",
			ActiveSubNav:   "locations-active",
			HeaderTitle:    name,
			HeaderSubtitle: address,
			HeaderIcon:     "icon-map-pin",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "location-detail-content",
		Labels:          deps.Labels,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		LocationName:    name,
		LocationAddress: address,
		LocationDesc:    description,
		LocationStatus:  locationStatus,
		StatusVariant:   statusVariant,
		EditDetailURL:   route.ResolveURL(deps.Routes.EditDetailURL, "id", id),
	}

	return pageData, nil
}

func buildTabItems(id string, labels entydad.LocationLabels, routes entydad.LocationRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: labels.Detail.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "users", Label: labels.Detail.Tabs.Users, Href: base + "?tab=users", HxGet: action + "users", Icon: "icon-users", Count: 0, Disabled: false},
		{Key: "pricelists", Label: labels.Detail.Tabs.PriceLists, Href: base + "?tab=pricelists", HxGet: action + "pricelists", Icon: "icon-tag", Count: 0, Disabled: false},
		{Key: "audit", Label: labels.Detail.Tabs.AuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock", Count: 0, Disabled: false},
	}
}
