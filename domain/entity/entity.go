// Package entity is the facade for the entydad `entity` esqyma domain.
//
// Consumers (block/, service-admin) use the E-prefixed names exported here
// (e.g. ClientLabels, DefaultClientRoutes, WorkspaceUserAddURL) which resolve
// to the entity-local types in the sub-context sub-packages below. This keeps
// consumer call sites unchanged while the internals follow the entity-local
// naming convention (one `Labels`/`Routes`/`DefaultRoutes` per entity package).
//
// Import-cycle rule (contract D): entity packages MUST NEVER import this
// facade. Cross-entity references go DIRECT to the sibling package. Only this
// facade and the sub-context module files import the entity packages.
//
// Sub-contexts (identity / party / location / commerce) are navigation-only
// folder groupings; there is exactly ONE Go contract package per entity, and
// exactly ONE import namespace for consumers: `entity.*`.
package entity

import (
	paymentterm "github.com/erniealice/entydad-golang/domain/entity/commerce/payment_term"
	permission "github.com/erniealice/entydad-golang/domain/entity/identity/permission"
	role "github.com/erniealice/entydad-golang/domain/entity/identity/role"
	user "github.com/erniealice/entydad-golang/domain/entity/identity/user"
	workspace "github.com/erniealice/entydad-golang/domain/entity/identity/workspace"
	workspaceuser "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user"
	workspaceuserrole "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user_role"
	locationpkg "github.com/erniealice/entydad-golang/domain/entity/location/location"
	locationarea "github.com/erniealice/entydad-golang/domain/entity/location/location_area"
	client "github.com/erniealice/entydad-golang/domain/entity/party/client"
	clienttag "github.com/erniealice/entydad-golang/domain/entity/party/client_tag"
	supplier "github.com/erniealice/entydad-golang/domain/entity/party/supplier"
	suppliertag "github.com/erniealice/entydad-golang/domain/entity/party/supplier_tag"
)

// ---------------------------------------------------------------------------
// Client (party/client)
// ---------------------------------------------------------------------------

type ClientLabels = client.Labels
type ClientPageLabels = client.PageLabels
type ClientButtonLabels = client.ButtonLabels
type ClientColumnLabels = client.ColumnLabels
type ClientEmptyLabels = client.EmptyLabels
type ClientFormLabels = client.FormLabels
type ClientDetailLabels = client.DetailLabels
type ClientCompanyDetailLabels = client.CompanyDetailLabels
type ClientDetailSectionLabels = client.DetailSectionLabels
type ClientDetailTagLabels = client.DetailTagLabels
type ClientPurchaseHistoryLabels = client.PurchaseHistoryLabels
type ClientDetailTabLabels = client.DetailTabLabels
type ClientPriceSchedulesLabels = client.PriceSchedulesLabels
type ClientSubscriptionLabels = client.SubscriptionLabels
type ClientStatementLabels = client.StatementLabels
type ClientOutstandingTableLabels = client.OutstandingTableLabels
type ClientOutstandingTableColumnLabels = client.OutstandingTableColumnLabels
type ClientOutstandingTableEmptyLabels = client.OutstandingTableEmptyLabels
type ClientRevenueRunLabels = client.RevenueRunLabels
type ClientRevenueRunErrorLabels = client.RevenueRunErrorLabels
type ClientDetailActionLabels = client.DetailActionLabels
type ClientBulkActionLabels = client.BulkActionLabels
type ClientDashboardLabels = client.DashboardLabels

type ClientRoutes = client.Routes

func DefaultClientRoutes() ClientRoutes { return client.DefaultRoutes() }

const (
	ClientDashboardURL        = client.DashboardURL
	ClientListURL             = client.ListURL
	ClientTableURL            = client.TableURL
	ClientAddURL              = client.AddURL
	ClientEditURL             = client.EditURL
	ClientDeleteURL           = client.DeleteURL
	ClientBulkDeleteURL       = client.BulkDeleteURL
	ClientDetailURL           = client.DetailURL
	ClientTabActionURL        = client.TabActionURL
	ClientAttachmentUploadURL = client.AttachmentUploadURL
	ClientAttachmentDeleteURL = client.AttachmentDeleteURL
	ClientSetStatusURL        = client.SetStatusURL
	ClientBulkSetStatusURL    = client.BulkSetStatusURL
	ClientSearchURL           = client.SearchURL
	ClientStatementExportURL  = client.StatementExportURL
	ClientRevenueRunURL       = client.RevenueRunURL
)

