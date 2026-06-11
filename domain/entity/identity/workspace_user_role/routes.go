package workspace_user_role

// routes.go — WorkspaceUserRole route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, identity sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: WorkspaceUserRoleRoutes
// -> Routes, DefaultWorkspaceUserRoleRoutes -> DefaultRoutes,
// WorkspaceUserRole<Xxx>URL -> <Xxx>URL.

// WorkspaceUserRole — Phase 3 assignment drawer routes.
const (
	AddURL         = "/action/workspace_user_role/add"
	DeleteURL      = "/action/workspace_user_role/delete/{id}"
	PermissionsURL = "/action/workspace_user_role/permissions"
	SearchRolesURL = "/action/workspace_user_role/search-roles"
)

// Routes holds all route paths for the workspace_user_role
// assignment drawer (Phase 3).
type Routes struct {
	AddURL         string `json:"add_url"`
	DeleteURL      string `json:"delete_url"`
	PermissionsURL string `json:"permissions_url"`
	SearchRolesURL string `json:"search_roles_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants.
func DefaultRoutes() Routes {
	return Routes{
		AddURL:         AddURL,
		DeleteURL:      DeleteURL,
		PermissionsURL: PermissionsURL,
		SearchRolesURL: SearchRolesURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"workspace_user_role.add":          r.AddURL,
		"workspace_user_role.delete":       r.DeleteURL,
		"workspace_user_role.permissions":  r.PermissionsURL,
		"workspace_user_role.search_roles": r.SearchRolesURL,
	}
}
