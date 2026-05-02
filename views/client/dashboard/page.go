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
	Dashboard       entydad.ClientDashboardLabels
	Routes          entydad.ClientRoutes
	CommonLabels    pyeza.CommonLabels
}

// PageData holds the data for the client dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the client dashboard view.
//
// Phase 1b refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// Stat values, chart series, and activity items remain dummy until Phase 4
// wires real Client repository aggregate methods.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Dashboard

		trend := &types.ChartData{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
			Series: []types.ChartSeries{{
				Name:   l.ClientGrowth,
				Values: []float64{180, 195, 210, 218, 225, 232, 240, 248, 255, 263, 268, 275},
				Color:  "terracotta",
			}},
		}
		trend.AutoScale()

		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-user-plus", Label: l.QuickNew, Href: deps.Routes.AddURL, Variant: "primary", TestID: "client-action-new"},
				{Icon: "icon-list", Label: l.QuickViewAll, Href: deps.Routes.ListURL, TestID: "client-action-list"},
				{Icon: "icon-tag", Label: l.QuickTags, Href: deps.Routes.ClientTagListURL, TestID: "client-action-tags"},
				{Icon: "icon-clock", Label: l.QuickPaymentTerms, Href: deps.Routes.PaymentTermsURL, TestID: "client-action-payment-terms"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-users", Value: "248", Label: l.TotalClients, Trend: "+12%", TrendUp: true, Color: "terracotta", TestID: "client-stat-total"},
				{Icon: "icon-user-check", Value: "210", Label: l.Active, Trend: "+8%", TrendUp: true, Color: "sage", TestID: "client-stat-active"},
				{Icon: "icon-user-minus", Value: "38", Label: l.Inactive, Trend: "-3%", TrendUp: false, Color: "navy", TestID: "client-stat-inactive"},
				{Icon: "icon-user-plus", Value: "24", Label: l.NewThisMonth, Trend: "+5", TrendUp: true, Color: "amber", TestID: "client-stat-new"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "growth", Title: l.ClientGrowth, Type: "chart", ChartKind: "line",
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
						{IconName: "icon-user-plus", IconVariant: "client", Title: l.ClientAdded, Description: "Acme Corporation joined", Time: "2m ago", TestID: "client-activity-added"},
						{IconName: "icon-check-circle", IconVariant: "quote", Title: l.ClientActivated, Description: "TechCorp Inc. set to active", Time: "1h ago", TestID: "client-activity-activated"},
						{IconName: "icon-edit", IconVariant: "award", Title: l.ProfileUpdated, Description: "Global Solutions details changed", Time: "3h ago", TestID: "client-activity-updated"},
						{IconName: "icon-tag", IconVariant: "integration", Title: l.TagAssigned, Description: "VIP tag added to 5 clients", Time: "5h ago", TestID: "client-activity-tagged"},
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.DashboardLabels.ClientTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "client",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.DashboardLabels.ClientTitle,
				HeaderIcon:   "icon-users",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "client-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("client-dashboard", pageData)
	})
}
