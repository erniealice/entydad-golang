package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"
)

// classifySignupError maps an adapter error from Register to a short code
// consumed by the signup02 page. Mirrors classifyChangePasswordError /
// classifyResetError: codes only, never raw err.Error() in the URL.
//
// Codes:
//   - "email_taken"     — the email is already registered
//   - "weak_password"   — the password is below the minimum length
//   - "invalid_email"   — the email is malformed
//   - "generic"         — any other failure
func classifySignupError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "already"),
		strings.Contains(msg, "exists"),
		strings.Contains(msg, "duplicate"):
		return "email_taken"
	case strings.Contains(msg, "at least"),
		strings.Contains(msg, "too short"),
		strings.Contains(msg, "minimum length"):
		return "weak_password"
	case strings.Contains(msg, "invalid email"),
		strings.Contains(msg, "email format"),
		strings.Contains(msg, "malformed"):
		return "invalid_email"
	default:
		return "generic"
	}
}

// classifyResetError maps an adapter error from ExecutePasswordReset to a
// short code consumed by the reset-password page handler. Same contract as
// classifyChangePasswordError: codes only, never raw err.Error().
//
// Codes:
//   - "expired_token"  — the reset token has expired
//   - "invalid_token"  — the reset token is malformed / signature mismatch / unknown
//   - "weak_password"  — the new password is below the minimum length
//   - "generic"        — any other failure
func classifyResetError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "expired"):
		return "expired_token"
	case strings.Contains(msg, "token"),
		strings.Contains(msg, "signature"),
		strings.Contains(msg, "hmac"):
		return "invalid_token"
	case strings.Contains(msg, "at least"),
		strings.Contains(msg, "too short"),
		strings.Contains(msg, "minimum length"):
		return "weak_password"
	default:
		return "generic"
	}
}

// classifyChangePasswordError maps an adapter error to a short code that the
// change-password page handler resolves to a lyngua-loaded label. Keeps raw
// err.Error() strings out of the rendered HTML.
//
// Codes:
//   - "incorrect" — the user typed the wrong current password
//   - "too_short" — the new password is below the minimum length
//   - "generic"   — any other failure
func classifyChangePasswordError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "current password is incorrect"),
		strings.Contains(msg, "invalid password"):
		return "incorrect"
	case strings.Contains(msg, "at least"),
		strings.Contains(msg, "too short"),
		strings.Contains(msg, "minimum length"):
		return "too_short"
	default:
		return "generic"
	}
}

// iconForPrincipalKind returns a pyeza icon template name for a given
// principal type. Used by the select-workspace-role page card list. Keep these
// in sync with packages/pyeza-golang/web/templates/components/icons/ —
// adding a new principal type means adding (or aliasing) an icon.
func iconForPrincipalKind(t PrincipalType) string {
	switch t {
	case PrincipalTypeOperatorOwner:
		return "icon-shield-check"
	case PrincipalTypeOperatorStaff:
		return "icon-users"
	case PrincipalTypeClient:
		return "icon-user"
	case PrincipalTypeClientDelegate:
		return "icon-user-group"
	case PrincipalTypeSupplier:
		return "icon-briefcase"
	case PrincipalTypeSupplierDelegate:
		return "icon-users"
	}
	return ""
}

// homeURLForWorkspaceID constructs the post-login home URL: /w/{slug}/home.
// Falls back to /me/inbox (workspace-less personal landing per Q-WS-7 → B)
// if the lookup fails or returns empty — that path always works regardless
// of slug state.
//
// Added 2026-05-22 to land users directly on their workspace-keyed home
// after login + principal-switch (was landing on /me/inbox unconditionally).
func (m *AuthModule) homeURLForWorkspaceID(ctx context.Context, workspaceID string) string {
	const fallback = "/me/inbox"
	if m.deps.WorkspaceSlugResolver == nil || workspaceID == "" {
		return fallback
	}
	slug := m.deps.WorkspaceSlugResolver(ctx, workspaceID)
	if slug == "" {
		return fallback
	}
	return "/w/" + slug + "/home"
}

// renderAuthView renders a view.ViewResult to the response for auth-shell
// pages (login, select-workspace-role). We bypass the main ViewAdapter pipeline
// because (a) the chooser sits behind /auth/ which is session-middleware
// excluded and so doesn't have a workspace/permissions context, and (b)
// the auth-shell views don't carry a Sidebar field for the adapter to
// fill.
//
// The renderer is reused from the deps so all template namespaces
// (login02-content, fonts, alert, etc.) are available.
func (m *AuthModule) renderAuthView(w http.ResponseWriter, r *http.Request, result view.ViewResult) {
	if result.Error != nil {
		log.Printf("[AUTH] auth view error: %v", result.Error)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if result.Redirect != "" {
		http.Redirect(w, r, result.Redirect, http.StatusSeeOther)
		return
	}
	if result.Template == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if m.deps.Renderer == nil {
		http.Error(w, "renderer not available", http.StatusInternalServerError)
		return
	}
	status := result.StatusCode
	if status == 0 {
		status = http.StatusOK
	}
	// This bypass path skips the ViewAdapter (which is what sets Content-Type on
	// the normal pipeline), so set it explicitly here BEFORE WriteHeader. Without
	// it, WriteHeader commits the response with no Content-Type and the
	// X-Content-Type-Options: nosniff security header makes the browser render the
	// HTML source as plain text. (The inline-HTML handlers below set it the same way.)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := m.deps.Renderer.Render(w, result.Template, result.Data); err != nil {
		log.Printf("[AUTH] render %q failed: %v", result.Template, err)
	}
}

// findPrincipalByID searches a slice of Principal values for one matching the
// given id. An optional kindHint (case-insensitive) disambiguates when
// different principal types might share an id value. Returns the matching
// Principal and true, or a zero value and false.
func findPrincipalByID(principals []Principal, id, kindHint string) (Principal, bool) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Principal{}, false
	}
	kindHint = strings.TrimSpace(kindHint)
	for _, p := range principals {
		if p.PrincipalID == id {
			// Case-insensitive: PrincipalType.String() returns lowercase ("client_delegate"),
			// picker templates POST uppercase ("CLIENT_DELEGATE"). Both must match.
			if kindHint != "" && !strings.EqualFold(PrincipalTypeString(p.Type), kindHint) {
				continue
			}
			return p, true
		}
	}
	return Principal{}, false
}
