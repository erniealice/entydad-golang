// Package block — typed wiring contract for entydad.Block.
//
// This file declares what entydad's Block() needs from outside.
// Service-admin's composition layer constructs a *UseCases value from
// espyna's consumer container; entydad's Block() consumes only this
// typed shape.
//
// Shape this struct by what ENTYDAD needs, NOT by mirroring espyna's
// *consumer.UseCases. Service-admin's adapter is the only place that
// knows both vocabularies. If espyna restructures its container, only
// that adapter changes.
package block

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	conversationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/communication/conversation"
	conversationpostpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/communication/conversation_post"
	conversationreadreceiptpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/communication/conversation_read_receipt"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientcatpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	locationareapb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location_area"
	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	rolepermissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role_permission"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	suppliercatpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier_category"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	wurpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	taxregistrationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_registration"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	stmtspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/statements"

	locationdashboard "github.com/erniealice/entydad-golang/domain/entity/location/location/dashboard"
	admindashboard "github.com/erniealice/entydad-golang/service/dashboard/views/admin/dashboard"
)

// UseCases declares everything entydad's Block() needs from outside.
// Construction is service-admin's job; entydad only declares the shape.
//
// Naming conventions (plan.md §2):
//  1. Field names are SINGULAR matching the proto folder name.
//  2. Group struct types use the `<Entity>UseCases` suffix.
//  3. Closure signatures use proto request/response types (no block-local
//     transport types). Ex-orchestrators are regular fields after Phase 0.
//  4. Nested entity ops mirror proto nesting (e.g. Client.Category).
type UseCases struct {
	// GetWorkspaceIDFromCtx extracts the workspace ID from a request context.
	// Wired by service-admin as consumer.GetWorkspaceIDFromContext.
	// Used by the dashboard wiring helpers in wiring.go.
	GetWorkspaceIDFromCtx func(ctx context.Context) string

	// Domain CRUD + use-case groups (singular field, XxxUseCases type)
	Client            ClientUseCases
	User              UserUseCases
	Role              RoleUseCases
	RolePermission    RolePermissionUseCases
	Permission        PermissionUseCases
	Location          LocationUseCases
	LocationArea      LocationAreaUseCases
	Workspace         WorkspaceUseCases
	WorkspaceUser     WorkspaceUserUseCases
	WorkspaceUserRole WorkspaceUserRoleUseCases
	Supplier          SupplierUseCases
	Subscription      SubscriptionUseCases
	Revenue           RevenueUseCases
	Collection        CollectionUseCases
	Category          CategoryUseCases
	PriceSchedule     PriceScheduleUseCases
	PricePlan         PricePlanUseCases
	PurchaseOrder     PurchaseOrderUseCases
	TaxRegistration   TaxRegistrationUseCases
	Conversation      ConversationUseCases

	// Reports — service-driven report use case closures consumed by the
	// client/supplier detail + list views. Wave B P1.E.4 (statements).
	Reports ReportsUseCases

	// Dashboard closures — view-layer typed; service-admin builds these
	// by calling the espyna use cases and mapping their internal response
	// types to entydad's view types.
	GetLocationDashboardPageData func(ctx context.Context) (*locationdashboard.LocationDashboardData, error)
	GetAdminDashboardPageData    func(ctx context.Context) (*admindashboard.AdminDashboardData, error)
}

// ReportsUseCases aggregates the service-driven report use case closures
// the entydad views need. Wave B P1.E.4 — statements + balances migrated
// out of `ctx.LedgerReportingSvc` / `entydad/block.LedgerReportingService`
// (the duck interface returning maps) into the proto-shaped service layer.
type ReportsUseCases struct {
	Statements StatementsUseCases
}

