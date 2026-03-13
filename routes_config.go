package entydad

// Route configuration structs for the entydad domain.
//
// This file implements a three-level routing system:
//
//   Level 1 — Generic defaults from Go consts (this file).
//     DefaultXxxRoutes() constructors populate structs from the route constants
//     defined in routes.go. These are sensible defaults that work out of the box.
//
//   Level 2 — Industry-specific overrides via JSON (loaded by consumer apps).
//     Consumer apps can load a JSON configuration file and unmarshal it into
//     route structs, overriding some or all of the default paths. JSON tags on
//     every field enable this.
//
//   Level 3 — App-specific overrides via Go field assignment (optional).
//     After constructing or loading routes, the consumer app can directly set
//     individual fields to customize specific routes for its own needs.
//
// Each route struct also provides a RouteMap() method that returns a
// map[string]string keyed by dot-notation identifiers (e.g. "client.list"),
// suitable for template rendering and reverse-routing lookups.

// ---------------------------------------------------------------------------
// ClientRoutes
// ---------------------------------------------------------------------------

// ClientRoutes holds all route paths for client management, including
// client tags and dashboard views.
type ClientRoutes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	SearchURL        string `json:"search_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Report routes
	ReceivablesAgingURL string `json:"receivables_aging_url"`
}

// DefaultClientRoutes returns a ClientRoutes populated from the package-level
// route constants.
func DefaultClientRoutes() ClientRoutes {
	return ClientRoutes{
		DashboardURL:     ClientDashboardURL,
		ListURL:          ClientListURL,
		TableURL:         ClientTableURL,
		AddURL:           ClientAddURL,
		EditURL:          ClientEditURL,
		DeleteURL:        ClientDeleteURL,
		BulkDeleteURL:    ClientBulkDeleteURL,
		DetailURL:        ClientDetailURL,
		TabActionURL:     ClientTabActionURL,
		SetStatusURL:     ClientSetStatusURL,
		BulkSetStatusURL: ClientBulkSetStatusURL,
		SearchURL:        ClientSearchURL,

		AttachmentUploadURL: ClientAttachmentUploadURL,
		AttachmentDeleteURL: ClientAttachmentDeleteURL,

		ReceivablesAgingURL: ReceivablesAgingURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r ClientRoutes) RouteMap() map[string]string {
	return map[string]string{
		"client.dashboard":      r.DashboardURL,
		"client.list":           r.ListURL,
		"client.table":          r.TableURL,
		"client.add":            r.AddURL,
		"client.edit":           r.EditURL,
		"client.delete":         r.DeleteURL,
		"client.bulk_delete":    r.BulkDeleteURL,
		"client.detail":         r.DetailURL,
		"client.tab_action":     r.TabActionURL,
		"client.set_status":     r.SetStatusURL,
		"client.bulk_set_status": r.BulkSetStatusURL,
		"client.search":          r.SearchURL,

		"client.attachment.upload": r.AttachmentUploadURL,
		"client.attachment.delete": r.AttachmentDeleteURL,

		"client.receivables_aging": r.ReceivablesAgingURL,
	}
}

// ---------------------------------------------------------------------------
// UserRoutes
// ---------------------------------------------------------------------------

// UserRoutes holds all route paths for user management, including user-role
// associations and legacy /manage/ paths.
type UserRoutes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	DetailURL         string `json:"detail_url"`
	TabActionURL      string `json:"tab_action_url"`
	ResetPasswordURL  string `json:"reset_password_url"`

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
}

// DefaultUserRoutes returns a UserRoutes populated from the package-level
// route constants.
func DefaultUserRoutes() UserRoutes {
	return UserRoutes{
		DashboardURL:     UserDashboardURL,
		ListURL:          UserListURL,
		TableURL:         UserTableURL,
		AddURL:           UserAddURL,
		EditURL:          UserEditURL,
		DeleteURL:        UserDeleteURL,
		BulkDeleteURL:    UserBulkDeleteURL,
		SetStatusURL:     UserSetStatusURL,
		BulkSetStatusURL: UserBulkSetStatusURL,
		DetailURL:         UserDetailURL,
		TabActionURL:      UserTabActionURL,
		ResetPasswordURL:  UserResetPasswordURL,

		AttachmentUploadURL: UserAttachmentUploadURL,
		AttachmentDeleteURL: UserAttachmentDeleteURL,

		// Legacy /manage/ routes
		RolesURL:       UserRolesURL,
		RolesTableURL:  UserRolesTableURL,
		RolesAssignURL: UserRolesAssignURL,
		RolesRemoveURL: UserRolesRemoveURL,

		// Migrated /detail/ routes
		DetailRolesURL:       UserDetailRolesURL,
		DetailRolesTableURL:  UserDetailRolesTableURL,
		DetailRolesAssignURL: UserDetailRolesAssignURL,
		DetailRolesRemoveURL: UserDetailRolesRemoveURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r UserRoutes) RouteMap() map[string]string {
	return map[string]string{
		"user.dashboard":      r.DashboardURL,
		"user.list":           r.ListURL,
		"user.table":          r.TableURL,
		"user.add":            r.AddURL,
		"user.edit":           r.EditURL,
		"user.delete":         r.DeleteURL,
		"user.bulk_delete":    r.BulkDeleteURL,
		"user.set_status":     r.SetStatusURL,
		"user.bulk_set_status": r.BulkSetStatusURL,
		"user.detail":         r.DetailURL,
		"user.tab_action":     r.TabActionURL,

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

// ---------------------------------------------------------------------------
// RoleRoutes
// ---------------------------------------------------------------------------

// RoleRoutes holds all route paths for role management, including
// role-permission and role-user associations, plus legacy /manage/ paths.
type RoleRoutes struct {
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

	// Migrated /detail/ role-permissions routes
	DetailPermissionsURL       string `json:"detail_permissions_url"`
	DetailPermissionsTableURL  string `json:"detail_permissions_table_url"`
	DetailPermissionsAssignURL string `json:"detail_permissions_assign_url"`
	DetailPermissionsRemoveURL string `json:"detail_permissions_remove_url"`
}

// DefaultRoleRoutes returns a RoleRoutes populated from the package-level
// route constants.
func DefaultRoleRoutes() RoleRoutes {
	return RoleRoutes{
		ListURL:          RoleListURL,
		TableURL:         RoleTableURL,
		AddURL:           RoleAddURL,
		EditURL:          RoleEditURL,
		DeleteURL:        RoleDeleteURL,
		BulkDeleteURL:    RoleBulkDeleteURL,
		SetStatusURL:     RoleSetStatusURL,
		BulkSetStatusURL: RoleBulkSetStatusURL,
		DetailURL:        RoleDetailURL,
		TabActionURL:     RoleTabActionURL,

		AttachmentUploadURL: RoleAttachmentUploadURL,
		AttachmentDeleteURL: RoleAttachmentDeleteURL,

		// Legacy /manage/ routes
		PermissionsURL:       RolePermissionsURL,
		PermissionsTableURL:  RolePermissionsTableURL,
		PermissionsAssignURL: RolePermissionsAssignURL,
		PermissionsRemoveURL: RolePermissionsRemoveURL,

		// Role-users routes
		UsersURL:       RoleUsersURL,
		UsersTableURL:  RoleUsersTableURL,
		UsersAssignURL: RoleUsersAssignURL,
		UsersRemoveURL: RoleUsersRemoveURL,

		// Migrated /detail/ routes
		DetailPermissionsURL:       RoleDetailPermissionsURL,
		DetailPermissionsTableURL:  RoleDetailPermissionsTableURL,
		DetailPermissionsAssignURL: RoleDetailPermissionsAssignURL,
		DetailPermissionsRemoveURL: RoleDetailPermissionsRemoveURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r RoleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"role.list":           r.ListURL,
		"role.table":          r.TableURL,
		"role.add":            r.AddURL,
		"role.edit":           r.EditURL,
		"role.delete":         r.DeleteURL,
		"role.bulk_delete":    r.BulkDeleteURL,
		"role.set_status":     r.SetStatusURL,
		"role.bulk_set_status": r.BulkSetStatusURL,
		"role.detail":         r.DetailURL,
		"role.tab_action":     r.TabActionURL,

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

		// Migrated /detail/ routes
		"role.detail_permission.list":   r.DetailPermissionsURL,
		"role.detail_permission.table":  r.DetailPermissionsTableURL,
		"role.detail_permission.assign": r.DetailPermissionsAssignURL,
		"role.detail_permission.remove": r.DetailPermissionsRemoveURL,
	}
}

// ---------------------------------------------------------------------------
// LocationRoutes
// ---------------------------------------------------------------------------

// LocationRoutes holds all route paths for location management.
type LocationRoutes struct {
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`
	EditDetailURL    string `json:"edit_detail_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultLocationRoutes returns a LocationRoutes populated from the
// package-level route constants.
func DefaultLocationRoutes() LocationRoutes {
	return LocationRoutes{
		ListURL:          LocationListURL,
		DetailURL:        LocationDetailURL,
		TableURL:         LocationTableURL,
		AddURL:           LocationAddURL,
		EditURL:          LocationEditURL,
		DeleteURL:        LocationDeleteURL,
		BulkDeleteURL:    LocationBulkDeleteURL,
		SetStatusURL:     LocationSetStatusURL,
		BulkSetStatusURL: LocationBulkSetStatusURL,
		TabActionURL:     LocationTabActionURL,
		EditDetailURL:    LocationEditDetailURL,

		AttachmentUploadURL: LocationAttachmentUploadURL,
		AttachmentDeleteURL: LocationAttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r LocationRoutes) RouteMap() map[string]string {
	return map[string]string{
		"location.list":            r.ListURL,
		"location.detail":          r.DetailURL,
		"location.table":           r.TableURL,
		"location.add":             r.AddURL,
		"location.edit":            r.EditURL,
		"location.delete":          r.DeleteURL,
		"location.bulk_delete":     r.BulkDeleteURL,
		"location.set_status":      r.SetStatusURL,
		"location.bulk_set_status": r.BulkSetStatusURL,
		"location.tab_action":      r.TabActionURL,
		"location.edit_detail":     r.EditDetailURL,

		"location.attachment.upload": r.AttachmentUploadURL,
		"location.attachment.delete": r.AttachmentDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// PermissionRoutes
// ---------------------------------------------------------------------------

// PermissionRoutes holds all route paths for permission management.
type PermissionRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultPermissionRoutes returns a PermissionRoutes populated from the
// package-level route constants.
func DefaultPermissionRoutes() PermissionRoutes {
	return PermissionRoutes{
		ListURL:          PermissionListURL,
		TableURL:         PermissionTableURL,
		AddURL:           PermissionAddURL,
		EditURL:          PermissionEditURL,
		DeleteURL:        PermissionDeleteURL,
		BulkDeleteURL:    PermissionBulkDeleteURL,
		SetStatusURL:     PermissionSetStatusURL,
		BulkSetStatusURL: PermissionBulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r PermissionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"permission.list":           r.ListURL,
		"permission.table":          r.TableURL,
		"permission.add":            r.AddURL,
		"permission.edit":           r.EditURL,
		"permission.delete":         r.DeleteURL,
		"permission.bulk_delete":    r.BulkDeleteURL,
		"permission.set_status":     r.SetStatusURL,
		"permission.bulk_set_status": r.BulkSetStatusURL,
	}
}

// ---------------------------------------------------------------------------
// WorkspaceRoutes
// ---------------------------------------------------------------------------

// WorkspaceRoutes holds all route paths for workspace management.
type WorkspaceRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultWorkspaceRoutes returns a WorkspaceRoutes populated from the
// package-level route constants.
func DefaultWorkspaceRoutes() WorkspaceRoutes {
	return WorkspaceRoutes{
		ListURL:          WorkspaceListURL,
		TableURL:         WorkspaceTableURL,
		AddURL:           WorkspaceAddURL,
		EditURL:          WorkspaceEditURL,
		DeleteURL:        WorkspaceDeleteURL,
		BulkDeleteURL:    WorkspaceBulkDeleteURL,
		SetStatusURL:     WorkspaceSetStatusURL,
		BulkSetStatusURL: WorkspaceBulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r WorkspaceRoutes) RouteMap() map[string]string {
	return map[string]string{
		"workspace.list":           r.ListURL,
		"workspace.table":          r.TableURL,
		"workspace.add":            r.AddURL,
		"workspace.edit":           r.EditURL,
		"workspace.delete":         r.DeleteURL,
		"workspace.bulk_delete":    r.BulkDeleteURL,
		"workspace.set_status":     r.SetStatusURL,
		"workspace.bulk_set_status": r.BulkSetStatusURL,
	}
}

// ---------------------------------------------------------------------------
// SupplierRoutes
// ---------------------------------------------------------------------------

// SupplierRoutes holds all route paths for supplier management.
type SupplierRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Report routes
	PayablesAgingURL string `json:"payables_aging_url"`
}

// DefaultSupplierRoutes returns a SupplierRoutes populated from the
// package-level route constants.
func DefaultSupplierRoutes() SupplierRoutes {
	return SupplierRoutes{
		ListURL:          SupplierListURL,
		TableURL:         SupplierTableURL,
		AddURL:           SupplierAddURL,
		EditURL:          SupplierEditURL,
		DeleteURL:        SupplierDeleteURL,
		BulkDeleteURL:    SupplierBulkDeleteURL,
		DetailURL:        SupplierDetailURL,
		TabActionURL:     SupplierTabActionURL,
		SetStatusURL:     SupplierSetStatusURL,
		BulkSetStatusURL: SupplierBulkSetStatusURL,

		AttachmentUploadURL: SupplierAttachmentUploadURL,
		AttachmentDeleteURL: SupplierAttachmentDeleteURL,

		PayablesAgingURL: PayablesAgingURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r SupplierRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier.list":            r.ListURL,
		"supplier.table":           r.TableURL,
		"supplier.add":             r.AddURL,
		"supplier.edit":            r.EditURL,
		"supplier.delete":          r.DeleteURL,
		"supplier.bulk_delete":     r.BulkDeleteURL,
		"supplier.detail":          r.DetailURL,
		"supplier.tab_action":      r.TabActionURL,
		"supplier.set_status":      r.SetStatusURL,
		"supplier.bulk_set_status": r.BulkSetStatusURL,

		"supplier.attachment.upload": r.AttachmentUploadURL,
		"supplier.attachment.delete": r.AttachmentDeleteURL,

		"supplier.payables_aging": r.PayablesAgingURL,
	}
}

// ---------------------------------------------------------------------------
// ClientTagRoutes
// ---------------------------------------------------------------------------

// ClientTagRoutes holds all route paths for client tag (category) management.
type ClientTagRoutes struct {
	ListURL       string `json:"list_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
}

