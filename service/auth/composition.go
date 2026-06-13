package auth

import (
	"log"
	"net/http"
	"sync"

	entydad "github.com/erniealice/entydad-golang"
	changepasswordmod "github.com/erniealice/entydad-golang/service/auth/views/change-password"
	login02mod "github.com/erniealice/entydad-golang/service/auth/views/login02"
	resetpassword02mod "github.com/erniealice/entydad-golang/service/auth/views/reset-password02"
	signup02mod "github.com/erniealice/entydad-golang/service/auth/views/signup02"
	pyeza "github.com/erniealice/pyeza-golang"
)

// AuthLabels holds all label structs the auth views need.
type AuthLabels struct {
	Login02         entydad.Login02Labels
	Signup02        entydad.Signup02Labels
	ResetPassword02 entydad.ResetPassword02Labels
	ChangePassword  entydad.ChangePasswordLabels
	Common          pyeza.CommonLabels
	Messages        map[string]string
}

// Deps holds all dependencies the auth module needs from the host application.
// The host's composition layer constructs this once and passes it to NewAuthModule.
type Deps struct {
	// Credential operations
	AuthAdapter    AuthAdapter
	SessionManager SessionManager

	// Principal resolution and switching
	PrincipalResolver PrincipalResolver
	PrincipalSwitcher PrincipalSwitcher

	// CSRF
	CSRFSecret []byte
	CSRFIssuer CSRFIssuer

	// Rendering (auth-shell bypass path)
	Renderer Renderer

	// Data lookups (injected closures replace raw *sql.DB access)
	UserIDByEmail         UserIDByEmail
	WorkspaceSlugResolver WorkspaceSlugResolver

	// Labels
	Labels AuthLabels

	// App chrome
	LogoText       string
	LogoIcon       string
	CarouselSlides []CarouselSlide // nil = use DefaultCarouselSlides()

	// Test mode
	AuthProvider string // e.g. "password"
	TestMode     bool

	// Cookie policy
	SecureCookies func() bool

	// Session cookie name — defaults to espyna's consumer.DefaultSessionCookieName
	// ("ichizen_session") when empty.
	SessionCookieName string

	// Workspace CSRF cookie name — defaults to "ws_csrf" when empty.
	WorkspaceCSRFCookieName string

	// GetUserIDFromContext extracts the authenticated user ID from the request
	// context. Defaults to identity.FromContext(r.Context()).UserID when nil.
	GetUserIDFromContext func(r *http.Request) string
}

// AuthModule is the assembled auth service ready to register routes.
type AuthModule struct {
	deps *Deps
	// lastResetTokens stores raw HMAC reset tokens keyed by user_id for test
	// environments where PASSWORD_AUTH_TEST_MODE=true. This map is populated by
	// the POST /auth/reset-password handler and consumed by the test-only
	// GET /test/last-reset-token endpoint so E2E specs can obtain the raw token
	// without a real email delivery pipeline.
	lastResetTokens sync.Map
}

// NewAuthModule validates deps and returns a ready-to-register module.
func NewAuthModule(deps *Deps) *AuthModule {
	if deps.CarouselSlides == nil {
		deps.CarouselSlides = DefaultCarouselSlides()
	}
	if deps.SessionCookieName == "" {
		deps.SessionCookieName = "ichizen_session"
	}
	if deps.WorkspaceCSRFCookieName == "" {
		deps.WorkspaceCSRFCookieName = "ws_csrf"
	}
	return &AuthModule{deps: deps}
}

