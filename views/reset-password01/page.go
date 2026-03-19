package resetpassword01

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies for the reset-password01 page.
type Deps struct {
	Labels         entydad.ResetPasswordLabels
	CommonLabels   pyeza.CommonLabels
	LoginURL       string // back to login link (default: /auth/login)
	ResetPostURL   string // form action for request step (default: /auth/reset-password)
	ConfirmPostURL string // form action for confirm step (default: /auth/reset-password/confirm)
}

// PageData holds the data for the reset-password01 page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.ResetPasswordLabels
	LoginURL        string
	ResetPostURL    string
	ConfirmPostURL  string
	Step            string // "request" or "confirm"
	Token           string // populated from query param when Step="confirm"
	Success         bool   // true after successful reset request or password change
	Error           string // validation error message
}

// NewView creates the reset-password01 page view (GET /auth/reset-password).
// It reads query params to determine which step to show:
//   - ?token=xxx → Step = "confirm", Token = xxx
//   - (no token) → Step = "request"
func NewView(deps *Deps) view.View {
	loginURL := deps.LoginURL
	if loginURL == "" {
		loginURL = entydad.AuthLoginURL
	}
	resetPostURL := deps.ResetPostURL
	if resetPostURL == "" {
		resetPostURL = entydad.AuthResetPasswordPostURL
	}
	confirmPostURL := deps.ConfirmPostURL
	if confirmPostURL == "" {
		confirmPostURL = entydad.AuthResetConfirmPostURL
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		step := "request"
		token := ""

		if viewCtx.Request != nil {
			token = viewCtx.Request.URL.Query().Get("token")
			if token != "" {
				step = "confirm"
			}
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "reset-password01-content",
			Labels:          deps.Labels,
			LoginURL:        loginURL,
			ResetPostURL:    resetPostURL,
			ConfirmPostURL:  confirmPostURL,
			Step:            step,
			Token:           token,
		}

		return view.OK("reset-password01", pageData)
	})
}
