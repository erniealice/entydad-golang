package entydad

// Default route constants for entydad ROOT LEFTOVERS.
//
// The migrated entity URL constants (Client/User/Location/Role/Permission/
// Workspace/Supplier/Tag/PaymentTerm) now live in their domain/entity/**
// <entity>/routes.go files; the entity facade re-exports the prefixed names.
// What remains here is the admin dashboard, cross-entity report links, auth
// service surface, and the still-root-resident tax/conversation routes.
const (
	// Admin app dashboard — composite app spanning permission/role/workspace/
	// workspace_user/workspace_user_role. Lives under service/dashboard/views/admin.
	AdminDashboardURL = "/admin/dashboard"

	// Client report routes (cross-entity reports surface)
	ReceivablesAgingURL = "/reports/receivables-aging"

	// Supplier report routes (cross-entity reports surface)
	PayablesAgingURL = "/suppliers/reports/payables-aging"

	// Auth routes (all prefixed with /auth/)
	AuthLoginURL             = "/auth/login"
	AuthLoginPostURL         = "/auth/login"
	AuthSignupURL            = "/auth/signup"
	AuthSignupPostURL        = "/auth/signup"
	AuthResetPasswordURL     = "/auth/reset-password"
	AuthResetPasswordPostURL = "/auth/reset-password"
	AuthResetConfirmURL      = "/auth/reset-password/confirm"
	AuthResetConfirmPostURL  = "/auth/reset-password/confirm"
	AuthChangePasswordURL    = "/auth/change-password"
	AuthLogoutURL            = "/auth/logout"

	// Legacy login routes (redirect to /auth/login)
	LoginURL     = "/login"
	LoginPostURL = "/login"

	// DefaultAppRedirectURL is the default post-login redirect path.
	// Consumer apps should set Deps.RedirectURL to override this.
	//
	// Post-P12 (2026-05-22) of docs/plan/20260521-workspace-keyed-routing:
	// /app/* is gone. Post-login redirects should use
	// composition.homeURLForWorkspaceID() to land on /w/{slug}/home; this
	// constant is the workspace-less fallback (/me/inbox) for callers that
	// can't resolve a workspace slug at redirect time (password reset, etc.).
	DefaultAppRedirectURL = "/me/inbox"

	// TaxRegistration — polymorphic (client + workspace party types in v1)
	// URL convention: party_type + party_id come from the parent detail page context.
	ClientTaxRegistrationListURL   = "/clients/detail/{id}/tax-registrations"
	ClientTaxRegistrationAddURL    = "/action/client/{id}/tax-registration/add"
	ClientTaxRegistrationEditURL   = "/action/client/{id}/tax-registration/edit/{reg_id}"
	ClientTaxRegistrationDeleteURL = "/action/client/{id}/tax-registration/delete"

	WorkspaceTaxRegistrationListURL   = "/workspace/settings/tax-registrations"
	WorkspaceTaxRegistrationAddURL    = "/action/workspace/tax-registration/add"
	WorkspaceTaxRegistrationEditURL   = "/action/workspace/tax-registration/edit/{reg_id}"
	WorkspaceTaxRegistrationDeleteURL = "/action/workspace/tax-registration/delete"

	// Conversation — secure messaging / ticketing (Plan-4, 2026-06-03).
	// Staff surface: /app/conversations/* (pages) + /action/conversation* (HTMX).
	// Client portal surface: /portal/conversations — gated behind AUTHZ_ENFORCE
	// + the inherited 20260601 Phase-4 acting_as_client_id prerequisite (see block.go).
	ConversationListURL        = "/app/conversations/list/{status}"
	ConversationTableURL       = "/action/conversation/table/{status}"
	ConversationDetailURL      = "/app/conversations/detail/{id}"
	ConversationAddURL         = "/action/conversation/add"
	ConversationAssignURL      = "/action/conversation/assign"
	ConversationSetStatusURL   = "/action/conversation/set-status"
	ConversationPostsURL       = "/action/conversation/posts"
	ConversationSendURL        = "/action/conversation_post/send"
	ConversationMarkReadURL    = "/action/conversation/mark-read"
	ConversationPortalListURL  = "/portal/conversations"
	ConversationPortalPostsURL = "/action/conversation/portal-posts"
)
