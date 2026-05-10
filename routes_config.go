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

	// Statement export
	StatementExportURL string `json:"statement_export_url"`

	// Report routes
	ReceivablesAgingURL string `json:"receivables_aging_url"`

	// Settings routes
	PaymentTermsURL  string `json:"payment_terms_url"`
	ClientTagListURL string `json:"client_tag_list_url"` // cross-app link to client-tag list (dashboard quick-action)

	// RevenueRunURL is the per-client "Run Invoices" drawer endpoint.
	RevenueRunURL string `json:"revenue_run_url"`
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

		StatementExportURL: ClientStatementExportURL,

		ReceivablesAgingURL: ReceivablesAgingURL,

		PaymentTermsURL:  PaymentTermListURL,
		ClientTagListURL: ClientTagListURL,
		RevenueRunURL:    ClientRevenueRunURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r ClientRoutes) RouteMap() map[string]string {
	return map[string]string{
		"client.dashboard":       r.DashboardURL,
		"client.list":            r.ListURL,
		"client.table":           r.TableURL,
		"client.add":             r.AddURL,
		"client.edit":            r.EditURL,
		"client.delete":          r.DeleteURL,
		"client.bulk_delete":     r.BulkDeleteURL,
		"client.detail":          r.DetailURL,
		"client.tab_action":      r.TabActionURL,
		"client.set_status":      r.SetStatusURL,
		"client.bulk_set_status": r.BulkSetStatusURL,
		"client.search":          r.SearchURL,

		"client.attachment.upload": r.AttachmentUploadURL,
		"client.attachment.delete": r.AttachmentDeleteURL,

		"client.statement_export": r.StatementExportURL,

		"client.receivables_aging": r.ReceivablesAgingURL,

		"client.payment_terms":  r.PaymentTermsURL,
		"client.revenue_run":    r.RevenueRunURL,
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
		DetailURL:        UserDetailURL,
		TabActionURL:     UserTabActionURL,
		ResetPasswordURL: UserResetPasswordURL,

		SearchTimezonesURL: UserSearchTimezonesURL,

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

		// Cross-app quick-action links for the user dashboard.
		RoleListURL:       RoleListURL,
		PermissionListURL: "/app/permissions/list/active", // PermissionListURL with {status}=active
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r UserRoutes) RouteMap() map[string]string {
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
	UsersSearchURL string `json:"users_search_url"`

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
		UsersSearchURL: RoleUsersSearchURL,

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

// ---------------------------------------------------------------------------
// LocationRoutes
// ---------------------------------------------------------------------------

// LocationRoutes holds all route paths for location management.
type LocationRoutes struct {
	DashboardURL     string `json:"dashboard_url"`
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
		DashboardURL:     LocationDashboardURL,
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
		"location.dashboard":       r.DashboardURL,
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
// LocationAreaRoutes
// ---------------------------------------------------------------------------

// LocationAreaRoutes holds all route paths for location area management.
type LocationAreaRoutes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultLocationAreaRoutes returns a LocationAreaRoutes populated from the
// package-level route constants.
func DefaultLocationAreaRoutes() LocationAreaRoutes {
	return LocationAreaRoutes{
		DashboardURL:     LocationAreaDashboardURL,
		ListURL:          LocationAreaListURL,
		TableURL:         LocationAreaTableURL,
		DetailURL:        LocationAreaDetailURL,
		AddURL:           LocationAreaAddURL,
		EditURL:          LocationAreaEditURL,
		DeleteURL:        LocationAreaDeleteURL,
		BulkDeleteURL:    LocationAreaBulkDeleteURL,
		SetStatusURL:     LocationAreaSetStatusURL,
		BulkSetStatusURL: LocationAreaBulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r LocationAreaRoutes) RouteMap() map[string]string {
	return map[string]string{
		"location_area.dashboard":       r.DashboardURL,
		"location_area.list":            r.ListURL,
		"location_area.table":           r.TableURL,
		"location_area.detail":          r.DetailURL,
		"location_area.add":             r.AddURL,
		"location_area.edit":            r.EditURL,
		"location_area.delete":          r.DeleteURL,
		"location_area.bulk_delete":     r.BulkDeleteURL,
		"location_area.set_status":      r.SetStatusURL,
		"location_area.bulk_set_status": r.BulkSetStatusURL,
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
		"permission.list":            r.ListURL,
		"permission.table":           r.TableURL,
		"permission.add":             r.AddURL,
		"permission.edit":            r.EditURL,
		"permission.delete":          r.DeleteURL,
		"permission.bulk_delete":     r.BulkDeleteURL,
		"permission.set_status":      r.SetStatusURL,
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
	SwitchURL        string `json:"switch_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
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
		SwitchURL:        WorkspaceSwitchURL,
		DetailURL:        WorkspaceDetailURL,
		TabActionURL:     WorkspaceTabActionURL,

		AttachmentUploadURL: WorkspaceAttachmentUploadURL,
		AttachmentDeleteURL: WorkspaceAttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r WorkspaceRoutes) RouteMap() map[string]string {
	return map[string]string{
		"workspace.list":            r.ListURL,
		"workspace.table":           r.TableURL,
		"workspace.add":             r.AddURL,
		"workspace.edit":            r.EditURL,
		"workspace.delete":          r.DeleteURL,
		"workspace.bulk_delete":     r.BulkDeleteURL,
		"workspace.set_status":      r.SetStatusURL,
		"workspace.bulk_set_status": r.BulkSetStatusURL,
		"workspace.switch_url":      r.SwitchURL,
		"workspace.detail":          r.DetailURL,
		"workspace.tab_action":      r.TabActionURL,

		"workspace.attachment.upload": r.AttachmentUploadURL,
		"workspace.attachment.delete": r.AttachmentDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// WorkspaceUserRoutes
// ---------------------------------------------------------------------------

// WorkspaceUserRoutes holds all route paths for workspace-user nested detail management.
type WorkspaceUserRoutes struct {
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	AddURL       string `json:"add_url"`
	DeleteURL    string `json:"delete_url"`
	SetStatusURL string `json:"set_status_url"`
	SearchURL    string `json:"search_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultWorkspaceUserRoutes returns a WorkspaceUserRoutes populated from the
// package-level route constants.
func DefaultWorkspaceUserRoutes() WorkspaceUserRoutes {
	return WorkspaceUserRoutes{
		ListURL:      WorkspaceUserListURL,
		DetailURL:    WorkspaceUserDetailURL,
		TabActionURL: WorkspaceUserTabActionURL,
		AddURL:       WorkspaceUserAddURL,
		DeleteURL:    WorkspaceUserDeleteURL,
		SetStatusURL: WorkspaceUserSetStatusURL,
		SearchURL:    WorkspaceUserSearchURL,

		AttachmentUploadURL: WorkspaceUserAttachmentUploadURL,
		AttachmentDeleteURL: WorkspaceUserAttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r WorkspaceUserRoutes) RouteMap() map[string]string {
	return map[string]string{
		"workspace_user.list":       r.ListURL,
		"workspace_user.detail":     r.DetailURL,
		"workspace_user.tab_action": r.TabActionURL,
		"workspace_user.add":        r.AddURL,
		"workspace_user.delete":     r.DeleteURL,
		"workspace_user.set_status": r.SetStatusURL,
		"workspace_user.search":     r.SearchURL,

		"workspace_user.attachment.upload": r.AttachmentUploadURL,
		"workspace_user.attachment.delete": r.AttachmentDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// WorkspaceUserRoleRoutes
// ---------------------------------------------------------------------------

// WorkspaceUserRoleRoutes holds all route paths for the workspace_user_role
// assignment drawer (Phase 3).
type WorkspaceUserRoleRoutes struct {
	AddURL          string `json:"add_url"`
	DeleteURL       string `json:"delete_url"`
	PermissionsURL  string `json:"permissions_url"`
	SearchRolesURL  string `json:"search_roles_url"`
}

// DefaultWorkspaceUserRoleRoutes returns a WorkspaceUserRoleRoutes populated from
// the package-level route constants.
func DefaultWorkspaceUserRoleRoutes() WorkspaceUserRoleRoutes {
	return WorkspaceUserRoleRoutes{
		AddURL:         WorkspaceUserRoleAddURL,
		DeleteURL:      WorkspaceUserRoleDeleteURL,
		PermissionsURL: WorkspaceUserRolePermissionsURL,
		SearchRolesURL: WorkspaceUserRoleSearchRolesURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r WorkspaceUserRoleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"workspace_user_role.add":          r.AddURL,
		"workspace_user_role.delete":       r.DeleteURL,
		"workspace_user_role.permissions":  r.PermissionsURL,
		"workspace_user_role.search_roles": r.SearchRolesURL,
	}
}

// ---------------------------------------------------------------------------
// AdminDashboardRoutes
// ---------------------------------------------------------------------------

// AdminDashboardRoutes holds the route path for the admin app dashboard.
//
// Admin is a *composite* app: its sidebar block spans permission, role,
// workspace, workspace_user, and workspace_user_role entities. The dashboard
// sits at the app level (not the entity level), so it gets its own thin
// route struct rather than being attached to PermissionRoutes / RoleRoutes /
// WorkspaceRoutes — none of which would be the obvious owner.
type AdminDashboardRoutes struct {
	DashboardURL string `json:"dashboard_url"`
}

// DefaultAdminDashboardRoutes returns an AdminDashboardRoutes populated from
// the package-level route constants.
func DefaultAdminDashboardRoutes() AdminDashboardRoutes {
	return AdminDashboardRoutes{
		DashboardURL: AdminDashboardURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r AdminDashboardRoutes) RouteMap() map[string]string {
	return map[string]string{
		"admin.dashboard": r.DashboardURL,
	}
}

// ---------------------------------------------------------------------------
// SupplierRoutes
// ---------------------------------------------------------------------------

// SupplierRoutes holds all route paths for supplier management.
type SupplierRoutes struct {
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

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Statement export
	StatementExportURL string `json:"statement_export_url"`

	// Report routes
	PayablesAgingURL string `json:"payables_aging_url"`

	// Settings routes
	PaymentTermsURL string `json:"payment_terms_url"`
}

// DefaultSupplierRoutes returns a SupplierRoutes populated from the
// package-level route constants.
func DefaultSupplierRoutes() SupplierRoutes {
	return SupplierRoutes{
		DashboardURL:     SupplierDashboardURL,
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

		StatementExportURL: SupplierStatementExportURL,

		PayablesAgingURL: PayablesAgingURL,

		PaymentTermsURL: SupplierPaymentTermListURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r SupplierRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier.dashboard":       r.DashboardURL,
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

		"supplier.statement_export": r.StatementExportURL,

		"supplier.payables_aging": r.PayablesAgingURL,

		"supplier.payment_terms": r.PaymentTermsURL,
	}
}

// ---------------------------------------------------------------------------
// PaymentTermRoutes
// ---------------------------------------------------------------------------

// PaymentTermRoutes holds all route paths for payment term management.
type PaymentTermRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultPaymentTermRoutes returns a PaymentTermRoutes populated from the
// package-level route constants.
func DefaultPaymentTermRoutes() PaymentTermRoutes {
	return PaymentTermRoutes{
		ListURL:          PaymentTermListURL,
		TableURL:         PaymentTermTableURL,
		AddURL:           PaymentTermAddURL,
		EditURL:          PaymentTermEditURL,
		DeleteURL:        PaymentTermDeleteURL,
		BulkDeleteURL:    PaymentTermBulkDeleteURL,
		SetStatusURL:     PaymentTermSetStatusURL,
		BulkSetStatusURL: PaymentTermBulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r PaymentTermRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payment_term.list":             r.ListURL,
		"payment_term.table":            r.TableURL,
		"payment_term.add":              r.AddURL,
		"payment_term.edit":             r.EditURL,
		"payment_term.delete":           r.DeleteURL,
		"payment_term.bulk_delete":      r.BulkDeleteURL,
		"payment_term.set_status":       r.SetStatusURL,
		"payment_term.bulk_set_status":  r.BulkSetStatusURL,
	}
}

// ---------------------------------------------------------------------------
// SupplierPaymentTermRoutes
// ---------------------------------------------------------------------------

// SupplierPaymentTermRoutes holds payment term route paths for the supplier context.
// These routes show only payment terms with entity_scope IN ('supplier', 'both').
type SupplierPaymentTermRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultSupplierPaymentTermRoutes returns a SupplierPaymentTermRoutes from package-level constants.
func DefaultSupplierPaymentTermRoutes() SupplierPaymentTermRoutes {
	return SupplierPaymentTermRoutes{
		ListURL:          SupplierPaymentTermListURL,
		TableURL:         SupplierPaymentTermTableURL,
		AddURL:           SupplierPaymentTermAddURL,
		EditURL:          SupplierPaymentTermEditURL,
		DeleteURL:        SupplierPaymentTermDeleteURL,
		BulkDeleteURL:    SupplierPaymentTermBulkDeleteURL,
		SetStatusURL:     SupplierPaymentTermSetStatusURL,
		BulkSetStatusURL: SupplierPaymentTermBulkSetStatusURL,
	}
}

// ToPaymentTermRoutes converts SupplierPaymentTermRoutes to a PaymentTermRoutes,
// allowing the payment term module to be reused with supplier-context paths.
func (r SupplierPaymentTermRoutes) ToPaymentTermRoutes() PaymentTermRoutes {
	return PaymentTermRoutes{
		ListURL:          r.ListURL,
		TableURL:         r.TableURL,
		AddURL:           r.AddURL,
		EditURL:          r.EditURL,
		DeleteURL:        r.DeleteURL,
		BulkDeleteURL:    r.BulkDeleteURL,
		SetStatusURL:     r.SetStatusURL,
		BulkSetStatusURL: r.BulkSetStatusURL,
	}
}

// ---------------------------------------------------------------------------
// ClientTagRoutes
// ---------------------------------------------------------------------------

// ClientTagRoutes holds all route paths for client tag (category) management.
type ClientTagRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultClientTagRoutes returns a ClientTagRoutes populated from the
// package-level route constants.
func DefaultClientTagRoutes() ClientTagRoutes {
	return ClientTagRoutes{
		ListURL:          ClientTagListURL,
		TableURL:         ClientTagTableURL,
		AddURL:           ClientTagAddURL,
		EditURL:          ClientTagEditURL,
		DeleteURL:        ClientTagDeleteURL,
		BulkDeleteURL:    ClientTagBulkDeleteURL,
		SetStatusURL:     ClientTagSetStatusURL,
		BulkSetStatusURL: ClientTagBulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r ClientTagRoutes) RouteMap() map[string]string {
	return map[string]string{
		"client_tag.list":             r.ListURL,
		"client_tag.table":            r.TableURL,
		"client_tag.add":              r.AddURL,
		"client_tag.edit":             r.EditURL,
		"client_tag.delete":           r.DeleteURL,
		"client_tag.bulk_delete":      r.BulkDeleteURL,
		"client_tag.set_status":       r.SetStatusURL,
		"client_tag.bulk_set_status":  r.BulkSetStatusURL,
	}
}

// SupplierTagRoutes holds all route paths for supplier tag (category) management.
type SupplierTagRoutes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

func DefaultSupplierTagRoutes() SupplierTagRoutes {
	return SupplierTagRoutes{
		ListURL:          SupplierTagListURL,
		TableURL:         SupplierTagTableURL,
		AddURL:           SupplierTagAddURL,
		EditURL:          SupplierTagEditURL,
		DeleteURL:        SupplierTagDeleteURL,
		BulkDeleteURL:    SupplierTagBulkDeleteURL,
		SetStatusURL:     SupplierTagSetStatusURL,
		BulkSetStatusURL: SupplierTagBulkSetStatusURL,
	}
}

func (r SupplierTagRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_tag.list":            r.ListURL,
		"supplier_tag.table":           r.TableURL,
		"supplier_tag.add":             r.AddURL,
		"supplier_tag.edit":            r.EditURL,
		"supplier_tag.delete":          r.DeleteURL,
		"supplier_tag.bulk_delete":     r.BulkDeleteURL,
		"supplier_tag.set_status":      r.SetStatusURL,
		"supplier_tag.bulk_set_status": r.BulkSetStatusURL,
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

// ---------------------------------------------------------------------------
// AuthRoutes
// ---------------------------------------------------------------------------

// AuthRoutes holds all route paths for authentication views (signup, reset, logout).
type AuthRoutes struct {
	LoginURL             string `json:"login_url"`
	LoginPostURL         string `json:"login_post_url"`
	SignupURL            string `json:"signup_url"`
	SignupPostURL        string `json:"signup_post_url"`
	ResetPasswordURL     string `json:"reset_password_url"`
	ResetPasswordPostURL string `json:"reset_password_post_url"`
	ResetConfirmURL      string `json:"reset_confirm_url"`
	ResetConfirmPostURL  string `json:"reset_confirm_post_url"`
	LogoutURL            string `json:"logout_url"`
}

// DefaultAuthRoutes returns an AuthRoutes populated from the package-level
// URL constants defined in routes.go.
func DefaultAuthRoutes() AuthRoutes {
	return AuthRoutes{
		LoginURL:             AuthLoginURL,
		LoginPostURL:         AuthLoginPostURL,
		SignupURL:            AuthSignupURL,
		SignupPostURL:        AuthSignupPostURL,
		ResetPasswordURL:     AuthResetPasswordURL,
		ResetPasswordPostURL: AuthResetPasswordPostURL,
		ResetConfirmURL:      AuthResetConfirmURL,
		ResetConfirmPostURL:  AuthResetConfirmPostURL,
		LogoutURL:            AuthLogoutURL,
	}
}

// RouteMap returns a map of route keys to URL paths for AuthRoutes.
func (r AuthRoutes) RouteMap() map[string]string {
	return map[string]string{
		"auth.login.page":          r.LoginURL,
		"auth.login.post":          r.LoginPostURL,
		"auth.signup.page":         r.SignupURL,
		"auth.signup.post":         r.SignupPostURL,
		"auth.reset-password.page": r.ResetPasswordURL,
		"auth.reset-password.post": r.ResetPasswordPostURL,
		"auth.reset-confirm.page":  r.ResetConfirmURL,
		"auth.reset-confirm.post":  r.ResetConfirmPostURL,
		"auth.logout":              r.LogoutURL,
	}
}

// ---------------------------------------------------------------------------
// TaxRegistrationRoutes
// ---------------------------------------------------------------------------

// TaxRegistrationRoutes holds route paths for the polymorphic TaxRegistration
// views. v1 surfaces client + workspace party types only.
// The AddURL and DeleteURL are party-scoped (include party_id in path).
type TaxRegistrationRoutes struct {
	// Client-scoped routes
	ClientListURL   string `json:"client_list_url"`
	ClientAddURL    string `json:"client_add_url"`
	ClientEditURL   string `json:"client_edit_url"`
	ClientDeleteURL string `json:"client_delete_url"`

	// Workspace-scoped routes
	WorkspaceListURL   string `json:"workspace_list_url"`
	WorkspaceAddURL    string `json:"workspace_add_url"`
	WorkspaceEditURL   string `json:"workspace_edit_url"`
	WorkspaceDeleteURL string `json:"workspace_delete_url"`

	// AddURL and DeleteURL are the active-context URLs (set by the view wiring).
	// For client views: same as ClientAddURL / ClientDeleteURL.
	// For workspace views: same as WorkspaceAddURL / WorkspaceDeleteURL.
	AddURL    string `json:"add_url"`
	DeleteURL string `json:"delete_url"`
}

// DefaultTaxRegistrationRoutes returns a TaxRegistrationRoutes populated from
// the package-level route constants.
func DefaultTaxRegistrationRoutes() TaxRegistrationRoutes {
	return TaxRegistrationRoutes{
		ClientListURL:   ClientTaxRegistrationListURL,
		ClientAddURL:    ClientTaxRegistrationAddURL,
		ClientEditURL:   ClientTaxRegistrationEditURL,
		ClientDeleteURL: ClientTaxRegistrationDeleteURL,

		WorkspaceListURL:   WorkspaceTaxRegistrationListURL,
		WorkspaceAddURL:    WorkspaceTaxRegistrationAddURL,
		WorkspaceEditURL:   WorkspaceTaxRegistrationEditURL,
		WorkspaceDeleteURL: WorkspaceTaxRegistrationDeleteURL,

		// Default to client context; override at wiring time for workspace.
		AddURL:    ClientTaxRegistrationAddURL,
		DeleteURL: ClientTaxRegistrationDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r TaxRegistrationRoutes) RouteMap() map[string]string {
	return map[string]string{
		"tax_registration.client.list":      r.ClientListURL,
		"tax_registration.client.add":       r.ClientAddURL,
		"tax_registration.client.edit":      r.ClientEditURL,
		"tax_registration.client.delete":    r.ClientDeleteURL,
		"tax_registration.workspace.list":   r.WorkspaceListURL,
		"tax_registration.workspace.add":    r.WorkspaceAddURL,
		"tax_registration.workspace.edit":   r.WorkspaceEditURL,
		"tax_registration.workspace.delete": r.WorkspaceDeleteURL,
	}
}
