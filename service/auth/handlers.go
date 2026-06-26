package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	entydad "github.com/erniealice/entydad-golang"
	selectWorkspaceRole "github.com/erniealice/entydad-golang/service/auth/views/login02/select-workspace-role"
	"github.com/erniealice/espyna-golang/shared/identity"
	"github.com/erniealice/pyeza-golang/view"
)

// ctxKeyAuthUserIDType is a typed context key for the cookie-validated
// user_id installed by the chooser GET handler. Typed (not string) so it
// can't collide with arbitrary "user_id" keys elsewhere in the request
// stack.
type ctxKeyAuthUserIDType struct{}

var ctxKeyAuthUserID = ctxKeyAuthUserIDType{}

// handleLogin returns the POST /auth/login handler — multi-principal aware.
//
// Sequence (the auth-cycle UX flow):
//
//	 1. Validate credentials via authAdapter.Login (existing path).
//	 2. Resolve all active principal bindings for the user via
//	    principalLoader.Resolve. Three branches follow:
//	      a. 0 principals    → /auth/no-access (signed-in but no access)
//	      b. 1 principal     → mint principal-scoped session row,
//	                            redirect to that principal's home
//	      c. 2+ principals   → keep the (principal-less) session cookie
//	                            from step 1, redirect to chooser
//
//	The 1-principal branch ROTATES the session: the token from step 1
//	is invalidated and a fresh, principal-stamped session row is
//	inserted in its place. The fresh token becomes the cookie. This
//	is the same security invariant as the cross-principal switch
//	(see principal_switch.go) — the cookie that ever sees a
//	principal-less authenticated state has a different value from
//	the cookie that sees authenticated+principal-resolved state.
func (m *AuthModule) handleLogin() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter
	sessionMw := m.deps.SessionManager
	principalLoader := m.deps.PrincipalResolver

	return func(w http.ResponseWriter, r *http.Request) {
		if authAdapter == nil || sessionMw == nil {
			// mock_auth: no real adapter — just redirect to app
			http.Redirect(w, r, entydad.DefaultAppRedirectURL, http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		token, identity, err := authAdapter.Login(r.Context(), email, password)
		if err != nil {
			log.Printf("[AUTH] login failed for %s: %v", email, err)
			http.Redirect(w, r, entydad.AuthLoginURL+"?error=invalid", http.StatusSeeOther)
			return
		}
		sessionMw.SetSessionCookie(w, token)
		// C2: issue CSRF cookie carrying the session token claim.
		// workspace_id is empty at login (no workspace chosen yet) — the
		// claim will be refreshed on the next GET once workspace_path
		// resolves the workspace.
		m.deps.CSRFIssuer(w, m.deps.CSRFSecret, token, "")

		// If principal resolution isn't wired (no sqlDB / no loader),
		// keep the legacy behaviour. The fail-closed permission model
		// in view_adapter.go is the safety net for unauthenticated
		// reads — it just means the chooser/portal flow is skipped.
		if principalLoader == nil || !principalLoader.IsEnabled() {
			http.Redirect(w, r, entydad.DefaultAppRedirectURL, http.StatusSeeOther)
			return
		}

		userID := ""
		if identity != nil {
			userID = identity.GetId()
		}
		if userID == "" {
			// Fall back to a lookup by email — the identity proto's
			// Subject field is normally the user_id but mock providers
			// occasionally leave it blank.
			if m.deps.UserIDByEmail != nil {
				userID = m.deps.UserIDByEmail(r.Context(), email)
			}
		}
		if userID == "" {
			log.Printf("[AUTH] login succeeded but no user_id resolvable for %s", email)
			http.Redirect(w, r, "/auth/no-access", http.StatusSeeOther)
			return
		}

		// Shared post-auth tail: resolve principals + auto-route/chooser.
		// routePrincipals performs any session rotation + CSRF refresh as a
		// side effect and RETURNS the redirect target; the password form
		// answers with a 303. The Firebase endpoint reuses the exact same
		// routePrincipals (answering with JSON), so the delegate-guard /
		// fail-closed logic lives in ONE place.
		http.Redirect(w, r, m.routePrincipals(w, r, token, userID), http.StatusSeeOther)
	}
}

// routePrincipals is the SHARED post-authentication tail for every login path
// (password form + Firebase ID token). The caller has already minted a session
// `token` for `userID` and set the session cookie + CSRF. This resolves the
// user's active principal bindings and decides where to send them, performing
// any session rotation / CSRF refresh as a SIDE EFFECT (on w). It RETURNS the
// redirect target URL rather than writing the navigation itself, so the
// password path can answer with a 303 and the Firebase fetch path with a JSON
// redirect — both reuse the identical, security-critical principal-resolution
// + delegate-guard logic. Behaviour for the password path is byte-identical to
// the prior inline switch (same URLs, same cookie writes, same 303).
func (m *AuthModule) routePrincipals(w http.ResponseWriter, r *http.Request, token, userID string) string {
	authAdapter := m.deps.AuthAdapter
	sessionMw := m.deps.SessionManager
	principalLoader := m.deps.PrincipalResolver

	principals, presolveErr := principalLoader.Resolve(r.Context(), userID)
	if presolveErr != nil {
		log.Printf("[AUTH] principal resolve failed for user %s: %v", userID, presolveErr)
		return "/auth/no-access"
	}
	switch len(principals) {
	case 0:
		// Authenticated but has no active grants — sign them out
		// and surface the no-access page. The cookie is cleared so
		// they don't sit on a permissionless session.
		if invalidErr := authAdapter.InvalidateSession(r.Context(), token); invalidErr != nil {
			log.Printf("[AUTH] no-access path: failed to invalidate session: %v", invalidErr)
		}
		sessionMw.ClearSessionCookie(w)
		return "/auth/no-access"

	case 1:
		// codex RBC#1 — High-1 (2026-06-02). A multi-target delegate
		// (len(ActingAsTargets) > 1) reaching this auto-route would be
		// switched into WITHOUT an explicit acting-as id, persisting a
		// delegate principal with an EMPTY acting_as_* (→ fail-closed
		// to zero permissions). Even though there is exactly ONE
		// principal, the ACT-AS identity is ambiguous — fail closed by
		// routing to the chooser so the user picks a target, never
		// calling executePrincipalSwitch with an empty acting-as.
		// Single-target delegates and non-delegate principals pass.
		if !DelegateActingAsResolved(principals[0], "", "") {
			log.Printf("[AUTH] auto-route blocked: multi-target delegate for user %s requires target selection — routing to chooser", userID)
			return "/auth/select-workspace-role"
		}
		// Auto-route. Rotate the session so the cookie token
		// reflects the principal-stamped row. This is the
		// post-login "exactly one principal" branch; we always
		// rotate to mint a principal-scoped session row.
		result, err := m.deps.PrincipalSwitcher(r.Context(), PrincipalSwitchInput{
			UserID:          userID,
			Token:           token,
			TargetPrincipal: principals[0],
			// Explicit-form caller (this is the chooser GET that
			// auto-routes); RequireAudit=false preserves dev-mode
			// best-effort behavior. UseCase per red-team X-2.
			UseCase:      "switch_explicit_rotate",
			RequestURL:   r.URL.Path,
			Referer:      r.Header.Get("Referer"),
			SecFetchSite: r.Header.Get("Sec-Fetch-Site"),
			UserAgent:    r.Header.Get("User-Agent"),
			RequireAudit: false,
		})
		if err != nil {
			log.Printf("[AUTH] auto-route principal switch failed for user %s: %v", userID, err)
			return "/auth/no-access"
		}
		if result.NewToken != "" {
			sessionMw.SetSessionCookie(w, result.NewToken)
		}
		// C2: always refresh the CSRF cookie after a successful
		// principal switch — even when NewToken is empty (in-place
		// session update, same workspace). The initial login CSRF
		// cookie carries workspaceID="" which becomes stale once the
		// session is stamped with a workspace. Use the effective
		// session token (new if rotated, original if in-place).
		effectiveToken := result.NewToken
		if effectiveToken == "" {
			effectiveToken = token
		}
		m.deps.CSRFIssuer(w, m.deps.CSRFSecret,
			effectiveToken, principals[0].WorkspaceID)
		return m.homeURLForWorkspaceID(r.Context(), principals[0].WorkspaceID)

	default:
		// 2+ principals. If all are the same kind (e.g. one user
		// has OPERATOR_STAFF in multiple workspaces), auto-route to
		// the first one — the workspace switcher in the sidebar
		// handles same-kind workspace switching. Only show the
		// chooser when principals span different kinds (e.g.
		// OPERATOR_STAFF + CLIENT), where the UX role differs.
		allSameKind := true
		firstKind := principals[0].Type
		for _, p := range principals[1:] {
			if p.Type != firstKind {
				allSameKind = false
				break
			}
		}
		if allSameKind {
			// codex RBC#1 — High-1 (2026-06-02). Same-kind auto-route
			// silently lands on principals[0]. That is safe for staff
			// (the sidebar workspace switcher handles multi-WU), but a
			// delegate principal holding N>1 acting-as targets would be
			// switched into with an EMPTY acting_as_* (→ fail-closed to
			// zero permissions). The act-as identity is ambiguous —
			// fail closed by routing to the chooser rather than calling
			// executePrincipalSwitch with an empty acting-as. A
			// single-target delegate is unambiguous and proceeds.
			if !DelegateActingAsResolved(principals[0], "", "") {
				log.Printf("[AUTH] same-kind auto-route blocked: multi-target delegate for user %s requires target selection — routing to chooser", userID)
				return "/auth/select-workspace-role"
			}
			result, err := m.deps.PrincipalSwitcher(r.Context(), PrincipalSwitchInput{
				UserID:          userID,
				Token:           token,
				TargetPrincipal: principals[0],
				// Same-kind auto-route after login (multi-WU staff
				// across N workspaces lands in the first one).
				// Explicit-form caller; RequireAudit=false.
				UseCase:      "switch_explicit_rotate",
				RequestURL:   r.URL.Path,
				Referer:      r.Header.Get("Referer"),
				SecFetchSite: r.Header.Get("Sec-Fetch-Site"),
				UserAgent:    r.Header.Get("User-Agent"),
				RequireAudit: false,
			})
			if err != nil {
				log.Printf("[AUTH] auto-route same-kind principal switch failed for user %s: %v", userID, err)
				return "/auth/no-access"
			}
			if result.NewToken != "" {
				sessionMw.SetSessionCookie(w, result.NewToken)
			}
			// C2: always refresh CSRF cookie (see case-1 comment).
			effectiveToken := result.NewToken
			if effectiveToken == "" {
				effectiveToken = token
			}
			m.deps.CSRFIssuer(w, m.deps.CSRFSecret,
				effectiveToken, principals[0].WorkspaceID)
			return m.homeURLForWorkspaceID(r.Context(), principals[0].WorkspaceID)
		}
		// Different kinds: present the chooser. The principal-less
		// session stays — view_adapter's fail-closed default keeps
		// the chooser page itself accessible (it doesn't gate on
		// perms) but blocks everything else.
		return "/auth/select-workspace-role"
	}
}

// handleFirebaseLogin returns the POST /auth/firebase handler — the Firebase
// ID-token login path (Microsoft / Google / etc. via the Firebase JS SDK). The
// browser signs in with the Firebase SDK, obtains an ID token, and POSTs it
// here as `id_token`. The server VERIFIES the token (FirebaseVerifier),
// enforces the optional sign-in-method allow-list, resolves the DB user by
// EMAIL (migrated users have no firebase_uid, so token.email -> user.email is
// the join key), mints a server-side session (SessionMinter — provider-agnostic
// after the session-decoupling), and reuses routePrincipals for the identical
// principal-resolution / delegate-guard logic as the password path.
//
// Mounted at /auth/firebase — under the session-middleware exclude prefix
// "/auth/" (like /auth/login) so this pre-session POST is not bounced to login,
// and outside the /action/* CSRF surface. Only registered when FirebaseVerifier
// + SessionMinter are wired (i.e. CONFIG_AUTH_PROVIDER=firebase).
//
// Answered as JSON {"redirect": "..."} because the client is a fetch(), which
// cannot follow a 303 as a navigation — the browser JS does
// window.location = redirect.
func (m *AuthModule) handleFirebaseLogin() http.HandlerFunc {
	verifier := m.deps.FirebaseVerifier
	minter := m.deps.SessionMinter
	sessionMw := m.deps.SessionManager
	principalLoader := m.deps.PrincipalResolver
	allowed := m.deps.AllowedSignInMethods

	return func(w http.ResponseWriter, r *http.Request) {
		if verifier == nil || minter == nil || sessionMw == nil {
			writeFirebaseJSON(w, http.StatusNotFound, "firebase sign-in not enabled")
			return
		}
		if err := r.ParseForm(); err != nil {
			writeFirebaseJSON(w, http.StatusBadRequest, "invalid_request")
			return
		}
		idToken := strings.TrimSpace(r.FormValue("id_token"))
		if idToken == "" {
			writeFirebaseJSON(w, http.StatusBadRequest, "missing_token")
			return
		}
		email, signInProvider, err := verifier(r.Context(), idToken)
		if err != nil || strings.TrimSpace(email) == "" {
			log.Printf("[AUTH] firebase verify failed (provider=%s): %v", signInProvider, err)
			writeFirebaseError(w, http.StatusUnauthorized, "invalid")
			return
		}
		// Layer 5: enforce the configured sign-in-method allow-list. Empty list
		// = allow any verified method. The token's sign_in_provider claim is the
		// source of truth (e.g. "microsoft.com", "google.com", "password").
		if len(allowed) > 0 && !signInMethodAllowed(allowed, signInProvider) {
			log.Printf("[AUTH] firebase sign-in method %q not in allow-list for %s", signInProvider, email)
			writeFirebaseError(w, http.StatusForbidden, "method_not_allowed")
			return
		}
		// Resolve the DB user by email (case-tolerant).
		userID := ""
		if m.deps.UserIDByEmail != nil {
			userID = m.deps.UserIDByEmail(r.Context(), email)
			if userID == "" {
				userID = m.deps.UserIDByEmail(r.Context(), strings.ToLower(email))
			}
		}
		if userID == "" {
			log.Printf("[AUTH] firebase: no DB user maps to email %s", email)
			writeFirebaseError(w, http.StatusForbidden, "no_account")
			return
		}
		// Mint a server-side session (workspace-less; routePrincipals rotates it
		// to a principal-stamped row, exactly like the password path).
		token, err := minter(r.Context(), userID)
		if err != nil || token == "" {
			log.Printf("[AUTH] firebase: mint session failed for user %s: %v", userID, err)
			writeFirebaseError(w, http.StatusInternalServerError, "session")
			return
		}
		sessionMw.SetSessionCookie(w, token)
		m.deps.CSRFIssuer(w, m.deps.CSRFSecret, token, "")
		// Observability: a success line for the firebase login (failures are
		// already logged above). userID + method, never the token.
		log.Printf("[AUTH] firebase login OK: user=%s method=%s", userID, signInProvider)

		if principalLoader == nil || !principalLoader.IsEnabled() {
			writeFirebaseRedirect(w, entydad.DefaultAppRedirectURL)
			return
		}
		writeFirebaseRedirect(w, m.routePrincipals(w, r, token, userID))
	}
}

// signInMethodAllowed reports whether the Firebase sign_in_provider claim is in
// the configured allow-list, case-insensitively.
func signInMethodAllowed(allowed []string, method string) bool {
	for _, a := range allowed {
		if strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(method)) {
			return true
		}
	}
	return false
}

