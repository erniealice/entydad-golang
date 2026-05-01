package changepassword

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies for the change-password page.
type Deps struct {
	Labels       entydad.ChangePasswordLabels
	CommonLabels pyeza.CommonLabels
	// PostURL is the form action URL (default: /auth/change-password).
	PostURL string
	// BackURL is where the user navigates after a successful change (default: /).
	BackURL string
}

// PageData holds the template data for the change-password page.
type PageData struct {
	types.PageData
	Labels   entydad.ChangePasswordLabels
	PostURL  string
	BackURL  string
	Error    string // error message shown in the error banner
	Success  bool   // true after a successful password change
}

// NewView creates the change-password GET view (GET /auth/change-password).
// Query params:
//   - ?error=<message>  → shows the error banner
//   - ?success=1        → shows the success banner
func NewView(deps *Deps) view.View {
	postURL := deps.PostURL
	if postURL == "" {
		postURL = entydad.AuthChangePasswordURL
	}
	backURL := deps.BackURL
	if backURL == "" {
		backURL = entydad.DefaultAppRedirectURL
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		errorMsg := ""
		success := false

		if viewCtx.Request != nil {
			q := viewCtx.Request.URL.Query()
			errorMsg = q.Get("error")
			success = q.Get("success") == "1"
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			Labels:  deps.Labels,
			PostURL: postURL,
			BackURL: backURL,
			Error:   errorMsg,
			Success: success,
		}

		return view.OK("change-password", pageData)
	})
}
