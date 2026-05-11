// Package action provides handlers for workspace_user mutations.
// Handles: Add (assign user to workspace), Delete, SetStatus.
// Search (JSON autocomplete for user selection on the add form).
package action

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"

	"github.com/erniealice/entydad-golang/views/workspace_user/form"
)

// Deps holds dependencies for workspace_user action handlers.
type Deps struct {
	Routes                 entydad.WorkspaceUserRoutes
	CreateWorkspaceUser    func(ctx context.Context, req *workspaceuserpb.CreateWorkspaceUserRequest) (*workspaceuserpb.CreateWorkspaceUserResponse, error)
	DeleteWorkspaceUser    func(ctx context.Context, req *workspaceuserpb.DeleteWorkspaceUserRequest) (*workspaceuserpb.DeleteWorkspaceUserResponse, error)
	SetWorkspaceUserActive func(ctx context.Context, id string, active bool) error
	// ListUsers is used by the user search endpoint to find users for autocomplete.
	ListUsers func(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error)
}

// searchOption is the JSON shape returned by the user search handler.
type searchOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// NewAddAction creates the workspace_user add action (GET = form, POST = create).
// GET: renders the "Add user to workspace" drawer form.
// POST: creates the workspace_user junction record.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("workspace_user", "create") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			workspaceID := viewCtx.Request.URL.Query().Get("workspace_id")
			return view.OK("workspace-user-add-form", &form.Data{
				FormAction:    deps.Routes.AddURL,
				WorkspaceID:   workspaceID,
				Labels:        entydad.WorkspaceUserLabels{},
				UserSearchURL: deps.Routes.SearchURL,
				CommonLabels:  nil,
			})
		}

		// POST — create workspace_user
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		workspaceID := r.FormValue("workspace_id")
		userID := r.FormValue("user_id")

		if workspaceID == "" || userID == "" {
			return entydad.HTMXError("workspace_id and user_id are required")
		}

		_, err := deps.CreateWorkspaceUser(ctx, &workspaceuserpb.CreateWorkspaceUserRequest{
			Data: &workspaceuserpb.WorkspaceUser{
				WorkspaceId: workspaceID,
				UserId:      userID,
				Active:      true,
			},
		})
		if err != nil {
			log.Printf("Failed to create workspace_user: %v", err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("workspace-users-table")
	})
}

// NewDeleteAction creates the workspace_user delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("workspace_user", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			id = viewCtx.Request.URL.Query().Get("id")
		}
		if id == "" {
			return entydad.HTMXError("id is required")
		}

		_, err := deps.DeleteWorkspaceUser(ctx, &workspaceuserpb.DeleteWorkspaceUserRequest{
			Data: &workspaceuserpb.WorkspaceUser{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete workspace_user %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("workspace-users-table")
	})
}

// NewSetStatusAction creates the workspace_user set-status action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("workspace_user", "update") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			id = viewCtx.Request.URL.Query().Get("id")
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}
		active := viewCtx.Request.FormValue("active") == "true"

		if err := deps.SetWorkspaceUserActive(ctx, id, active); err != nil {
			log.Printf("Failed to set workspace_user active %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("workspace-users-table")
	})
}

// NewUserSearchAction returns an http.HandlerFunc that searches users for the
// "Add user to workspace" autocomplete input.
// GET /action/workspace_user/search?q=term
// Returns JSON: [{"value":"user_id","label":"First Last (email)"}, ...]
func NewUserSearchAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

		if query == "" {
			writeSearchJSON(w, []searchOption{})
			return
		}

		if deps.ListUsers == nil {
			writeSearchJSON(w, []searchOption{})
			return
		}

		resp, err := deps.ListUsers(ctx, &userpb.ListUsersRequest{})
		if err != nil {
			log.Printf("workspace_user search: failed to list users: %v", err)
			writeSearchJSON(w, []searchOption{})
			return
		}

		var results []searchOption
		for _, u := range resp.GetData() {
			if !u.GetActive() {
				continue
			}
			name := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			email := u.GetEmailAddress()
			label := name
			if email != "" {
				label = label + " (" + email + ")"
			}
			if !strings.Contains(strings.ToLower(label), query) {
				continue
			}
			results = append(results, searchOption{
				Value: u.GetId(),
				Label: label,
			})
		}

		if results == nil {
			results = []searchOption{}
		}
		writeSearchJSON(w, results)
	}
}

// writeSearchJSON marshals data as JSON and writes it to the response writer.
func writeSearchJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("workspace_user search: failed to encode JSON response: %v", err)
	}
}
