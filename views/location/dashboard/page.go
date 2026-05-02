// Package dashboard is the view layer for the location app dashboard
// (Phase 4 of the Pyeza dashboard plan). It calls the
// GetLocationDashboardPageData use case via a callback in Deps and projects
// the response into a typed pyeza DashboardData rendered by the shared
// "dashboard" component block.
package dashboard

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
)

// LocationAreaCount is one row of the "top areas by location count" widget.
//
// Mirrors the use case's typed row, kept here so the view layer does not
// import the espyna use-case package directly. The Deps callback returns
// LocationAreaCount values populated from the adapter aggregate.
type LocationAreaCount struct {
	LocationAreaID   string
	LocationAreaName string
	LocationCount    int64
}

// LocationDashboardData is the response shape consumed by the view.
//
// The orchestrator's callback (wired in service-admin) projects the espyna
// use-case response onto this struct so the view stays free of espyna
// imports — same pattern as the user/inventory/procurement dashboards.
type LocationDashboardData struct {
	TotalLocations    int64
	ActiveLocations   int64
	RegionsCount      int64
	AreasCount        int64
	LocationsByRegion map[string]int64
	TopAreas          []LocationAreaCount
	RecentLocations   []*locationpb.Location
}

// Deps holds view dependencies.
type Deps struct {
	DashboardLabels entydad.DashboardLabels
	Dashboard       entydad.LocationDashboardLabels
	Routes          entydad.LocationRoutes
	CommonLabels    pyeza.CommonLabels

	// GetDashboardData is the workspace-scoped page-data fetch. The
	// container constructs this by calling the
	// GetLocationDashboardPageDataUseCase. nil-safe: when missing, the
	// view renders empty-state widgets.
	GetDashboardData func(ctx context.Context) (*LocationDashboardData, error)
}

// PageData holds the data for the location dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the location dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Dashboard

		var data *LocationDashboardData
		if deps.GetDashboardData != nil {
			d, err := deps.GetDashboardData(ctx)
			if err != nil {
				log.Printf("location dashboard: failed to load page data: %v", err)
			}
			data = d
		}
		if data == nil {
			data = &LocationDashboardData{}
		}

		// Stats — 4 cards.
		stats := []types.StatCardData{
			{
				Icon: "icon-map-pin", Value: strconv.FormatInt(data.TotalLocations, 10),
				Label: l.TotalLocations, Color: "terracotta",
				TestID: "location-stat-total",
			},
			{
				Icon: "icon-check-circle", Value: strconv.FormatInt(data.ActiveLocations, 10),
				Label: l.Active, Color: "sage",
				TestID: "location-stat-active",
			},
			{
				Icon: "icon-globe", Value: strconv.FormatInt(data.RegionsCount, 10),
				Label: l.Regions, Color: "navy",
				TestID: "location-stat-regions",
			},
			{
				Icon: "icon-layers", Value: strconv.FormatInt(data.AreasCount, 10),
				Label: l.AreasCount, Color: "amber",
				TestID: "location-stat-areas",
			},
		}

		// Bar chart: locations per region/area.
		regionLabels, regionValues := projectRegions(data.LocationsByRegion)
		regionChart := &types.ChartData{
			Labels: regionLabels,
			Series: []types.ChartSeries{{
				Name:   l.LocationsByRegion,
				Values: regionValues,
				Color:  "terracotta",
			}},
		}
		regionChart.AutoScale()

		// Table widget — top areas by location count. Pyeza has no typed
		// generic table widget shape that fits this aggregate cleanly, so
		// we follow the procurement-dashboard pattern (Type:"custom" with
		// raw template.HTML), which the dashboard plan calls out as the
		// fallback for tables that don't slot into TableConfig.
		topAreasHTML := buildTopAreasHTML(data.TopAreas, l)

		// Activity list — recent additions.
		recentItems := buildRecentList(data.RecentLocations, l)

		dash := types.DashboardData{
			Title:    deps.DashboardLabels.LocationTitle,
			Icon:     "icon-map-pin",
			Subtitle: l.LocationsByRegion,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNewLocation, Href: deps.Routes.AddURL, Variant: "primary", TestID: "location-action-new"},
				{Icon: "icon-layers", Label: l.QuickNewArea, Href: "/app/location-areas/dashboard", TestID: "location-action-new-area"},
				{Icon: "icon-truck", Label: l.QuickMoveStock, Href: "#", TestID: "location-action-move-stock"},
				{Icon: "icon-map", Label: l.QuickLocationMap, Href: "#", TestID: "location-action-map"},
			},
			Stats: stats,
			Widgets: []types.DashboardWidget{
				{
					ID: "by-region", Title: l.LocationsByRegion,
					Type: "chart", ChartKind: "bar",
					ChartData: regionChart, Span: 2,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-globe",
						Title: l.LocationsByRegion,
					},
				},
				{
					ID: "top-areas", Title: l.TopLocationsByArea,
					Type: "custom", Span: 1,
					Custom: topAreasHTML,
				},
				{
					ID: "recent", Title: l.RecentAdditions,
					Type: "list", Span: 3,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ListURL},
					},
					ListItems: recentItems,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-clock",
						Title: l.RecentAdditions,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.DashboardLabels.LocationTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "location",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.DashboardLabels.LocationTitle,
				HeaderIcon:   "icon-map-pin",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "location-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("location-dashboard", pageData)
	})
}

