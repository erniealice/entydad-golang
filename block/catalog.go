// Package block — composition-v2 catalog for the entydad entity domain.
//
// Each binder function constructs a compose.Unit with a typed Mount closure
// that wires the entity module's Deps from UseCases + Infra, then calls
// NewXxxModule(deps).RegisterRoutes(mc.Routes). The structure mirrors the
// three sub-context files in this package (identity.go, party.go, commerce.go)
// plus the tax domain.
//
// Aggregator: AllUnits(uc, infra) returns the full ordered slice for the
// compose engine, mirroring the registration order in Block().
package block

import (
	"context"
	"fmt"
	"log"

	entity "github.com/erniealice/entydad-golang/domain/entity"
	commerce "github.com/erniealice/entydad-golang/domain/entity/commerce"
	entitypaymentterm "github.com/erniealice/entydad-golang/domain/entity/commerce/payment_term"
	identity "github.com/erniealice/entydad-golang/domain/entity/identity"
	entitypermission "github.com/erniealice/entydad-golang/domain/entity/identity/permission"
	entityrole "github.com/erniealice/entydad-golang/domain/entity/identity/role"
	roleusers "github.com/erniealice/entydad-golang/domain/entity/identity/role/users"
	entityuser "github.com/erniealice/entydad-golang/domain/entity/identity/user"
	entityworkspace "github.com/erniealice/entydad-golang/domain/entity/identity/workspace"
	workspaceaction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/action"
	entityworkspaceuser "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user"
	entityworkspaceuserrole "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user_role"
	location "github.com/erniealice/entydad-golang/domain/entity/location"
	entitylocation "github.com/erniealice/entydad-golang/domain/entity/location/location"
	locationaction "github.com/erniealice/entydad-golang/domain/entity/location/location/action"
	entitylocationarea "github.com/erniealice/entydad-golang/domain/entity/location/location_area"
	locationareaaction "github.com/erniealice/entydad-golang/domain/entity/location/location_area/action"
	locationarealist "github.com/erniealice/entydad-golang/domain/entity/location/location_area/list"
	party "github.com/erniealice/entydad-golang/domain/entity/party"
	entityclient "github.com/erniealice/entydad-golang/domain/entity/party/client"
	clientdetail "github.com/erniealice/entydad-golang/domain/entity/party/client/detail"
	entityclienttag "github.com/erniealice/entydad-golang/domain/entity/party/client_tag"
	entitysupplier "github.com/erniealice/entydad-golang/domain/entity/party/supplier"
	entitysuppliertag "github.com/erniealice/entydad-golang/domain/entity/party/supplier_tag"
	tax "github.com/erniealice/entydad-golang/domain/tax"
	taxregistration "github.com/erniealice/entydad-golang/domain/tax/tax_registration"
	"github.com/erniealice/entydad-golang/service/auth"
	"github.com/erniealice/espyna-golang/consumer/compose"
	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	locationareapb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location_area"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"google.golang.org/protobuf/proto"
)

// ---------------------------------------------------------------------------
// Party sub-context
// ---------------------------------------------------------------------------

func ClientUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entityclient.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entityclient.Routes)
		l := u.Labels.(*entityclient.Labels)

		// Resolve the timezone search URL from the user unit if available.
		searchTimezonesURL := ""
		if ur, ok := compose.RoutesOf[*entityuser.Routes](mc, "entity.user"); ok {
			searchTimezonesURL = ur.SearchTimezonesURL
		}

		deps := &party.ClientModuleDeps{
			Routes:               *r,
			SearchTimezonesURL:   searchTimezonesURL,
			CommonLabels:         mc.Common,
			SharedLabels:         infra.SharedLabels,
			Labels:               *l,
			DashboardLabels:      infra.ClientDashboardLabels,
			DashboardTitleLabels: infra.DashboardTitleLabels,
			TableLabels:          mc.Table,
			GetListPageData:      uc.Client.GetListPageData,
			GetInUseIDs:          infra.RefChecker.GetClientInUseIDs,
			CreateClient:         uc.Client.Create,
			ReadClient:           uc.Client.Read,
			UpdateClient:         uc.Client.Update,
			DeleteClient:         uc.Client.Delete,
			SetStatus: setStatusClosure(uc, "client", func(status string) bool {
				return status != "inactive"
			}),
			ListPaymentTerms: func(fctx context.Context) ([]*party.ClientPaymentTermOption, error) {
				rows := listPaymentTermRows(fctx, uc)
				opts := make([]*party.ClientPaymentTermOption, 0, len(rows))
				for _, row := range rows {
					id := row.GetId()
					if id == "" {
						continue
					}
					scope := row.GetEntityScope()
					if scope != "client" && scope != "both" {
						continue
					}
					opts = append(opts, &party.ClientPaymentTermOption{Id: id, Name: row.GetName()})
				}
				return opts, nil
			},
			SubscriptionAddURL:               infra.SubscriptionRoutes.AddURL,
			SubscriptionDetailURL:            infra.SubscriptionRoutes.DetailURL,
			SubscriptionUnderClientDetailURL: infra.SubscriptionRoutes.UnderClientDetailURL,
			SubscriptionEditURL:              infra.SubscriptionRoutes.EditURL,
			SubscriptionDeleteURL:            infra.SubscriptionRoutes.DeleteURL,
			UploadFile:                       infra.UploadFile,
			ListAttachments:                  infra.ListAttachments,
			CreateAttachment:                 infra.CreateAttachment,
			DeleteAttachment:                 infra.DeleteAttachment,
			NewID:                            infra.NewAttachmentID,
		}
		if uc.Category.List != nil {
			deps.ListCategories = uc.Category.List
		}
		if uc.Client.Category.List != nil {
			deps.ListClientCategories = uc.Client.Category.List
			deps.CreateClientCategory = uc.Client.Category.Create
			deps.DeleteClientCategory = uc.Client.Category.Delete
		}
		if uc.Subscription.List != nil {
			deps.ListSubscriptions = uc.Subscription.List
			deps.GetSubscriptionListPageData = uc.Subscription.GetListPageData
		}
		if uc.PriceSchedule.List != nil {
			listPriceSchedules := uc.PriceSchedule.List
			var listPricePlans func(context.Context, *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
			if uc.PricePlan.List != nil {
				listPricePlans = uc.PricePlan.List
			}
			psDetailURL := infra.PriceScheduleRoutes.DetailURL
			deps.ListClientPriceSchedules = func(fctx context.Context, clientID string) ([]clientdetail.ClientPriceScheduleRow, error) {
				resp, err := listPriceSchedules(fctx, &priceschedulepb.ListPriceSchedulesRequest{})
				if err != nil {
					return nil, err
				}
				planCountBySchedule := map[string]int{}
				if listPricePlans != nil {
					ppResp, ppErr := listPricePlans(fctx, &priceplanpb.ListPricePlansRequest{})
					if ppErr != nil {
						log.Printf("entydad catalog: failed to list price plans for plan-count column: %v", ppErr)
					} else {
						for _, pp := range ppResp.GetData() {
							if sid := pp.GetPriceScheduleId(); sid != "" {
								planCountBySchedule[sid]++
							}
						}
					}
				}
				tz := pyezatypes.LocationFromContext(fctx)
				var rows []clientdetail.ClientPriceScheduleRow
				for _, s := range resp.GetData() {
					if s.GetClientId() != clientID {
						continue
					}
					detailURL := ""
					if psDetailURL != "" {
						detailURL = route.ResolveURL(psDetailURL, "id", s.GetId())
					}
					startDate, startTime := pyezatypes.FormatTimestampSplitInTZ(s.GetDateTimeStart(), tz)
					endDate, endTime := pyezatypes.FormatTimestampSplitInTZ(s.GetDateTimeEnd(), tz)
					rows = append(rows, clientdetail.ClientPriceScheduleRow{
						ID:            s.GetId(),
						Name:          s.GetName(),
						DateStartDate: startDate,
						DateStartTime: startTime,
						DateEndDate:   endDate,
						DateEndTime:   endTime,
						PlanCount:     planCountBySchedule[s.GetId()],
						DetailURL:     detailURL,
					})
				}
				return rows, nil
			}
			deps.PriceScheduleAddURL = infra.PriceScheduleRoutes.AddURL
		}
		if uc.Workspace.Read != nil {
			readWorkspace := uc.Workspace.Read
			wsID := getDefaultWorkspaceID()
			deps.GetFunctionalCurrency = func(fctx context.Context) string {
				resp, err := readWorkspace(fctx, &workspacepb.ReadWorkspaceRequest{
					Data: &workspacepb.Workspace{Id: wsID},
				})
				if err != nil {
					return ""
				}
				if data := resp.GetData(); len(data) > 0 {
					return data[0].GetFunctionalCurrency()
				}
				return ""
			}
		}
		if uc.Reports.Statements.ListClientBalancesAsMap != nil {
			deps.GetClientBalances = uc.Reports.Statements.ListClientBalancesAsMap
		}
		if uc.Reports.Statements.GetClientStatement != nil {
			getCSvc := uc.Reports.Statements.GetClientStatement
			deps.GetClientStatement = func(fctx context.Context, req *clientstmtpb.ClientStatementRequest) (*clientstmtpb.ClientStatementResponse, error) {
				resp, err := getCSvc(fctx, translateClientStatementReq(req))
				if err != nil {
					return nil, err
				}
				return translateClientStatementResp(resp), nil
			}
		}
		if uc.Subscription.CountActiveByClientIDs != nil {
			countActive := uc.Subscription.CountActiveByClientIDs
			deps.GetActiveSubscriptionCounts = func(fctx context.Context) (map[string]int32, error) {
				resp, err := countActive(fctx, &subscriptionpb.CountActiveByClientIdsRequest{})
				if err != nil {
					return nil, err
				}
				return resp.GetCounts(), nil
			}
		}
		if uc.Revenue.List != nil {
			listRevenues := uc.Revenue.List
			deps.ListRevenuesByClient = func(fctx context.Context, clientID string) ([]*revenuepb.Revenue, error) {
				resp, err := listRevenues(fctx, &revenuepb.ListRevenuesRequest{
					Filters: &categorypb.FilterRequest{
						Filters: []*categorypb.TypedFilter{
							{
								Field: "client_id",
								FilterType: &categorypb.TypedFilter_StringFilter{
									StringFilter: &categorypb.StringFilter{
										Value:         clientID,
										Operator:      categorypb.StringOperator_STRING_EQUALS,
										CaseSensitive: true,
									},
								},
							},
						},
					},
				})
				if err != nil {
					return nil, err
				}
				return resp.GetData(), nil
			}
		}
		if uc.Collection.ListByClient != nil {
			listByClient := uc.Collection.ListByClient
			deps.ListCollectionsByClient = func(fctx context.Context, clientID string) ([]*collectionpb.Collection, error) {
				resp, err := listByClient(fctx, &collectionpb.ListByClientRequest{ClientId: clientID})
				if err != nil {
					return nil, err
				}
				return resp.GetData(), nil
			}
		}
		if uc.Revenue.ListRevenueRunCandidates != nil && uc.Revenue.GenerateRevenueRun != nil {
			listCandidates := uc.Revenue.ListRevenueRunCandidates
			generateRun := uc.Revenue.GenerateRevenueRun
			deps.ListRevenueRunCandidates = func(fctx context.Context, scope clientdetail.RevenueRunScope) ([]clientdetail.RevenueRunCandidate, string, error) {
				resp, err := listCandidates(fctx, &revrunpb.ListRevenueRunCandidatesRequest{
					Scope: &revrunpb.RevenueRunScope{
						WorkspaceId:    proto.String(scope.WorkspaceID),
						ClientId:       proto.String(scope.ClientID),
						SubscriptionId: proto.String(scope.SubscriptionID),
						AsOfDate:       proto.String(scope.AsOfDate),
					},
					IncludeAdvanceCollections: proto.Bool(true),
				})
				if err != nil || resp == nil {
					return nil, "", err
				}
				out := make([]clientdetail.RevenueRunCandidate, 0, len(resp.GetData()))
				for _, c := range resp.GetData() {
					amtDisplay := fmt.Sprintf("%.2f", float64(c.GetAmount())/100)
					out = append(out, clientdetail.RevenueRunCandidate{
						SubscriptionID:                 c.GetSubscriptionId(),
						SubscriptionName:               c.GetSubscriptionName(),
						ClientID:                       c.GetClientId(),
						ClientName:                     c.GetClientName(),
						PlanName:                       c.GetPlanName(),
						BillingCycleLabel:              c.GetBillingCycleLabel(),
						Currency:                       c.GetCurrency(),
						PeriodStart:                    c.GetPeriodStart(),
						PeriodEnd:                      c.GetPeriodEnd(),
						PeriodLabel:                    c.GetPeriodLabel(),
						PeriodMarker:                   c.GetPeriodMarker(),
						Amount:                         c.GetAmount(),
						AmountDisplay:                  amtDisplay,
						LineItemCount:                  int(c.GetLineItemCount()),
						Eligible:                       c.GetEligible(),
						BlockerReason:                  c.GetBlockerReason(),
						SourceKind:                     c.GetSourceKind().String(),
						AdvanceCollectionID:            c.GetAdvanceCollectionId(),
						SuppressingAdvanceCollectionID: c.GetSuppressingAdvanceCollectionId(),
					})
				}
				return out, resp.GetNextCursor(), nil
			}
			deps.GenerateRevenueRun = func(fctx context.Context, scope clientdetail.RevenueRunScope, selections clientdetail.RevenueRunSelections) (*clientdetail.RevenueRunResult, error) {
				protoSelections := &revrunpb.RevenueRunSelections{
					FilterToken: proto.String(selections.FilterToken),
				}
				for _, s := range selections.ExplicitList {
					protoSel := &revrunpb.SelectedRevenueRunCandidate{
						SubscriptionId: s.SubscriptionID,
						PeriodStart:    s.PeriodStart,
						PeriodEnd:      s.PeriodEnd,
						PeriodMarker:   s.PeriodMarker,
					}
					if s.SourceKind == "REVENUE_RUN_SOURCE_KIND_ADVANCE_COLLECTION" || s.SourceKind == "ADVANCE_COLLECTION" {
						protoSel.SourceKind = revrunpb.RevenueRunSourceKind_REVENUE_RUN_SOURCE_KIND_ADVANCE_COLLECTION
						if s.AdvanceCollectionID != "" {
							protoSel.AdvanceCollectionId = proto.String(s.AdvanceCollectionID)
						}
					}
					protoSelections.ExplicitList = append(protoSelections.ExplicitList, protoSel)
				}
				resp, err := generateRun(fctx, &revrunpb.GenerateRevenueRunRequest{
					Scope: &revrunpb.RevenueRunScope{
						WorkspaceId:    proto.String(scope.WorkspaceID),
						ClientId:       proto.String(scope.ClientID),
						SubscriptionId: proto.String(scope.SubscriptionID),
						AsOfDate:       proto.String(scope.AsOfDate),
					},
					Selections: protoSelections,
				})
				if err != nil || resp == nil {
					return nil, err
				}
				run := resp.GetRun()
				runID, runStatus := "", ""
				if run != nil {
					runID = run.GetId()
					runStatus = run.GetStatus().String()
				}
				var created, skipped, errored int32
				for _, a := range resp.GetAttempts() {
					switch a.GetOutcome().String() {
					case "REVENUE_RUN_ATTEMPT_OUTCOME_CREATED":
						created++
					case "REVENUE_RUN_ATTEMPT_OUTCOME_SKIPPED":
						skipped++
					default:
						errored++
					}
				}
				return &clientdetail.RevenueRunResult{
					RunID:   runID,
					Status:  runStatus,
					Created: created,
					Skipped: skipped,
					Errored: errored,
				}, nil
			}
		}
		party.NewClientModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func SupplierUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entitysupplier.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entitysupplier.Routes)
		l := u.Labels.(*entitysupplier.Labels)

		// Resolve timezone search URL from the user unit if available.
		searchTimezonesURL := ""
		if ur, ok := compose.RoutesOf[*entityuser.Routes](mc, "entity.user"); ok {
			searchTimezonesURL = ur.SearchTimezonesURL
		}
		// Resolve supplier_tag routes for dashboard quick-action links.
		supplierTagRoutes := entitysuppliertag.DefaultRoutes()
		if str, ok := compose.RoutesOf[*entitysuppliertag.Routes](mc, "entity.supplier_tag"); ok {
			supplierTagRoutes = *str
		}

		deps := &party.SupplierModuleDeps{
			Routes:               *r,
			SupplierTagRoutes:    supplierTagRoutes,
			SearchTimezonesURL:   searchTimezonesURL,
			CommonLabels:         mc.Common,
			SharedLabels:         infra.SharedLabels,
			Labels:               *l,
			DashboardLabels:      infra.SupplierDashboardLabels,
			DashboardTitleLabels: infra.DashboardTitleLabels,
			TableLabels:          mc.Table,
			GetInUseIDs:          infra.RefChecker.GetSupplierInUseIDs,
			SetStatus: setStatusClosure(uc, "supplier", func(string) bool {
				return true
			}),
			ListPaymentTerms: func(fctx context.Context) ([]*party.SupplierPaymentTermOption, error) {
				rows := listPaymentTermRows(fctx, uc)
				opts := make([]*party.SupplierPaymentTermOption, 0, len(rows))
				for _, row := range rows {
					id := row.GetId()
					if id == "" {
						continue
					}
					scope := row.GetEntityScope()
					if scope != "supplier" && scope != "both" {
						continue
					}
					opts = append(opts, &party.SupplierPaymentTermOption{Id: id, Name: row.GetName()})
				}
				return opts, nil
			},
			UploadFile:       infra.UploadFile,
			ListAttachments:  infra.ListAttachments,
			CreateAttachment: infra.CreateAttachment,
			DeleteAttachment: infra.DeleteAttachment,
			NewID:            infra.NewAttachmentID,
		}
		if uc.Supplier.GetListPageData != nil {
			deps.GetListPageData = uc.Supplier.GetListPageData
			deps.CreateSupplier = uc.Supplier.Create
			deps.ReadSupplier = uc.Supplier.Read
			deps.UpdateSupplier = uc.Supplier.Update
			deps.DeleteSupplier = uc.Supplier.Delete
		}
		if uc.PurchaseOrder.List != nil {
			deps.ListPurchaseOrders = uc.PurchaseOrder.List
		}
		if uc.Reports.Statements.GetSupplierStatement != nil {
			getSvcSupp := uc.Reports.Statements.GetSupplierStatement
			deps.GetSupplierStatement = func(fctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error) {
				resp, err := getSvcSupp(fctx, translateSupplierStatementReq(req))
				if err != nil {
					return nil, err
				}
				return translateSupplierStatementResp(resp), nil
			}
		}
		if uc.Reports.Statements.ListSupplierBalancesAsMap != nil {
			deps.GetSupplierBalances = uc.Reports.Statements.ListSupplierBalancesAsMap
		}
		if uc.Category.List != nil {
			deps.ListCategories = uc.Category.List
		}
		if uc.Supplier.Category.List != nil {
			deps.ListSupplierCategories = uc.Supplier.Category.List
			deps.CreateSupplierCategory = uc.Supplier.Category.Create
			deps.DeleteSupplierCategory = uc.Supplier.Category.Delete
		}
		party.NewSupplierModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func ClientTagUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entityclienttag.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entityclienttag.Routes)
		l := u.Labels.(*entityclienttag.Labels)

		deps := &party.ClientTagModuleDeps{
			Routes:            *r,
			Labels:            *l,
			SharedLabels:      infra.SharedLabels,
			CommonLabels:      mc.Common,
			TableLabels:       mc.Table,
			GetInUseIDs:       infra.RefChecker.GetCategoryInUseIDs,
			SetCategoryActive: setActiveClosure(uc, "category"),
		}
		if uc.Category.List != nil {
			deps.ListCategories = uc.Category.List
			deps.CreateCategory = uc.Category.Create
			deps.ReadCategory = uc.Category.Read
			deps.UpdateCategory = uc.Category.Update
			deps.DeleteCategory = uc.Category.Delete
		}
		if uc.Category.GetListPageData != nil {
			deps.GetCategoryListPageData = uc.Category.GetListPageData
		}
		if uc.Client.Category.List != nil {
			deps.ListClientCategories = uc.Client.Category.List
		}
		party.NewClientTagModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func SupplierTagUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entitysuppliertag.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entitysuppliertag.Routes)
		l := u.Labels.(*entitysuppliertag.Labels)

		deps := &party.SupplierTagModuleDeps{
			Routes:            *r,
			Labels:            *l,
			SharedLabels:      infra.SharedLabels,
			CommonLabels:      mc.Common,
			TableLabels:       mc.Table,
			GetInUseIDs:       infra.RefChecker.GetCategoryInUseIDs,
			SetCategoryActive: setActiveClosure(uc, "category"),
		}
		if uc.Category.List != nil {
			deps.ListCategories = uc.Category.List
			deps.CreateCategory = uc.Category.Create
			deps.ReadCategory = uc.Category.Read
			deps.UpdateCategory = uc.Category.Update
			deps.DeleteCategory = uc.Category.Delete
		}
		if uc.Category.GetListPageData != nil {
			deps.GetCategoryListPageData = uc.Category.GetListPageData
		}
		if uc.Supplier.Category.List != nil {
			deps.ListSupplierCategories = uc.Supplier.Category.List
		}
		party.NewSupplierTagModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Identity sub-context
