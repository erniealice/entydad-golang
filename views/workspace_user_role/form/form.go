package form

import (
	"github.com/erniealice/entydad-golang"
)

// Data is the template data for the "Assign role" drawer form.
type Data struct {
	FormAction         string
	WorkspaceID         string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	WorkspaceUserID    string
	WorkspaceUserName  string
	WorkspaceUserEmail string
	SearchRolesURL     string
	PermissionsURL     string
	Labels             entydad.WorkspaceUserRoleLabels
	CommonLabels       any
}

// PermissionsData is the template data for the reactive permissions partial.
type PermissionsData struct {
	Permissions []PermissionItem
}

// PermissionItem holds a single permission code for display.
type PermissionItem struct {
	Code string
}
