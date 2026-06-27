package block

import (
	"context"
	"log"

	locationdashboardview "github.com/erniealice/entydad-golang/domain/entity/location/location/dashboard"
	admindashboardview "github.com/erniealice/entydad-golang/service/dashboard/views/admin/dashboard"

	"github.com/erniealice/espyna-golang/consumer"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	admindashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/admin"
	locationdashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/location"
	stmtspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/statements"
)

// ---------------------------------------------------------------------------
// entydad use-case mapper (relocated from app composition/adapters_entydad.go,
// Wave B D2a). Maps espyna's *consumer.UseCases to the entydad block's typed
// *UseCases shape. Imports only espyna consumer + esqyma proto + entydad's own
// domain dashboard view packages — all legal inside the entydad block package.
// ---------------------------------------------------------------------------

// balanceRowsToMap converts the typed `[]*BalanceRow` response shape
// (Wave B P1.E.4 Statements proto) into the legacy `map[string]int64`
// shape consumed by entydad's client/supplier list views. The new
// `ListClient/SupplierBalances` use cases return ordered rows; the map
// shim preserves the legacy view behavior without touching the views.
func balanceRowsToMap(rows []*stmtspb.BalanceRow) map[string]int64 {
	if len(rows) == 0 {
		return nil
	}
	out := make(map[string]int64, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		out[r.GetCounterpartyId()] = r.GetAmountCentavos()
	}
	return out
}

// entydadDBOps is the capability-narrow operations surface the ops-backed
// entydad closures (SetActive / SetStatus) need from the concrete database
// adapter passed via ctx.DB (a *consumer.DatabaseAdapter). It is the typed
// successor to the deleted entydad DataSource/UpdateableSource duck --
// service-admin asserts ctx.DB to it once and binds only the closures that
// genuinely need generic-collection ops (active/status are plain column writes
// the typed proto Update can't express because proto3 omits zero values).
// Signatures MUST match *consumer.DatabaseAdapter exactly.
// 20260612-datasource-typed-path (entydad duck delete).
type entydadDBOps interface {
	Update(ctx context.Context, collection string, id string, data map[string]any) (map[string]any, error)
}

