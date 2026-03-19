package signup02

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

// SocialProvider holds data for a social signup button.
type SocialProvider struct {
	Name    string // e.g. "Google", "Apple"
	IconSVG string // raw SVG markup
}

// Deps holds view dependencies for the signup02 page.
type Deps struct {
	Labels          entydad.Signup02Labels
	CommonLabels    pyeza.CommonLabels
	LogoText        string           // brand name displayed on form side
	LogoIcon        string           // icon template name for logo mark
	SignupPostURL   string           // form action URL (default: /auth/signup)
	LoginURL        string           // sign-in link URL (default: /auth/login)
	TermsURL        string           // terms & conditions URL (default: /terms)
	Slides          []CarouselSlide  // carousel slides (left panel)
	SocialProviders []SocialProvider // social signup buttons
}

// PageData holds the data for the signup02 page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.Signup02Labels
	LogoText        string
	LogoIcon        string
	SignupPostURL   string
	LoginURL        string
	TermsURL        string
	Slides          []CarouselSlide
	SocialProviders []SocialProvider
}

// NewView creates the signup02 page view (GET /auth/signup).
func NewView(deps *Deps) view.View {
	signupPostURL := deps.SignupPostURL
	if signupPostURL == "" {
		signupPostURL = entydad.AuthSignupPostURL
	}
	loginURL := deps.LoginURL
	if loginURL == "" {
		loginURL = entydad.AuthLoginURL
	}
	termsURL := deps.TermsURL
	if termsURL == "" {
		termsURL = "/terms"
	}

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Title,
				CurrentPath:  viewCtx.CurrentPath,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "signup02-content",
			Labels:          deps.Labels,
			LogoText:        deps.LogoText,
			LogoIcon:        deps.LogoIcon,
			SignupPostURL:   signupPostURL,
			LoginURL:        loginURL,
			TermsURL:        termsURL,
			Slides:          deps.Slides,
			SocialProviders: deps.SocialProviders,
		}

		return view.OK("signup02", pageData)
	})
}
