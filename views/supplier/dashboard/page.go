package dashboard

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	DashboardLabels entydad.DashboardLabels
	Dashboard       entydad.SupplierDashboardLabels
	Routes          entydad.SupplierRoutes
	CommonLabels    pyeza.CommonLabels
}

// PageData holds the data for the supplier dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the supplier dashboard view.
//
// Phase 1b refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// Stat values, chart series, and activity items remain dummy until Phase 4
// wires real Supplier repository aggregate methods.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Dashboard

		trend := &types.ChartData{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
			Series: []types.ChartSeries{{
				Name:   l.SupplierActivity,
				Values: []float64{110, 118, 124, 132, 138, 142, 146, 150, 153, 154, 155, 156},
				Color:  "terracotta",
			}},
		}
		trend.AutoScale()

		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-truck", Label: l.QuickNew, Href: deps.Routes.AddURL, Variant: "primary", TestID: "supplier-action-new"},
				{Icon: "icon-list", Label: l.QuickViewAll, Href: deps.Routes.ListURL, TestID: "supplier-action-list"},
				{Icon: "icon-tag", Label: l.QuickTags, Href: "/app/suppliers/settings/tags/list", TestID: "supplier-action-tags"},
				{Icon: "icon-folder", Label: l.QuickCategories, Href: "/app/suppliers/settings/tags/list", TestID: "supplier-action-categories"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-truck", Value: "156", Label: l.TotalSuppliers, Trend: "+8%", TrendUp: true, Color: "terracotta", TestID: "supplier-stat-total"},
				{Icon: "icon-check-circle", Value: "132", Label: l.Active, Trend: "+5%", TrendUp: true, Color: "sage", TestID: "supplier-stat-active"},
				{Icon: "icon-slash", Value: "18", Label: l.Blocked, Trend: "-2%", TrendUp: false, Color: "navy", TestID: "supplier-stat-blocked"},
				{Icon: "icon-pause-circle", Value: "6", Label: l.OnHold, Trend: "+1", TrendUp: false, Color: "amber", TestID: "supplier-stat-on-hold"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "activity", Title: l.SupplierActivity, Type: "chart", ChartKind: "line",
					ChartData: trend, Span: 2,
					HeaderActions: []types.QuickAction{
						{Label: l.FilterWeek, Href: "#"},
						{Label: l.FilterMonth, Href: "#"},
						{Label: l.FilterYear, Href: "#", Variant: "primary"},
					},
				},
				{
					ID: "recent", Title: l.RecentActivity, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ListURL},
					},
					ListItems: []types.ActivityItem{
						{IconName: "icon-truck", IconVariant: "client", Title: l.SupplierAdded, Description: "Global Materials Inc. onboarded", Time: "2m ago", TestID: "supplier-activity-added"},
						{IconName: "icon-check-circle", IconVariant: "quote", Title: l.SupplierActivated, Description: "TechParts Ltd. set to active", Time: "1h ago", TestID: "supplier-activity-activated"},
						{IconName: "icon-edit", IconVariant: "award", Title: l.DetailsUpdated, Description: "Acme Supplies payment terms changed", Time: "3h ago", TestID: "supplier-activity-updated"},
						{IconName: "icon-tag", IconVariant: "integration", Title: l.TagAssigned, Description: "Preferred tag added to 3 suppliers", Time: "5h ago", TestID: "supplier-activity-tagged"},
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.DashboardLabels.SupplierTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "supplier",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.DashboardLabels.SupplierTitle,
				HeaderIcon:   "icon-truck",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "supplier-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("supplier-dashboard", pageData)
	})
}
