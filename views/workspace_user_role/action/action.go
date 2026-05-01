// Package action provides handlers for workspace_user_role mutations and partials.
// Phase 3 of the bootstrap-auth plan.
//
// Routes served:
//   GET  /action/workspace_user_role/add?workspace_user_id={wu}  — drawer form
//   POST /action/workspace_user_role/add                         — create junction row
//   GET  /action/workspace_user_role/permissions?role_id={id}    — reactive permissions partial
//   GET  /action/workspace_user_role/search-roles?q={q}&workspace_id={ws} — autocomplete JSON
//   GET  /action/workspace_user_role/delete/{id}                 — confirm (unused; form-based)
//   POST /action/workspace_user_role/delete/{id}                 — soft-delete
package action

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
)

// Deps holds dependencies for workspace_user_role action handlers.
type Deps struct {
	Routes entydad.WorkspaceUserRoleRoutes
	// GetWorkspaceUserItemPageData loads a WorkspaceUser (with nested user) by ID.
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	// CreateWorkspaceUserRole creates the junction row.
	CreateWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.CreateWorkspaceUserRoleRequest) (*workspaceuserrolepb.CreateWorkspaceUserRoleResponse, error)
	// DeleteWorkspaceUserRole soft-deletes a workspace_user_role row.
	DeleteWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.DeleteWorkspaceUserRoleRequest) (*workspaceuserrolepb.DeleteWorkspaceUserRoleResponse, error)
	// ListRoles lists all roles (used for search-roles autocomplete).
	ListRoles func(ctx context.Context, req *rolepb.ListRolesRequest) (*rolepb.ListRolesResponse, error)
	// Labels provides i18n strings for the drawer form.
	Labels entydad.WorkspaceUserRoleLabels
	// CommonLabels provides shared labels (for sheet footer).
	CommonLabels any
}

// AssignFormData is the template data for the "Assign role" drawer form.
type AssignFormData struct {
	FormAction          string
	WorkspaceUserID     string
	WorkspaceUserName   string
	WorkspaceUserEmail  string
	SearchRolesURL      string
	PermissionsURL      string
	Labels              entydad.WorkspaceUserRoleLabels
	CommonLabels        any
}

// PermissionsData is the template data for the reactive permissions partial.
type PermissionsData struct {
	Permissions []PermissionItem
}

// PermissionItem holds a single permission code for display.
type PermissionItem struct {
	Code string
}

// searchOption is the JSON shape returned by the search-roles handler.
type searchOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// NewAddAction creates the workspace_user_role add action (GET = form, POST = create).
// GET: renders the "Assign role" drawer form with read-only workspace user label.
// POST: creates the workspace_user_role junction record.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("workspace_user_role", "create") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			workspaceUserID := viewCtx.Request.URL.Query().Get("workspace_user_id")

			// Load workspace user to populate the read-only label.
			userName := ""
			email := ""
			if workspaceUserID != "" && deps.GetWorkspaceUserItemPageData != nil {
				resp, err := deps.GetWorkspaceUserItemPageData(ctx, &workspaceuserpb.GetWorkspaceUserItemPageDataRequest{
					WorkspaceUserId: workspaceUserID,
				})
				if err != nil {
					log.Printf("workspace_user_role add form: failed to load workspace_user %s: %v", workspaceUserID, err)
				} else if wu := resp.GetWorkspaceUser(); wu != nil {
					if u := wu.GetUser(); u != nil {
						userName = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
						email = u.GetEmailAddress()
					}
				}
			}

			return view.OK("wur-assign-form", &AssignFormData{
				FormAction:         deps.Routes.AddURL,
				WorkspaceUserID:    workspaceUserID,
				WorkspaceUserName:  userName,
				WorkspaceUserEmail: email,
				SearchRolesURL:     deps.Routes.SearchRolesURL,
				PermissionsURL:     deps.Routes.PermissionsURL,
				Labels:             deps.Labels,
				CommonLabels:       deps.CommonLabels,
			})
		}

		// POST — create workspace_user_role
		if err := viewCtx.Request.ParseForm(); err != nil {
			return entydad.HTMXError(viewCtx.T("shared.errors.invalidFormData"))
		}

		r := viewCtx.Request
		workspaceUserID := r.FormValue("workspace_user_id")
		roleID := r.FormValue("role_id")

		if workspaceUserID == "" || roleID == "" {
			return entydad.HTMXError("workspace_user_id and role_id are required")
		}

		if deps.CreateWorkspaceUserRole == nil {
			return entydad.HTMXError("CreateWorkspaceUserRole not wired")
		}

		_, err := deps.CreateWorkspaceUserRole(ctx, &workspaceuserrolepb.CreateWorkspaceUserRoleRequest{
			Data: &workspaceuserrolepb.WorkspaceUserRole{
				WorkspaceUserId: workspaceUserID,
				RoleId:          roleID,
				Active:          true,
			},
		})
		if err != nil {
			log.Printf("Failed to create workspace_user_role: %v", err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("workspace-user-roles-table")
	})
}

