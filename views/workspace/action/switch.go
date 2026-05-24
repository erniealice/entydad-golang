package action

import (
	"context"
	"net/http"

	"github.com/erniealice/espyna-golang/consumer"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
)

// SecureSwitchResult is the outcome of an override-driven workspace switch.
// When NewToken is non-empty the handler will call SetSessionCookie before
// emitting HX-Redirect to RedirectURL.
type SecureSwitchResult struct {
	NewToken    string
	RedirectURL string
}

// SecureSwitchFn is the optional override that bypasses the legacy in-place
// SwitchWorkspace use case. When wired (A1 fix, WKR-P0-1, 2026-05-22) it
// MUST route the switch through the rotation primitive so cross-workspace
// switches rotate the session token and write an audit row. The provider
// (service-admin) is responsible for cookie management via SetSessionCookie
// (the handler owns the http.ResponseWriter and calls SetSessionCookie when
// NewToken is non-empty, then redirects).
type SecureSwitchFn func(ctx context.Context, userID, sessionToken, targetWorkspaceID string) (*SecureSwitchResult, error)

// SwitchWorkspaceDeps holds dependencies for the switch workspace handler.
type SwitchWorkspaceDeps struct {
	// SecureSwitch (A1 fix WKR-P0-1, 2026-05-22): when non-nil, the handler
	// uses this override INSTEAD of the legacy SwitchWorkspace use case.
	// Service-admin wires this to executePrincipalSwitch so that
	// /action/admin/switch-workspace (sidebar workspace-switcher) honours
	// the workspace-boundary rotation invariant and writes an audit row.
	// Legacy callers that haven't migrated still get the unrotated in-place
	// SwitchWorkspace behavior via the fallback branch below.
	SecureSwitch SecureSwitchFn
	// ResolveUserID extracts the authenticated user_id from the request.
	// Required when SecureSwitch is non-nil; ignored otherwise.
	ResolveUserID func(r *http.Request) string
	// SetSessionCookie writes the post-rotation session cookie to the
	// response. Required when SecureSwitch is non-nil and may return a new
	// token; ignored otherwise.
	SetSessionCookie func(w http.ResponseWriter, token string)

	SwitchWorkspace func(ctx context.Context, req *workspacepb.SwitchWorkspaceRequest) (*workspacepb.SwitchWorkspaceResponse, error)
	// HomeURLForWorkspaceID resolves the post-switch redirect URL given the
	// newly-active workspace_id. Optional — when nil the handler falls back
	// to HomeURL (or "/home" if both unset). Per Q-WS-1 the redirect should
	// land on /w/{slug}/home so the URL reflects the active workspace.
	HomeURLForWorkspaceID func(ctx context.Context, workspaceID string) string
	// HomeURL is the static fallback when HomeURLForWorkspaceID is nil.
	// Defaults to "/home" (post-P12 of workspace-keyed-routing plan;
	// "/app/home" is gone). The bare /home handler reads workspace from the
	// (post-switch) session.
	HomeURL string
}

// NewSwitchWorkspaceHandler creates an http.HandlerFunc that switches the
// active workspace for the current session and issues an HTMX full-page
// redirect to /app/home so the new workspace context is picked up.
//
// POST /action/admin/switch-workspace
// Form field: workspace_id (required)
func NewSwitchWorkspaceHandler(deps *SwitchWorkspaceDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		workspaceID := r.FormValue("workspace_id")
		if workspaceID == "" {
			http.Error(w, "workspace_id required", http.StatusBadRequest)
			return
		}

		// Get session token from cookie (try production name first, then dev name)
		cookie, err := r.Cookie(consumer.DefaultSessionCookieName)
		if err != nil {
			cookie, err = r.Cookie("session_token")
		}
		if err != nil {
			http.Error(w, "no session", http.StatusUnauthorized)
			return
		}

		// A1 fix (WKR-P0-1, 2026-05-22): when the host app has wired the
		// SecureSwitch override (service-admin does, via
		// executePrincipalSwitch), prefer it. The legacy SwitchWorkspace
		// use case performs an in-place mutation of session.workspace_id
		// with no rotation and no audit — violating the
		// workspace-boundary rotation invariant (Q-WS-13) and the
		// audit-on-every-switch invariant (red-team A-4 / X-2).
		if deps.SecureSwitch != nil {
			if deps.ResolveUserID == nil || deps.SetSessionCookie == nil {
				http.Error(w, "switch workspace misconfigured", http.StatusInternalServerError)
				return
			}
			userID := deps.ResolveUserID(r)
			if userID == "" {
				http.Error(w, "no session", http.StatusUnauthorized)
				return
			}
			result, sErr := deps.SecureSwitch(r.Context(), userID, cookie.Value, workspaceID)
			if sErr != nil || result == nil {
				msg := "failed to switch workspace"
				if sErr != nil {
					msg = sErr.Error()
				}
				http.Error(w, msg, http.StatusForbidden)
				return
			}
			if result.NewToken != "" {
				deps.SetSessionCookie(w, result.NewToken)
			}
			homeURL := result.RedirectURL
			if homeURL == "" {
				if deps.HomeURLForWorkspaceID != nil {
					homeURL = deps.HomeURLForWorkspaceID(r.Context(), workspaceID)
				}
				if homeURL == "" {
					homeURL = deps.HomeURL
				}
				if homeURL == "" {
					homeURL = "/home"
				}
			}
			w.Header().Set("HX-Redirect", homeURL)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Legacy fallback: hosts that haven't migrated to the rotation
		// primitive still get the unrotated in-place behavior. New hosts
		// MUST wire SecureSwitch.
		resp, err := deps.SwitchWorkspace(r.Context(), &workspacepb.SwitchWorkspaceRequest{
			WorkspaceId:  workspaceID,
			SessionToken: cookie.Value,
		})
		if err != nil || !resp.GetSuccess() {
			msg := "failed to switch workspace"
			if resp != nil && resp.GetError() != nil {
				msg = resp.GetError().GetMessage()
			}
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		// HTMX redirect to home (full page reload to pick up new workspace context).
		// Prefer /w/{slug}/home via HomeURLForWorkspaceID so the URL reflects the
		// new workspace (Q-WS-1 → A / Q-WS-13). Fall back to bare /home which
		// reads workspace from the just-switched session.
		var homeURL string
		if deps.HomeURLForWorkspaceID != nil {
			homeURL = deps.HomeURLForWorkspaceID(r.Context(), workspaceID)
		}
		if homeURL == "" {
			homeURL = deps.HomeURL
		}
		if homeURL == "" {
			homeURL = "/home"
		}
		w.Header().Set("HX-Redirect", homeURL)
		w.WriteHeader(http.StatusOK)
	}
}
