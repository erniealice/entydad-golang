package entydad

// Default route constants for entydad views.
// Consumer apps can use these or define their own.
const (
	ClientDashboardURL = "/app/clients/dashboard"
	ClientListURL      = "/app/clients/list/{status}"
	ClientTableURL     = "/action/clients/table/{status}"
	ClientAddURL        = "/action/clients/add"
	ClientEditURL       = "/action/clients/edit/{id}"
	ClientDeleteURL        = "/action/clients/delete"
	ClientBulkDeleteURL    = "/action/clients/bulk-delete"
	ClientDetailURL     = "/app/clients/detail/{id}"
	ClientTabActionURL          = "/action/clients/{id}/tab/{tab}"
	ClientAttachmentUploadURL   = "/action/clients/{id}/attachments/upload"
	ClientAttachmentDeleteURL   = "/action/clients/{id}/attachments/delete"
	ClientSetStatusURL     = "/action/clients/set-status"
	ClientBulkSetStatusURL = "/action/clients/bulk-set-status"
	ClientSearchURL        = "/action/clients/search"

	UserDashboardURL     = "/app/users/dashboard"
	UserListURL          = "/app/users/list/{status}"
	UserTableURL         = "/action/users/table/{status}"
	UserAddURL           = "/action/users/add"
	UserEditURL          = "/action/users/edit/{id}"
	UserDeleteURL        = "/action/users/delete"
	UserBulkDeleteURL    = "/action/users/bulk-delete"
	UserSetStatusURL     = "/action/users/set-status"
	UserBulkSetStatusURL = "/action/users/bulk-set-status"

	LocationDetailURL        = "/app/locations/{id}"
	LocationListURL          = "/app/locations/list/{status}"
	LocationTableURL         = "/action/locations/table/{status}"
	LocationAddURL           = "/action/locations/add"
	LocationEditURL          = "/action/locations/edit/{id}"
	LocationDeleteURL        = "/action/locations/delete"
	LocationBulkDeleteURL    = "/action/locations/bulk-delete"
	LocationSetStatusURL     = "/action/locations/set-status"
	LocationBulkSetStatusURL = "/action/locations/bulk-set-status"
	LocationTabActionURL         = "/action/locations/{id}/tab/{tab}"
	LocationAttachmentUploadURL  = "/action/locations/{id}/attachments/upload"
	LocationAttachmentDeleteURL  = "/action/locations/{id}/attachments/delete"
	LocationEditDetailURL    = "/action/locations/edit-detail/{id}"

	UserDetailURL         = "/app/users/detail/{id}"
	UserTabActionURL          = "/action/users/{id}/tab/{tab}"
	UserAttachmentUploadURL   = "/action/users/{id}/attachments/upload"
	UserAttachmentDeleteURL   = "/action/users/{id}/attachments/delete"
	UserResetPasswordURL  = "/action/users/reset-password/{id}"

	RoleDetailURL        = "/app/roles/detail/{id}"
	RoleTabActionURL         = "/action/roles/{id}/tab/{tab}"
	RoleAttachmentUploadURL  = "/action/roles/{id}/attachments/upload"
	RoleAttachmentDeleteURL  = "/action/roles/{id}/attachments/delete"
	RoleListURL          = "/app/roles/list"
	RoleTableURL         = "/action/roles/table"
	RoleAddURL           = "/action/roles/add"
	RoleEditURL          = "/action/roles/edit/{id}"
	RoleDeleteURL        = "/action/roles/delete"
	RoleBulkDeleteURL    = "/action/roles/bulk-delete"
	RoleSetStatusURL     = "/action/roles/set-status"
	RoleBulkSetStatusURL = "/action/roles/bulk-set-status"

	PermissionListURL          = "/app/permissions/list/{status}"
	PermissionTableURL         = "/action/permissions/table/{status}"
	PermissionAddURL           = "/action/permissions/add"
	PermissionEditURL          = "/action/permissions/edit/{id}"
	PermissionDeleteURL        = "/action/permissions/delete"
	PermissionBulkDeleteURL    = "/action/permissions/bulk-delete"
	PermissionSetStatusURL     = "/action/permissions/set-status"
	PermissionBulkSetStatusURL = "/action/permissions/bulk-set-status"

	RolePermissionsURL       = "/app/manage/roles/{id}/permissions"
	RolePermissionsTableURL  = "/action/manage/roles/{id}/permissions/table"
	RolePermissionsAssignURL = "/action/manage/roles/{id}/permissions/assign"
	RolePermissionsRemoveURL = "/action/manage/roles/{id}/permissions/remove"

	UserRolesURL       = "/app/manage/users/{id}/roles"
	UserRolesTableURL  = "/action/manage/users/{id}/roles/table"
	UserRolesAssignURL = "/action/manage/users/{id}/roles/assign"
	UserRolesRemoveURL = "/action/manage/users/{id}/roles/remove"

	RoleUsersURL       = "/app/roles/detail/{id}/users"
	RoleUsersTableURL  = "/action/roles/detail/{id}/users/table"
	RoleUsersAssignURL = "/action/roles/detail/{id}/users/assign"
	RoleUsersRemoveURL = "/action/roles/detail/{id}/users/remove"

	// Migrated route constants: /detail/ pattern for user-roles and role-permissions
	// Old /manage/ constants kept above for backward compatibility
	UserDetailRolesURL       = "/app/users/detail/{id}/roles"
	UserDetailRolesTableURL  = "/action/users/detail/{id}/roles/table"
	UserDetailRolesAssignURL = "/action/users/detail/{id}/roles/assign"
	UserDetailRolesRemoveURL = "/action/users/detail/{id}/roles/remove"

	RoleDetailPermissionsURL       = "/app/roles/detail/{id}/permissions"
	RoleDetailPermissionsTableURL  = "/action/roles/detail/{id}/permissions/table"
	RoleDetailPermissionsAssignURL = "/action/roles/detail/{id}/permissions/assign"
	RoleDetailPermissionsRemoveURL = "/action/roles/detail/{id}/permissions/remove"

	WorkspaceListURL          = "/app/workspaces/list/{status}"
	WorkspaceTableURL         = "/action/workspaces/table/{status}"
	WorkspaceAddURL           = "/action/workspaces/add"
	WorkspaceEditURL          = "/action/workspaces/edit/{id}"
	WorkspaceDeleteURL        = "/action/workspaces/delete"
	WorkspaceBulkDeleteURL    = "/action/workspaces/bulk-delete"
	WorkspaceSetStatusURL     = "/action/workspaces/set-status"
	WorkspaceBulkSetStatusURL = "/action/workspaces/bulk-set-status"

	// Client report routes
	ReceivablesAgingURL = "/app/clients/reports/receivables-aging"

	// Supplier routes
	SupplierListURL          = "/app/suppliers/list/{status}"
	SupplierTableURL         = "/action/suppliers/table/{status}"
	SupplierAddURL           = "/action/suppliers/add"
	SupplierEditURL          = "/action/suppliers/edit/{id}"
	SupplierDeleteURL        = "/action/suppliers/delete"
	SupplierBulkDeleteURL    = "/action/suppliers/bulk-delete"
	SupplierDetailURL        = "/app/suppliers/detail/{id}"
	SupplierTabActionURL         = "/action/suppliers/{id}/tab/{tab}"
	SupplierAttachmentUploadURL  = "/action/suppliers/{id}/attachments/upload"
	SupplierAttachmentDeleteURL  = "/action/suppliers/{id}/attachments/delete"
	SupplierSetStatusURL     = "/action/suppliers/set-status"
	SupplierBulkSetStatusURL = "/action/suppliers/bulk-set-status"

	// Supplier report routes
	PayablesAgingURL = "/app/suppliers/reports/payables-aging"

	// Client Tag (Category) routes
	ClientTagListURL       = "/app/clients/settings/tags/list"
	ClientTagAddURL        = "/action/clients/tags/add"
	ClientTagEditURL       = "/action/clients/tags/edit/{id}"
	ClientTagDeleteURL     = "/action/clients/tags/delete"
	ClientTagBulkDeleteURL = "/action/clients/tags/bulk-delete"

	// Login routes
	LoginURL     = "/login"
	LoginPostURL = "/login"

	// DefaultAppRedirectURL is the default post-login redirect path.
	// Consumer apps should set Deps.RedirectURL to override this.
	DefaultAppRedirectURL = "/app/"
)
