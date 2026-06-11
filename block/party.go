// party.go — block sub-context lift (B, block-go-anatomy).
//
// wirePartyModule registers the party sub-context entity modules
// (client, supplier, client_tag, supplier_tag) into the app router. It is a
// PURE code-move of the corresponding `if cfg.enableAll || cfg.X { ... }`
// blocks from block.go's Block() — same construction order, registration
// order, callbacks, and nil-checks. No behaviour change.
//
// All deps the lifted bodies need from Block()'s scope are carried on
// partyWiring (block-go-anatomy: >6 deps → struct).
package block

import (
	"context"
	"fmt"
	"log"

	party "github.com/erniealice/entydad-golang/domain/entity/party"
	clientdetail "github.com/erniealice/entydad-golang/domain/entity/party/client/detail"
	"github.com/erniealice/espyna-golang/reference"
	"github.com/erniealice/espyna-golang/registry"
	entityid "github.com/erniealice/espyna-golang/registry/entityid"
	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"google.golang.org/protobuf/proto"
)

// partyWiring carries everything the party cluster needs from Block()'s
// scope. Implementation detail of the wiring; never re-exported.
type partyWiring struct {
	cfg        *blockConfig
	uc         *UseCases
	db         UpdateableSource
	labels     blockLabels
	routes     blockRoutes
	refChecker reference.Checker

	getClientStatement   func(ctx context.Context, req *clientstmtpb.ClientStatementRequest) (*clientstmtpb.ClientStatementResponse, error)
	getClientBalances    func(ctx context.Context) (map[string]int64, error)
	getSupplierStatement func(ctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error)
	getSupplierBalances  func(ctx context.Context) (map[string]int64, error)

	uploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	listAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID  func() string
}