// ---------------------------------------------------------------------------
// Supplier (party/supplier)
// ---------------------------------------------------------------------------

type SupplierLabels = supplier.Labels
type SupplierPageLabels = supplier.PageLabels
type SupplierButtonLabels = supplier.ButtonLabels
type SupplierColumnLabels = supplier.ColumnLabels
type SupplierEmptyLabels = supplier.EmptyLabels
type SupplierFormLabels = supplier.FormLabels
type SupplierDetailLabels = supplier.DetailLabels
type SupplierDetailSectionLabels = supplier.DetailSectionLabels
type SupplierPurchaseOrdersLabels = supplier.PurchaseOrdersLabels
type SupplierActionLabels = supplier.ActionLabels
type SupplierDashboardLabels = supplier.DashboardLabels

type SupplierRoutes = supplier.Routes

func DefaultSupplierRoutes() SupplierRoutes { return supplier.DefaultRoutes() }

const (
	SupplierDashboardURL            = supplier.DashboardURL
	SupplierListURL                 = supplier.ListURL
	SupplierTableURL                = supplier.TableURL
	SupplierAddURL                  = supplier.AddURL
	SupplierEditURL                 = supplier.EditURL
	SupplierDeleteURL               = supplier.DeleteURL
	SupplierBulkDeleteURL           = supplier.BulkDeleteURL
	SupplierDetailURL               = supplier.DetailURL
	SupplierTabActionURL            = supplier.TabActionURL
	SupplierAttachmentUploadURL     = supplier.AttachmentUploadURL
	SupplierAttachmentDeleteURL     = supplier.AttachmentDeleteURL
	SupplierSetStatusURL             = supplier.SetStatusURL
	SupplierBulkSetStatusURL         = supplier.BulkSetStatusURL
	SupplierStatementExportURL       = supplier.StatementExportURL
	SupplierExpenseRecognitionRunURL = supplier.ExpenseRecognitionRunURL
)

// ---------------------------------------------------------------------------
// ClientTag (party/client_tag)
// ---------------------------------------------------------------------------

type ClientTagLabels = clienttag.Labels
type ClientTagPageLabels = clienttag.PageLabels
type ClientTagButtonLabels = clienttag.ButtonLabels
type ClientTagColumnLabels = clienttag.ColumnLabels
type ClientTagEmptyLabels = clienttag.EmptyLabels
type ClientTagActionLabels = clienttag.ActionLabels
type ClientTagConfirmLabels = clienttag.ConfirmLabels

type ClientTagRoutes = clienttag.Routes

func DefaultClientTagRoutes() ClientTagRoutes { return clienttag.DefaultRoutes() }

// ---------------------------------------------------------------------------
// SupplierTag (party/supplier_tag)
// ---------------------------------------------------------------------------

type SupplierTagLabels = suppliertag.Labels
type SupplierTagPageLabels = suppliertag.PageLabels
type SupplierTagButtonLabels = suppliertag.ButtonLabels
type SupplierTagColumnLabels = suppliertag.ColumnLabels
type SupplierTagEmptyLabels = suppliertag.EmptyLabels
type SupplierTagActionLabels = suppliertag.ActionLabels
type SupplierTagConfirmLabels = suppliertag.ConfirmLabels

type SupplierTagRoutes = suppliertag.Routes

func DefaultSupplierTagRoutes() SupplierTagRoutes { return suppliertag.DefaultRoutes() }

// ---------------------------------------------------------------------------
// PaymentTerm (commerce/payment_term)
// ---------------------------------------------------------------------------

type PaymentTermLabels = paymentterm.Labels
type PaymentTermPageLabels = paymentterm.PageLabels
type PaymentTermButtonLabels = paymentterm.ButtonLabels
type PaymentTermColumnLabels = paymentterm.ColumnLabels
type PaymentTermEmptyLabels = paymentterm.EmptyLabels
type PaymentTermFormLabels = paymentterm.FormLabels
type PaymentTermActionLabels = paymentterm.ActionLabels

type PaymentTermRoutes = paymentterm.Routes
type SupplierPaymentTermRoutes = paymentterm.SupplierRoutes

