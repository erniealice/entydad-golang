package user

// routes.go — User route struct, URL consts, and constructors (including the
// user-attached user↔role junction routes).
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, identity sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: UserRoutes -> Routes,
// DefaultUserRoutes -> DefaultRoutes, User<Xxx>URL -> <Xxx>URL. The cross-entity
// RoleListURL reference resolves DIRECT to the sibling role package
// (child→parent DAG; never through the facade).

import (
	role "github.com/erniealice/entydad-golang/domain/entity/identity/role"
)

// Default route constants for the user view.
const (
	DashboardURL       = "/users/dashboard"
	ListURL            = "/users/list/{status}"
	TableURL           = "/action/user/table/{status}"
	AddURL             = "/action/user/add"
	EditURL            = "/action/user/edit/{id}"
	DeleteURL          = "/action/user/delete"
	BulkDeleteURL      = "/action/user/bulk-delete"
	SetStatusURL       = "/action/user/set-status"
	BulkSetStatusURL   = "/action/user/bulk-set-status"
	SearchTimezonesURL = "/action/user/search-timezones"

	DetailURL           = "/users/detail/{id}"
	TabActionURL        = "/action/user/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/user/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/user/{id}/attachments/delete"
	ResetPasswordURL    = "/action/user/reset-password/{id}"

	// Legacy /manage/ user-roles routes
	RolesURL       = "/manage/users/{id}/roles"
	RolesTableURL  = "/action/manage/users/{id}/roles/table"
	RolesAssignURL = "/action/manage/users/{id}/roles/assign"
	RolesRemoveURL = "/action/manage/users/{id}/roles/remove"

	// Migrated /detail/ user-roles routes
	DetailRolesURL       = "/users/detail/{id}/roles"
	DetailRolesTableURL  = "/action/user/detail/{id}/roles/table"
	DetailRolesAssignURL = "/action/user/detail/{id}/roles/assign"
	DetailRolesRemoveURL = "/action/user/detail/{id}/roles/remove"
)

// Routes holds all route paths for user management.
type Routes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`
	ResetPasswordURL string `json:"reset_password_url"`

	// Timezone autocomplete search endpoint (returns JSON [{value,label}, ...])
	SearchTimezonesURL string `json:"search_timezones_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Legacy /manage/ user-roles routes (kept for backward compatibility)
	RolesURL       string `json:"roles_url"`
	RolesTableURL  string `json:"roles_table_url"`
	RolesAssignURL string `json:"roles_assign_url"`
	RolesRemoveURL string `json:"roles_remove_url"`

	// Migrated /detail/ user-roles routes
	DetailRolesURL       string `json:"detail_roles_url"`
	DetailRolesTableURL  string `json:"detail_roles_table_url"`
	DetailRolesAssignURL string `json:"detail_roles_assign_url"`
	DetailRolesRemoveURL string `json:"detail_roles_remove_url"`

	// Cross-app links used by the user dashboard quick-actions.
	RoleListURL       string `json:"role_list_url"`
	PermissionListURL string `json:"permission_list_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		TableURL:         TableURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
		DetailURL:        DetailURL,
		TabActionURL:     TabActionURL,
		ResetPasswordURL: ResetPasswordURL,

		SearchTimezonesURL: SearchTimezonesURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		// Legacy /manage/ routes
		RolesURL:       RolesURL,
		RolesTableURL:  RolesTableURL,
		RolesAssignURL: RolesAssignURL,
		RolesRemoveURL: RolesRemoveURL,

		// Migrated /detail/ routes
		DetailRolesURL:       DetailRolesURL,
		DetailRolesTableURL:  DetailRolesTableURL,
		DetailRolesAssignURL: DetailRolesAssignURL,
		DetailRolesRemoveURL: DetailRolesRemoveURL,

		// Cross-app quick-action links for the user dashboard.
		RoleListURL:       role.ListURL,
		PermissionListURL: "/app/permissions/list/active", // PermissionListURL with {status}=active
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"user.dashboard":       r.DashboardURL,
		"user.list":            r.ListURL,
		"user.table":           r.TableURL,
		"user.add":             r.AddURL,
		"user.edit":            r.EditURL,
		"user.delete":          r.DeleteURL,
		"user.bulk_delete":     r.BulkDeleteURL,
		"user.set_status":      r.SetStatusURL,
		"user.bulk_set_status": r.BulkSetStatusURL,
		"user.detail":          r.DetailURL,
		"user.tab_action":      r.TabActionURL,

		"user.search_timezones": r.SearchTimezonesURL,

		"user.attachment.upload": r.AttachmentUploadURL,
		"user.attachment.delete": r.AttachmentDeleteURL,

		// Legacy /manage/ routes
		"user.role.list":   r.RolesURL,
		"user.role.table":  r.RolesTableURL,
		"user.role.assign": r.RolesAssignURL,
		"user.role.remove": r.RolesRemoveURL,

		// Migrated /detail/ routes
		"user.detail_role.list":   r.DetailRolesURL,
		"user.detail_role.table":  r.DetailRolesTableURL,
		"user.detail_role.assign": r.DetailRolesAssignURL,
		"user.detail_role.remove": r.DetailRolesRemoveURL,
	}
}
