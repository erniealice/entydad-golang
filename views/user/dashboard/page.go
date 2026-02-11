package dashboard

import (
	"context"
	"fmt"
	"html/template"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// DashboardStats holds count values for stat cards.
type DashboardStats struct {
	TotalUsers    int
	ActiveUsers   int
	InactiveUsers int
	TotalRoles    int
}

// ActivityItem represents a single entry in the recent activity feed.
type ActivityItem struct {
	IconHTML    template.HTML
	Title      string
	Description string
	TimeAgo    string
}

// ChartData holds data for the activity chart.
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
	CommonLabels     pyeza.CommonLabels
	GetDashboardData func(ctx context.Context) (*DashboardData, error)
}

// PageData holds the data for the user dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Stats           DashboardStats
	RecentActivity  []ActivityItem
	ChartData       ChartData
	ActivePeriod    string
	ChartSVGPath    string
	ChartFillPath   string
}

// NewView creates the user dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		period := viewCtx.Request.URL.Query().Get("period")
		if period == "" {
			period = "year"
		}

		var stats DashboardStats
		var recentActivity []ActivityItem
		var chartData ChartData
		var chartSVGPath, chartFillPath string

		if deps.GetDashboardData != nil {
			data, err := deps.GetDashboardData(ctx)
			if err != nil {
				log.Printf("Failed to get dashboard data: %v", err)
			} else if data != nil {
				stats = data.Stats
				recentActivity = data.RecentActivity
				chartData = data.Chart
				chartSVGPath, chartFillPath = buildChartPaths(chartData.Values)
			}
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        "Users Dashboard",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "users",
				ActiveSubNav: "dashboard",
				HeaderTitle:  "Users Dashboard",
				HeaderIcon:   "icon-shield",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "user-dashboard-content",
			Stats:           stats,
			RecentActivity:  recentActivity,
			ChartData:       chartData,
			ActivePeriod:    period,
			ChartSVGPath:    chartSVGPath,
			ChartFillPath:   chartFillPath,
		}

		return view.OK("user-dashboard", pageData)
	})
}

// buildChartPaths generates SVG path strings from chart values.
// Returns (linePath, fillPath). If values is empty, returns empty strings.
func buildChartPaths(values []int) (string, string) {
	if len(values) == 0 {
		return "", ""
	}

	maxVal := 0
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	width := 400.0
	height := 200.0
	padding := 10.0
	usableH := height - 2*padding

	stepX := width / float64(len(values)-1)

	points := make([][2]float64, len(values))
	for i, v := range values {
		x := float64(i) * stepX
		y := padding + usableH*(1-float64(v)/float64(maxVal))
		points[i] = [2]float64{x, y}
	}

	// Build line path
	linePath := fmt.Sprintf("M%.0f,%.0f", points[0][0], points[0][1])
	for i := 1; i < len(points); i++ {
		linePath += fmt.Sprintf(" L%.0f,%.0f", points[i][0], points[i][1])
	}

	// Build fill path (line + close to bottom)
	fillPath := linePath + fmt.Sprintf(" L%.0f,%.0f L%.0f,%.0f Z", width, height, 0.0, height)

	return linePath, fillPath
}
