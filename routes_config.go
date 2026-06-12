package entydad

// Route configuration structs for the entydad domain.
//
// This file implements a three-level routing system:
//
//   Level 1 — Generic defaults from Go consts (this file).
//     DefaultXxxRoutes() constructors populate structs from the route constants
//     defined in routes.go. These are sensible defaults that work out of the box.
//
//   Level 2 — Industry-specific overrides via JSON (loaded by consumer apps).
//     Consumer apps can load a JSON configuration file and unmarshal it into
//     route structs, overriding some or all of the default paths. JSON tags on
//     every field enable this.
//
//   Level 3 — App-specific overrides via Go field assignment (optional).
//     After constructing or loading routes, the consumer app can directly set
//     individual fields to customize specific routes for its own needs.
//
// Each route struct also provides a RouteMap() method that returns a
// map[string]string keyed by dot-notation identifiers (e.g. "client.list"),
// suitable for template rendering and reverse-routing lookups.

// ---------------------------------------------------------------------------
// ClientRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// UserRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// RoleRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// LocationRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// LocationAreaRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// PermissionRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// WorkspaceRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// WorkspaceUserRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// WorkspaceUserRoleRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// AdminDashboardRoutes
// ---------------------------------------------------------------------------

// AdminDashboardRoutes holds the route path for the admin app dashboard.
//
// Admin is a *composite* app: its sidebar block spans permission, role,
// workspace, workspace_user, and workspace_user_role entities. The dashboard
// sits at the app level (not the entity level), so it gets its own thin
// route struct rather than being attached to PermissionRoutes / RoleRoutes /
// WorkspaceRoutes — none of which would be the obvious owner.
type AdminDashboardRoutes struct {
	DashboardURL string `json:"dashboard_url"`
}

// DefaultAdminDashboardRoutes returns an AdminDashboardRoutes populated from
// the package-level route constants.
func DefaultAdminDashboardRoutes() AdminDashboardRoutes {
	return AdminDashboardRoutes{
		DashboardURL: AdminDashboardURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r AdminDashboardRoutes) RouteMap() map[string]string {
	return map[string]string{
		"admin.dashboard": r.DashboardURL,
	}
}

// ---------------------------------------------------------------------------
// SupplierRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// PaymentTermRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// SupplierPaymentTermRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// ClientTagRoutes
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// LoginRoutes
// ---------------------------------------------------------------------------

// LoginRoutes holds all route paths for login/authentication views.
type LoginRoutes struct {
	LoginURL     string `json:"login_url"`
	LoginPostURL string `json:"login_post_url"`
}

// DefaultLoginRoutes returns a LoginRoutes populated from the package-level
// route constants.
func DefaultLoginRoutes() LoginRoutes {
	return LoginRoutes{
		LoginURL:     LoginURL,
		LoginPostURL: LoginPostURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r LoginRoutes) RouteMap() map[string]string {
	return map[string]string{
		"login.page": r.LoginURL,
		"login.post": r.LoginPostURL,
	}
}

// ---------------------------------------------------------------------------
// AuthRoutes
// ---------------------------------------------------------------------------

// AuthRoutes holds all route paths for authentication views (signup, reset, logout).
type AuthRoutes struct {
	LoginURL             string `json:"login_url"`
	LoginPostURL         string `json:"login_post_url"`
	SignupURL            string `json:"signup_url"`
	SignupPostURL        string `json:"signup_post_url"`
	ResetPasswordURL     string `json:"reset_password_url"`
	ResetPasswordPostURL string `json:"reset_password_post_url"`
	ResetConfirmURL      string `json:"reset_confirm_url"`
	ResetConfirmPostURL  string `json:"reset_confirm_post_url"`
	LogoutURL            string `json:"logout_url"`
}

// DefaultAuthRoutes returns an AuthRoutes populated from the package-level
// URL constants defined in routes.go.
func DefaultAuthRoutes() AuthRoutes {
	return AuthRoutes{
		LoginURL:             AuthLoginURL,
		LoginPostURL:         AuthLoginPostURL,
		SignupURL:            AuthSignupURL,
		SignupPostURL:        AuthSignupPostURL,
		ResetPasswordURL:     AuthResetPasswordURL,
		ResetPasswordPostURL: AuthResetPasswordPostURL,
		ResetConfirmURL:      AuthResetConfirmURL,
		ResetConfirmPostURL:  AuthResetConfirmPostURL,
		LogoutURL:            AuthLogoutURL,
	}
}

// RouteMap returns a map of route keys to URL paths for AuthRoutes.
func (r AuthRoutes) RouteMap() map[string]string {
	return map[string]string{
		"auth.login.page":          r.LoginURL,
		"auth.login.post":          r.LoginPostURL,
		"auth.signup.page":         r.SignupURL,
		"auth.signup.post":         r.SignupPostURL,
		"auth.reset-password.page": r.ResetPasswordURL,
		"auth.reset-password.post": r.ResetPasswordPostURL,
		"auth.reset-confirm.page":  r.ResetConfirmURL,
		"auth.reset-confirm.post":  r.ResetConfirmPostURL,
		"auth.logout":              r.LogoutURL,
	}
}
