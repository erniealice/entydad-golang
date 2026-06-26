package login02

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

// SocialProvider holds data for a social login button.
type SocialProvider struct {
	Name    string // e.g. "Google", "Microsoft"
	IconSVG string // raw SVG markup
	Method  string // Firebase provider id, e.g. "google.com", "microsoft.com"
}

// FirebaseConfig is the PUBLIC browser config for the Firebase JS SDK. When
// non-nil the login page renders in "firebase mode": the email/password form
// and every social button sign in via the Firebase SDK and POST the resulting
// ID token to FirebasePostURL (/auth/firebase). The legacy password POST
// (/auth/login) is NOT used in firebase mode. Nil = classic password mode
// (the page is byte-identical to before).
type FirebaseConfig struct {
	APIKey          string
	AuthDomain      string
	ProjectID       string
	EmulatorHost    string // optional FIREBASE_AUTH_EMULATOR_HOST
	MicrosoftTenant string // optional: pin the Azure-AD tenant for microsoft.com (single-tenant apps reject /common — AADSTS50194)
	FirebasePostURL string // where the client POSTs the verified ID token
}

// Deps holds view dependencies for the login02 page.
type Deps struct {
	Labels          entydad.Login02Labels
	CommonLabels    pyeza.CommonLabels
	RedirectURL     string           // where to send after login (default: /app/)
	LogoText        string           // brand name displayed on form side
	LogoIcon        string           // icon template name for logo mark
	LoginPostURL    string           // form action URL (default: /login)
	RegisterURL     string           // sign-up link URL (default: /register)
	ForgotURL       string           // forgot password URL (default: /auth/reset-password)
	Slides          []CarouselSlide  // carousel slides (left panel)
	SocialProviders []SocialProvider // social login buttons
	// FirebaseConfig non-nil ⇒ firebase mode (see type doc). ShowPasswordForm
	// controls whether the email/password form renders; it is forced true in
	// classic password mode (FirebaseConfig == nil) so legacy is unchanged.
	FirebaseConfig   *FirebaseConfig
	ShowPasswordForm bool
	// AllowSignups renders the "no account? sign up" footer link when true.
	AllowSignups bool
}

// PageData holds the data for the login02 page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.Login02Labels
	RedirectURL     string
	LogoText        string
	LogoIcon        string
	LoginPostURL    string
	RegisterURL     string
	ForgotURL       string
	Slides           []CarouselSlide
	SocialProviders  []SocialProvider
	FirebaseConfig   *FirebaseConfig
	ShowPasswordForm bool
	AllowSignups     bool
	Error            string // non-empty when login failed (e.g. ?error=invalid)
}

// NewView creates the login02 page view (GET /login).
func NewView(deps *Deps) view.View {
	redirectURL := deps.RedirectURL
	if redirectURL == "" {
		redirectURL = entydad.DefaultAppRedirectURL
	}
	loginPostURL := deps.LoginPostURL
	if loginPostURL == "" {
		loginPostURL = "/login"
	}
	registerURL := deps.RegisterURL
	if registerURL == "" {
		registerURL = "/register"
	}
	forgotURL := deps.ForgotURL
	if forgotURL == "" {
		forgotURL = "/auth/reset-password"
	}
	// Classic password mode (no firebase config) ALWAYS shows the password form
	// — preserves the legacy page exactly. In firebase mode the composition
	// layer decides via ShowPasswordForm (password ∈ allowed methods).
	showPasswordForm := deps.ShowPasswordForm
	if deps.FirebaseConfig == nil {
		showPasswordForm = true
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		errorMsg := ""
		if viewCtx.Request != nil {
			if viewCtx.Request.URL.Query().Get("error") != "" {
				errorMsg = deps.Labels.Error
			}
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "login02-content",
			Labels:          deps.Labels,
			RedirectURL:     redirectURL,
			LogoText:        deps.LogoText,
			LogoIcon:        deps.LogoIcon,
			LoginPostURL:    loginPostURL,
			RegisterURL:     registerURL,
			ForgotURL:       forgotURL,
			Slides:           deps.Slides,
			SocialProviders:  deps.SocialProviders,
			FirebaseConfig:   deps.FirebaseConfig,
			ShowPasswordForm: showPasswordForm,
			AllowSignups:     deps.AllowSignups,
			Error:            errorMsg,
		}

		return view.OK("login02", pageData)
	})
}