// NewDeleteAction creates the workspace_user_role delete action.
// GET: renders a confirmation view (reuses sheet pattern).
// POST: soft-deletes the workspace_user_role row.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("workspace_user_role", "delete") {
			return entydad.HTMXError(viewCtx.T("shared.errors.permissionDenied"))
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			id = viewCtx.Request.URL.Query().Get("id")
		}
		if id == "" {
			return entydad.HTMXError("id is required")
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Return a simple confirmation form inside the sheet.
			return view.OK("wur-delete-confirm", map[string]string{
				"ID":        id,
				"DeleteURL": deps.Routes.DeleteURL,
			})
		}

		// POST — soft-delete
		if deps.DeleteWorkspaceUserRole == nil {
			return entydad.HTMXError("DeleteWorkspaceUserRole not wired")
		}

		_, err := deps.DeleteWorkspaceUserRole(ctx, &workspaceuserrolepb.DeleteWorkspaceUserRoleRequest{
			Data: &workspaceuserrolepb.WorkspaceUserRole{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete workspace_user_role %s: %v", id, err)
			return entydad.HTMXError(err.Error())
		}

		return entydad.HTMXSuccess("workspace-user-roles-table")
	})
}

// NewPermissionsAction returns a view.View for the reactive permissions partial.
// GET /action/workspace_user_role/permissions?role_id={id}
// Returns HTML: <ul data-testid="wur-permissions-list">…</ul>
// Empty list when role_id is blank.
func NewPermissionsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		roleID := viewCtx.Request.URL.Query().Get("role_id")

		data := &PermissionsData{
			Permissions: []PermissionItem{},
		}

		if roleID == "" || deps.ListRoles == nil {
			return view.OK("wur-permissions-list", data)
		}

		resp, err := deps.ListRoles(ctx, &rolepb.ListRolesRequest{})
		if err != nil {
			log.Printf("workspace_user_role permissions: failed to list roles: %v", err)
			return view.OK("wur-permissions-list", data)
		}

		for _, role := range resp.GetData() {
			if role.GetId() != roleID {
				continue
			}
			for _, rp := range role.GetRolePermissions() {
				code := ""
				if p := rp.GetPermission(); p != nil {
					code = p.GetPermissionCode()
				}
				if code == "" {
					continue
				}
				data.Permissions = append(data.Permissions, PermissionItem{Code: code})
			}
			break
		}

		return view.OK("wur-permissions-list", data)
	})
}

// NewSearchRolesAction returns an http.HandlerFunc for the role search autocomplete.
// GET /action/workspace_user_role/search-roles?q={query}&workspace_id={ws}
// Scope: roles where workspace_id = $ws OR workspace_id IS NULL (global roles).
// Returns JSON: [{"value":"role_id","label":"Role Name"}, ...]
func NewSearchRolesAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		workspaceID := r.URL.Query().Get("workspace_id")

		if deps.ListRoles == nil {
			writeSearchJSON(w, []searchOption{})
			return
		}

		resp, err := deps.ListRoles(ctx, &rolepb.ListRolesRequest{})
		if err != nil {
			log.Printf("workspace_user_role search-roles: failed to list roles: %v", err)
			writeSearchJSON(w, []searchOption{})
			return
		}

		var results []searchOption
		for _, role := range resp.GetData() {
			if !role.GetActive() {
				continue
			}

			// Scope: include if workspace_id matches OR role has no workspace (global).
			roleWorkspaceID := role.GetWorkspaceId()
			if workspaceID != "" && roleWorkspaceID != "" && roleWorkspaceID != workspaceID {
				continue
			}

			name := role.GetName()
			if query != "" && !strings.Contains(strings.ToLower(name), query) {
				continue
			}

			results = append(results, searchOption{
				Value: role.GetId(),
				Label: name,
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
		log.Printf("workspace_user_role search: failed to encode JSON response: %v", err)
	}
}
