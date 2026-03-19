package signup02

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ActionDeps holds dependencies for the signup02 form handler.
type ActionDeps struct {
	RedirectURL string // where to send after signup (default: /auth/login)
	LoginURL    string // link back to login in error responses
}

// NewAction creates the signup form handler (POST /auth/signup).
// The actual account creation logic is delegated to the app via middleware/handlers.
// This view handles the form rendering response for validation errors.
func NewAction(deps *ActionDeps) view.View {
	redirectURL := deps.RedirectURL
	if redirectURL == "" {
		redirectURL = entydad.AuthLoginURL
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// Redirect to login on successful signup (auto-login will be wired in Phase 3)
		return view.Redirect(redirectURL)
	})
}
