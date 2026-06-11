// workspace_user_role_module.go provides the view module for the workspace_user_role
// assignment drawer and reactive permissions partial.
// This is Phase 3 of the bootstrap-auth plan.
//
// Routes registered:
//
//	GET  /action/workspace_user_role/add?workspace_user_id={wu}  — drawer form
//	POST /action/workspace_user_role/add                         — create junction
//	GET  /action/workspace_user_role/permissions?role_id={id}    — permissions partial
//	GET  /action/workspace_user_role/search-roles?q={q}          — autocomplete JSON
//	GET  /action/workspace_user_role/delete/{id}                 — delete confirm
//	POST /action/workspace_user_role/delete/{id}                 — soft-delete
package identity

import (
	"context"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	workspaceuserrole "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user_role"
	wuaction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user_role/action"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
)

// WorkspaceUserRoleModuleDeps holds all dependencies for the workspace_user_role module.
type WorkspaceUserRoleModuleDeps struct {
	Routes       workspaceuserrole.Routes
	Labels       workspaceuserrole.Labels
	CommonLabels any

	// GetWorkspaceUserItemPageData loads a WorkspaceUser (with nested user) by ID.
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	// CreateWorkspaceUserRole creates the junction row.
	CreateWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.CreateWorkspaceUserRoleRequest) (*workspaceuserrolepb.CreateWorkspaceUserRoleResponse, error)
	// DeleteWorkspaceUserRole soft-deletes a workspace_user_role row.
	DeleteWorkspaceUserRole func(ctx context.Context, req *workspaceuserrolepb.DeleteWorkspaceUserRoleRequest) (*workspaceuserrolepb.DeleteWorkspaceUserRoleResponse, error)
	// ListRoles lists roles for the search-roles autocomplete.
	ListRoles func(ctx context.Context, req *rolepb.ListRolesRequest) (*rolepb.ListRolesResponse, error)
}

// WorkspaceUserRoleModule holds all constructed workspace_user_role views.
type WorkspaceUserRoleModule struct {
	routes      workspaceuserrole.Routes
	Add         view.View
	Delete      view.View
	Permissions view.View
	SearchRoles http.HandlerFunc
}

// NewWorkspaceUserRoleModule constructs all workspace_user_role views from deps.
func NewWorkspaceUserRoleModule(deps *WorkspaceUserRoleModuleDeps) *WorkspaceUserRoleModule {
	actionDeps := &wuaction.Deps{
		Routes:                       deps.Routes,
		GetWorkspaceUserItemPageData: deps.GetWorkspaceUserItemPageData,
		CreateWorkspaceUserRole:      deps.CreateWorkspaceUserRole,
		DeleteWorkspaceUserRole:      deps.DeleteWorkspaceUserRole,
		ListRoles:                    deps.ListRoles,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
	}

	return &WorkspaceUserRoleModule{
		routes:      deps.Routes,
		Add:         wuaction.NewAddAction(actionDeps),
		Delete:      wuaction.NewDeleteAction(actionDeps),
		Permissions: wuaction.NewPermissionsAction(actionDeps),
		SearchRoles: wuaction.NewSearchRolesAction(actionDeps),
	}
}

// RegisterRoutes registers all workspace_user_role routes into the app router.
func (m *WorkspaceUserRoleModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.routes.AddURL != "" {
		r.GET(m.routes.AddURL, m.Add)
		r.POST(m.routes.AddURL, m.Add)
	}
	if m.routes.DeleteURL != "" {
		r.GET(m.routes.DeleteURL, m.Delete)
		r.POST(m.routes.DeleteURL, m.Delete)
	}
	if m.routes.PermissionsURL != "" {
		r.GET(m.routes.PermissionsURL, m.Permissions)
	}
	// SearchRoles returns JSON — uses HandleFunc if available.
	if m.routes.SearchRolesURL != "" && m.SearchRoles != nil {
		if full, ok := r.(identityRouteRegistrarFull); ok {
			full.HandleFunc("GET", m.routes.SearchRolesURL, m.SearchRoles)
		}
	}
}
