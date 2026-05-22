package action

import (
	"context"
	"net/http"

	"github.com/erniealice/espyna-golang/consumer"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
)

// SwitchWorkspaceDeps holds dependencies for the switch workspace handler.
type SwitchWorkspaceDeps struct {
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
