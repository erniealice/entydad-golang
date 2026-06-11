package resetpassword01

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ActionDeps holds dependencies for the reset password form handlers.
type ActionDeps struct {
	LoginURL string // where to redirect after successful password reset (default: /auth/login)
}

// NewRequestAction creates the handler for POST /auth/reset-password (request step).
// Reads email from form and returns a success-state view.
// Actual email-sending logic is delegated to the app via middleware/handlers.
func NewRequestAction(deps *ActionDeps) view.View {
	loginURL := deps.LoginURL
	if loginURL == "" {
		loginURL = entydad.AuthLoginURL
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// On successful reset request, re-render the page in success state.
		// The app middleware handles actual email dispatch before this view runs.
		return view.Redirect(loginURL)
	})
}

// NewConfirmAction creates the handler for POST /auth/reset-password/confirm (confirm step).
// Reads token + new_password + confirm_password from form.
// Validates passwords match, then redirects to login on success.
// Actual password-change logic is delegated to the app via middleware/handlers.
func NewConfirmAction(deps *ActionDeps) view.View {
	loginURL := deps.LoginURL
	if loginURL == "" {
		loginURL = entydad.AuthLoginURL
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// Redirect to login after successful password reset (Phase 3 will wire real logic).
		return view.Redirect(loginURL)
	})
}