// RegisterRoutes registers all auth GET/POST handlers on the given registrar.
// Mirrors the pattern used by entydad's portal/profile/billing modules.
func (m *AuthModule) RegisterRoutes(routes RouteRegistrar) {
	deps := m.deps
	logoText := deps.LogoText
	carouselSlides := deps.CarouselSlides
	login02Slides := toLogin02Slides(carouselSlides)
	signup02Slides := toSignup02Slides(carouselSlides)
	resetpassword02Slides := toResetPassword02Slides(carouselSlides)

	// Login (GET + POST)
	routes.GET(entydad.AuthLoginURL, login02mod.NewView(&login02mod.Deps{
		Labels:       deps.Labels.Login02,
		CommonLabels: deps.Labels.Common,
		LogoText:     logoText,
		LogoIcon:     deps.LogoIcon,
		LoginPostURL: entydad.AuthLoginPostURL,
		RegisterURL:  entydad.AuthSignupURL,
		ForgotURL:    entydad.AuthResetPasswordURL,
		Slides:       login02Slides,
	}))

	// POST /auth/login
	routes.HandleFunc("POST", entydad.AuthLoginPostURL, m.handleLogin())

	// Signup (GET + POST)
	routes.GET(entydad.AuthSignupURL, signup02mod.NewView(&signup02mod.Deps{
		Labels:       deps.Labels.Signup02,
		CommonLabels: deps.Labels.Common,
		LogoText:     logoText,
		LogoIcon:     deps.LogoIcon,
		LoginURL:     entydad.AuthLoginURL,
		Slides:       signup02Slides,
	}))
	routes.HandleFunc("POST", entydad.AuthSignupPostURL, m.handleSignup())

	// Reset password (GET + POST request step + GET/POST confirm step)
	resetPasswordDeps := &resetpassword02mod.Deps{
		Labels:       deps.Labels.ResetPassword02,
		CommonLabels: deps.Labels.Common,
		LogoText:     logoText,
		LogoIcon:     deps.LogoIcon,
		LoginURL:     entydad.AuthLoginURL,
		Slides:       resetpassword02Slides,
	}
	routes.GET(entydad.AuthResetPasswordURL, resetpassword02mod.NewView(resetPasswordDeps))
	routes.HandleFunc("POST", entydad.AuthResetPasswordPostURL, m.handleResetPasswordRequest())
	routes.GET(entydad.AuthResetConfirmURL, resetpassword02mod.NewView(resetPasswordDeps))
	routes.HandleFunc("POST", entydad.AuthResetConfirmPostURL, m.handleResetPasswordConfirm())

	// Change password (GET + POST)
	routes.GET(entydad.AuthChangePasswordURL, changepasswordmod.NewView(&changepasswordmod.Deps{
		Labels:       deps.Labels.ChangePassword,
		CommonLabels: deps.Labels.Common,
		PostURL:      entydad.AuthChangePasswordURL,
		BackURL:      entydad.DefaultAppRedirectURL,
	}))
	routes.HandleFunc("POST", entydad.AuthChangePasswordURL, m.handleChangePassword())

	// Multi-principal: chooser page + switch action + no-access fallback
	routes.HandleFunc("GET", "/auth/select-workspace-role", m.handleSelectWorkspaceRole())
	routes.HandleFunc("POST", "/action/auth/switch-principal", m.handleSwitchPrincipal())

	// GET /auth/no-access
	routes.HandleFunc("GET", "/auth/no-access", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(noAccessHTML))
	})

	// Logout
	logoutHandler := m.handleLogout()
	routes.HandleFunc("POST", "/action/auth/logout", logoutHandler)
	routes.HandleFunc("POST", entydad.AuthLogoutURL, logoutHandler)
	routes.HandleFunc("GET", entydad.AuthLogoutURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(logoutLoadingHTML))
	})

	// Test-only endpoint: GET /test/last-reset-token?user_id=...
	if deps.AuthProvider == "password" && deps.TestMode {
		routes.HandleFunc("GET", "/test/last-reset-token", m.handleTestLastResetToken())
		log.Println("  ✓ Test-only endpoint mounted: GET /test/last-reset-token (live — reads in-process sync.Map)")
	}

	log.Println("  ✓ Auth screens initialized (login, signup, reset-password, change-password, logout)")
}
