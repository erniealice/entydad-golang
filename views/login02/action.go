package login02

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"
)

// ActionDeps holds dependencies for the login02 form handler.
type ActionDeps struct {
	RedirectURL string // where to send after login (default: /app/)
}

// NewAction creates the login form handler (POST /login).
// The actual authentication logic is delegated to the app via middleware/handlers.
// This view handles the form rendering response for validation errors.
func NewAction(deps *ActionDeps) view.View {
	redirectURL := deps.RedirectURL
	if redirectURL == "" {
		redirectURL = "/app/"
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		return view.Redirect(redirectURL)
	})
}