// StatementsUseCases — service-driven counterparty statement + balance
// closures consumed by client/supplier detail/list views. Migrated
// 2026-05-21 (Wave B P1.E.4) out of `entydad/block.LedgerReportingService`.
//
// **Map shim:** the view layer historically accepted
// `map[string]int64` (counterparty_id → centavo balance) for
// GetClientBalances/GetSupplierBalances. The new service-driven use cases
// return `[]*BalanceRow` (typed proto rows) per Q-SDM-MAP-SHAPES. The
// `ListClient/SupplierBalancesAsMap` closures expose the legacy shape so
// the view files keep their map-based table cell lookup unchanged; the
// typed pivots are also surfaced for any future migration.
type StatementsUseCases struct {
	GetClientStatement   func(context.Context, *stmtspb.GetClientStatementRequest) (*stmtspb.GetClientStatementResponse, error)
	GetSupplierStatement func(context.Context, *stmtspb.GetSupplierStatementRequest) (*stmtspb.GetSupplierStatementResponse, error)
	ListClientBalances   func(context.Context, *stmtspb.ListClientBalancesRequest) (*stmtspb.ListClientBalancesResponse, error)
	ListSupplierBalances func(context.Context, *stmtspb.ListSupplierBalancesRequest) (*stmtspb.ListSupplierBalancesResponse, error)
	// Map-returning shims for backwards compatibility with the entydad
	// client/supplier list views, which still consume the legacy
	// `map[string]int64` shape directly into their per-row TableCell
	// lookups. Service-admin's adapter wires these on top of
	// ListClient/SupplierBalances by converting `[]*BalanceRow` to the map.
	ListClientBalancesAsMap   func(context.Context) (map[string]int64, error)
	ListSupplierBalancesAsMap func(context.Context) (map[string]int64, error)
}