func DefaultPaymentTermRoutes() PaymentTermRoutes { return paymentterm.DefaultRoutes() }
func DefaultSupplierPaymentTermRoutes() SupplierPaymentTermRoutes {
	return paymentterm.DefaultSupplierRoutes()
}

// ---------------------------------------------------------------------------
// User (identity/user)
// ---------------------------------------------------------------------------

type UserLabels = user.Labels
type UserPageLabels = user.PageLabels
type UserButtonLabels = user.ButtonLabels
type UserColumnLabels = user.ColumnLabels
type UserEmptyLabels = user.EmptyLabels
type UserFormLabels = user.FormLabels
type UserActionLabels = user.ActionLabels
type UserDetailLabels = user.DetailLabels
type UserDetailSecurityLabels = user.DetailSecurityLabels
type UserDetailEmptyStateLabels = user.DetailEmptyStateLabels
type UserDetailBasicInfoLabels = user.DetailBasicInfoLabels
type UserDetailTabLabels = user.DetailTabLabels
type UserDashboardLabels = user.DashboardLabels
type UserRoleLabels = user.RoleLabels
type UserRolePageLabels = user.RolePageLabels
type UserRoleButtonLabels = user.RoleButtonLabels
type UserRoleColumnLabels = user.RoleColumnLabels
type UserRoleEmptyLabels = user.RoleEmptyLabels
type UserRoleFormLabels = user.RoleFormLabels
type UserRoleActionLabels = user.RoleActionLabels

type UserRoutes = user.Routes

func DefaultUserRoutes() UserRoutes { return user.DefaultRoutes() }

// ---------------------------------------------------------------------------
// Role (identity/role)
// ---------------------------------------------------------------------------

type RoleLabels = role.Labels
type RolePageLabels = role.PageLabels
type RoleButtonLabels = role.ButtonLabels
type RoleColumnLabels = role.ColumnLabels
type RoleEmptyLabels = role.EmptyLabels
type RoleFormLabels = role.FormLabels
type RoleActionLabels = role.ActionLabels
type RoleDetailLabels = role.DetailLabels
type RoleDetailTabLabels = role.DetailTabLabels
type RoleDetailInfoLabels = role.DetailInfoLabels
type RolePermissionLabels = role.PermissionLabels
type RolePermissionPageLabels = role.PermissionPageLabels
type RolePermissionButtonLabels = role.PermissionButtonLabels
type RolePermissionColumnLabels = role.PermissionColumnLabels
type RolePermissionEmptyLabels = role.PermissionEmptyLabels
type RolePermissionFormLabels = role.PermissionFormLabels
type RolePermissionActionLabels = role.PermissionActionLabels
type RoleUserLabels = role.UserLabels
type RoleUserPageLabels = role.UserPageLabels
type RoleUserButtonLabels = role.UserButtonLabels
type RoleUserColumnLabels = role.UserColumnLabels
type RoleUserEmptyLabels = role.UserEmptyLabels
type RoleUserFormLabels = role.UserFormLabels
type RoleUserActionLabels = role.UserActionLabels

type RoleRoutes = role.Routes

func DefaultRoleRoutes() RoleRoutes { return role.DefaultRoutes() }

// ---------------------------------------------------------------------------
// Permission (identity/permission)
// ---------------------------------------------------------------------------

type PermissionLabels = permission.Labels
type PermissionPageLabels = permission.PageLabels
type PermissionButtonLabels = permission.ButtonLabels
type PermissionColumnLabels = permission.ColumnLabels
type PermissionEmptyLabels = permission.EmptyLabels
type PermissionFormLabels = permission.FormLabels
type PermissionActionLabels = permission.ActionLabels

type PermissionRoutes = permission.Routes

func DefaultPermissionRoutes() PermissionRoutes { return permission.DefaultRoutes() }

// ---------------------------------------------------------------------------
// Workspace (identity/workspace)
// ---------------------------------------------------------------------------

type WorkspaceLabels = workspace.Labels
type WorkspaceDetailLabels = workspace.DetailLabels
type WorkspaceDetailTaxRegLabels = workspace.DetailTaxRegLabels
type WorkspaceDetailInfoLabels = workspace.DetailInfoLabels
type WorkspaceDetailTabLabels = workspace.DetailTabLabels
type WorkspaceDetailUserLabels = workspace.DetailUserLabels
type WorkspacePageLabels = workspace.PageLabels
type WorkspaceButtonLabels = workspace.ButtonLabels
type WorkspaceColumnLabels = workspace.ColumnLabels
type WorkspaceEmptyLabels = workspace.EmptyLabels
type WorkspaceFormLabels = workspace.FormLabels
type WorkspaceActionLabels = workspace.ActionLabels

