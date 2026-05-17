package resetpassword02

import (
	"context"

	entydad "github.com/erniealice/entydad-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// CarouselSlide holds data for a single carousel slide.
type CarouselSlide struct {
	Title       string
	Description string
}

// Deps holds view dependencies for the reset-password02 split-screen page.
type Deps struct {
	Labels         entydad.ResetPassword02Labels
	CommonLabels   pyeza.CommonLabels
	LogoText       string          // brand name displayed on form side
	LogoIcon       string          // icon template name for logo mark
	LoginURL       string          // back to login link (default: /auth/login)
	ResetPostURL   string          // form action for request step (default: /auth/reset-password)
	ConfirmPostURL string          // form action for confirm step (default: /auth/reset-password/confirm)
	Slides         []CarouselSlide // carousel slides (left panel)
}

// PageData holds the data for the reset-password02 page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.ResetPassword02Labels
	LogoText        string
	LogoIcon        string
	LoginURL        string
	ResetPostURL    string
	ConfirmPostURL  string
	Slides          []CarouselSlide
	Step            string // "request" or "confirm"
	Token           string // populated from query param when Step="confirm"
	Success         bool   // true after successful reset request or password change
	Error           string // validation error message
}

// NewView creates the reset-password02 page view (GET /auth/reset-password).
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
		errorMsg := ""
		success := false

		if viewCtx.Request != nil {
			q := viewCtx.Request.URL.Query()
			token = q.Get("token")
			if token != "" {
				step = "confirm"
			}
			if q.Get("sent") == "true" {
				success = true
			}
			if code := q.Get("error"); code != "" {
				// Map the short error code from the action handler to a lyngua-
				// loaded label. Anything unrecognized falls through to the
				// generic Error label. Raw err.Error() strings are never
				// displayed — they may leak internals and aren't localisable.
				errorMsg = resolveErrorLabel(code, deps.Labels)
			}
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "reset-password02-content",
			Labels:          deps.Labels,
			LogoText:        deps.LogoText,
			LogoIcon:        deps.LogoIcon,
			LoginURL:        loginURL,
			ResetPostURL:    resetPostURL,
			ConfirmPostURL:  confirmPostURL,
			Slides:          deps.Slides,
			Step:            step,
			Token:           token,
			Success:         success,
			Error:           errorMsg,
		}

		return view.OK("reset-password02", pageData)
	})
}

// resolveErrorLabel maps a short error code from the action handler to the
// matching localised label on ResetPassword02Labels. Anything unrecognized
// returns the generic Error label so the user always sees a meaningful
// message (and never a raw Go error string).
func resolveErrorLabel(code string, l entydad.ResetPassword02Labels) string {
	switch code {
	case "mismatch":
		if l.ErrorMismatch != "" {
			return l.ErrorMismatch
		}
	case "invalid_token":
		if l.ErrorInvalidToken != "" {
			return l.ErrorInvalidToken
		}
	case "expired_token":
		if l.ErrorExpiredToken != "" {
			return l.ErrorExpiredToken
		}
	case "weak_password":
		if l.ErrorWeakPassword != "" {
			return l.ErrorWeakPassword
		}
	}
	return l.Error
}