// ClientUseCases — direct CRUD on the Client entity + nested ClientCategory ops.
// Category (singular) under Client mirrors how proto nests client_category under entity/.
type ClientUseCases struct {
	GetListPageData func(context.Context, *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
	Create          func(context.Context, *clientpb.CreateClientRequest) (*clientpb.CreateClientResponse, error)
	Read            func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	Update          func(context.Context, *clientpb.UpdateClientRequest) (*clientpb.UpdateClientResponse, error)
	Delete          func(context.Context, *clientpb.DeleteClientRequest) (*clientpb.DeleteClientResponse, error)
	Category        ClientCategoryUseCases
}

type ClientCategoryUseCases struct {
	List   func(context.Context, *clientcatpb.ListClientCategoriesRequest) (*clientcatpb.ListClientCategoriesResponse, error)
	Create func(context.Context, *clientcatpb.CreateClientCategoryRequest) (*clientcatpb.CreateClientCategoryResponse, error)
	Delete func(context.Context, *clientcatpb.DeleteClientCategoryRequest) (*clientcatpb.DeleteClientCategoryResponse, error)
}

type UserUseCases struct {
	GetListPageData func(context.Context, *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error)
	Create          func(context.Context, *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
	Read            func(context.Context, *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)
	Update          func(context.Context, *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error)
	Delete          func(context.Context, *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error)
	List            func(context.Context, *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error)
}

type RoleUseCases struct {
	GetListPageData func(context.Context, *rolepb.GetRoleListPageDataRequest) (*rolepb.GetRoleListPageDataResponse, error)
	Create          func(context.Context, *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error)
	Read            func(context.Context, *rolepb.ReadRoleRequest) (*rolepb.ReadRoleResponse, error)
	Update          func(context.Context, *rolepb.UpdateRoleRequest) (*rolepb.UpdateRoleResponse, error)
	Delete          func(context.Context, *rolepb.DeleteRoleRequest) (*rolepb.DeleteRoleResponse, error)
	GetItemPageData func(context.Context, *rolepb.GetRoleItemPageDataRequest) (*rolepb.GetRoleItemPageDataResponse, error)
	List            func(context.Context, *rolepb.ListRolesRequest) (*rolepb.ListRolesResponse, error)
}

type RolePermissionUseCases struct {
	Create func(context.Context, *rolepermissionpb.CreateRolePermissionRequest) (*rolepermissionpb.CreateRolePermissionResponse, error)
	Delete func(context.Context, *rolepermissionpb.DeleteRolePermissionRequest) (*rolepermissionpb.DeleteRolePermissionResponse, error)
}

type PermissionUseCases struct {
	GetListPageData func(context.Context, *permissionpb.GetPermissionListPageDataRequest) (*permissionpb.GetPermissionListPageDataResponse, error)
	Create          func(context.Context, *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error)
	Read            func(context.Context, *permissionpb.ReadPermissionRequest) (*permissionpb.ReadPermissionResponse, error)
	Update          func(context.Context, *permissionpb.UpdatePermissionRequest) (*permissionpb.UpdatePermissionResponse, error)
	Delete          func(context.Context, *permissionpb.DeletePermissionRequest) (*permissionpb.DeletePermissionResponse, error)
	List            func(context.Context, *permissionpb.ListPermissionsRequest) (*permissionpb.ListPermissionsResponse, error)
}

type LocationUseCases struct {
	GetListPageData func(context.Context, *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error)
	Create          func(context.Context, *locationpb.CreateLocationRequest) (*locationpb.CreateLocationResponse, error)
	Read            func(context.Context, *locationpb.ReadLocationRequest) (*locationpb.ReadLocationResponse, error)
	Update          func(context.Context, *locationpb.UpdateLocationRequest) (*locationpb.UpdateLocationResponse, error)
	Delete          func(context.Context, *locationpb.DeleteLocationRequest) (*locationpb.DeleteLocationResponse, error)
}

// LocationAreaUseCases — typed proto closures for the location_area entity.
// Replaces the former duck path (block-local CRUDSource.ListSimple/Create/
// Read/Update/Delete against the "location_area" collection string) with the
// proto-shaped espyna use cases. Wired by service-admin's adapter from
// uc.Entity.LocationArea.* (mirrors the Location group). Closures are
// nil-safe at the call site (commerce.go) so a partially-wired or mock build
// renders empty rather than panicking.
type LocationAreaUseCases struct {
	List   func(context.Context, *locationareapb.ListLocationAreasRequest) (*locationareapb.ListLocationAreasResponse, error)
	Create func(context.Context, *locationareapb.CreateLocationAreaRequest) (*locationareapb.CreateLocationAreaResponse, error)
	Read   func(context.Context, *locationareapb.ReadLocationAreaRequest) (*locationareapb.ReadLocationAreaResponse, error)
	Update func(context.Context, *locationareapb.UpdateLocationAreaRequest) (*locationareapb.UpdateLocationAreaResponse, error)
	Delete func(context.Context, *locationareapb.DeleteLocationAreaRequest) (*locationareapb.DeleteLocationAreaResponse, error)
}

type WorkspaceUseCases struct {
	GetListPageData func(context.Context, *workspacepb.GetWorkspaceListPageDataRequest) (*workspacepb.GetWorkspaceListPageDataResponse, error)
	Create          func(context.Context, *workspacepb.CreateWorkspaceRequest) (*workspacepb.CreateWorkspaceResponse, error)
	Read            func(context.Context, *workspacepb.ReadWorkspaceRequest) (*workspacepb.ReadWorkspaceResponse, error)
	Update          func(context.Context, *workspacepb.UpdateWorkspaceRequest) (*workspacepb.UpdateWorkspaceResponse, error)
	Delete          func(context.Context, *workspacepb.DeleteWorkspaceRequest) (*workspacepb.DeleteWorkspaceResponse, error)
	Switch          func(context.Context, *workspacepb.SwitchWorkspaceRequest) (*workspacepb.SwitchWorkspaceResponse, error)
}

type WorkspaceUserUseCases struct {
	GetListPageData func(context.Context, *workspaceuserpb.GetWorkspaceUserListPageDataRequest) (*workspaceuserpb.GetWorkspaceUserListPageDataResponse, error)
	GetItemPageData func(context.Context, *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	Create          func(context.Context, *workspaceuserpb.CreateWorkspaceUserRequest) (*workspaceuserpb.CreateWorkspaceUserResponse, error)
	Delete          func(context.Context, *workspaceuserpb.DeleteWorkspaceUserRequest) (*workspaceuserpb.DeleteWorkspaceUserResponse, error)
	List            func(context.Context, *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
}

type WorkspaceUserRoleUseCases struct {
	Create          func(context.Context, *wurpb.CreateWorkspaceUserRoleRequest) (*wurpb.CreateWorkspaceUserRoleResponse, error)
	Delete          func(context.Context, *wurpb.DeleteWorkspaceUserRoleRequest) (*wurpb.DeleteWorkspaceUserRoleResponse, error)
	GetListPageData func(context.Context, *wurpb.GetWorkspaceUserRoleListPageDataRequest) (*wurpb.GetWorkspaceUserRoleListPageDataResponse, error)
}

// SupplierUseCases — direct CRUD + nested SupplierCategory ops.
// Category (singular) mirrors how proto nests supplier_category under entity/.
type SupplierUseCases struct {
	GetListPageData func(context.Context, *supplierpb.GetSupplierListPageDataRequest) (*supplierpb.GetSupplierListPageDataResponse, error)
	Create          func(context.Context, *supplierpb.CreateSupplierRequest) (*supplierpb.CreateSupplierResponse, error)
	Read            func(context.Context, *supplierpb.ReadSupplierRequest) (*supplierpb.ReadSupplierResponse, error)
	Update          func(context.Context, *supplierpb.UpdateSupplierRequest) (*supplierpb.UpdateSupplierResponse, error)
	Delete          func(context.Context, *supplierpb.DeleteSupplierRequest) (*supplierpb.DeleteSupplierResponse, error)
	Category        SupplierCategoryUseCases
}

type SupplierCategoryUseCases struct {
	List   func(context.Context, *suppliercatpb.ListSupplierCategoriesRequest) (*suppliercatpb.ListSupplierCategoriesResponse, error)
	Create func(context.Context, *suppliercatpb.CreateSupplierCategoryRequest) (*suppliercatpb.CreateSupplierCategoryResponse, error)
	Delete func(context.Context, *suppliercatpb.DeleteSupplierCategoryRequest) (*suppliercatpb.DeleteSupplierCategoryResponse, error)
}

type SubscriptionUseCases struct {
	List                   func(context.Context, *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)
	GetListPageData        func(context.Context, *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error)
	CountActiveByClientIDs func(context.Context, *subscriptionpb.CountActiveByClientIdsRequest) (*subscriptionpb.CountActiveByClientIdsResponse, error)
}

type RevenueUseCases struct {
	List func(context.Context, *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error)
	// Ex-helpers promoted to proto-defined use cases in Phase 0.
	ListRevenueRunCandidates func(context.Context, *revrunpb.ListRevenueRunCandidatesRequest) (*revrunpb.ListRevenueRunCandidatesResponse, error)
	GenerateRevenueRun       func(context.Context, *revrunpb.GenerateRevenueRunRequest) (*revrunpb.GenerateRevenueRunResponse, error)
}

type CollectionUseCases struct {
	ListByClient func(context.Context, *collectionpb.ListByClientRequest) (*collectionpb.ListByClientResponse, error)
}

// CategoryUseCases — generic common/category CRUD used by client-tag and supplier-tag modules.
type CategoryUseCases struct {
	List   func(context.Context, *commonpb.ListCategoriesRequest) (*commonpb.ListCategoriesResponse, error)
	Create func(context.Context, *commonpb.CreateCategoryRequest) (*commonpb.CreateCategoryResponse, error)
	Read   func(context.Context, *commonpb.ReadCategoryRequest) (*commonpb.ReadCategoryResponse, error)
	Update func(context.Context, *commonpb.UpdateCategoryRequest) (*commonpb.UpdateCategoryResponse, error)
	Delete func(context.Context, *commonpb.DeleteCategoryRequest) (*commonpb.DeleteCategoryResponse, error)
}

type PriceScheduleUseCases struct {
	List func(context.Context, *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
}

type PricePlanUseCases struct {
	List func(context.Context, *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
}

type PurchaseOrderUseCases struct {
	List func(context.Context, *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
}

type TaxRegistrationUseCases struct {
	List func(context.Context, *taxregistrationpb.ListTaxRegistrationsRequest) (*taxregistrationpb.ListTaxRegistrationsResponse, error)
}

// ConversationUseCases — secure-messaging surface (Plan-4, 2026-06-03).
//
// Closure signatures use the REAL espyna use-case request/response types:
// AssignConversation + SetConversationStatus consume UpdateConversationRequest
// (no distinct Assign/SetStatus proto message exists — the espyna use cases
// dispatch on the mutated field); SendConversationPost consumes
// CreateConversationPostRequest; MarkConversationRead consumes
// CreateConversationReadReceiptRequest.
//
// Client-portal scoping (acting_as_client_id) is applied inside the espyna use
// cases; the view/block layer never reads it directly.
type ConversationUseCases struct {
	List      func(context.Context, *conversationpb.ListConversationsRequest) (*conversationpb.ListConversationsResponse, error)
	Read      func(context.Context, *conversationpb.ReadConversationRequest) (*conversationpb.ReadConversationResponse, error)
	Create    func(context.Context, *conversationpb.CreateConversationRequest) (*conversationpb.CreateConversationResponse, error)
	Assign    func(context.Context, *conversationpb.UpdateConversationRequest) (*conversationpb.UpdateConversationResponse, error)
	SetStatus func(context.Context, *conversationpb.UpdateConversationRequest) (*conversationpb.UpdateConversationResponse, error)
	Post      ConversationPostUseCases
	Receipt   ConversationReadReceiptUseCases
}

// ConversationPostUseCases — post list + composer send.
type ConversationPostUseCases struct {
	List func(context.Context, *conversationpostpb.ListConversationPostsRequest) (*conversationpostpb.ListConversationPostsResponse, error)
	Send func(context.Context, *conversationpostpb.CreateConversationPostRequest) (*conversationpostpb.CreateConversationPostResponse, error)
}

// ConversationReadReceiptUseCases — read-receipt high-water-mark upsert.
type ConversationReadReceiptUseCases struct {
	MarkRead func(context.Context, *conversationreadreceiptpb.CreateConversationReadReceiptRequest) (*conversationreadreceiptpb.CreateConversationReadReceiptResponse, error)
}

// RequireFor returns an error listing every needed-but-nil field for cfg's
// enabled modules. Called at Block() entry; missing field → startup error.
//
// CRITICAL: this is the deterministic completeness check. Partial wiring
// is a startup error, not a runtime nil panic.
func (u *UseCases) RequireFor(cfg *blockConfig) error {
	if u == nil {
		return fmt.Errorf("entydad.Block: WithUseCases(...) was not supplied")
	}

	var missing []string
	check := func(ok bool, name string) {
		if !ok {
			missing = append(missing, name)
		}
	}

	if cfg.enableAll || cfg.client {
		check(u.Client.GetListPageData != nil, "UseCases.Client.GetListPageData")
		check(u.Client.Create != nil, "UseCases.Client.Create")
		check(u.Client.Read != nil, "UseCases.Client.Read")
		check(u.Client.Update != nil, "UseCases.Client.Update")
		check(u.Client.Delete != nil, "UseCases.Client.Delete")
		// Category and cross-domain deps are optional (nil-safe wiring)
	}

	if cfg.enableAll || cfg.user {
		check(u.User.GetListPageData != nil, "UseCases.User.GetListPageData")
		check(u.User.Create != nil, "UseCases.User.Create")
		check(u.User.Read != nil, "UseCases.User.Read")
		check(u.User.Update != nil, "UseCases.User.Update")
		check(u.User.Delete != nil, "UseCases.User.Delete")
	}

	if cfg.enableAll || cfg.role {
		check(u.Role.GetListPageData != nil, "UseCases.Role.GetListPageData")
		check(u.Role.Create != nil, "UseCases.Role.Create")
		check(u.Role.Read != nil, "UseCases.Role.Read")
		check(u.Role.Update != nil, "UseCases.Role.Update")
		check(u.Role.Delete != nil, "UseCases.Role.Delete")
		check(u.Role.GetItemPageData != nil, "UseCases.Role.GetItemPageData")
	}

	if cfg.enableAll || cfg.permission {
		check(u.Permission.GetListPageData != nil, "UseCases.Permission.GetListPageData")
		check(u.Permission.Create != nil, "UseCases.Permission.Create")
		check(u.Permission.Read != nil, "UseCases.Permission.Read")
		check(u.Permission.Update != nil, "UseCases.Permission.Update")
		check(u.Permission.Delete != nil, "UseCases.Permission.Delete")
	}

	if cfg.enableAll || cfg.location {
		check(u.Location.GetListPageData != nil, "UseCases.Location.GetListPageData")
		check(u.Location.Create != nil, "UseCases.Location.Create")
		check(u.Location.Read != nil, "UseCases.Location.Read")
		check(u.Location.Update != nil, "UseCases.Location.Update")
		check(u.Location.Delete != nil, "UseCases.Location.Delete")
	}

	if cfg.enableAll || cfg.workspace {
		check(u.Workspace.GetListPageData != nil, "UseCases.Workspace.GetListPageData")
		check(u.Workspace.Create != nil, "UseCases.Workspace.Create")
		check(u.Workspace.Read != nil, "UseCases.Workspace.Read")
		check(u.Workspace.Update != nil, "UseCases.Workspace.Update")
		check(u.Workspace.Delete != nil, "UseCases.Workspace.Delete")
	}

	if cfg.enableAll || cfg.conversation {
		check(u.Conversation.List != nil, "UseCases.Conversation.List")
		check(u.Conversation.Read != nil, "UseCases.Conversation.Read")
		check(u.Conversation.Create != nil, "UseCases.Conversation.Create")
		check(u.Conversation.Post.List != nil, "UseCases.Conversation.Post.List")
		check(u.Conversation.Post.Send != nil, "UseCases.Conversation.Post.Send")
		// Assign, SetStatus, MarkRead are optional (nil-safe: the assign /
		// set-status drawers refuse cleanly and mark-read becomes a no-op).
	}

	if len(missing) > 0 {
		return fmt.Errorf("entydad.Block: incomplete UseCases — missing %v", missing)
	}
	return nil
}

// MustValidate is the FAIL-CLOSED enforcement wrapper around RequireFor
// (architecture-roast burn #1). RequireFor decides WHICH fields gate (required
// vs optional); MustValidate decides HOW a gap fails — it adds posture, not
// policy.
//
// The bare `return RequireFor(...)` from Block() was fail-OPEN by convention: a
// returned error can be dropped (`_ =`, an ignored value, a future caller that
// doesn't check) and the block silently registers an empty feature — the exact
// nil-closure trap the architecture roast (burn #1) named. MustValidate removes
// that escape hatch:
//
//   - In dev/test (running under `go test`, OR ENTYDAD_BLOCK_STRICT truthy) a
//     missing REQUIRED closure PANICS with the full field list. A panic cannot
//     be silently dropped, prints a stack trace at the offending wiring site,
//     and fails the test/CI loudly. This is where a developer wiring a new
//     entity discovers a gap — at their desk, not in prod.
//   - In prod a missing REQUIRED closure logs a screaming FATAL line at the
//     seam (so even a caller that drops the returned error leaves an
//     unmissable log record) AND returns the error so Block() propagates it and
//     NewServiceAdmin halts boot with a clear "domain block failed" message.
//
// OPTIONAL ports (the Client/Conversation optional sub-ops, the un-asserted
// modules: admin, clientTag, supplierTag, paymentTerm, locationArea,
// workspaceUser, workspaceUserRole, supplier, taxRegistration, and the
// dashboard/report closures) are NEVER flagged — that required-vs-optional
// discrimination lives entirely in RequireFor, which only asserts a field when
// its enabling cfg module is on. MustValidate adds posture, not policy: it
// changes HOW a gap fails, not WHICH fields gate.
func (u *UseCases) MustValidate(cfg *blockConfig) error {
	err := u.RequireFor(cfg)
	if err == nil {
		return nil
	}
	if blockStrictMode() {
		// Dev/test: loud, uncatchable-by-accident, stack-traced.
		panic("FATAL: " + err.Error() + " — REQUIRED block wiring is nil. " +
			"Fix the closure assignment in service-admin's buildEntydadUseCases " +
			"(adapters.go) before this reaches prod.")
	}
	// Prod: scream at the seam, then return so boot halts. The log line is the
	// belt to the returned-error's suspenders (a dropped error still screams).
	log.Printf("FATAL: %v — refusing to register entydad modules with a nil "+
		"REQUIRED closure (fail-closed wiring).", err)
	return err
}

// blockStrictMode reports whether the fail-closed wiring guard should PANIC
// (dev/test) rather than return-and-log (prod) on a missing REQUIRED closure.
//
// True when running under `go test` (testing.Testing(), Go 1.21+ — zero env
// coupling, auto-on in every test + CI run) OR when ENTYDAD_BLOCK_STRICT is set
// to an explicit truthy value (the dev escape hatch for `go run` smoke tests).
// The env matching mirrors container.go's authzEnforceEnabled — anything else
// (unset, "", "0", "false") is prod posture.
func blockStrictMode() bool {
	if testing.Testing() {
		return true
	}
	switch os.Getenv("ENTYDAD_BLOCK_STRICT") {
	case "1", "true", "TRUE", "True", "yes", "on":
		return true
	default:
		return false
	}
}
