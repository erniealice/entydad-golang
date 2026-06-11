package role

// routes.go — Role route struct, URL consts, and constructors (including the
// role-attached role↔permission and role↔user junction routes).
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, identity sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: RoleRoutes -> Routes,
// DefaultRoleRoutes -> DefaultRoutes, Role<Xxx>URL -> <Xxx>URL.

// Default route constants for the role view.
const (
	DetailURL           = "/roles/detail/{id}"
	TabActionURL        = "/action/role/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/role/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/role/{id}/attachments/delete"
	ListURL             = "/roles/list"
	TableURL            = "/action/role/table"
	AddURL              = "/action/role/add"
	EditURL             = "/action/role/edit/{id}"
	DeleteURL           = "/action/role/delete"
	BulkDeleteURL       = "/action/role/bulk-delete"
	SetStatusURL        = "/action/role/set-status"
	BulkSetStatusURL    = "/action/role/bulk-set-status"

	// Legacy /manage/ role-permissions routes (kept for backward compatibility)
	PermissionsURL       = "/manage/roles/{id}/permissions"
	PermissionsTableURL  = "/action/manage/roles/{id}/permissions/table"
	PermissionsAssignURL = "/action/manage/roles/{id}/permissions/assign"
	PermissionsRemoveURL = "/action/manage/roles/{id}/permissions/remove"

	// Role-users routes
	UsersURL       = "/roles/detail/{id}/users"
	UsersTableURL  = "/action/role/detail/{id}/users/table"
	UsersAssignURL = "/action/role/detail/{id}/users/assign"
	UsersRemoveURL = "/action/role/detail/{id}/users/remove"
	UsersSearchURL = "/action/role/detail/{id}/users/search"

	// Migrated /detail/ role-permissions routes
	DetailPermissionsURL       = "/roles/detail/{id}/permissions"
	DetailPermissionsTableURL  = "/action/role/detail/{id}/permissions/table"
	DetailPermissionsAssignURL = "/action/role/detail/{id}/permissions/assign"
	DetailPermissionsRemoveURL = "/action/role/detail/{id}/permissions/remove"
)

// Routes holds all route paths for role management, including
// role-permission and role-user associations, plus legacy /manage/ paths.
type Routes struct {
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

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Legacy /manage/ role-permissions routes (kept for backward compatibility)
	PermissionsURL       string `json:"permissions_url"`
	PermissionsTableURL  string `json:"permissions_table_url"`
	PermissionsAssignURL string `json:"permissions_assign_url"`
	PermissionsRemoveURL string `json:"permissions_remove_url"`

	// Role-users routes
	UsersURL       string `json:"users_url"`
	UsersTableURL  string `json:"users_table_url"`
	UsersAssignURL string `json:"users_assign_url"`
	UsersRemoveURL string `json:"users_remove_url"`
	UsersSearchURL string `json:"users_search_url"`

	// Migrated /detail/ role-permissions routes
	DetailPermissionsURL       string `json:"detail_permissions_url"`
	DetailPermissionsTableURL  string `json:"detail_permissions_table_url"`
	DetailPermissionsAssignURL string `json:"detail_permissions_assign_url"`
	DetailPermissionsRemoveURL string `json:"detail_permissions_remove_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants.
func DefaultRoutes() Routes {
	return Routes{
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

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		// Legacy /manage/ routes
		PermissionsURL:       PermissionsURL,
		PermissionsTableURL:  PermissionsTableURL,
		PermissionsAssignURL: PermissionsAssignURL,
		PermissionsRemoveURL: PermissionsRemoveURL,

		// Role-users routes
		UsersURL:       UsersURL,
		UsersTableURL:  UsersTableURL,
		UsersAssignURL: UsersAssignURL,
		UsersRemoveURL: UsersRemoveURL,
		UsersSearchURL: UsersSearchURL,

		// Migrated /detail/ routes
		DetailPermissionsURL:       DetailPermissionsURL,
		DetailPermissionsTableURL:  DetailPermissionsTableURL,
		DetailPermissionsAssignURL: DetailPermissionsAssignURL,
		DetailPermissionsRemoveURL: DetailPermissionsRemoveURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"role.list":            r.ListURL,
		"role.table":           r.TableURL,
		"role.add":             r.AddURL,
		"role.edit":            r.EditURL,
		"role.delete":          r.DeleteURL,
		"role.bulk_delete":     r.BulkDeleteURL,
		"role.set_status":      r.SetStatusURL,
		"role.bulk_set_status": r.BulkSetStatusURL,
		"role.detail":          r.DetailURL,
		"role.tab_action":      r.TabActionURL,

		"role.attachment.upload": r.AttachmentUploadURL,
		"role.attachment.delete": r.AttachmentDeleteURL,

		// Legacy /manage/ routes
		"role.permission.list":   r.PermissionsURL,
		"role.permission.table":  r.PermissionsTableURL,
		"role.permission.assign": r.PermissionsAssignURL,
		"role.permission.remove": r.PermissionsRemoveURL,

		// Role-users routes
		"role.user.list":   r.UsersURL,
		"role.user.table":  r.UsersTableURL,
		"role.user.assign": r.UsersAssignURL,
		"role.user.remove": r.UsersRemoveURL,
		"role.user.search": r.UsersSearchURL,

		// Migrated /detail/ routes
		"role.detail_permission.list":   r.DetailPermissionsURL,
		"role.detail_permission.table":  r.DetailPermissionsTableURL,
		"role.detail_permission.assign": r.DetailPermissionsAssignURL,
		"role.detail_permission.remove": r.DetailPermissionsRemoveURL,
	}
}
