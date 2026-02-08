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

// PageData holds the data for the user dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
}

// NewView creates the user dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
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
		}

		return view.OK("user-dashboard", pageData)
	})
}
