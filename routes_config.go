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

// ---------------------------------------------------------------------------
// TaxRegistrationRoutes
// ---------------------------------------------------------------------------

// TaxRegistrationRoutes holds route paths for the polymorphic TaxRegistration
// views. v1 surfaces client + workspace party types only.
// The AddURL and DeleteURL are party-scoped (include party_id in path).
type TaxRegistrationRoutes struct {
	// Client-scoped routes
	ClientListURL   string `json:"client_list_url"`
	ClientAddURL    string `json:"client_add_url"`
	ClientEditURL   string `json:"client_edit_url"`
	ClientDeleteURL string `json:"client_delete_url"`

	// Workspace-scoped routes
	WorkspaceListURL   string `json:"workspace_list_url"`
	WorkspaceAddURL    string `json:"workspace_add_url"`
	WorkspaceEditURL   string `json:"workspace_edit_url"`
	WorkspaceDeleteURL string `json:"workspace_delete_url"`

	// AddURL and DeleteURL are the active-context URLs (set by the view wiring).
	// For client views: same as ClientAddURL / ClientDeleteURL.
	// For workspace views: same as WorkspaceAddURL / WorkspaceDeleteURL.
	AddURL    string `json:"add_url"`
	DeleteURL string `json:"delete_url"`
}

// DefaultTaxRegistrationRoutes returns a TaxRegistrationRoutes populated from
// the package-level route constants.
func DefaultTaxRegistrationRoutes() TaxRegistrationRoutes {
	return TaxRegistrationRoutes{
		ClientListURL:   ClientTaxRegistrationListURL,
		ClientAddURL:    ClientTaxRegistrationAddURL,
		ClientEditURL:   ClientTaxRegistrationEditURL,
		ClientDeleteURL: ClientTaxRegistrationDeleteURL,

		WorkspaceListURL:   WorkspaceTaxRegistrationListURL,
		WorkspaceAddURL:    WorkspaceTaxRegistrationAddURL,
		WorkspaceEditURL:   WorkspaceTaxRegistrationEditURL,
		WorkspaceDeleteURL: WorkspaceTaxRegistrationDeleteURL,

		// Default to client context; override at wiring time for workspace.
		AddURL:    ClientTaxRegistrationAddURL,
		DeleteURL: ClientTaxRegistrationDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r TaxRegistrationRoutes) RouteMap() map[string]string {
	return map[string]string{
		"tax_registration.client.list":      r.ClientListURL,
		"tax_registration.client.add":       r.ClientAddURL,
		"tax_registration.client.edit":      r.ClientEditURL,
		"tax_registration.client.delete":    r.ClientDeleteURL,
		"tax_registration.workspace.list":   r.WorkspaceListURL,
		"tax_registration.workspace.add":    r.WorkspaceAddURL,
		"tax_registration.workspace.edit":   r.WorkspaceEditURL,
		"tax_registration.workspace.delete": r.WorkspaceDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// ConversationRoutes
// ---------------------------------------------------------------------------

// ConversationRoutes holds all URL constants for the conversation domain view
// module (secure messaging / ticketing, Plan-4 2026-06-03).
//
// Staff routes follow the standard /app/ (pages) + /action/ (HTMX) split.
// Portal routes are client-facing and only registered behind the AUTHZ_ENFORCE
// gate + the inherited 20260601 Phase-4 acting_as_client_id prerequisite.
type ConversationRoutes struct {
	// Staff surface
	ListURL      string `json:"list_url"`       // /app/conversations/list/{status}
	TableURL     string `json:"table_url"`      // /action/conversation/table/{status}
	DetailURL    string `json:"detail_url"`     // /app/conversations/detail/{id}
	AddURL       string `json:"add_url"`        // /action/conversation/add
	AssignURL    string `json:"assign_url"`     // /action/conversation/assign
	SetStatusURL string `json:"set_status_url"` // /action/conversation/set-status
	PostsURL     string `json:"posts_url"`      // /action/conversation/posts
	SendURL      string `json:"send_url"`       // /action/conversation_post/send
	MarkReadURL  string `json:"mark_read_url"`  // /action/conversation/mark-read

	// Client portal surface (gated)
	PortalListURL  string `json:"portal_list_url"`  // /portal/conversations
	PortalPostsURL string `json:"portal_posts_url"` // /action/conversation/portal-posts
}

// DefaultConversationRoutes returns ConversationRoutes populated from the
// package-level route constants.
func DefaultConversationRoutes() ConversationRoutes {
	return ConversationRoutes{
		ListURL:        ConversationListURL,
		TableURL:       ConversationTableURL,
		DetailURL:      ConversationDetailURL,
		AddURL:         ConversationAddURL,
		AssignURL:      ConversationAssignURL,
		SetStatusURL:   ConversationSetStatusURL,
		PostsURL:       ConversationPostsURL,
		SendURL:        ConversationSendURL,
		MarkReadURL:    ConversationMarkReadURL,
		PortalListURL:  ConversationPortalListURL,
		PortalPostsURL: ConversationPortalPostsURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ConversationRoutes) RouteMap() map[string]string {
	return map[string]string{
		"conversation.list":         r.ListURL,
		"conversation.table":        r.TableURL,
		"conversation.detail":       r.DetailURL,
		"conversation.add":          r.AddURL,
		"conversation.assign":       r.AssignURL,
		"conversation.set_status":   r.SetStatusURL,
		"conversation.posts":        r.PostsURL,
		"conversation.send":         r.SendURL,
		"conversation.mark_read":    r.MarkReadURL,
		"conversation.portal_list":  r.PortalListURL,
		"conversation.portal_posts": r.PortalPostsURL,
	}
}
