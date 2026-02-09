package login01

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	entydad "github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies for the login page.
type Deps struct {
	Labels       entydad.LoginLabels
	CommonLabels pyeza.CommonLabels
	RedirectURL  string // where to send after login (default: /app/)
}

// PageData holds the data for the login page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.LoginLabels
	RedirectURL     string
}

// NewView creates the login page view (GET /login).
func NewView(deps *Deps) view.View {
	redirectURL := deps.RedirectURL
	if redirectURL == "" {
		redirectURL = "/app/"
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "login01-content",
			Labels:          deps.Labels,
			RedirectURL:     redirectURL,
		}

		return view.OK("login01", pageData)
	})
}