// buildEntydadUseCases maps espyna's *consumer.UseCases to entydad block's
// typed shape. All sub-group wiring is nil-safe.
//
// db is the concrete database adapter (ctx.DB, a *consumer.DatabaseAdapter). It
// backs the generic-collection closures (SetActive, SetStatus) the deleted
// entydad DataSource duck used to serve. When the assertion to entydadDBOps
// fails (mock build, nil DB) those closures are left nil -- entydad degrades
// gracefully (active/status toggles no-op).
// 20260612-datasource-typed-path (entydad duck delete).
func buildEntydadUseCases(uc *consumer.UseCases, db any) *UseCases {
	result := &UseCases{
		GetWorkspaceIDFromCtx: consumer.GetWorkspaceIDFromContext,
	}

	// Assert ctx.DB to the capability-narrow ops surface once. ok==false (mock
	// build / nil DB) -> the ops-backed closures below stay nil (nil-safe at the
	// entydad call sites: toggles no-op).
	if ops, opsOK := db.(entydadDBOps); opsOK {
		// SetActive — sets ONLY the `active` boolean on the named collection.
		// Typed proto3 Update can't clear `active` (omits false bools), so this
		// goes through the generic ops.Update with an explicit {"active": active}
		// map. Mirrors the deleted entydad duck's
		// Update(collection, id, {"active": active}).
		result.SetActive = func(ctx context.Context, collection string, id string, active bool) error {
			_, err := ops.Update(ctx, collection, id, map[string]any{"active": active})
			return err
		}
		// SetStatus — sets the `status` text column AND the `active` boolean
		// together on the named collection. The status→active derivation is done
		// at the entydad call site (client: status != "inactive"; supplier:
		// forced true) and BOTH values arrive here; this just persists them.
		// Mirrors the deleted duck's
		// Update(collection, id, {"status": status, "active": active}).
		result.SetStatus = func(ctx context.Context, collection string, id string, status string, active bool) error {
			_, err := ops.Update(ctx, collection, id, map[string]any{
				"status": status,
				"active": active,
			})
			return err
		}
	} else {
		log.Printf("buildEntydadUseCases: ctx.DB does not satisfy entydadDBOps — SetActive/SetStatus left unwired (active+status toggles no-op)")
	}

	e := uc.Entity
	if e == nil {
		return result
	}

	// Client
	if e.Client != nil {
		result.Client.GetListPageData = e.Client.GetClientListPageData.Execute
		result.Client.Create = e.Client.CreateClient.Execute
		result.Client.Read = e.Client.ReadClient.Execute
		result.Client.Update = e.Client.UpdateClient.Execute
		result.Client.Delete = e.Client.DeleteClient.Execute
		// List — typed ListClients replacing the deleted duck's
		// ListSimple("client") scan behind the conversation inbox's
		// ClientNameByID display-name resolver. nil-safe if the use case
		// is absent. 20260612-datasource-typed-path (entydad duck delete).
		if e.Client.ListClients != nil {
			result.Client.List = e.Client.ListClients.Execute
		}
	}
	if e.ClientCategory != nil {
		result.Client.Category.List = e.ClientCategory.ListClientCategories.Execute
		result.Client.Category.Create = e.ClientCategory.CreateClientCategory.Execute
		result.Client.Category.Delete = e.ClientCategory.DeleteClientCategory.Execute
	}

	// Delegate — workspace-scoped read adapter committed in espyna; view layer is
	// list + CRUD. ListDelegates is optional (not surfaced in the view yet).
	if e.Delegate != nil {
		result.Delegate.GetListPageData = e.Delegate.GetDelegateListPageData.Execute
		result.Delegate.Create = e.Delegate.CreateDelegate.Execute
		result.Delegate.Read = e.Delegate.ReadDelegate.Execute
		result.Delegate.Update = e.Delegate.UpdateDelegate.Execute
		result.Delegate.Delete = e.Delegate.DeleteDelegate.Execute
		if e.Delegate.ListDelegates != nil {
			result.Delegate.List = e.Delegate.ListDelegates.Execute
		}
	}

	// User
	if e.User != nil {
		result.User.GetListPageData = e.User.GetUserListPageData.Execute
		result.User.Create = e.User.CreateUser.Execute
		result.User.Read = e.User.ReadUser.Execute
		result.User.Update = e.User.UpdateUser.Execute
		result.User.Delete = e.User.DeleteUser.Execute
		result.User.List = e.User.ListUsers.Execute
		// Provider-abstracted admin user-lifecycle use cases (design §5/§6).
		// nil-safe: present only when espyna wired them (they are part of the
		// standard user use-case group, so they accompany the CRUD ops above).
		if e.User.DisableUser != nil {
			result.User.Disable = e.User.DisableUser.Execute
		}
		if e.User.EnableUser != nil {
			result.User.Enable = e.User.EnableUser.Execute
		}
		if e.User.AdminResetPassword != nil {
			result.User.ResetPassword = e.User.AdminResetPassword.Execute
		}
	}

	// Role
	if e.Role != nil {
		result.Role.GetListPageData = e.Role.GetRoleListPageData.Execute
		result.Role.Create = e.Role.CreateRole.Execute
		result.Role.Read = e.Role.ReadRole.Execute
		result.Role.Update = e.Role.UpdateRole.Execute
		result.Role.Delete = e.Role.DeleteRole.Execute
		result.Role.GetItemPageData = e.Role.GetRoleItemPageData.Execute
		result.Role.List = e.Role.ListRoles.Execute
	}

	// RolePermission
	if e.RolePermission != nil {
		result.RolePermission.Create = e.RolePermission.CreateRolePermission.Execute
		result.RolePermission.Delete = e.RolePermission.DeleteRolePermission.Execute
	}

	// Permission
	if e.Permission != nil {
		result.Permission.GetListPageData = e.Permission.GetPermissionListPageData.Execute
		result.Permission.Create = e.Permission.CreatePermission.Execute
		result.Permission.Read = e.Permission.ReadPermission.Execute
		result.Permission.Update = e.Permission.UpdatePermission.Execute
		result.Permission.Delete = e.Permission.DeletePermission.Execute
		result.Permission.List = e.Permission.ListPermissions.Execute
	}

	// Location
	if e.Location != nil {
		result.Location.GetListPageData = e.Location.GetLocationListPageData.Execute
		result.Location.Create = e.Location.CreateLocation.Execute
		result.Location.Read = e.Location.ReadLocation.Execute
		result.Location.Update = e.Location.UpdateLocation.Execute
		result.Location.Delete = e.Location.DeleteLocation.Execute
	}

	// LocationArea — typed proto closures replacing the former entydad duck
	// path (block-local CRUDSource.ListSimple/Create/Read/Update/Delete against
	// the "location_area" collection string). The block group field `List` maps
	// from the espyna ListLocationAreas use case; closures are nil-safe at the
	// commerce.go call sites. (20260612-datasource-typed-path W1.E.)
	if e.LocationArea != nil {
		result.LocationArea.List = e.LocationArea.ListLocationAreas.Execute
		result.LocationArea.Create = e.LocationArea.CreateLocationArea.Execute
		result.LocationArea.Read = e.LocationArea.ReadLocationArea.Execute
		result.LocationArea.Update = e.LocationArea.UpdateLocationArea.Execute
		result.LocationArea.Delete = e.LocationArea.DeleteLocationArea.Execute
	}

	// PaymentTerm — typed use cases replacing raw *sql.DB
	// registry.CreateRepository path. The espyna entity aggregate now exposes
	// payment_term use cases directly (E6 provider hardcode fix, 20260614).
	if e.PaymentTerm != nil {
		result.PaymentTerm.GetListPageData = e.PaymentTerm.GetPaymentTermListPageData.Execute
		result.PaymentTerm.CreatePaymentTerm = e.PaymentTerm.CreatePaymentTerm.Execute
		result.PaymentTerm.ReadPaymentTerm = e.PaymentTerm.ReadPaymentTerm.Execute
		result.PaymentTerm.UpdatePaymentTerm = e.PaymentTerm.UpdatePaymentTerm.Execute
		result.PaymentTerm.DeletePaymentTerm = e.PaymentTerm.DeletePaymentTerm.Execute
		if e.PaymentTerm.ListPaymentTerms != nil {
			result.PaymentTerm.ListPaymentTerms = e.PaymentTerm.ListPaymentTerms.Execute
		}
	}

	// Workspace
	if e.Workspace != nil {
		result.Workspace.GetListPageData = e.Workspace.GetWorkspaceListPageData.Execute
		result.Workspace.Create = e.Workspace.CreateWorkspace.Execute
		result.Workspace.Read = e.Workspace.ReadWorkspace.Execute
		result.Workspace.Update = e.Workspace.UpdateWorkspace.Execute
		result.Workspace.Delete = e.Workspace.DeleteWorkspace.Execute
		if e.Workspace.SwitchWorkspace != nil {
			result.Workspace.Switch = e.Workspace.SwitchWorkspace.Execute
		}
	}

	// WorkspaceUser
	if e.WorkspaceUser != nil {
		result.WorkspaceUser.GetListPageData = e.WorkspaceUser.GetWorkspaceUserListPageData.Execute
		result.WorkspaceUser.GetItemPageData = e.WorkspaceUser.GetWorkspaceUserItemPageData.Execute
		result.WorkspaceUser.Create = e.WorkspaceUser.CreateWorkspaceUser.Execute
		result.WorkspaceUser.Delete = e.WorkspaceUser.DeleteWorkspaceUser.Execute
		result.WorkspaceUser.List = e.WorkspaceUser.ListWorkspaceUsers.Execute
	}

	// WorkspaceUserRole
	if e.WorkspaceUserRole != nil {
		result.WorkspaceUserRole.Create = e.WorkspaceUserRole.CreateWorkspaceUserRole.Execute
		result.WorkspaceUserRole.Delete = e.WorkspaceUserRole.DeleteWorkspaceUserRole.Execute
		result.WorkspaceUserRole.GetListPageData = e.WorkspaceUserRole.GetWorkspaceUserRoleListPageData.Execute
	}

	// Supplier
	if e.Supplier != nil {
		result.Supplier.GetListPageData = e.Supplier.GetSupplierListPageData.Execute
		result.Supplier.Create = e.Supplier.CreateSupplier.Execute
		result.Supplier.Read = e.Supplier.ReadSupplier.Execute
		result.Supplier.Update = e.Supplier.UpdateSupplier.Execute
		result.Supplier.Delete = e.Supplier.DeleteSupplier.Execute
	}
	if e.SupplierCategory != nil {
		result.Supplier.Category.List = e.SupplierCategory.ListSupplierCategories.Execute
		result.Supplier.Category.Create = e.SupplierCategory.CreateSupplierCategory.Execute
		result.Supplier.Category.Delete = e.SupplierCategory.DeleteSupplierCategory.Execute
	}

	// Common categories (used by ClientTag and SupplierTag modules)
	if uc.Common != nil && uc.Common.Category != nil {
		result.Category.List = uc.Common.Category.ListCategories.Execute
		result.Category.Create = uc.Common.Category.CreateCategory.Execute
		result.Category.Read = uc.Common.Category.ReadCategory.Execute
		result.Category.Update = uc.Common.Category.UpdateCategory.Execute
		result.Category.Delete = uc.Common.Category.DeleteCategory.Execute
		// GetListPageData — settings list page needs all categories (active AND
		// inactive). ListCategories returns both, so shim it to the simpler
		// func(ctx) ([]*Category, error) shape the entydad block expects.
		// Replaces the raw *sql.DB registry.CreateRepository path (E6 fix).
		listCategoriesUC := uc.Common.Category.ListCategories
		result.Category.GetListPageData = func(ctx context.Context) ([]*commonpb.Category, error) {
			resp, err := listCategoriesUC.Execute(ctx, &commonpb.ListCategoriesRequest{})
			if err != nil {
				return nil, err
			}
			return resp.GetData(), nil
		}
	}

	// Subscription (for client detail subscription tab)
	if uc.Subscription != nil && uc.Subscription.Subscription != nil {
		result.Subscription.List = uc.Subscription.Subscription.ListSubscriptions.Execute
		result.Subscription.GetListPageData = uc.Subscription.Subscription.GetSubscriptionListPageData.Execute
		result.Subscription.CountActiveByClientIDs = uc.Subscription.Subscription.CountActiveByClientIds.Execute
	}

	// Revenue (for client detail revenue tab + revenue run)
	if uc.Revenue != nil && uc.Revenue.Revenue != nil {
		result.Revenue.List = uc.Revenue.Revenue.ListRevenues.Execute
		result.Revenue.ListRevenueRunCandidates = uc.Revenue.Revenue.ListRevenueRunCandidates.Execute
		result.Revenue.GenerateRevenueRun = uc.Revenue.Revenue.GenerateRevenueRun.Execute
	}

	// Collection (for client detail collections tab)
	if uc.Treasury != nil && uc.Treasury.Collection != nil {
		result.Collection.ListByClient = uc.Treasury.Collection.ListByClient.Execute
	}

	// PriceSchedule, PricePlan (for subscription/revenue selectors)
	if uc.Subscription != nil && uc.Subscription.PriceSchedule != nil {
		result.PriceSchedule.List = uc.Subscription.PriceSchedule.ListPriceSchedules.Execute
	}
	if uc.Subscription != nil && uc.Subscription.PricePlan != nil {
		result.PricePlan.List = uc.Subscription.PricePlan.ListPricePlans.Execute
	}

	// PurchaseOrder (for supplier detail purchase orders tab)
	if uc.Expenditure != nil && uc.Expenditure.PurchaseOrder != nil {
		result.PurchaseOrder.List = uc.Expenditure.PurchaseOrder.ListPurchaseOrders.Execute
	}

	// TaxRegistration (for workspace tax registration tab)
	if uc.Tax != nil && uc.Tax.TaxRegistration != nil {
		result.TaxRegistration.List = uc.Tax.TaxRegistration.ListTaxRegistrations.Execute
	}

	// Conversation — secure-messaging surface (Plan-4, 2026-06-03).
	// Bridges espyna's communication use cases to entydad's typed block shape.
	result.Conversation = buildConversationUseCases(uc)

	// 20260521 Wave B P1.E.4 — service-driven statements + balances
	// closures threaded into entydad's typed block UseCases. The legacy
	// `ctx.LedgerReportingSvc` assertion in entydad/block.go is removed;
	// these closures replace it.
	//
	// The map-returning balance shims wrap the typed `[]*BalanceRow`
	// closure to preserve the legacy `map[string]int64` shape consumed by
	// entydad client/supplier list views (per the documented entydad
	// balances shape decision in this commit).
	if uc.Service != nil && uc.Service.Reporting != nil && uc.Service.Reporting.Statements != nil {
		st := uc.Service.Reporting.Statements
		if st.GetClientStatement != nil {
			result.Reports.Statements.GetClientStatement = st.GetClientStatement.Execute
		}
		if st.GetSupplierStatement != nil {
			result.Reports.Statements.GetSupplierStatement = st.GetSupplierStatement.Execute
		}
		if st.ListClientBalances != nil {
			listClientBalances := st.ListClientBalances
			result.Reports.Statements.ListClientBalances = listClientBalances.Execute
			result.Reports.Statements.ListClientBalancesAsMap = func(ctx context.Context) (map[string]int64, error) {
				resp, err := listClientBalances.Execute(ctx, &stmtspb.ListClientBalancesRequest{})
				if err != nil {
					return nil, err
				}
				return balanceRowsToMap(resp.GetBalances()), nil
			}
		}
		if st.ListSupplierBalances != nil {
			listSupplierBalances := st.ListSupplierBalances
			result.Reports.Statements.ListSupplierBalances = listSupplierBalances.Execute
			result.Reports.Statements.ListSupplierBalancesAsMap = func(ctx context.Context) (map[string]int64, error) {
				resp, err := listSupplierBalances.Execute(ctx, &stmtspb.ListSupplierBalancesRequest{})
				if err != nil {
					return nil, err
				}
				return balanceRowsToMap(resp.GetBalances()), nil
			}
		}
	}

	// Dashboard closures — call use-case Execute directly with proto Request.
	// Returns nil when use case is not wired.
	//
	// Location dashboard relocated to service.Dashboard.Location per Wave B P1.C.2
	// (Q-SDM-DASHBOARD-LAYOUT / Q-SDM-DASHBOARD-DOWNSTREAM, 2026-05-20).
	if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Location != nil && uc.Service.Dashboard.Location.GetLocationDashboard != nil {
		locationDash := uc.Service.Dashboard.Location.GetLocationDashboard
		result.GetLocationDashboardPageData = func(ctx context.Context) (*locationdashboardview.LocationDashboardData, error) {
			wsID := consumer.GetWorkspaceIDFromContext(ctx)
			resp, err := locationDash.Execute(ctx, &locationdashpb.GetLocationDashboardRequest{
				WorkspaceId: wsID,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			out := &locationdashboardview.LocationDashboardData{
				TotalLocations:    resp.GetStats().GetTotalLocations(),
				ActiveLocations:   resp.GetStats().GetActiveLocations(),
				RegionsCount:      resp.GetStats().GetRegionsCount(),
				AreasCount:        resp.GetStats().GetAreasCount(),
				LocationsByRegion: resp.GetLocationsByRegion(),
				RecentLocations:   resp.GetRecentLocations(),
			}
			for _, a := range resp.GetTopAreas() {
				out.TopAreas = append(out.TopAreas, locationdashboardview.LocationAreaCount{
					LocationAreaID:   a.GetLocationAreaId(),
					LocationAreaName: a.GetLocationAreaName(),
					LocationCount:    a.GetLocationCount(),
				})
			}
			return out, nil
		}
	}

	// Admin dashboard relocated to service.Dashboard.Admin per Wave B P1.C.1
	// (Q-SDM-DASHBOARD-LAYOUT / Q-SDM-DASHBOARD-DOWNSTREAM, 2026-05-20).
	if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Admin != nil && uc.Service.Dashboard.Admin.GetAdminDashboard != nil {
		adminDash := uc.Service.Dashboard.Admin.GetAdminDashboard
		result.GetAdminDashboardPageData = func(ctx context.Context) (*admindashboardview.AdminDashboardData, error) {
			wsID := consumer.GetWorkspaceIDFromContext(ctx)
			resp, err := adminDash.Execute(ctx, &admindashpb.GetAdminDashboardRequest{
				WorkspaceId: wsID,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			out := &admindashboardview.AdminDashboardData{
				WorkspaceUsers:      resp.GetStats().GetWorkspaceUsers(),
				Roles:               resp.GetStats().GetRoles(),
				Permissions:         resp.GetStats().GetPermissions(),
				RecentRoleChanges7d: resp.GetStats().GetRecentRoleChanges7D(),
				UsersPerRole:        resp.GetUsersPerRole(),
				RecentAssignments:   resp.GetRecentAssignments(),
			}
			for _, r := range resp.GetTopRolesByPerms() {
				out.TopRolesByPerms = append(out.TopRolesByPerms, admindashboardview.RolePermissionCount{
					RoleID:          r.GetRoleId(),
					RoleName:        r.GetRoleName(),
					PermissionCount: r.GetPermissionCount(),
				})
			}
			return out, nil
		}
	}

	return result
}

// buildConversationUseCases bridges espyna's communication use cases to the
// entydad block's typed ConversationUseCases shape (Plan-4 secure messaging).
//
// All closures are the REAL espyna Execute signatures (verified against
// packages/espyna-golang/.../communication/*/usecases.go):
//   - Assign / SetStatus consume *conversationpb.UpdateConversationRequest
//     (the espyna use cases dispatch on the mutated field — there is no
//     distinct AssignConversationRequest / SetConversationStatusRequest).
//   - Send consumes *conversationpostpb.CreateConversationPostRequest.
//   - MarkRead consumes *conversationreadreceiptpb.CreateConversationReadReceiptRequest.
//
// Nil-safe everywhere: when espyna has not wired the communication aggregate
// (e.g. mock builds) the whole struct stays zero and entydad's RequireFor
// would reject WithConversation() at boot — so WithConversation() is only
// added to the block options when these use cases are present (container.go).
func buildConversationUseCases(uc *consumer.UseCases) ConversationUseCases {
	if uc == nil || uc.Communication == nil {
		return ConversationUseCases{}
	}
	comm := uc.Communication
	out := ConversationUseCases{}

	if comm.Conversation != nil {
		conv := comm.Conversation
		if conv.ListConversations != nil {
			out.List = conv.ListConversations.Execute
		}
		if conv.ReadConversation != nil {
			out.Read = conv.ReadConversation.Execute
		}
		if conv.CreateConversation != nil {
			out.Create = conv.CreateConversation.Execute
		}
		if conv.AssignConversation != nil {
			out.Assign = conv.AssignConversation.Execute
		}
		if conv.SetConversationStatus != nil {
			out.SetStatus = conv.SetConversationStatus.Execute
		}
	}

	if comm.ConversationPost != nil {
		post := comm.ConversationPost
		if post.ListConversationPosts != nil {
			out.Post.List = post.ListConversationPosts.Execute
		}
		if post.SendConversationPost != nil {
			out.Post.Send = post.SendConversationPost.Execute
		}
	}

	if comm.ConversationReadReceipt != nil {
		rr := comm.ConversationReadReceipt
		if rr.MarkConversationRead != nil {
			out.Receipt.MarkRead = rr.MarkConversationRead.Execute
		}
	}

	return out
}