func wirePartyModule(ctx *pyeza.AppContext, w partyWiring) {
	cfg := w.cfg
	uc := w.uc
	db := w.db
	labels := w.labels
	routes := w.routes
	refChecker := w.refChecker
	getClientStatement := w.getClientStatement
	getClientBalances := w.getClientBalances
	getSupplierStatement := w.getSupplierStatement
	getSupplierBalances := w.getSupplierBalances
	uploadFile := w.uploadFile
	listAttachments := w.listAttachments
	createAttachment := w.createAttachment
	deleteAttachment := w.deleteAttachment
	newAttachmentID := w.newAttachmentID

	if cfg.enableAll || cfg.client {
		clientDeps := &party.ClientModuleDeps{
			Routes: routes.Client,
			// User module owns the timezone search endpoint; the client
			// representative form reuses the same JSON handler.
			SearchTimezonesURL:   routes.User.SearchTimezonesURL,
			CommonLabels:         ctx.Common,
			SharedLabels:         labels.Shared,
			Labels:               labels.Client,
			DashboardLabels:      labels.ClientDashboard,
			DashboardTitleLabels: labels.Dashboard,
			TableLabels:          ctx.Table,
			GetListPageData:      uc.Client.GetListPageData,
			GetInUseIDs:          refChecker.GetClientInUseIDs,
			CreateClient:         uc.Client.Create,
			ReadClient:           uc.Client.Read,
			UpdateClient:         uc.Client.Update,
			DeleteClient:         uc.Client.Delete,
			SetStatus: func(fctx context.Context, id string, status string) error {
				// Keep the legacy `active` boolean in sync with the lifecycle
				// status so consumers that still filter on c.active see
				// consistent values: only "inactive" flips active=false.
				active := status != "inactive"
				_, err := db.Update(fctx, "client", id, map[string]any{
					"status": status,
					"active": active,
				})
				return err
			},
			ListPaymentTerms: func(fctx context.Context) ([]*party.ClientPaymentTermOption, error) {
				rows, err := db.ListSimple(fctx, "payment_term")
				if err != nil {
					return nil, err
				}
				opts := make([]*party.ClientPaymentTermOption, 0, len(rows))
				for _, row := range rows {
					id, _ := row["id"].(string)
					name, _ := row["name"].(string)
					scope, _ := row["entity_scope"].(string)
					if id == "" {
						continue
					}
					// Client context: show terms scoped to "client" or "both"
					if scope != "client" && scope != "both" {
						continue
					}
					opts = append(opts, &party.ClientPaymentTermOption{Id: id, Name: name})
				}
				return opts, nil
			},
			ListRevenues:                     db.ListSimple,
			GetClientStatement:               getClientStatement,
			SubscriptionAddURL:               routes.Subscription.AddURL,
			SubscriptionDetailURL:            routes.Subscription.DetailURL,
			SubscriptionUnderClientDetailURL: routes.Subscription.UnderClientDetailURL,
			SubscriptionEditURL:              routes.Subscription.EditURL,
			SubscriptionDeleteURL:            routes.Subscription.DeleteURL,
			UploadFile:                       uploadFile,
			ListAttachments:                  listAttachments,
			CreateAttachment:                 createAttachment,
			DeleteAttachment:                 deleteAttachment,
			NewID:                            newAttachmentID,
		}
		if uc.Category.List != nil {
			clientDeps.ListCategories = uc.Category.List
		}
		if uc.Client.Category.List != nil {
			clientDeps.ListClientCategories = uc.Client.Category.List
			clientDeps.CreateClientCategory = uc.Client.Category.Create
			clientDeps.DeleteClientCategory = uc.Client.Category.Delete
		}
		if uc.Subscription.List != nil {
			clientDeps.ListSubscriptions = uc.Subscription.List
			clientDeps.GetSubscriptionListPageData = uc.Subscription.GetListPageData
		}
		// Wire the PriceSchedules tab: list all price_schedules, filter to
		// client_id == clientID, then map each to a ClientPriceScheduleRow.
		// PlanCount is computed from a single secondary ListPricePlans call —
		// build a map[scheduleID]int once, then injected per row. DetailURL
		// resolved from the active PriceSchedule route config (lyngua overrides
		// per business type).
		if uc.PriceSchedule.List != nil {
			listPriceSchedules := uc.PriceSchedule.List
			var listPricePlans func(context.Context, *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
			if uc.PricePlan.List != nil {
				listPricePlans = uc.PricePlan.List
			}
			psDetailURL := routes.PriceSchedule.DetailURL
			clientDeps.ListClientPriceSchedules = func(fctx context.Context, clientID string) ([]clientdetail.ClientPriceScheduleRow, error) {
				resp, err := listPriceSchedules(fctx, &priceschedulepb.ListPriceSchedulesRequest{})
				if err != nil {
					return nil, err
				}
				// Build plan count by schedule_id: one ListPricePlans call,
				// tally GetPriceScheduleId() per row. Failures are non-fatal
				// (counts default to 0).
				planCountBySchedule := map[string]int{}
				if listPricePlans != nil {
					ppResp, ppErr := listPricePlans(fctx, &priceplanpb.ListPricePlansRequest{})
					if ppErr != nil {
						log.Printf("entydad.Block: failed to list price plans for plan-count column: %v", ppErr)
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
			clientDeps.PriceScheduleAddURL = routes.PriceSchedule.AddURL
		}
		if uc.Workspace.Read != nil {
			readWorkspace := uc.Workspace.Read
			wsID := getDefaultWorkspaceID()
			clientDeps.GetFunctionalCurrency = func(fctx context.Context) string {
				resp, err := readWorkspace(fctx, &workspacepb.ReadWorkspaceRequest{
					Data: &workspacepb.Workspace{Id: wsID},
				})
				if err != nil {
					return ""
				}
				data := resp.GetData()
				if len(data) == 0 {
					return ""
				}
				return data[0].GetFunctionalCurrency()
			}
		}
		if getClientBalances != nil {
			clientDeps.GetClientBalances = getClientBalances
		}
		if uc.Subscription.CountActiveByClientIDs != nil {
			countActive := uc.Subscription.CountActiveByClientIDs
			clientDeps.GetActiveSubscriptionCounts = func(fctx context.Context) (map[string]int32, error) {
				resp, err := countActive(fctx, &subscriptionpb.CountActiveByClientIdsRequest{})
				if err != nil {
					return nil, err
				}
				return resp.GetCounts(), nil
			}
		}
		// Wire ListRevenuesByClient using the typed Revenue use case with a
		// client_id StringFilter. The postgres adapter supports BuildFilterWhere
		// on the revenue table which includes client_id as a filterable column.
		if uc.Revenue.List != nil {
			listRevenues := uc.Revenue.List
			clientDeps.ListRevenuesByClient = func(fctx context.Context, clientID string) ([]*revenuepb.Revenue, error) {
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
			clientDeps.ListCollectionsByClient = func(fctx context.Context, clientID string) ([]*collectionpb.Collection, error) {
				resp, err := listByClient(fctx, &collectionpb.ListByClientRequest{ClientId: clientID})
				if err != nil {
					return nil, err
				}
				return resp.GetData(), nil
			}
		}

		// Wire Revenue Run drawer shims. Both use cases must be present for
		// either callback to be wired, ensuring the drawer is either fully
		// functional or fully absent.
		if uc.Revenue.ListRevenueRunCandidates != nil &&
			uc.Revenue.GenerateRevenueRun != nil {
			listCandidates := uc.Revenue.ListRevenueRunCandidates
			generateRun := uc.Revenue.GenerateRevenueRun
			clientDeps.ListRevenueRunCandidates = func(fctx context.Context, scope clientdetail.RevenueRunScope) ([]clientdetail.RevenueRunCandidate, string, error) {
				// Plan B Phase 5c — opt-in to advance Collection candidates by default.
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

			clientDeps.GenerateRevenueRun = func(fctx context.Context, scope clientdetail.RevenueRunScope, selections clientdetail.RevenueRunSelections) (*clientdetail.RevenueRunResult, error) {
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
					// Plan B Phase 5c — dispatch on source_kind.
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
				runID := ""
				runStatus := ""
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

		party.NewClientModule(clientDeps).RegisterRoutes(ctx.Routes)
	}

	if cfg.enableAll || cfg.supplier {
		supplierDeps := &party.SupplierModuleDeps{
			Routes:            routes.Supplier,
			SupplierTagRoutes: routes.SupplierTag,
			// User module owns the timezone search endpoint; the supplier
			// representative form reuses the same JSON handler.
			SearchTimezonesURL:   routes.User.SearchTimezonesURL,
			CommonLabels:         ctx.Common,
			SharedLabels:         labels.Shared,
			Labels:               labels.Supplier,
			DashboardLabels:      labels.SupplierDashboard,
			DashboardTitleLabels: labels.Dashboard,
			TableLabels:          ctx.Table,
			GetInUseIDs:          refChecker.GetSupplierInUseIDs,
			SetStatus: func(fctx context.Context, id string, status string) error {
				// `active` is the soft-delete flag; status transitions
				// (active/blocked/on_hold) must NOT flip it, otherwise
				// deactivated suppliers look identical to deleted ones.
				// Only DeleteSupplier should set active=false.
				_, err := db.Update(fctx, "supplier", id, map[string]any{
					"active": true,
					"status": status,
				})
				return err
			},
			ListPaymentTerms: func(fctx context.Context) ([]*party.SupplierPaymentTermOption, error) {
				rows, err := db.ListSimple(fctx, "payment_term")
				if err != nil {
					return nil, err
				}
				opts := make([]*party.SupplierPaymentTermOption, 0, len(rows))
				for _, row := range rows {
					id, _ := row["id"].(string)
					name, _ := row["name"].(string)
					scope, _ := row["entity_scope"].(string)
					if id == "" {
						continue
					}
					// Supplier context: show terms scoped to "supplier" or "both"
					if scope != "supplier" && scope != "both" {
						continue
					}
					opts = append(opts, &party.SupplierPaymentTermOption{Id: id, Name: name})
				}
				return opts, nil
			},
			UploadFile:       uploadFile,
			ListAttachments:  listAttachments,
			CreateAttachment: createAttachment,
			DeleteAttachment: deleteAttachment,
			NewID:            newAttachmentID,
		}
		if uc.Supplier.GetListPageData != nil {
			supplierDeps.GetListPageData = uc.Supplier.GetListPageData
			supplierDeps.CreateSupplier = uc.Supplier.Create
			supplierDeps.ReadSupplier = uc.Supplier.Read
			supplierDeps.UpdateSupplier = uc.Supplier.Update
			supplierDeps.DeleteSupplier = uc.Supplier.Delete
		}
		if uc.PurchaseOrder.List != nil {
			supplierDeps.ListPurchaseOrders = uc.PurchaseOrder.List
		}
		if getSupplierStatement != nil {
			supplierDeps.GetSupplierStatement = getSupplierStatement
		}
		if getSupplierBalances != nil {
			supplierDeps.GetSupplierBalances = getSupplierBalances
		}
		// Tag-related deps for supplier form multi-select
		if uc.Category.List != nil {
			supplierDeps.ListCategories = uc.Category.List
		}
		if uc.Supplier.Category.List != nil {
			supplierDeps.ListSupplierCategories = uc.Supplier.Category.List
			supplierDeps.CreateSupplierCategory = uc.Supplier.Category.Create
			supplierDeps.DeleteSupplierCategory = uc.Supplier.Category.Delete
		}
		party.NewSupplierModule(supplierDeps).RegisterRoutes(ctx.Routes)
	}

	if cfg.enableAll || cfg.clientTag {
		clienttagDeps := &party.ClientTagModuleDeps{
			Routes:       routes.ClientTag,
			Labels:       labels.ClientTag,
			SharedLabels: labels.Shared,
			CommonLabels: ctx.Common,
			TableLabels:  ctx.Table,
			GetInUseIDs:  refChecker.GetCategoryInUseIDs,
			SetCategoryActive: func(fctx context.Context, id string, active bool) error {
				_, err := db.Update(fctx, "category", id, map[string]any{"active": active})
				return err
			},
		}
		if uc.Category.List != nil {
			clienttagDeps.ListCategories = uc.Category.List
			clienttagDeps.CreateCategory = uc.Category.Create
			clienttagDeps.ReadCategory = uc.Category.Read
			clienttagDeps.UpdateCategory = uc.Category.Update
			clienttagDeps.DeleteCategory = uc.Category.Delete
		}
		if ctx.SqlDB != nil {
			repoAny, err := registry.CreateRepository("postgresql", entityid.Category, ctx.SqlDB, "category")
			if err == nil {
				if pgd, ok := repoAny.(categoryListPageDataGetter); ok {
					clienttagDeps.GetCategoryListPageData = pgd.GetCategoryListPageData
				}
			}
		}
		if uc.Client.Category.List != nil {
			clienttagDeps.ListClientCategories = uc.Client.Category.List
		}
		party.NewClientTagModule(clienttagDeps).RegisterRoutes(ctx.Routes)
	}

	if cfg.enableAll || cfg.supplierTag {
		suppliertagDeps := &party.SupplierTagModuleDeps{
			Routes:       routes.SupplierTag,
			Labels:       labels.SupplierTag,
			SharedLabels: labels.Shared,
			CommonLabels: ctx.Common,
			TableLabels:  ctx.Table,
			GetInUseIDs:  refChecker.GetCategoryInUseIDs,
			SetCategoryActive: func(fctx context.Context, id string, active bool) error {
				_, err := db.Update(fctx, "category", id, map[string]any{"active": active})
				return err
			},
		}
		if uc.Category.List != nil {
			suppliertagDeps.ListCategories = uc.Category.List
			suppliertagDeps.CreateCategory = uc.Category.Create
			suppliertagDeps.ReadCategory = uc.Category.Read
			suppliertagDeps.UpdateCategory = uc.Category.Update
			suppliertagDeps.DeleteCategory = uc.Category.Delete
		}
		if ctx.SqlDB != nil {
			repoAny, err := registry.CreateRepository("postgresql", entityid.Category, ctx.SqlDB, "category")
			if err == nil {
				if pgd, ok := repoAny.(categoryListPageDataGetter); ok {
					suppliertagDeps.GetCategoryListPageData = pgd.GetCategoryListPageData
				}
			}
		}
		if uc.Supplier.Category.List != nil {
			suppliertagDeps.ListSupplierCategories = uc.Supplier.Category.List
		}
		party.NewSupplierTagModule(suppliertagDeps).RegisterRoutes(ctx.Routes)
	}
}