// DefaultClientTagRoutes returns a ClientTagRoutes populated from the
// package-level route constants.
func DefaultClientTagRoutes() ClientTagRoutes {
	return ClientTagRoutes{
		ListURL:       ClientTagListURL,
		AddURL:        ClientTagAddURL,
		EditURL:       ClientTagEditURL,
		DeleteURL:     ClientTagDeleteURL,
		BulkDeleteURL: ClientTagBulkDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r ClientTagRoutes) RouteMap() map[string]string {
	return map[string]string{
		"client_tag.list":        r.ListURL,
		"client_tag.add":         r.AddURL,
		"client_tag.edit":        r.EditURL,
		"client_tag.delete":      r.DeleteURL,
		"client_tag.bulk_delete": r.BulkDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// LoginRoutes
// ---------------------------------------------------------------------------

// LoginRoutes holds all route paths for login/authentication views.
type LoginRoutes struct {
	LoginURL     string `json:"login_url"`
	LoginPostURL string `json:"login_post_url"`
}

// DefaultLoginRoutes returns a LoginRoutes populated from the package-level
// route constants.
func DefaultLoginRoutes() LoginRoutes {
	return LoginRoutes{
		LoginURL:     LoginURL,
		LoginPostURL: LoginPostURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r LoginRoutes) RouteMap() map[string]string {
	return map[string]string{
		"login.page": r.LoginURL,
		"login.post": r.LoginPostURL,
	}
}