// projectRegions converts the use case's region map into chart-ready
// parallel slices, preserving descending count order.
func projectRegions(byRegion map[string]int64) ([]string, []float64) {
	if len(byRegion) == 0 {
		return []string{"-"}, []float64{0}
	}
	type pair struct {
		label string
		count int64
	}
	pairs := make([]pair, 0, len(byRegion))
	for k, v := range byRegion {
		pairs = append(pairs, pair{label: k, count: v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].label < pairs[j].label
	})
	labels := make([]string, len(pairs))
	values := make([]float64, len(pairs))
	for i, p := range pairs {
		labels[i] = p.label
		values[i] = float64(p.count)
	}
	return labels, values
}

// buildTopAreasHTML renders the top-areas table as a small HTML fragment
// for the Type=custom widget body. Mirrors the procurement-dashboard
// pattern: a thin <table> with testid hooks for E2E selectors.
func buildTopAreasHTML(rows []LocationAreaCount, l entydad.LocationDashboardLabels) template.HTML {
	if len(rows) == 0 {
		return template.HTML(`<div class="empty-state empty-state--inline" data-testid="location-top-areas-empty">` +
			template.HTMLEscapeString(l.TopLocationsByArea) +
			`</div>`)
	}
	var sb strings.Builder
	sb.WriteString(`<table class="dashboard-mini-table" data-testid="location-top-areas-table"><thead><tr>`)
	sb.WriteString(`<th>` + template.HTMLEscapeString(l.ColumnLocation) + `</th>`)
	sb.WriteString(`<th class="num">` + template.HTMLEscapeString(l.ColumnAreas) + `</th>`)
	sb.WriteString(`</tr></thead><tbody>`)
	for _, r := range rows {
		rowID := template.HTMLEscapeString(r.LocationAreaID)
		sb.WriteString(`<tr data-testid="location-top-areas-row-` + rowID + `">`)
		sb.WriteString(`<td>` + template.HTMLEscapeString(r.LocationAreaName) + `</td>`)
		sb.WriteString(`<td class="num">` + strconv.FormatInt(r.LocationCount, 10) + `</td>`)
		sb.WriteString(`</tr>`)
	}
	sb.WriteString(`</tbody></table>`)
	return template.HTML(sb.String())
}

// buildRecentList projects recent locations into the activity-list shape.
func buildRecentList(locations []*locationpb.Location, l entydad.LocationDashboardLabels) []types.ActivityItem {
	if len(locations) == 0 {
		return nil
	}
	out := make([]types.ActivityItem, 0, len(locations))
	for i, loc := range locations {
		desc := loc.GetAddress()
		if desc == "" {
			desc = "—"
		}
		t := ""
		if loc.GetDateCreated() != 0 {
			t = formatRelative(time.UnixMilli(loc.GetDateCreated()))
		}
		out = append(out, types.ActivityItem{
			IconName:    "icon-map-pin",
			IconVariant: "client",
			Title:       fmt.Sprintf("%s — %s", l.LocationAdded, loc.GetName()),
			Description: desc,
			Time:        t,
			TestID:      fmt.Sprintf("location-activity-%d", i+1),
		})
	}
	return out
}

// formatRelative renders a coarse "Nh ago" / "Nd ago" / date label suitable
// for the activity-list time slot. We deliberately keep this view-local
// rather than depending on a shared formatter, mirroring the existing
// per-package dashboards.
func formatRelative(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	delta := time.Since(t)
	switch {
	case delta < time.Minute:
		return "just now"
	case delta < time.Hour:
		return fmt.Sprintf("%dm ago", int(delta.Minutes()))
	case delta < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(delta.Hours()))
	case delta < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(delta.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
}