type WorkspaceRoutes = workspace.Routes

func DefaultWorkspaceRoutes() WorkspaceRoutes { return workspace.DefaultRoutes() }

const (
	WorkspaceDetailURL = workspace.DetailURL
)

// ---------------------------------------------------------------------------
// WorkspaceUser (identity/workspace_user)
// ---------------------------------------------------------------------------

type WorkspaceUserLabels = workspaceuser.Labels
type WorkspaceUserPageLabels = workspaceuser.PageLabels
type WorkspaceUserColumnLabels = workspaceuser.ColumnLabels
type WorkspaceUserDetailLabels = workspaceuser.DetailLabels
type WorkspaceUserDetailInfoLabels = workspaceuser.DetailInfoLabels
type WorkspaceUserDetailTabLabels = workspaceuser.DetailTabLabels
type WorkspaceUserDetailRolesLabels = workspaceuser.DetailRolesLabels
type WorkspaceUserFormLabels = workspaceuser.FormLabels
type WorkspaceUserActionLabels = workspaceuser.ActionLabels

type WorkspaceUserRoutes = workspaceuser.Routes

func DefaultWorkspaceUserRoutes() WorkspaceUserRoutes { return workspaceuser.DefaultRoutes() }

const (
	WorkspaceUserAddURL    = workspaceuser.AddURL
	WorkspaceUserDetailURL = workspaceuser.DetailURL
)

// ---------------------------------------------------------------------------
// WorkspaceUserRole (identity/workspace_user_role)
// ---------------------------------------------------------------------------

type WorkspaceUserRoleLabels = workspaceuserrole.Labels
type WorkspaceUserRoleFormLabels = workspaceuserrole.FormLabels
type WorkspaceUserRoleButtonLabels = workspaceuserrole.ButtonLabels

type WorkspaceUserRoleRoutes = workspaceuserrole.Routes

func DefaultWorkspaceUserRoleRoutes() WorkspaceUserRoleRoutes {
	return workspaceuserrole.DefaultRoutes()
}

const (
	WorkspaceUserRoleAddURL    = workspaceuserrole.AddURL
	WorkspaceUserRoleDeleteURL = workspaceuserrole.DeleteURL
)

// ---------------------------------------------------------------------------
// Location (location/location)
// ---------------------------------------------------------------------------

type LocationLabels = locationpkg.Labels
type LocationPageLabels = locationpkg.PageLabels
type LocationButtonLabels = locationpkg.ButtonLabels
type LocationColumnLabels = locationpkg.ColumnLabels
type LocationEmptyLabels = locationpkg.EmptyLabels
type LocationFormLabels = locationpkg.FormLabels
type LocationActionLabels = locationpkg.ActionLabels
type LocationDetailLabels = locationpkg.DetailLabels
type LocationDetailBasicInfoLabels = locationpkg.DetailBasicInfoLabels
type LocationDetailTabLabels = locationpkg.DetailTabLabels
type LocationDetailEmptyLabels = locationpkg.DetailEmptyLabels
type LocationDashboardLabels = locationpkg.DashboardLabels

type LocationRoutes = locationpkg.Routes

func DefaultLocationRoutes() LocationRoutes { return locationpkg.DefaultRoutes() }

// ---------------------------------------------------------------------------
// LocationArea (location/location_area)
// ---------------------------------------------------------------------------

type LocationAreaLabels = locationarea.Labels
type LocationAreaPageLabels = locationarea.PageLabels
type LocationAreaButtonLabels = locationarea.ButtonLabels
type LocationAreaColumnLabels = locationarea.ColumnLabels
type LocationAreaEmptyLabels = locationarea.EmptyLabels
type LocationAreaFormLabels = locationarea.FormLabels
type LocationAreaActionLabels = locationarea.ActionLabels
type LocationAreaErrorLabels = locationarea.ErrorLabels

func DefaultLocationAreaLabels() LocationAreaLabels { return locationarea.DefaultLabels() }

type LocationAreaRoutes = locationarea.Routes

func DefaultLocationAreaRoutes() LocationAreaRoutes { return locationarea.DefaultRoutes() }
