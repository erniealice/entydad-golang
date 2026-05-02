package dashboard

import (
	"context"
	"html/template"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
)

// DashboardStats holds count values for stat cards.
//
// Kept compatible with the existing service-admin GetDashboardData closure;
// values are projected into typed StatCardData inside the view.
type DashboardStats struct {
	TotalUsers    int
	ActiveUsers   int
	InactiveUsers int
	TotalRoles    int
}

// ActivityItem represents a single entry in the recent activity feed.
//
// IconHTML is kept for backward compatibility with the existing closure;
// the view falls back to it when IconName is unset. New callers should
// populate IconName + IconVariant so the pyeza activity-list renders the
// shared icon chip.
type ActivityItem struct {
	IconHTML    template.HTML
	IconName    string
	IconVariant string
	Title       string
	Description string
	TimeAgo     string
}

// ChartData holds data for the activity chart.
//
// Values stays []int for compatibility with the existing closure that
// counts user creations per month.
type ChartData struct {
	Labels []string
	Values []int
	Period string
}

// DashboardData is the combined result from the data provider.
type DashboardData struct {
	Stats          DashboardStats
	RecentActivity []ActivityItem
	Chart          ChartData
}

// Deps holds view dependencies.
type Deps struct {
	DashboardLabels  entydad.DashboardLabels
	Dashboard        entydad.UserDashboardLabels
	Routes           entydad.UserRoutes
	CommonLabels     pyeza.CommonLabels
	GetDashboardData func(ctx context.Context) (*DashboardData, error)
}

// PageData holds the data for the user dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the user dashboard view.
//
// Phase 1b refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// If GetDashboardData is set, its values feed the stats + chart + activity
// list; otherwise the view renders dummy values consistent with the other
// entydad dashboards (matching the DUMB-phase classification per the
// dashboard plan). Phase 2+ may replace the closure with a real use case.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Dashboard

		// Default dummy values mirror the supplier/client shape so all three
		// entydad dashboards look consistent in the Phase 1b state.
		statTotal := "184"
		statActive := "162"
		statInactive := "22"
		statRoles := "8"

		chartLabels := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		chartValues := []float64{120, 128, 134, 142, 150, 156, 162, 168, 172, 176, 180, 184}

		var activityItems []types.ActivityItem
		if deps.GetDashboardData != nil {
			if data, err := deps.GetDashboardData(ctx); err == nil && data != nil {
				if data.Stats.TotalUsers > 0 || data.Stats.ActiveUsers > 0 {
					statTotal = strconv.Itoa(data.Stats.TotalUsers)
					statActive = strconv.Itoa(data.Stats.ActiveUsers)
					statInactive = strconv.Itoa(data.Stats.InactiveUsers)
					statRoles = strconv.Itoa(data.Stats.TotalRoles)
				}
				if len(data.Chart.Labels) > 0 {
					chartLabels = data.Chart.Labels
				}
				if len(data.Chart.Values) > 0 {
					chartValues = chartValues[:0]
					for _, v := range data.Chart.Values {
						chartValues = append(chartValues, float64(v))
					}
				}
				for i, item := range data.RecentActivity {
					iconName := item.IconName
					iconVariant := item.IconVariant
					if iconName == "" {
						// Closures that pre-render IconHTML can't be projected
						// onto the activity-list shape — fall back to a generic
						// icon so the row still renders.
						iconName = "icon-info"
					}
					if iconVariant == "" {
						iconVariant = "client"
					}
					activityItems = append(activityItems, types.ActivityItem{
						IconName:    iconName,
						IconVariant: iconVariant,
						Title:       item.Title,
						Description: item.Description,
						Time:        item.TimeAgo,
						TestID:      "user-activity-" + strconv.Itoa(i+1),
					})
				}
			}
		}

		if len(activityItems) == 0 {
			activityItems = []types.ActivityItem{
				{IconName: "icon-shield", IconVariant: "client", Title: l.UserAdded, Description: "alice@example.com joined", Time: "2m ago", TestID: "user-activity-added"},
				{IconName: "icon-check-circle", IconVariant: "quote", Title: l.UserActivated, Description: "bob@example.com set to active", Time: "1h ago", TestID: "user-activity-activated"},
				{IconName: "icon-edit", IconVariant: "award", Title: l.RoleAssigned, Description: "Manager role granted to 2 users", Time: "3h ago", TestID: "user-activity-role"},
				{IconName: "icon-user", IconVariant: "integration", Title: l.ProfileUpdated, Description: "carol@example.com updated profile", Time: "5h ago", TestID: "user-activity-updated"},
			}
		}

		trend := &types.ChartData{
			Labels: chartLabels,
			Series: []types.ChartSeries{{
				Name:   l.UserActivity,
				Values: chartValues,
				Color:  "terracotta",
			}},
		}
		trend.AutoScale()

		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-user-plus", Label: l.QuickNew, Href: deps.Routes.AddURL, Variant: "primary", TestID: "user-action-new"},
				{Icon: "icon-list", Label: l.QuickViewAll, Href: deps.Routes.ListURL, TestID: "user-action-list"},
				{Icon: "icon-shield", Label: l.QuickRoles, Href: deps.Routes.RoleListURL, TestID: "user-action-roles"},
				{Icon: "icon-key", Label: l.QuickPermissions, Href: deps.Routes.PermissionListURL, TestID: "user-action-permissions"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-shield", Value: statTotal, Label: l.TotalUsers, Color: "terracotta", TestID: "user-stat-total"},
				{Icon: "icon-user-check", Value: statActive, Label: l.Active, Color: "sage", TestID: "user-stat-active"},
				{Icon: "icon-user-minus", Value: statInactive, Label: l.Inactive, Color: "navy", TestID: "user-stat-inactive"},
				{Icon: "icon-shield", Value: statRoles, Label: l.Roles, Color: "amber", TestID: "user-stat-roles"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "activity", Title: l.UserActivity, Type: "chart", ChartKind: "line",
					ChartData: trend, Span: 2,
					HeaderActions: []types.QuickAction{
						{Label: l.FilterWeek, Href: "?period=week"},
						{Label: l.FilterMonth, Href: "?period=month"},
						{Label: l.FilterYear, Href: "?period=year", Variant: "primary"},
					},
				},
				{
					ID: "recent", Title: l.RecentActivity, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ListURL},
					},
					ListItems: activityItems,
					EmptyState: &types.EmptyStateData{
						Title: l.NoRecentActivity,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.DashboardLabels.UserTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "user",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.DashboardLabels.UserTitle,
				HeaderIcon:   "icon-shield",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "user-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("user-dashboard", pageData)
	})
}
