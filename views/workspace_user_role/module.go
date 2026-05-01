// Package workspace_user_role provides the view module for the workspace_user_role
// assignment drawer and reactive permissions partial.
// This is Phase 3 of the bootstrap-auth plan.
//
// Routes registered:
//   GET  /action/workspace_user_role/add?workspace_user_id={wu}  — drawer form
//   POST /action/workspace_user_role/add                         — create junction
//   GET  /action/workspace_user_role/permissions?role_id={id}    — permissions partial
//   GET  /action/workspace_user_role/search-roles?q={q}          — autocomplete JSON
//   GET  /action/workspace_user_role/delete/{id}                 — delete confirm
//   POST /action/workspace_user_role/delete/{id}                 — soft-delete
package workspace_user_role

import (
	"context"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	wuaction "github.com/erniealice/entydad-golang/views/workspace_user_role/action"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
)

// routeRegistrarFull extends view.RouteRegistrar with HandleFunc support.
type routeRegistrarFull interface {
	view.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// ModuleDeps holds all dependencies for the workspace_user_role module.
type ModuleDeps struct {
	Routes       entydad.WorkspaceUserRoleRoutes
	Labels       entydad.WorkspaceUserRoleLabels
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

// Module holds all constructed workspace_user_role views.
type Module struct {
	routes      entydad.WorkspaceUserRoleRoutes
	Add         view.View
	Delete      view.View
	Permissions view.View
	SearchRoles http.HandlerFunc
}

// NewModule constructs all workspace_user_role views from deps.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &wuaction.Deps{
		Routes:                      deps.Routes,
		GetWorkspaceUserItemPageData: deps.GetWorkspaceUserItemPageData,
		CreateWorkspaceUserRole:     deps.CreateWorkspaceUserRole,
		DeleteWorkspaceUserRole:     deps.DeleteWorkspaceUserRole,
		ListRoles:                   deps.ListRoles,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
	}

	return &Module{
		routes:      deps.Routes,
		Add:         wuaction.NewAddAction(actionDeps),
		Delete:      wuaction.NewDeleteAction(actionDeps),
		Permissions: wuaction.NewPermissionsAction(actionDeps),
		SearchRoles: wuaction.NewSearchRolesAction(actionDeps),
	}
}

// RegisterRoutes registers all workspace_user_role routes into the app router.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
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
		if full, ok := r.(routeRegistrarFull); ok {
			full.HandleFunc("GET", m.routes.SearchRolesURL, m.SearchRoles)
		}
	}
}
