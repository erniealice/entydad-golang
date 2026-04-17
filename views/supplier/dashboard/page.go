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
	CommonLabels    pyeza.CommonLabels
}

// PageData holds the data for the supplier dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          supplierDashboardPageLabels
}

// supplierDashboardPageLabels holds labels exposed to the supplier dashboard template.
type supplierDashboardPageLabels struct {
	Dashboard entydad.SupplierDashboardLabels
}

// NewView creates the supplier dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
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
			Labels: supplierDashboardPageLabels{
				Dashboard: deps.Dashboard,
			},
		}

		return view.OK("supplier-dashboard", pageData)
	})
}
