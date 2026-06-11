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
	Labels  entydad.ChangePasswordLabels
	PostURL string
	BackURL string
	Error   string // error message shown in the error banner
	Success bool   // true after a successful password change
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
			if code := q.Get("error"); code != "" {
				// Map the short error code from the action handler to a lyngua-
				// loaded label. Anything unrecognized falls through to the
				// generic Error label. Raw err.Error() strings are never
				// displayed — they may leak internals and aren't localisable.
				errorMsg = resolveErrorLabel(code, deps.Labels)
			}
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

// resolveErrorLabel maps a short error code from the action handler to the
// matching localised label on ChangePasswordLabels. Anything unrecognized
// returns the generic Error label so the user always sees a meaningful
// message (and never a raw Go error string).
func resolveErrorLabel(code string, l entydad.ChangePasswordLabels) string {
	switch code {
	case "mismatch":
		if l.ErrorMismatch != "" {
			return l.ErrorMismatch
		}
	case "incorrect":
		if l.ErrorCurrentIncorrect != "" {
			return l.ErrorCurrentIncorrect
		}
	case "too_short":
		if l.ErrorTooShort != "" {
			return l.ErrorTooShort
		}
	}
	return l.Error
}