// ---------------------------------------------------------------------------

func UserUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entityuser.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entityuser.Routes)
		l := u.Labels.(*entityuser.Labels)

		identity.NewUserModule(&identity.UserModuleDeps{
			Routes:                       *r,
			CommonLabels:                 mc.Common,
			SharedLabels:                 infra.SharedLabels,
			Labels:                       *l,
			DashboardLabels:              infra.UserDashboardLabels,
			DashboardTitleLabels:         infra.DashboardTitleLabels,
			UserRoleLabels:               infra.UserRoleLabels,
			TableLabels:                  mc.Table,
			GetListPageData:              uc.User.GetListPageData,
			GetUserWorkspacesMap:         infra.GetUserWorkspacesMap,
			CreateUser:                   uc.User.Create,
			ReadUser:                     uc.User.Read,
			UpdateUser:                   uc.User.Update,
			DeleteUser:                   uc.User.Delete,
			SetActive:                    setActiveClosure(uc, "user"),
			CreateWorkspaceUser:          uc.WorkspaceUser.Create,
			ListWorkspaceUsers:           uc.WorkspaceUser.List,
			GetWorkspaceUserItemPageData: uc.WorkspaceUser.GetItemPageData,
			DefaultWorkspaceID:           getDefaultWorkspaceID(),
			CreateWorkspaceUserRole:      uc.WorkspaceUserRole.Create,
			DeleteWorkspaceUserRole:      uc.WorkspaceUserRole.Delete,
			ListRoles:                    uc.Role.List,
			GetDashboardData:             infra.GetDashboardData,
			HashPassword:                 infra.HashPassword,
			UploadFile:                   infra.UploadFile,
			ListAttachments:              infra.ListAttachments,
			CreateAttachment:             infra.CreateAttachment,
			DeleteAttachment:             infra.DeleteAttachment,
			NewID:                        infra.NewAttachmentID,
		}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func RoleUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entityrole.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entityrole.Routes)
		l := u.Labels.(*entityrole.Labels)

		identity.NewRoleModule(&identity.RoleModuleDeps{
			Routes:                  *r,
			CommonLabels:            mc.Common,
			SharedLabels:            infra.SharedLabels,
			Labels:                  *l,
			RolePermissionLabels:    infra.RolePermissionLabels,
			RoleUserLabels:          infra.RoleUserLabels,
			TableLabels:             mc.Table,
			GetListPageData:         uc.Role.GetListPageData,
			GetInUseIDs:             infra.RefChecker.GetRoleInUseIDs,
			CreateRole:              uc.Role.Create,
			ReadRole:                uc.Role.Read,
			UpdateRole:              uc.Role.Update,
			DeleteRole:              uc.Role.Delete,
			SetActive:               setActiveClosure(uc, "role"),
			GetItemPageData:         uc.Role.GetItemPageData,
			CreateRolePermission:    uc.RolePermission.Create,
			DeleteRolePermission:    uc.RolePermission.Delete,
			ListPermissions:         uc.Permission.List,
			GetUsersByRoleID:        infra.GetUsersByRoleID,
			ListWorkspaceUsers:      uc.WorkspaceUser.List,
			CreateWorkspaceUserRole: uc.WorkspaceUserRole.Create,
			DeleteWorkspaceUserRole: uc.WorkspaceUserRole.Delete,
			UploadFile:              infra.UploadFile,
			ListAttachments:         infra.ListAttachments,
			CreateAttachment:        infra.CreateAttachment,
			DeleteAttachment:        infra.DeleteAttachment,
			NewID:                   infra.NewAttachmentID,
		}).RegisterRoutes(mc.Routes)

		// Role-User search (raw http.HandlerFunc endpoint).
		compose.HandleFunc(mc.Routes, "GET", r.UsersSearchURL, roleusers.NewSearchUsersAction(&roleusers.SearchDeps{
			ListWorkspaceUsers: uc.WorkspaceUser.List,
		}))
		return nil
	}
	return u
}

func PermissionUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entitypermission.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entitypermission.Routes)
		l := u.Labels.(*entitypermission.Labels)

		identity.NewPermissionModule(&identity.PermissionModuleDeps{
			Routes:           *r,
			CommonLabels:     mc.Common,
			SharedLabels:     infra.SharedLabels,
			Labels:           *l,
			TableLabels:      mc.Table,
			GetListPageData:  uc.Permission.GetListPageData,
			CreatePermission: uc.Permission.Create,
			ReadPermission:   uc.Permission.Read,
			UpdatePermission: uc.Permission.Update,
			DeletePermission: uc.Permission.Delete,
			SetActive:        setActiveClosure(uc, "permission"),
		}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func WorkspaceUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entityworkspace.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entityworkspace.Routes)
		l := u.Labels.(*entityworkspace.Labels)

		wsMod := &identity.WorkspaceModuleDeps{
			Routes:                 *r,
			CommonLabels:           mc.Common,
			SharedLabels:           infra.SharedLabels,
			Labels:                 *l,
			TableLabels:            mc.Table,
			GetListPageData:        uc.Workspace.GetListPageData,
			CreateWorkspace:        uc.Workspace.Create,
			ReadWorkspace:          uc.Workspace.Read,
			UpdateWorkspace:        uc.Workspace.Update,
			DeleteWorkspace:        uc.Workspace.Delete,
			SetActive:              setActiveClosure(uc, "workspace"),
			WorkspaceUserDetailURL: entity.WorkspaceUserDetailURL,
			WorkspaceUserAddURL:    entity.WorkspaceUserAddURL,
			UploadFile:             infra.UploadFile,
			ListAttachments:        infra.ListAttachments,
			CreateAttachment:       infra.CreateAttachment,
			DeleteAttachment:       infra.DeleteAttachment,
			NewID:                  infra.NewAttachmentID,
		}
		if uc.WorkspaceUser.GetListPageData != nil {
			wsMod.GetWorkspaceUserListPageData = uc.WorkspaceUser.GetListPageData
		}
		identity.NewWorkspaceModule(wsMod).RegisterRoutes(mc.Routes)

		// Switch-workspace raw POST handler (HandleFunc).
		if uc.Workspace.Switch != nil || infra.SecureSwitch != nil {
			compose.HandleFunc(mc.Routes, "POST", r.SwitchURL, workspaceaction.NewSwitchWorkspaceHandler(&workspaceaction.SwitchWorkspaceDeps{
				SecureSwitch:          infra.SecureSwitch,
				ResolveUserID:         infra.SecureSwitchResolveUser,
				SetSessionCookie:      infra.SecureSwitchSetCookie,
				SwitchWorkspace:       uc.Workspace.Switch,
				HomeURLForWorkspaceID: infra.HomeURLForWorkspaceID,
			}))
		}
		return nil
	}
	return u
}

func WorkspaceUserUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entityworkspaceuser.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entityworkspaceuser.Routes)
		l := u.Labels.(*entityworkspaceuser.Labels)

		if uc.WorkspaceUser.GetListPageData == nil {
			log.Println("entydad catalog: warning: workspace_user use cases not initialized — workspace_user detail routes will be unavailable")
			return nil
		}
		wuMod := &identity.WorkspaceUserModuleDeps{
			Routes:                       *r,
			WorkspaceDetailURL:           entity.WorkspaceDetailURL,
			CommonLabels:                 mc.Common,
			Labels:                       *l,
			TableLabels:                  mc.Table,
			GetListPageData:              uc.WorkspaceUser.GetListPageData,
			GetWorkspaceUserItemPageData: uc.WorkspaceUser.GetItemPageData,
			CreateWorkspaceUser:          uc.WorkspaceUser.Create,
			DeleteWorkspaceUser:          uc.WorkspaceUser.Delete,
			SetWorkspaceUserActive:       setActiveClosure(uc, "workspace_user"),
			WorkspaceUserRoleAddURL:      entity.WorkspaceUserRoleAddURL,
			WorkspaceUserRoleDeleteURL:   entity.WorkspaceUserRoleDeleteURL,
			UploadFile:                   infra.UploadFile,
			ListAttachments:              infra.ListAttachments,
			CreateAttachment:             infra.CreateAttachment,
			DeleteAttachment:             infra.DeleteAttachment,
			NewID:                        infra.NewAttachmentID,
		}
		if uc.User.List != nil {
			wuMod.ListUsers = uc.User.List
		}
		if uc.WorkspaceUserRole.GetListPageData != nil {
			wuMod.GetWorkspaceUserRoleListPageData = uc.WorkspaceUserRole.GetListPageData
		}
		identity.NewWorkspaceUserModule(wuMod).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func WorkspaceUserRoleUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entityworkspaceuserrole.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entityworkspaceuserrole.Routes)
		l := u.Labels.(*entityworkspaceuserrole.Labels)

		if uc.WorkspaceUserRole.Create == nil {
			log.Println("entydad catalog: warning: workspace_user_role use cases not initialized — workspace_user_role drawer routes will be unavailable")
			return nil
		}
		wurMod := &identity.WorkspaceUserRoleModuleDeps{
			Routes:                  *r,
			Labels:                  *l,
			CommonLabels:            mc.Common,
			CreateWorkspaceUserRole: uc.WorkspaceUserRole.Create,
			DeleteWorkspaceUserRole: uc.WorkspaceUserRole.Delete,
		}
		if uc.WorkspaceUser.GetItemPageData != nil {
			wurMod.GetWorkspaceUserItemPageData = uc.WorkspaceUser.GetItemPageData
		}
		if uc.Role.List != nil {
			wurMod.ListRoles = uc.Role.List
		}
		identity.NewWorkspaceUserRoleModule(wurMod).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Commerce / location sub-context
