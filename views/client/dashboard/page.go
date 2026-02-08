package dashboard

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels pyeza.CommonLabels
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
				Title:        "Clients Dashboard",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "clients",
				ActiveSubNav: "dashboard",
				HeaderTitle:  "Clients Dashboard",
				HeaderIcon:   "icon-users",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "client-dashboard-content",
		}

		return view.OK("client-dashboard", pageData)
	})
}
