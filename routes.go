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
)
