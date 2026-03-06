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
	CommonLabels    pyeza.CommonLabels
}

// PageData holds the data for the client dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
}

// NewView creates the client dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.DashboardLabels.ClientTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "clients",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.DashboardLabels.ClientTitle,
				HeaderIcon:   "icon-users",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "client-dashboard-content",
		}

		return view.OK("client-dashboard", pageData)
	})
}
