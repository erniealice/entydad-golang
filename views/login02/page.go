package login02

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	entydad "github.com/erniealice/entydad-golang"
)

// CarouselSlide holds data for a single carousel slide.
type CarouselSlide struct {
	Title       string
	Description string
}

// SocialProvider holds data for a social login button.
type SocialProvider struct {
	Name    string // e.g. "Google", "Apple"
	IconSVG string // raw SVG markup
}

// Deps holds view dependencies for the login02 page.
type Deps struct {
	Labels          entydad.Login02Labels
	CommonLabels    pyeza.CommonLabels
	RedirectURL     string          // where to send after login (default: /app/)
	LogoText        string          // brand name displayed on form side
	LogoIcon        string          // icon template name for logo mark
	LoginPostURL    string          // form action URL (default: /login)
	RegisterURL     string          // sign-up link URL (default: /register)
	ForgotURL       string          // forgot password URL (default: /auth/reset-password)
	Slides          []CarouselSlide // carousel slides (left panel)
	SocialProviders []SocialProvider // social login buttons
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
	Slides          []CarouselSlide
	SocialProviders []SocialProvider
}

// NewView creates the login02 page view (GET /login).
func NewView(deps *Deps) view.View {
	redirectURL := deps.RedirectURL
	if redirectURL == "" {
		redirectURL = "/app/"
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

	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
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
			Slides:          deps.Slides,
			SocialProviders: deps.SocialProviders,
		}

		return view.OK("login02", pageData)
	})
}
