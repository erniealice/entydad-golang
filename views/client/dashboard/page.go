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
	CommonLabels    pyeza.CommonLabels
}

// PageData holds the data for the client dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          clientDashboardPageLabels
}

// clientDashboardPageLabels holds labels exposed to the client dashboard template.
type clientDashboardPageLabels struct {
	Dashboard entydad.ClientDashboardLabels
}

// NewView creates the client dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
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
			Labels: clientDashboardPageLabels{
				Dashboard: deps.Dashboard,
			},
		}

		return view.OK("client-dashboard", pageData)
	})
}
