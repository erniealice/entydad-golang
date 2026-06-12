package entydad

// Default route constants for entydad ROOT LEFTOVERS.
//
// The migrated entity URL constants (Client/User/Location/Role/Permission/
// Workspace/Supplier/Tag/PaymentTerm) now live in their domain/entity/**
// <entity>/routes.go files; the entity facade re-exports the prefixed names.
// What remains here is the admin dashboard, cross-entity report links, and the
// auth service surface. (Tax routes relocated to domain/tax/tax_registration —
// thread TX; conversation routes relocated to hybra views/conversation/model —
// thread TC; both 2026-06-12.)
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
)

// TaxRegistration route constants relocated to domain/tax/tax_registration
// (routes.go) — entity-local under the tax domain (fork E4 / thread TX,
// 2026-06-12). Conversation route constants relocated to hybra
// views/conversation/model (contract.go) — cross-cutting communication surface
// (view-package-placement.md OCID / thread TC, 2026-06-12).
