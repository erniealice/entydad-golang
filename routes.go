package entydad

// Default route constants for entydad views.
// Consumer apps can use these or define their own.
const (
	ClientListURL       = "/app/clients/list/{status}"
	ClientAddURL        = "/action/clients/add"
	ClientEditURL       = "/action/clients/edit/{id}"
	ClientDeleteURL     = "/action/clients/delete"
	ClientBulkDeleteURL = "/action/clients/bulk-delete"

	UserListURL          = "/app/users/list/{status}"
	UserTableURL         = "/action/users/table/{status}"
	UserAddURL           = "/action/users/add"
	UserEditURL          = "/action/users/edit/{id}"
	UserDeleteURL        = "/action/users/delete"
	UserBulkDeleteURL    = "/action/users/bulk-delete"
	UserSetStatusURL     = "/action/users/set-status"
	UserBulkSetStatusURL = "/action/users/bulk-set-status"

	LocationListURL          = "/app/locations/list/{status}"
	LocationTableURL         = "/action/locations/table/{status}"
	LocationAddURL           = "/action/locations/add"
	LocationEditURL          = "/action/locations/edit/{id}"
	LocationDeleteURL        = "/action/locations/delete"
	LocationBulkDeleteURL    = "/action/locations/bulk-delete"
	LocationSetStatusURL     = "/action/locations/set-status"
	LocationBulkSetStatusURL = "/action/locations/bulk-set-status"

	RoleListURL          = "/app/roles/list/{status}"
	RoleTableURL         = "/action/roles/table/{status}"
	RoleAddURL           = "/action/roles/add"
	RoleEditURL          = "/action/roles/edit/{id}"
	RoleDeleteURL        = "/action/roles/delete"
	RoleBulkDeleteURL    = "/action/roles/bulk-delete"
	RoleSetStatusURL     = "/action/roles/set-status"
	RoleBulkSetStatusURL = "/action/roles/bulk-set-status"
)
