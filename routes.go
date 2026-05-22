package entydad

// Default route constants for entydad views.
// Consumer apps can use these or define their own.
const (
	ClientDashboardURL        = "/clients/dashboard"
	ClientListURL             = "/clients/list/{status}"
	ClientTableURL            = "/action/client/table/{status}"
	ClientAddURL              = "/action/client/add"
	ClientEditURL             = "/action/client/edit/{id}"
	ClientDeleteURL           = "/action/client/delete"
	ClientBulkDeleteURL       = "/action/client/bulk-delete"
	ClientDetailURL           = "/clients/detail/{id}"
	ClientTabActionURL        = "/action/client/{id}/tab/{tab}"
	ClientAttachmentUploadURL = "/action/client/{id}/attachments/upload"
	ClientAttachmentDeleteURL = "/action/client/{id}/attachments/delete"
	ClientSetStatusURL        = "/action/client/set-status"
	ClientBulkSetStatusURL    = "/action/client/bulk-set-status"
	ClientSearchURL           = "/action/client/search"

	UserDashboardURL       = "/users/dashboard"
	UserListURL            = "/users/list/{status}"
	UserTableURL           = "/action/user/table/{status}"
	UserAddURL             = "/action/user/add"
	UserEditURL            = "/action/user/edit/{id}"
	UserDeleteURL          = "/action/user/delete"
	UserBulkDeleteURL      = "/action/user/bulk-delete"
	UserSetStatusURL       = "/action/user/set-status"
	UserBulkSetStatusURL   = "/action/user/bulk-set-status"
	UserSearchTimezonesURL = "/action/user/search-timezones"

	LocationDashboardURL        = "/locations/dashboard"
	LocationDetailURL           = "/locations/detail/{id}"
	LocationListURL             = "/locations/list/{status}"
	LocationTableURL            = "/action/location/table/{status}"
	LocationAddURL              = "/action/location/add"
	LocationEditURL             = "/action/location/edit/{id}"
	LocationDeleteURL           = "/action/location/delete"
	LocationBulkDeleteURL       = "/action/location/bulk-delete"
	LocationSetStatusURL        = "/action/location/set-status"
	LocationBulkSetStatusURL    = "/action/location/bulk-set-status"
	LocationTabActionURL        = "/action/location/{id}/tab/{tab}"
	LocationAttachmentUploadURL = "/action/location/{id}/attachments/upload"
	LocationAttachmentDeleteURL = "/action/location/{id}/attachments/delete"
	LocationEditDetailURL       = "/action/location/edit-detail/{id}"

	LocationAreaDashboardURL     = "/location-areas/dashboard"
	LocationAreaListURL          = "/location-areas/list/{status}"
	LocationAreaTableURL         = "/action/location-area/table/{status}"
	LocationAreaDetailURL        = "/location-areas/detail/{id}"
	LocationAreaAddURL           = "/action/location-area/add"
	LocationAreaEditURL          = "/action/location-area/edit/{id}"
	LocationAreaDeleteURL        = "/action/location-area/delete"
	LocationAreaBulkDeleteURL    = "/action/location-area/bulk-delete"
	LocationAreaSetStatusURL     = "/action/location-area/set-status"
	LocationAreaBulkSetStatusURL = "/action/location-area/bulk-set-status"

	UserDetailURL           = "/users/detail/{id}"
	UserTabActionURL        = "/action/user/{id}/tab/{tab}"
	UserAttachmentUploadURL = "/action/user/{id}/attachments/upload"
	UserAttachmentDeleteURL = "/action/user/{id}/attachments/delete"
	UserResetPasswordURL    = "/action/user/reset-password/{id}"

	RoleDetailURL           = "/roles/detail/{id}"
	RoleTabActionURL        = "/action/role/{id}/tab/{tab}"
	RoleAttachmentUploadURL = "/action/role/{id}/attachments/upload"
	RoleAttachmentDeleteURL = "/action/role/{id}/attachments/delete"
	RoleListURL             = "/roles/list"
	RoleTableURL            = "/action/role/table"
	RoleAddURL              = "/action/role/add"
	RoleEditURL             = "/action/role/edit/{id}"
	RoleDeleteURL           = "/action/role/delete"
	RoleBulkDeleteURL       = "/action/role/bulk-delete"
	RoleSetStatusURL        = "/action/role/set-status"
	RoleBulkSetStatusURL    = "/action/role/bulk-set-status"

	PermissionListURL          = "/permissions/list/{status}"
	PermissionTableURL         = "/action/permission/table/{status}"
	PermissionAddURL           = "/action/permission/add"
	PermissionEditURL          = "/action/permission/edit/{id}"
	PermissionDeleteURL        = "/action/permission/delete"
	PermissionBulkDeleteURL    = "/action/permission/bulk-delete"
	PermissionSetStatusURL     = "/action/permission/set-status"
	PermissionBulkSetStatusURL = "/action/permission/bulk-set-status"

	RolePermissionsURL       = "/manage/roles/{id}/permissions"
	RolePermissionsTableURL  = "/action/manage/roles/{id}/permissions/table"
	RolePermissionsAssignURL = "/action/manage/roles/{id}/permissions/assign"
	RolePermissionsRemoveURL = "/action/manage/roles/{id}/permissions/remove"

	UserRolesURL       = "/manage/users/{id}/roles"
	UserRolesTableURL  = "/action/manage/users/{id}/roles/table"
	UserRolesAssignURL = "/action/manage/users/{id}/roles/assign"
	UserRolesRemoveURL = "/action/manage/users/{id}/roles/remove"

	RoleUsersURL       = "/roles/detail/{id}/users"
	RoleUsersTableURL  = "/action/role/detail/{id}/users/table"
	RoleUsersAssignURL = "/action/role/detail/{id}/users/assign"
	RoleUsersRemoveURL = "/action/role/detail/{id}/users/remove"
	RoleUsersSearchURL = "/action/role/detail/{id}/users/search"

	// Migrated route constants: /detail/ pattern for user-roles and role-permissions
	// Old /manage/ constants kept above for backward compatibility
	UserDetailRolesURL       = "/users/detail/{id}/roles"
	UserDetailRolesTableURL  = "/action/user/detail/{id}/roles/table"
	UserDetailRolesAssignURL = "/action/user/detail/{id}/roles/assign"
	UserDetailRolesRemoveURL = "/action/user/detail/{id}/roles/remove"

	RoleDetailPermissionsURL       = "/roles/detail/{id}/permissions"
	RoleDetailPermissionsTableURL  = "/action/role/detail/{id}/permissions/table"
	RoleDetailPermissionsAssignURL = "/action/role/detail/{id}/permissions/assign"
	RoleDetailPermissionsRemoveURL = "/action/role/detail/{id}/permissions/remove"

	WorkspaceListURL             = "/workspaces/list/{status}"
	WorkspaceTableURL            = "/action/workspace/table/{status}"
	WorkspaceAddURL              = "/action/workspace/add"
	WorkspaceEditURL             = "/action/workspace/edit/{id}"
	WorkspaceDeleteURL           = "/action/workspace/delete"
	WorkspaceBulkDeleteURL       = "/action/workspace/bulk-delete"
	WorkspaceSetStatusURL        = "/action/workspace/set-status"
	WorkspaceBulkSetStatusURL    = "/action/workspace/bulk-set-status"
	WorkspaceSwitchURL           = "/action/admin/switch-workspace"
	WorkspaceDetailURL           = "/workspaces/detail/{id}"
	WorkspaceTabActionURL        = "/action/workspace/{id}/tab/{tab}"
	WorkspaceAttachmentUploadURL = "/action/workspace/{id}/attachments/upload"
	WorkspaceAttachmentDeleteURL = "/action/workspace/{id}/attachments/delete"

	WorkspaceUserListURL             = "/workspace-users/list/{status}"
	WorkspaceUserDetailURL           = "/workspace-users/detail/{id}"
	WorkspaceUserTabActionURL        = "/action/workspace_user/{id}/tab/{tab}"
	WorkspaceUserAddURL              = "/action/workspace_user/add"
	WorkspaceUserDeleteURL           = "/action/workspace_user/delete/{id}"
	WorkspaceUserSetStatusURL        = "/action/workspace_user/set-status/{id}"
	WorkspaceUserSearchURL           = "/action/workspace_user/search"
	WorkspaceUserAttachmentUploadURL = "/action/workspace_user/{id}/attachments/upload"
	WorkspaceUserAttachmentDeleteURL = "/action/workspace_user/{id}/attachments/delete"

	// WorkspaceUserRole — Phase 3 assignment drawer routes.
	WorkspaceUserRoleAddURL         = "/action/workspace_user_role/add"
	WorkspaceUserRoleDeleteURL      = "/action/workspace_user_role/delete/{id}"
	WorkspaceUserRolePermissionsURL = "/action/workspace_user_role/permissions"
	WorkspaceUserRoleSearchRolesURL = "/action/workspace_user_role/search-roles"

	// Admin app dashboard — composite app spanning permission/role/workspace/
	// workspace_user/workspace_user_role. Lives under entydad/views/admin/dashboard.
	AdminDashboardURL = "/admin/dashboard"

	// Client report routes
	ReceivablesAgingURL = "/reports/receivables-aging"

	// Client statement export
	ClientStatementExportURL = "/action/client/{id}/statement/export"

	// ClientRevenueRunURL is the per-client "Run Invoices" drawer endpoint.
	// Static verb segment "revenue-run" comes before {id} so the Go ServeMux
	// does not conflict with the /action/client/table/{status} pattern.
	ClientRevenueRunURL = "/action/client/revenue-run/{id}"

	// Supplier routes
	SupplierDashboardURL        = "/suppliers/dashboard"
	SupplierListURL             = "/suppliers/list/{status}"
	SupplierTableURL            = "/action/supplier/table/{status}"
	SupplierAddURL              = "/action/supplier/add"
	SupplierEditURL             = "/action/supplier/edit/{id}"
	SupplierDeleteURL           = "/action/supplier/delete"
	SupplierBulkDeleteURL       = "/action/supplier/bulk-delete"
	SupplierDetailURL           = "/suppliers/detail/{id}"
	SupplierTabActionURL        = "/action/supplier/{id}/tab/{tab}"
	SupplierAttachmentUploadURL = "/action/supplier/{id}/attachments/upload"
	SupplierAttachmentDeleteURL = "/action/supplier/{id}/attachments/delete"
	SupplierSetStatusURL        = "/action/supplier/set-status"
	SupplierBulkSetStatusURL    = "/action/supplier/bulk-set-status"

	// Supplier statement export
	SupplierStatementExportURL = "/action/supplier/{id}/statement/export"

	// Plan A 20260517-expense-run — Surface A per-supplier drawer URL.
	// "Run Recognitions" CTA on the Statement tab opens this drawer.
	SupplierExpenseRecognitionRunURL = "/action/supplier/expense-recognition-run/{id}"

	// Supplier report routes
	PayablesAgingURL = "/suppliers/reports/payables-aging"

	// Client Tag (Category) routes
	ClientTagListURL          = "/clients/settings/tags/list"
	ClientTagTableURL         = "/action/client/tags/table"
	ClientTagAddURL           = "/action/client/tags/add"
	ClientTagEditURL          = "/action/client/tags/edit/{id}"
	ClientTagDeleteURL        = "/action/client/tags/delete"
	ClientTagBulkDeleteURL    = "/action/client/tags/bulk-delete"
	ClientTagSetStatusURL     = "/action/client/tags/set-status"
	ClientTagBulkSetStatusURL = "/action/client/tags/bulk-set-status"

	// Supplier Tag (Category) routes
	SupplierTagListURL          = "/suppliers/settings/tags/list"
	SupplierTagTableURL         = "/action/supplier/tags/table"
	SupplierTagAddURL           = "/action/supplier/tags/add"
	SupplierTagEditURL          = "/action/supplier/tags/edit/{id}"
	SupplierTagDeleteURL        = "/action/supplier/tags/delete"
	SupplierTagBulkDeleteURL    = "/action/supplier/tags/bulk-delete"
	SupplierTagSetStatusURL     = "/action/supplier/tags/set-status"
	SupplierTagBulkSetStatusURL = "/action/supplier/tags/bulk-set-status"

	// Payment Term routes — client context (shows client + both scopes)
	PaymentTermListURL          = "/clients/settings/payment-terms/list"
	PaymentTermTableURL         = "/action/client/settings/payment-terms/table"
	PaymentTermAddURL           = "/action/client/settings/payment-terms/add"
	PaymentTermEditURL          = "/action/client/settings/payment-terms/edit/{id}"
	PaymentTermDeleteURL        = "/action/client/settings/payment-terms/delete"
	PaymentTermBulkDeleteURL    = "/action/client/settings/payment-terms/bulk-delete"
	PaymentTermSetStatusURL     = "/action/client/settings/payment-terms/set-status"
	PaymentTermBulkSetStatusURL = "/action/client/settings/payment-terms/bulk-set-status"

	// Payment Term routes — supplier context (shows supplier + both scopes)
	SupplierPaymentTermListURL          = "/suppliers/settings/payment-terms/list"
	SupplierPaymentTermTableURL         = "/action/supplier/settings/payment-terms/table"
	SupplierPaymentTermAddURL           = "/action/supplier/settings/payment-terms/add"
	SupplierPaymentTermEditURL          = "/action/supplier/settings/payment-terms/edit/{id}"
	SupplierPaymentTermDeleteURL        = "/action/supplier/settings/payment-terms/delete"
	SupplierPaymentTermBulkDeleteURL    = "/action/supplier/settings/payment-terms/bulk-delete"
	SupplierPaymentTermSetStatusURL     = "/action/supplier/settings/payment-terms/set-status"
	SupplierPaymentTermBulkSetStatusURL = "/action/supplier/settings/payment-terms/bulk-set-status"

	// Auth routes (all prefixed with /auth/)
	AuthLoginURL             = "/auth/login"
	AuthLoginPostURL         = "/auth/login"
	AuthSignupURL            = "/auth/signup"
	AuthSignupPostURL        = "/auth/signup"
	AuthResetPasswordURL     = "/auth/reset-password"
	AuthResetPasswordPostURL = "/auth/reset-password"
	AuthResetConfirmURL      = "/auth/reset-password/confirm"
	AuthResetConfirmPostURL  = "/auth/reset-password/confirm"
	AuthChangePasswordURL    = "/auth/change-password"
	AuthLogoutURL            = "/auth/logout"

	// Legacy login routes (redirect to /auth/login)
	LoginURL     = "/login"
	LoginPostURL = "/login"

	// DefaultAppRedirectURL is the default post-login redirect path.
	// Consumer apps should set Deps.RedirectURL to override this.
	//
	// Post-P12 (2026-05-22) of docs/plan/20260521-workspace-keyed-routing:
	// /app/* is gone. Post-login redirects should use
	// composition.homeURLForWorkspaceID() to land on /w/{slug}/home; this
	// constant is the workspace-less fallback (/me/inbox) for callers that
	// can't resolve a workspace slug at redirect time (password reset, etc.).
	DefaultAppRedirectURL = "/me/inbox"

	// TaxRegistration — polymorphic (client + workspace party types in v1)
	// URL convention: party_type + party_id come from the parent detail page context.
	ClientTaxRegistrationListURL   = "/clients/detail/{id}/tax-registrations"
	ClientTaxRegistrationAddURL    = "/action/client/{id}/tax-registration/add"
	ClientTaxRegistrationEditURL   = "/action/client/{id}/tax-registration/edit/{reg_id}"
	ClientTaxRegistrationDeleteURL = "/action/client/{id}/tax-registration/delete"

	WorkspaceTaxRegistrationListURL   = "/workspace/settings/tax-registrations"
	WorkspaceTaxRegistrationAddURL    = "/action/workspace/tax-registration/add"
	WorkspaceTaxRegistrationEditURL   = "/action/workspace/tax-registration/edit/{reg_id}"
	WorkspaceTaxRegistrationDeleteURL = "/action/workspace/tax-registration/delete"
)