// ---------------------------------------------------------------------------

func LocationUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entitylocation.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entitylocation.Routes)
		l := u.Labels.(*entitylocation.Labels)

		// Resolve location_area routes for dashboard quick-action links.
		locationAreaRoutes := entitylocationarea.DefaultRoutes()
		if lar, ok := compose.RoutesOf[*entitylocationarea.Routes](mc, "entity.location_area"); ok {
			locationAreaRoutes = *lar
		}

		locationDeps := &location.LocationModuleDeps{
			Routes:               *r,
			LocationAreaRoutes:   locationAreaRoutes,
			CommonLabels:         mc.Common,
			SharedLabels:         infra.SharedLabels,
			Labels:               *l,
			DashboardTitleLabels: infra.DashboardTitleLabels,
			TableLabels:          mc.Table,
			GetListPageData:      uc.Location.GetListPageData,
			GetInUseIDs:          infra.RefChecker.GetLocationInUseIDs,
			CreateLocation:       uc.Location.Create,
			ReadLocation:         uc.Location.Read,
			UpdateLocation:       uc.Location.Update,
			DeleteLocation:       uc.Location.Delete,
			SetActive:            setActiveClosure(uc, "location"),
			UploadFile:           infra.UploadFile,
			ListAttachments:      infra.ListAttachments,
			CreateAttachment:     infra.CreateAttachment,
			DeleteAttachment:     infra.DeleteAttachment,
			NewID:                infra.NewAttachmentID,
		}
		if uc.LocationArea.List != nil {
			listLocationAreas := uc.LocationArea.List
			locationDeps.ListLocationAreas = func(fctx context.Context) ([]locationaction.LocationAreaOption, error) {
				resp, err := listLocationAreas(fctx, &locationareapb.ListLocationAreasRequest{})
				if err != nil {
					return nil, err
				}
				rows := resp.GetData()
				opts := make([]locationaction.LocationAreaOption, 0, len(rows))
				for _, row := range rows {
					if !row.GetActive() {
						continue
					}
					id := row.GetId()
					if id == "" {
						continue
					}
					opts = append(opts, locationaction.LocationAreaOption{ID: id, Name: row.GetName()})
				}
				return opts, nil
			}
		}
		if uc.GetLocationDashboardPageData != nil {
			locationDeps.GetLocationDashboardPageData = uc.GetLocationDashboardPageData
		}
		location.NewLocationModule(locationDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func LocationAreaUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entitylocationarea.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entitylocationarea.Routes)
		l := u.Labels.(*entitylocationarea.Labels)

		la := uc.LocationArea
		if la.List == nil || la.Create == nil || la.Read == nil || la.Update == nil || la.Delete == nil {
			log.Println("entydad catalog: warning: LocationArea use cases not wired — skipping location_area module")
			return nil
		}
		location.NewLocationAreaModule(&location.LocationAreaModuleDeps{
			Routes:       *r,
			CommonLabels: mc.Common,
			SharedLabels: infra.SharedLabels,
			Labels:       *l,
			TableLabels:  mc.Table,
			GetListPageData: func(fctx context.Context, status string, search string, page, pageSize int) (*locationarealist.LocationAreaListResult, error) {
				resp, err := la.List(fctx, &locationareapb.ListLocationAreasRequest{})
				if err != nil {
					return nil, err
				}
				rows := resp.GetData()
				items := make([]*locationarealist.LocationAreaItem, 0, len(rows))
				for _, row := range rows {
					active := row.GetActive()
					recordStatus := "active"
					if !active {
						recordStatus = "inactive"
					}
					if recordStatus != status {
						continue
					}
					items = append(items, &locationarealist.LocationAreaItem{
						ID:          row.GetId(),
						Name:        row.GetName(),
						Description: row.GetDescription(),
						Active:      active,
						DateCreated: row.GetDateCreatedString(),
					})
				}
				return &locationarealist.LocationAreaListResult{Items: items, TotalItems: len(items)}, nil
			},
			GetInUseIDs: infra.RefChecker.GetLocationAreaInUseIDs,
			CreateLocationArea: func(fctx context.Context, name, description string, active bool) (string, error) {
				resp, err := la.Create(fctx, &locationareapb.CreateLocationAreaRequest{
					Data: &locationareapb.LocationArea{Name: name, Description: description, Active: active},
				})
				if err != nil {
					return "", err
				}
				if data := resp.GetData(); len(data) > 0 {
					return data[0].GetId(), nil
				}
				return "", nil
			},
			ReadLocationArea: func(fctx context.Context, id string) (*locationareaaction.LocationAreaRecord, error) {
				resp, err := la.Read(fctx, &locationareapb.ReadLocationAreaRequest{
					Data: &locationareapb.LocationArea{Id: id},
				})
				if err != nil {
					return nil, err
				}
				data := resp.GetData()
				if len(data) == 0 {
					return nil, nil
				}
				row := data[0]
				return &locationareaaction.LocationAreaRecord{
					ID:          row.GetId(),
					Name:        row.GetName(),
					Description: row.GetDescription(),
					Active:      row.GetActive(),
				}, nil
			},
			UpdateLocationArea: func(fctx context.Context, id, name, description string, active bool) error {
				_, err := la.Update(fctx, &locationareapb.UpdateLocationAreaRequest{
					Data: &locationareapb.LocationArea{Id: id, Name: name, Description: description, Active: active},
				})
				return err
			},
			DeleteLocationArea: func(fctx context.Context, id string) error {
				_, err := la.Delete(fctx, &locationareapb.DeleteLocationAreaRequest{
					Data: &locationareapb.LocationArea{Id: id},
				})
				return err
			},
			SetLocationAreaActive: func(fctx context.Context, id string, active bool) error {
				// Read-modify-write: UpdateLocationArea requires Name (NOT NULL).
				readResp, err := la.Read(fctx, &locationareapb.ReadLocationAreaRequest{
					Data: &locationareapb.LocationArea{Id: id},
				})
				if err != nil {
					return err
				}
				name, description := "", ""
				if data := readResp.GetData(); len(data) > 0 {
					name = data[0].GetName()
					description = data[0].GetDescription()
				}
				_, err = la.Update(fctx, &locationareapb.UpdateLocationAreaRequest{
					Data: &locationareapb.LocationArea{Id: id, Name: name, Description: description, Active: active},
				})
				return err
			},
		}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func PaymentTermUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := entitypaymentterm.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*entitypaymentterm.Routes)
		l := u.Labels.(*entitypaymentterm.Labels)

		if uc.PaymentTerm.GetListPageData == nil {
			log.Println("entydad catalog: warning: PaymentTerm use cases not wired — skipping payment_term module")
			return nil
		}
		setPaymentTermActive := setActiveClosure(uc, "payment_term")
		sharedDeps := &commerce.PaymentTermModuleDeps{
			CommonLabels:         mc.Common,
			SharedLabels:         infra.SharedLabels,
			Labels:               *l,
			TableLabels:          mc.Table,
			GetListPageData:      uc.PaymentTerm.GetListPageData,
			GetInUseIDs:          infra.RefChecker.GetPaymentTermInUseIDs,
			CreatePaymentTerm:    uc.PaymentTerm.CreatePaymentTerm,
			ReadPaymentTerm:      uc.PaymentTerm.ReadPaymentTerm,
			UpdatePaymentTerm:    uc.PaymentTerm.UpdatePaymentTerm,
			DeletePaymentTerm:    uc.PaymentTerm.DeletePaymentTerm,
			SetPaymentTermActive: setPaymentTermActive,
		}
		// Client-context payment term list (entity_scope IN ('client', 'both')).
		clientDeps := *sharedDeps
		clientDeps.Routes = *r
		clientDeps.Scope = "client"
		commerce.NewPaymentTermModule(&clientDeps).RegisterRoutes(mc.Routes)
		// Supplier-context payment term list (entity_scope IN ('supplier', 'both')).
		// DefaultSupplierRoutes provides the supplier-scoped URL constants; no
		// separate compose.Unit exists for supplier_payment_term.
		supplierRoutes := entitypaymentterm.DefaultSupplierRoutes()
		supplierDeps := *sharedDeps
		supplierDeps.Routes = supplierRoutes.ToRoutes()
		supplierDeps.Scope = "supplier"
		commerce.NewPaymentTermModule(&supplierDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Tax domain
// ---------------------------------------------------------------------------

func TaxRegistrationUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := taxregistration.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*taxregistration.Routes)
		l := u.Labels.(*taxregistration.Labels)

		taxRegDeps := &tax.TaxRegistrationModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		if uc.TaxRegistration.List != nil {
			taxRegDeps.ListTaxRegistrations = uc.TaxRegistration.List
		}
		tax.NewTaxRegistrationModule(taxRegDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Service: Auth
// ---------------------------------------------------------------------------

// AuthUnit returns a compose.Unit that mounts the auth service module (login,
// signup, reset-password, change-password, logout, multi-principal chooser).
//
// Auth is NOT a standard entity module with Describe() — it has no
// Routes/Labels/Templates loaded via the compose engine. Instead, its
// routes are hardcoded in entydad/service/auth, and labels are injected
// through auth.Deps by the host app's composition layer.
//
// The Mount closure type-asserts MountContext.Routes to auth.RouteRegistrar
// (view.RouteRegistrar + HandleFunc). The assertion always succeeds for
// service-admin's chi-based RouteRegistry; other hosts that lack HandleFunc
// degrade gracefully with a log warning (matching compose.HandleFunc's
// existing behavior across all blocks).
func AuthUnit(infra *Infra) compose.Unit {
	return compose.Unit{
		Key: "service.auth",
		Mount: func(mc *compose.MountContext) error {
			if infra.AuthDeps == nil {
				log.Println("entydad catalog: warning: AuthDeps is nil — skipping auth module")
				return nil
			}

			// The auth module's RegisterRoutes requires auth.RouteRegistrar
			// (view.RouteRegistrar + HandleFunc). Assert the compose registrar
			// to that interface.
			authRoutes, ok := mc.Routes.(auth.RouteRegistrar)
			if !ok {
				log.Println("entydad catalog: warning: RouteRegistrar does not implement auth.RouteRegistrar (HandleFunc) — skipping auth module")
				return nil
			}

			module := auth.NewAuthModule(infra.AuthDeps)
			module.RegisterRoutes(authRoutes)

			log.Println("  ✓ Auth module delegated to entydad/service/auth (compose Unit)")
			return nil
		},
	}
}

// ---------------------------------------------------------------------------
// Aggregator
// ---------------------------------------------------------------------------

// AllUnits returns the complete curated unit list for all entydad entity
// domains in the same registration order as Block(): party → identity →
// commerce → tax → service.auth.
func AllUnits(uc *UseCases, infra *Infra) []compose.Unit {
	units := []compose.Unit{
		// Party sub-context
		ClientUnit(uc, infra),
		SupplierUnit(uc, infra),
		ClientTagUnit(uc, infra),
		SupplierTagUnit(uc, infra),
		// Identity sub-context
		UserUnit(uc, infra),
		RoleUnit(uc, infra),
		PermissionUnit(uc, infra),
		WorkspaceUnit(uc, infra),
		WorkspaceUserUnit(uc, infra),
		WorkspaceUserRoleUnit(uc, infra),
		// Commerce / location sub-context
		LocationUnit(uc, infra),
		LocationAreaUnit(uc, infra),
		PaymentTermUnit(uc, infra),
		// Tax domain
		TaxRegistrationUnit(uc, infra),
		// Service: auth (login, signup, reset-password, etc.)
		AuthUnit(infra),
	}
	return units
}