// writeFirebaseRedirect answers the Firebase login fetch with the post-login
// navigation target.
func writeFirebaseRedirect(w http.ResponseWriter, target string) {
	writeFirebaseBody(w, http.StatusOK, map[string]string{"redirect": target})
}

// writeFirebaseError answers the Firebase login fetch with a short error code
// (never the raw error — not localisable, may leak internals).
func writeFirebaseError(w http.ResponseWriter, status int, code string) {
	writeFirebaseBody(w, status, map[string]string{"error": code})
}

// writeFirebaseJSON is the bare-message form used before a code is resolved.
func writeFirebaseJSON(w http.ResponseWriter, status int, msg string) {
	writeFirebaseBody(w, status, map[string]string{"error": msg})
}

func writeFirebaseBody(w http.ResponseWriter, status int, body map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// handleSignup returns the POST /auth/signup handler.
func (m *AuthModule) handleSignup() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter
	sessionMw := m.deps.SessionManager

	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH_FIREBASE_ALLOW_SIGNUPS=false → self-signup disabled: reject the
		// endpoint, not just hide the link (defense-in-depth).
		if !m.deps.AllowSignups {
			http.Redirect(w, r, entydad.AuthLoginURL, http.StatusSeeOther)
			return
		}
		if authAdapter == nil || sessionMw == nil {
			// mock_auth: no real adapter — redirect to login
			http.Redirect(w, r, entydad.AuthLoginURL, http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		if password != confirmPassword {
			http.Redirect(w, r, entydad.AuthSignupURL+"?error=mismatch", http.StatusSeeOther)
			return
		}
		if _, err := authAdapter.Register(r.Context(), email, password, firstName, lastName, ""); err != nil {
			log.Printf("[AUTH] register failed for %s: %v", email, err)
			code := classifySignupError(err)
			http.Redirect(w, r, entydad.AuthSignupURL+"?error="+code, http.StatusSeeOther)
			return
		}
		// Auto-login after successful registration
		token, _, err := authAdapter.Login(r.Context(), email, password)
		if err != nil {
			// Registered but login failed — send to login page with success hint
			http.Redirect(w, r, entydad.AuthLoginURL+"?registered=true", http.StatusSeeOther)
			return
		}
		sessionMw.SetSessionCookie(w, token)
		// C2: signup has no workspace_id yet; workspace claim filled on
		// first GET after principal resolution.
		m.deps.CSRFIssuer(w, m.deps.CSRFSecret, token, "")
		http.Redirect(w, r, entydad.DefaultAppRedirectURL, http.StatusSeeOther)
	}
}

// handleResetPasswordRequest returns the POST /auth/reset-password handler
// (request step — sends the reset email).
func (m *AuthModule) handleResetPasswordRequest() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter

	return func(w http.ResponseWriter, r *http.Request) {
		if authAdapter == nil {
			// mock_auth: no-op — redirect to confirm sent page
			http.Redirect(w, r, entydad.AuthResetPasswordURL+"?sent=true", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		email := r.FormValue("email")
		// Ignore errors to prevent email enumeration
		resetToken, err := authAdapter.RequestPasswordReset(r.Context(), email)
		if err == nil && resetToken != "" {
			log.Printf("[AUTH] Password reset token for %s: %s", email, resetToken)
			log.Printf("[AUTH] Reset URL: %s?token=%s", entydad.AuthResetConfirmURL, resetToken)

			// In test mode, store the raw token keyed by user_id so the
			// test-only GET /test/last-reset-token endpoint can return it.
			if m.deps.TestMode && m.deps.UserIDByEmail != nil {
				userID := m.deps.UserIDByEmail(r.Context(), email)
				if userID != "" {
					m.lastResetTokens.Store(userID, resetToken)
					log.Printf("[AUTH] [TEST] stored reset token for user %s", userID)
				}
			}
		}
		// Always show success regardless of outcome
		http.Redirect(w, r, entydad.AuthResetPasswordURL+"?sent=true", http.StatusSeeOther)
	}
}

// handleResetPasswordConfirm returns the POST /auth/reset-password/confirm
// handler (confirm step — sets the new password).
func (m *AuthModule) handleResetPasswordConfirm() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter

	return func(w http.ResponseWriter, r *http.Request) {
		if authAdapter == nil {
			// mock_auth: no-op — redirect to login
			http.Redirect(w, r, entydad.AuthLoginURL+"?reset=true", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		token := r.FormValue("token")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")
		if newPassword != confirmPassword {
			http.Redirect(w, r, entydad.AuthResetConfirmURL+"?token="+url.QueryEscape(token)+"&error=mismatch", http.StatusSeeOther)
			return
		}
		if err := authAdapter.ExecutePasswordReset(r.Context(), token, newPassword); err != nil {
			log.Printf("[AUTH] password reset confirm failed: %v", err)
			// Map adapter error to a short code so the page handler can pick
			// the right lyngua-loaded label. Never emit err.Error() — it's
			// not localisable and may leak internals.
			code := classifyResetError(err)
			http.Redirect(w, r, entydad.AuthResetConfirmURL+"?token="+url.QueryEscape(token)+"&error="+code, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, entydad.AuthLoginURL+"?reset=true", http.StatusSeeOther)
	}
}

// handleChangePassword returns the POST /auth/change-password handler.
func (m *AuthModule) handleChangePassword() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter

	return func(w http.ResponseWriter, r *http.Request) {
		// Federated / IdP-managed deployments have no local password to change
		// (AUTH_FIREBASE_PASSWORD_CHANGE_ENABLED / derived) — reject the endpoint,
		// not just hide the link (defense-in-depth).
		if !m.deps.AllowPasswordChange {
			http.Redirect(w, r, entydad.DefaultAppRedirectURL, http.StatusSeeOther)
			return
		}
		if authAdapter == nil {
			// mock_auth: no-op — redirect back to app
			http.Redirect(w, r, entydad.DefaultAppRedirectURL, http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")
		if newPassword != confirmPassword {
			http.Redirect(w, r, entydad.AuthChangePasswordURL+"?error=mismatch", http.StatusSeeOther)
			return
		}
		id, ok := identity.FromContext(r.Context())
		if !ok || id.UserID == "" {
			http.Redirect(w, r, entydad.AuthLoginURL, http.StatusSeeOther)
			return
		}
		userID := id.UserID
		if err := authAdapter.ChangePassword(r.Context(), userID, oldPassword, newPassword); err != nil {
			log.Printf("[AUTH] change password failed for user %s: %v", userID, err)
			// Map adapter error to a short code so the page handler can pick
			// the right lyngua-loaded label. Never emit err.Error() — it's
			// not localisable and may leak internals.
			code := classifyChangePasswordError(err)
			http.Redirect(w, r, entydad.AuthChangePasswordURL+"?error="+code, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, entydad.AuthChangePasswordURL+"?success=1", http.StatusSeeOther)
	}
}

// handleSelectWorkspaceRole returns the GET /auth/select-workspace-role handler.
// Rendered when login resolves to 2+ principal bindings. The page reuses the
// login02 auth shell (no app shell, no sidebar).
func (m *AuthModule) handleSelectWorkspaceRole() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter
	principalLoader := m.deps.PrincipalResolver
	deps := m.deps

	// resolveUserFromCookie is shared between the chooser GET and the
	// switch POST. Both run under the /auth/ exclude-prefix so the
	// session middleware does NOT inject the user into context — we
	// re-validate the cookie token ourselves.
	resolveUserFromCookie := func(r *http.Request) string {
		if authAdapter == nil {
			return ""
		}
		cookie, err := r.Cookie(deps.SessionCookieName)
		if err != nil || cookie.Value == "" {
			return ""
		}
		userID, err := authAdapter.ValidateSession(r.Context(), cookie.Value)
		if err != nil {
			return ""
		}
		return userID
	}

	chooseDeps := &selectWorkspaceRole.Deps{
		Labels:        selectWorkspaceRole.DefaultLabels(),
		Login02:       deps.Labels.Login02,
		CommonLabels:  deps.Labels.Common,
		LogoText:      deps.LogoText,
		LogoIcon:      deps.LogoIcon,
		SwitchPostURL: "/action/auth/switch-principal",
		LogoutURL:     entydad.AuthLogoutURL,
		ResolveCards: func(ctx context.Context) []selectWorkspaceRole.PrincipalCard {
			if principalLoader == nil || !principalLoader.IsEnabled() {
				return nil
			}
			// Prefer ctx-injected user (test paths) before falling
			// back to cookie validation done by the wrapping handler.
			var userID string
			if id, ok := identity.FromContext(ctx); ok {
				userID = id.UserID
			}
			if userID == "" {
				if v, ok := ctx.Value(ctxKeyAuthUserID).(string); ok {
					userID = v
				}
			}
			if userID == "" {
				return nil
			}
			principals, perr := principalLoader.Resolve(ctx, userID)
			if perr != nil {
				log.Printf("[AUTH] select-workspace-role resolve error: %v", perr)
				return nil
			}
			cards := make([]selectWorkspaceRole.PrincipalCard, 0, len(principals))
			for _, p := range principals {
				cards = append(cards, selectWorkspaceRole.PrincipalCard{
					Kind:        PrincipalTypeString(p.Type),
					PrincipalID: p.PrincipalID,
					DisplayName: p.DisplayName,
					IconName:    iconForPrincipalKind(p.Type),
				})
			}
			return cards
		},
	}
	chooseViewObj := selectWorkspaceRole.NewView(chooseDeps)

	return func(w http.ResponseWriter, r *http.Request) {
		userID := resolveUserFromCookie(r)
		if userID == "" {
			// No valid session → bounce to login. We never render the
			// chooser to an unauthenticated visitor.
			http.Redirect(w, r, entydad.AuthLoginURL, http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKeyAuthUserID, userID)
		result := chooseViewObj.Handle(ctx, &view.ViewContext{
			Request:     r,
			CurrentPath: r.URL.Path,
		})
		m.renderAuthView(w, r, result)
	}
}

// handleSwitchPrincipal returns the POST /action/auth/switch-principal handler
// — the security-critical switch handler.
//
// Form fields:
//
//	principal_id          (required) — id of the grant row to switch to
//	principal_kind        (optional) — hint to disambiguate id collisions
//	acting_as_client_id   (optional) — for delegate-of-N>1
//	acting_as_supplier_id (optional) — same
//
// The handler re-runs principalLoader.Resolve and locates the
// matching principal in the loader's authoritative list. A user
// trying to forge a principal_id they don't actually hold is
// rejected with a 403-equivalent redirect to the chooser.
func (m *AuthModule) handleSwitchPrincipal() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter
	sessionMw := m.deps.SessionManager
	principalLoader := m.deps.PrincipalResolver
	deps := m.deps

	// resolveUserFromCookie — same closure as selectWorkspaceRole.
	resolveUserFromCookie := func(r *http.Request) string {
		if authAdapter == nil {
			return ""
		}
		cookie, err := r.Cookie(deps.SessionCookieName)
		if err != nil || cookie.Value == "" {
			return ""
		}
		userID, err := authAdapter.ValidateSession(r.Context(), cookie.Value)
		if err != nil {
			return ""
		}
		return userID
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if authAdapter == nil || sessionMw == nil || principalLoader == nil {
			http.Redirect(w, r, entydad.AuthLoginURL, http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		var userID string
		if id, ok := identity.FromContext(ctx); ok {
			userID = id.UserID
		}
		if userID == "" {
			// /action/* is normally session-middleware-protected, but
			// defensively re-validate via cookie when context didn't
			// carry the user — keeps the handler robust to middleware
			// reordering.
			userID = resolveUserFromCookie(r)
		}
		if userID == "" {
			http.Redirect(w, r, entydad.AuthLoginURL, http.StatusSeeOther)
			return
		}

		principalID := strings.TrimSpace(r.FormValue("principal_id"))
		kindHint := strings.TrimSpace(r.FormValue("principal_kind"))
		actingAsClientID := strings.TrimSpace(r.FormValue("acting_as_client_id"))
		actingAsSupplierID := strings.TrimSpace(r.FormValue("acting_as_supplier_id"))

		principals, perr := principalLoader.Resolve(ctx, userID)
		if perr != nil {
			log.Printf("[AUTH] switch-principal: resolve failed for user %s: %v", userID, perr)
			http.Redirect(w, r, "/auth/select-workspace-role?error=resolve", http.StatusSeeOther)
			return
		}
		target, ok := findPrincipalByID(principals, principalID, kindHint)
		if !ok {
			log.Printf("[AUTH] switch-principal: rejected — user %s does not hold principal_id=%s (kind=%s)", userID, principalID, kindHint)
			http.Redirect(w, r, "/auth/select-workspace-role?error=forbidden", http.StatusSeeOther)
			return
		}

		// codex RBC#1 — High-1 (2026-06-02). The chooser form posts only
		// principal_id + principal_kind (no acting-as fields today), so a
		// multi-target delegate (len(ActingAsTargets) > 1) chosen here
		// would be switched into with an EMPTY acting_as_* — persisting a
		// delegate principal with no act-as identity (→ fail-closed to zero
		// permissions). Fail closed: when the chosen principal is a
		// multi-target delegate and no acting-as id (matching one of its
		// own targets) was supplied, re-render the chooser with an explicit
		// ambiguous-target outcome rather than calling
		// executePrincipalSwitch. This shares the SAME guard as the login
		// auto-route, so neither caller can persist an empty acting-as. A
		// full per-acting-as-target chooser UX (a second pick step) is a
		// SEPARATE feature — FLAGGED, not built here. Single-target
		// delegates and non-delegate principals are unaffected.
		if !DelegateActingAsResolved(target, actingAsClientID, actingAsSupplierID) {
			log.Printf("[AUTH] switch-principal: blocked — user %s chose multi-target delegate principal_id=%s with no acting-as target; routing back to chooser", userID, principalID)
			http.Redirect(w, r, "/auth/select-workspace-role?error=select_target", http.StatusSeeOther)
			return
		}

		currentToken := ""
		if cookie, err := r.Cookie(deps.SessionCookieName); err == nil {
			currentToken = cookie.Value
		}

		// Pick the explicit-form use_case discriminator. When the form
		// supplies an explicit acting-as target, tag as acting_as;
		// otherwise tag generically as rotate (the primitive will
		// demote to in-place automatically when workspace matches —
		// the audit row's reason field will carry `rotated:false` for
		// that case, so the taxonomy stays honest).
		explicitUseCase := "switch_explicit_rotate"
		if actingAsClientID != "" || actingAsSupplierID != "" {
			explicitUseCase = "switch_explicit_acting_as"
		}
		result, sErr := m.deps.PrincipalSwitcher(ctx, PrincipalSwitchInput{
			UserID:             userID,
			Token:              currentToken,
			TargetPrincipal:    target,
			ActingAsClientID:   actingAsClientID,
			ActingAsSupplierID: actingAsSupplierID,
			// Explicit-form caller (POST /action/auth/switch-principal).
			// RequireAudit=false preserves dev-mode best-effort behavior.
			UseCase:      explicitUseCase,
			RequestURL:   r.URL.Path,
			Referer:      r.Header.Get("Referer"),
			SecFetchSite: r.Header.Get("Sec-Fetch-Site"),
			UserAgent:    r.Header.Get("User-Agent"),
			RequireAudit: false,
		})
		if sErr != nil {
			log.Printf("[AUTH] switch-principal: execute failed for user %s: %v", userID, sErr)
			http.Redirect(w, r, "/auth/select-workspace-role?error=switch", http.StatusSeeOther)
			return
		}
		if result.NewToken != "" {
			sessionMw.SetSessionCookie(w, result.NewToken)
		}
		// C2: always refresh CSRF cookie after principal switch — even
		// when NewToken is empty (in-place, same workspace). An in-place
		// switch may change principal_type / acting_as without changing
		// the session token, but the CSRF workspace claim must still
		// reflect the current workspace. Use original cookie token when
		// no rotation occurred.
		effectiveToken := result.NewToken
		if effectiveToken == "" {
			effectiveToken = currentToken
		}
		m.deps.CSRFIssuer(w, m.deps.CSRFSecret,
			effectiveToken, target.WorkspaceID)
		http.Redirect(w, r, m.homeURLForWorkspaceID(r.Context(), target.WorkspaceID), http.StatusSeeOther)
	}
}

// handleLogout returns the POST /auth/logout and /action/auth/logout handler.
func (m *AuthModule) handleLogout() http.HandlerFunc {
	authAdapter := m.deps.AuthAdapter
	deps := m.deps

	return func(w http.ResponseWriter, r *http.Request) {
		if authAdapter != nil {
			if cookie, err := r.Cookie(deps.SessionCookieName); err == nil && cookie.Value != "" {
				if invalidErr := authAdapter.InvalidateSession(r.Context(), cookie.Value); invalidErr != nil {
					log.Printf("[AUTH] logout: failed to invalidate session: %v", invalidErr)
				}
			}
		}
		// Q-SEC-3 (2026-05-31): clear the session cookie with SameSite=Strict
		// + Secure — the locked "pin logout SameSite=Strict" posture. The
		// espyna ClearSessionCookie default is SameSite=Lax; logout is always
		// an already-same-site terminal flow, so Strict is the safe stricter
		// choice, and Secure (following the COOKIE_SECURE policy) ensures
		// browsers honor the deletion of a Secure session cookie. Practical
		// cross-site logout-CSRF is already blocked because the Lax-issued
		// session cookie is not sent on cross-site POSTs; this hardens the
		// deletion side to match. The cookie name matches espyna's default
		// (never overridden in container.go).
		http.SetCookie(w, &http.Cookie{
			Name:     deps.SessionCookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   deps.SecureCookies(),
			SameSite: http.SameSiteStrictMode,
		})
		// Belt-and-suspenders: also expire the workspace CSRF cookie so no
		// stale ws_csrf token outlives the invalidated session. The session
		// invalidation above is the real guard; this just avoids a dangling
		// token lingering until the next GET refresh re-issues one.
		http.SetCookie(w, &http.Cookie{
			Name:     deps.WorkspaceCSRFCookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   deps.SecureCookies(),
			SameSite: http.SameSiteLaxMode,
		})
		// The sidebar logout form is intercepted by app-shell's
		// hx-boost="true". A plain http.Redirect would make HTMX
		// AJAX-follow the 303 and swap the full standalone /auth/login
		// document into #main-content — nesting the login UI inside
		// the now-defunct app shell. HX-Redirect tells HTMX to perform
		// a real browser navigation so the bare auth shell renders
		// cleanly. Same pattern as espyna's SessionMiddleware.redirectToLogin.
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", entydad.AuthLoginURL)
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Redirect(w, r, entydad.AuthLoginURL, http.StatusSeeOther)
	}
}

// handleTestLastResetToken returns the GET /test/last-reset-token handler.
// Returns the raw HMAC reset token for the given user so E2E specs can
// construct the confirm URL without a real email pipeline.
func (m *AuthModule) handleTestLastResetToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing user_id query param", http.StatusBadRequest)
			return
		}
		val, ok := m.lastResetTokens.Load(userID)
		if !ok {
			http.Error(w, "no reset token found for user", http.StatusNotFound)
			return
		}
		token, _ := val.(string)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(token))
	}
}
