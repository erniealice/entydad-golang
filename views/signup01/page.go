package signup01

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies for the signup01 page.
type Deps struct {
	Labels       entydad.SignupLabels
	CommonLabels pyeza.CommonLabels
	LoginURL     string // link back to login page (default: /auth/login)
}

// PageData holds the data for the signup01 page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.SignupLabels
	LoginURL        string
}

// NewView creates the signup01 page view (GET /auth/signup).
func NewView(deps *Deps) view.View {
	loginURL := deps.LoginURL
	if loginURL == "" {
		loginURL = entydad.AuthLoginURL
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "signup01-content",
			Labels:          deps.Labels,
			LoginURL:        loginURL,
		}

		return view.OK("signup01", pageData)
	})
}
