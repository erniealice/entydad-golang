// Package block provides the Block() function for registering all entydad
// entity modules into a pyeza app via AppOption composition.
//
// It lives in a sub-package to avoid the import cycle that would arise if it
// were placed in the root entydad package (the view sub-packages already import
// the root package for route/label types).
//
// Companion files in this directory:
//   - helpers.go       — DB interface types (UpdateableSource, CRUDSource,
//     categoryListPageDataGetter) and getDefaultWorkspaceID.
//   - route_loading.go — blockLabels / blockRoutes types and their lyngua loaders
//     (loadBlockLabels, loadBlockRoutes).
//   - wiring.go        — dashboard reflective wiring helpers (wireLocationDashboard,
//     wireAdminDashboard).
//
// Usage:
//
//	import entydadblock "github.com/erniealice/entydad-golang/block"
//
//	app, err := pyeza.NewApp(
//	    entydadblock.Block(),                           // all modules
//	    entydadblock.Block(entydadblock.WithClient()),  // client only
//	)
package block

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	appcontext "github.com/erniealice/espyna-golang/appcontext"

	"github.com/erniealice/entydad-golang"
	adminmod "github.com/erniealice/entydad-golang/views/admin"
	admindashboardroutes "github.com/erniealice/entydad-golang/views/admin/dashboard"
	clientmod "github.com/erniealice/entydad-golang/views/client"
	clientdetail "github.com/erniealice/entydad-golang/views/client/detail"
	clienttagmod "github.com/erniealice/entydad-golang/views/clienttag"
	conversationmod "github.com/erniealice/entydad-golang/views/conversation"
	locationmod "github.com/erniealice/entydad-golang/views/location"
	locationaction "github.com/erniealice/entydad-golang/views/location/action"
	locationareamod "github.com/erniealice/entydad-golang/views/location_area"
	locationareaaction "github.com/erniealice/entydad-golang/views/location_area/action"
	locationarealist "github.com/erniealice/entydad-golang/views/location_area/list"
	paymenttermmod "github.com/erniealice/entydad-golang/views/payment_term"
	permissionmod "github.com/erniealice/entydad-golang/views/permission"
	rolemod "github.com/erniealice/entydad-golang/views/role"
	roleusers "github.com/erniealice/entydad-golang/views/role/users"
	suppliermod "github.com/erniealice/entydad-golang/views/supplier"
	suppliertagmod "github.com/erniealice/entydad-golang/views/suppliertag"
	taxregistrationmod "github.com/erniealice/entydad-golang/views/tax_registration"
	usermod "github.com/erniealice/entydad-golang/views/user"
	userdashboard "github.com/erniealice/entydad-golang/views/user/dashboard"
	workspacemod "github.com/erniealice/entydad-golang/views/workspace"
	workspaceaction "github.com/erniealice/entydad-golang/views/workspace/action"
	workspaceusermod "github.com/erniealice/entydad-golang/views/workspace_user"
	workspaceuserrolemod "github.com/erniealice/entydad-golang/views/workspace_user_role"
	"github.com/erniealice/espyna-golang/reference"
	"github.com/erniealice/espyna-golang/registry"
	entityid "github.com/erniealice/espyna-golang/registry/entityid"
	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"google.golang.org/protobuf/proto"
)

// routeRegistrarFull — optional extension for raw http.HandlerFunc routes.
// Apps that implement HandleFunc (e.g., service-admin chi router wrapper) can
// register JSON endpoints and other raw HTTP handlers. Apps that don't will skip.
type routeRegistrarFull interface {
	pyeza.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// handleFunc registers an http.HandlerFunc route if the registrar supports it.
// Silently skips if the registrar does not implement HandleFunc.
func handleFunc(r pyeza.RouteRegistrar, method, path string, handler http.HandlerFunc) {
	if full, ok := r.(routeRegistrarFull); ok {
		full.HandleFunc(method, path, handler)
		return
	}
	log.Printf("entydad.Block: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
}

// BlockOption configures which entydad modules are registered by Block().
type BlockOption func(*blockConfig)

type blockConfig struct {
	enableAll         bool
	admin             bool
	client            bool
	clientTag         bool
	supplierTag       bool
	paymentTerm       bool
	user              bool
	role              bool
	location          bool
	locationArea      bool
	permission        bool
	workspace         bool
	workspaceUser     bool
	workspaceUserRole bool
	supplier          bool
	taxRegistration   bool
	conversation      bool
	useCases          *UseCases
	// homeURL is the static fallback URL to redirect to after a successful
	// workspace switch. Defaults to "/home" (post-P12 of workspace-keyed-routing
	// plan; "/app/home" is gone) when empty.
	homeURL string
	// homeURLForWorkspaceID resolves the post-switch redirect URL given the
	// newly-active workspace_id. When non-nil, it takes precedence over
	// homeURL so callers can land on /w/{slug}/home per Q-WS-1 → A.
	homeURLForWorkspaceID func(ctx context.Context, workspaceID string) string
	// secureSwitch (A1 fix WKR-P0-1, 2026-05-22): when set, the
	// switch-workspace handler routes through the host app's secure
	// rotation primitive instead of the legacy in-place use case. See
	// workspaceaction.SwitchWorkspaceDeps.SecureSwitch.
	secureSwitch             workspaceaction.SecureSwitchFn
	secureSwitchResolveUser  func(r *http.Request) string
	secureSwitchSetCookie    func(w http.ResponseWriter, token string)
}

// WithUseCases supplies the typed use-case closures to Block().
// Required: Block() returns an error if this option is not provided.
// Service-admin constructs the *UseCases via an adapter function that
// bridges espyna's consumer container to entydad's typed shape.
func WithUseCases(uc *UseCases) BlockOption {
	return func(c *blockConfig) { c.useCases = uc }
}

// WithAdmin enables the Admin dashboard module in Block().
// The admin module registers the /app/admin/dashboard route.
func WithAdmin() BlockOption { return func(c *blockConfig) { c.admin = true } }

// WithClient enables the Client module in Block().
func WithClient() BlockOption { return func(c *blockConfig) { c.client = true } }

// WithClientTag enables the ClientTag module in Block().
func WithClientTag() BlockOption { return func(c *blockConfig) { c.clientTag = true } }

// WithSupplierTag enables the SupplierTag module in Block().
func WithSupplierTag() BlockOption { return func(c *blockConfig) { c.supplierTag = true } }

// WithPaymentTerm enables the PaymentTerm module in Block().
func WithPaymentTerm() BlockOption { return func(c *blockConfig) { c.paymentTerm = true } }

// WithUser enables the User module in Block().
func WithUser() BlockOption { return func(c *blockConfig) { c.user = true } }

// WithRole enables the Role module in Block().
func WithRole() BlockOption { return func(c *blockConfig) { c.role = true } }

// WithLocation enables the Location module in Block().
func WithLocation() BlockOption { return func(c *blockConfig) { c.location = true } }

// WithLocationArea enables the LocationArea module in Block().
func WithLocationArea() BlockOption { return func(c *blockConfig) { c.locationArea = true } }

// WithPermission enables the Permission module in Block().
func WithPermission() BlockOption { return func(c *blockConfig) { c.permission = true } }

// WithWorkspace enables the Workspace module in Block().
func WithWorkspace() BlockOption { return func(c *blockConfig) { c.workspace = true } }

// WithWorkspaceUser enables the WorkspaceUser nested-detail module in Block().
// Phase 2: registers detail page, tab-action, add/delete/set-status, and user-search routes.
func WithWorkspaceUser() BlockOption { return func(c *blockConfig) { c.workspaceUser = true } }

// WithWorkspaceUserRole enables the WorkspaceUserRole assignment drawer module in Block().
// Phase 3: registers add, delete, permissions, and search-roles routes.
func WithWorkspaceUserRole() BlockOption {
	return func(c *blockConfig) { c.workspaceUserRole = true }
}

// WithSupplier enables the Supplier module in Block().
func WithSupplier() BlockOption { return func(c *blockConfig) { c.supplier = true } }

// WithTaxRegistration enables the TaxRegistration polymorphic list module in Block().
// Registers both client-scoped and workspace-scoped views.
func WithTaxRegistration() BlockOption { return func(c *blockConfig) { c.taxRegistration = true } }

// WithConversation enables the Conversation secure-messaging module in Block()
// (staff inbox + thread detail + composer; client portal is built but gated).
func WithConversation() BlockOption { return func(c *blockConfig) { c.conversation = true } }

// WithHomeURL sets the URL the switch-workspace handler redirects to after a
// successful workspace switch. Defaults to "/app/home" when not provided.
func WithHomeURL(url string) BlockOption { return func(c *blockConfig) { c.homeURL = url } }

// WithHomeURLForWorkspaceID supplies a resolver that returns /w/{slug}/home (or
// the appropriate workspace-keyed URL) given the workspace_id being switched to.
// Used by service-admin to land users on the URL-canonical workspace home after
// the workspace switcher fires. Per Q-WS-1 → A of docs/plan/20260521-workspace-keyed-routing.
func WithHomeURLForWorkspaceID(fn func(ctx context.Context, workspaceID string) string) BlockOption {
	return func(c *blockConfig) { c.homeURLForWorkspaceID = fn }
}

// WithSecureSwitch wires the rotation-aware workspace switch primitive into
// the /action/admin/switch-workspace handler. Service-admin passes its
// executePrincipalSwitch closure here so the sidebar workspace-switcher
// rotates the session token, locks the binding inside tx, and writes an
// audit row — matching the workspace-boundary rotation invariant
// (Q-WS-13). When unset, the handler falls back to the legacy in-place
// SwitchWorkspace use case (no rotation, no audit) for backward
// compatibility with hosts that haven't migrated. A1 fix WKR-P0-1
// (2026-05-22).
//
// All three parameters are required when SecureSwitch is to be active:
//   - switchFn       — the rotation primitive
//   - resolveUserID  — extract user_id from the request (post-session-mw)
//   - setSessionCookie — set the post-rotation session cookie
func WithSecureSwitch(
	switchFn workspaceaction.SecureSwitchFn,
	resolveUserID func(r *http.Request) string,
	setSessionCookie func(w http.ResponseWriter, token string),
) BlockOption {
	return func(c *blockConfig) {
		c.secureSwitch = switchFn
		c.secureSwitchResolveUser = resolveUserID
		c.secureSwitchSetCookie = setSessionCookie
	}
}

// Block returns a pyeza.AppOption that registers entydad entity modules into the app.
// When called with no options, all modules are registered (enableAll mode).
// When called with specific WithXxx() options, only those modules are registered.
//
// Expected ctx fields (type-asserted from any):
//   - ctx.DB           → UpdateableSource (entydad.DataSource + Update method)
//   - ctx.RefChecker   → reference.Checker
//   - ctx.Translations → *lynguaV1.TranslationProvider
//   - ctx.UploadFile, ctx.ListAttachments, ctx.CreateAttachment,
//     ctx.DeleteAttachment, ctx.NewAttachmentID — attachment funcs
//   - ctx.GetUsersByRoleID, ctx.GetDashboardData, ctx.HashPassword,
//     ctx.GetUserWorkspacesMap — user/workspace helpers
//   - ctx.Routes, ctx.Common, ctx.Table, ctx.BusinessType — from pyeza.AppContext
func Block(opts ...BlockOption) pyeza.AppOption {
	cfg := &blockConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	// "Enable all modules" is derived — true when no module-toggling option was
	// passed. Non-module options (WithUseCases, future config options) must NOT
	// flip this off, otherwise the service-admin adapter that calls
	// `Block(WithUseCases(...))` would silently register zero modules (which
	// breaks workspace switch + every other admin route).
	moduleSelected := cfg.admin || cfg.client || cfg.clientTag || cfg.supplierTag ||
		cfg.paymentTerm || cfg.user || cfg.role || cfg.location || cfg.locationArea ||
		cfg.permission || cfg.workspace || cfg.workspaceUser || cfg.workspaceUserRole ||
		cfg.supplier || cfg.taxRegistration || cfg.conversation
	cfg.enableAll = !moduleSelected

	return func(ctx *pyeza.AppContext) error {
		// --- typed UseCases supplied via WithUseCases() ---

		if cfg.useCases == nil {
			return fmt.Errorf("entydad.Block: WithUseCases(...) is required")
		}
		if err := cfg.useCases.RequireFor(cfg); err != nil {
			return err
		}
		uc := cfg.useCases // local alias for brevity

		db, ok := ctx.DB.(UpdateableSource)
		if !ok {
			return fmt.Errorf("entydad.Block: DB must implement block.UpdateableSource (DataSource + Update)")
		}

		refChecker, ok := ctx.RefChecker.(reference.Checker)
		if !ok {
			return fmt.Errorf("entydad.Block: RefChecker must be reference.Checker")
		}

		translations, ok := ctx.Translations.(*lynguaV1.TranslationProvider)
		if !ok {
			return fmt.Errorf("entydad.Block: Translations must be *lynguaV1.TranslationProvider")
		}

		// type-assert attachment operations (nil-safe — attachment funcs may be absent)
		uploadFile, _ := ctx.UploadFile.(func(ctx context.Context, bucket, key string, content []byte, contentType string) error)
		listAttachments, _ := ctx.ListAttachments.(func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error))
		createAttachment, _ := ctx.CreateAttachment.(func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		deleteAttachment, _ := ctx.DeleteAttachment.(func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error))
		newAttachmentID, _ := ctx.NewAttachmentID.(func() string)

		// type-assert user/role helpers (nil-safe)
		getUsersByRoleID, _ := ctx.GetUsersByRoleID.(func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error))
		getDashboardData, _ := ctx.GetDashboardData.(func(ctx context.Context) (*userdashboard.DashboardData, error))
		hashPassword, _ := ctx.HashPassword.(func(password string) (string, error))
		getUserWorkspacesMap, _ := ctx.GetUserWorkspacesMap.(func(ctx context.Context) (map[string][]pyezatypes.ChipData, error))

		// 20260521 Wave B P1.E.4 — statements/balances now flow through
		// the typed `uc.Reports.Statements.*` closures (service-driven).
		// The legacy `ctx.LedgerReportingSvc` assertion is removed; the
		// duck interface is no longer asserted by entydad. See the
		// statement* helper shims below.
		statements := uc.Reports.Statements
		var getClientStatement func(ctx context.Context, req *clientstmtpb.ClientStatementRequest) (*clientstmtpb.ClientStatementResponse, error)
		if statements.GetClientStatement != nil {
			getClientStatement = func(fctx context.Context, req *clientstmtpb.ClientStatementRequest) (*clientstmtpb.ClientStatementResponse, error) {
				resp, err := statements.GetClientStatement(fctx, translateClientStatementReq(req))
				if err != nil {
					return nil, err
				}
				return translateClientStatementResp(resp), nil
			}
		}
		var getSupplierStatement func(ctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error)
		if statements.GetSupplierStatement != nil {
			getSupplierStatement = func(fctx context.Context, req *suppstmtpb.SupplierStatementRequest) (*suppstmtpb.SupplierStatementResponse, error) {
				resp, err := statements.GetSupplierStatement(fctx, translateSupplierStatementReq(req))
				if err != nil {
					return nil, err
				}
				return translateSupplierStatementResp(resp), nil
			}
		}
		getClientBalances := statements.ListClientBalancesAsMap
		getSupplierBalances := statements.ListSupplierBalancesAsMap

		// --- load labels from lyngua ---
		labels := loadBlockLabels(translations, ctx.BusinessType)

		// --- load routes (defaults + lyngua JSON overrides) ---
		routes := loadBlockRoutes(translations, ctx.BusinessType)

		// --- register modules ---

		if cfg.enableAll || cfg.client {
			clientDeps := &clientmod.ModuleDeps{
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
				ListPaymentTerms: func(fctx context.Context) ([]*clientmod.PaymentTermOption, error) {
					rows, err := db.ListSimple(fctx, "payment_term")
					if err != nil {
						return nil, err
					}
					opts := make([]*clientmod.PaymentTermOption, 0, len(rows))
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
						opts = append(opts, &clientmod.PaymentTermOption{Id: id, Name: name})
					}
					return opts, nil
				},
				ListRevenues:       db.ListSimple,
				GetClientStatement: getClientStatement,
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

			clientmod.NewModule(clientDeps).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.user {
			usermod.NewModule(&usermod.ModuleDeps{
				Routes:               routes.User,
				CommonLabels:         ctx.Common,
				SharedLabels:         labels.Shared,
				Labels:               labels.User,
				DashboardLabels:      labels.UserDashboard,
				DashboardTitleLabels: labels.Dashboard,
				UserRoleLabels:       labels.UserRole,
				TableLabels:          ctx.Table,
				GetListPageData:      uc.User.GetListPageData,
				GetUserWorkspacesMap: getUserWorkspacesMap,
				CreateUser:           uc.User.Create,
				ReadUser:             uc.User.Read,
				UpdateUser:           uc.User.Update,
				DeleteUser:           uc.User.Delete,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "user", id, map[string]any{"active": active})
					return err
				},
				CreateWorkspaceUser:          uc.WorkspaceUser.Create,
				ListWorkspaceUsers:           uc.WorkspaceUser.List,
				GetWorkspaceUserItemPageData: uc.WorkspaceUser.GetItemPageData,
				DefaultWorkspaceID:           getDefaultWorkspaceID(),
				CreateWorkspaceUserRole:      uc.WorkspaceUserRole.Create,
				DeleteWorkspaceUserRole:      uc.WorkspaceUserRole.Delete,
				ListRoles:                    uc.Role.List,
				GetDashboardData:             getDashboardData,
				HashPassword:                 hashPassword,
				UploadFile:                   uploadFile,
				ListAttachments:              listAttachments,
				CreateAttachment:             createAttachment,
				DeleteAttachment:             deleteAttachment,
				NewID:                        newAttachmentID,
			}).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.role {
			rolemod.NewModule(&rolemod.ModuleDeps{
				Routes:               routes.Role,
				CommonLabels:         ctx.Common,
				SharedLabels:         labels.Shared,
				Labels:               labels.Role,
				RolePermissionLabels: labels.RolePermission,
				RoleUserLabels:       labels.RoleUser,
				TableLabels:          ctx.Table,
				GetListPageData:      uc.Role.GetListPageData,
				GetInUseIDs:          refChecker.GetRoleInUseIDs,
				CreateRole:           uc.Role.Create,
				ReadRole:             uc.Role.Read,
				UpdateRole:           uc.Role.Update,
				DeleteRole:           uc.Role.Delete,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "role", id, map[string]any{"active": active})
					return err
				},
				GetItemPageData:         uc.Role.GetItemPageData,
				CreateRolePermission:    uc.RolePermission.Create,
				DeleteRolePermission:    uc.RolePermission.Delete,
				ListPermissions:         uc.Permission.List,
				GetUsersByRoleID:        getUsersByRoleID,
				ListWorkspaceUsers:      uc.WorkspaceUser.List,
				CreateWorkspaceUserRole: uc.WorkspaceUserRole.Create,
				DeleteWorkspaceUserRole: uc.WorkspaceUserRole.Delete,
				UploadFile:              uploadFile,
				ListAttachments:         listAttachments,
				CreateAttachment:        createAttachment,
				DeleteAttachment:        deleteAttachment,
				NewID:                   newAttachmentID,
			}).RegisterRoutes(ctx.Routes)

			// Role-User search (http.HandlerFunc — uses HandleFunc, not GET)
			handleFunc(ctx.Routes, "GET", routes.Role.UsersSearchURL, roleusers.NewSearchUsersAction(&roleusers.SearchDeps{
				ListWorkspaceUsers: uc.WorkspaceUser.List,
			}))
		}

		if cfg.enableAll || cfg.location {
			locationDeps := &locationmod.ModuleDeps{
				Routes:             routes.Location,
				LocationAreaRoutes: routes.LocationArea,
				CommonLabels:       ctx.Common,
				SharedLabels:       labels.Shared,
				Labels:             labels.Location,
				TableLabels:        ctx.Table,
				GetListPageData:    uc.Location.GetListPageData,
				GetInUseIDs:        refChecker.GetLocationInUseIDs,
				CreateLocation:     uc.Location.Create,
				ReadLocation:       uc.Location.Read,
				UpdateLocation:     uc.Location.Update,
				DeleteLocation:     uc.Location.Delete,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "location", id, map[string]any{"active": active})
					return err
				},
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
			}
			if crudDB, hasCRUD := db.(CRUDSource); hasCRUD {
				locationDeps.ListLocationAreas = func(fctx context.Context) ([]locationaction.LocationAreaOption, error) {
					rows, err := crudDB.ListSimple(fctx, "location_area")
					if err != nil {
						return nil, err
					}
					opts := make([]locationaction.LocationAreaOption, 0, len(rows))
					for _, row := range rows {
						active, _ := row["active"].(bool)
						if !active {
							continue
						}
						id, _ := row["id"].(string)
						name, _ := row["name"].(string)
						if id == "" {
							continue
						}
						opts = append(opts, locationaction.LocationAreaOption{ID: id, Name: name})
					}
					return opts, nil
				}
			}
			if uc.GetLocationDashboardPageData != nil {
				locationDeps.GetLocationDashboardPageData = uc.GetLocationDashboardPageData
			}
			locationmod.NewModule(locationDeps).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.locationArea {
			crudDB, hasCRUD := db.(CRUDSource)
			if !hasCRUD {
				log.Println("entydad.Block: warning: DB does not implement CRUDSource — skipping location_area module")
			} else {
				locationareamod.NewModule(&locationareamod.ModuleDeps{
					Routes:       routes.LocationArea,
					CommonLabels: ctx.Common,
					SharedLabels: labels.Shared,
					Labels:       labels.LocationArea,
					TableLabels:  ctx.Table,
					GetListPageData: func(fctx context.Context, status string, search string, page, pageSize int) (*locationarealist.LocationAreaListResult, error) {
						rows, err := crudDB.ListSimple(fctx, "location_area")
						if err != nil {
							return nil, err
						}
						items := make([]*locationarealist.LocationAreaItem, 0, len(rows))
						for _, row := range rows {
							active, _ := row["active"].(bool)
							recordStatus := "active"
							if !active {
								recordStatus = "inactive"
							}
							if recordStatus != status {
								continue
							}
							id, _ := row["id"].(string)
							name, _ := row["name"].(string)
							description, _ := row["description"].(string)
							dateCreated, _ := row["date_created"].(string)
							items = append(items, &locationarealist.LocationAreaItem{
								ID:          id,
								Name:        name,
								Description: description,
								Active:      active,
								DateCreated: dateCreated,
							})
						}
						return &locationarealist.LocationAreaListResult{Items: items, TotalItems: len(items)}, nil
					},
					GetInUseIDs: refChecker.GetLocationAreaInUseIDs,
					CreateLocationArea: func(fctx context.Context, name, description string, active bool) (string, error) {
						row, err := crudDB.Create(fctx, "location_area", map[string]any{
							"name":        name,
							"description": description,
							"active":      active,
						})
						if err != nil {
							return "", err
						}
						id, _ := row["id"].(string)
						return id, nil
					},
					ReadLocationArea: func(fctx context.Context, id string) (*locationareaaction.LocationAreaRecord, error) {
						row, err := crudDB.Read(fctx, "location_area", id)
						if err != nil {
							return nil, err
						}
						name, _ := row["name"].(string)
						description, _ := row["description"].(string)
						active, _ := row["active"].(bool)
						return &locationareaaction.LocationAreaRecord{
							ID:          id,
							Name:        name,
							Description: description,
							Active:      active,
						}, nil
					},
					UpdateLocationArea: func(fctx context.Context, id, name, description string, active bool) error {
						_, err := crudDB.Update(fctx, "location_area", id, map[string]any{
							"name":        name,
							"description": description,
							"active":      active,
						})
						return err
					},
					DeleteLocationArea: func(fctx context.Context, id string) error {
						return crudDB.Delete(fctx, "location_area", id)
					},
					SetLocationAreaActive: func(fctx context.Context, id string, active bool) error {
						_, err := crudDB.Update(fctx, "location_area", id, map[string]any{"active": active})
						return err
					},
				}).RegisterRoutes(ctx.Routes)
			}
		}

		if cfg.enableAll || cfg.permission {
			permissionmod.NewModule(&permissionmod.ModuleDeps{
				Routes:           routes.Permission,
				CommonLabels:     ctx.Common,
				SharedLabels:     labels.Shared,
				Labels:           labels.Permission,
				TableLabels:      ctx.Table,
				GetListPageData:  uc.Permission.GetListPageData,
				CreatePermission: uc.Permission.Create,
				ReadPermission:   uc.Permission.Read,
				UpdatePermission: uc.Permission.Update,
				DeletePermission: uc.Permission.Delete,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "permission", id, map[string]any{"active": active})
					return err
				},
			}).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.workspace {
			wsMod := &workspacemod.ModuleDeps{
				Routes:          routes.Workspace,
				CommonLabels:    ctx.Common,
				SharedLabels:    labels.Shared,
				Labels:          labels.Workspace,
				TableLabels:     ctx.Table,
				GetListPageData: uc.Workspace.GetListPageData,
				CreateWorkspace: uc.Workspace.Create,
				ReadWorkspace:   uc.Workspace.Read,
				UpdateWorkspace: uc.Workspace.Update,
				DeleteWorkspace: uc.Workspace.Delete,
				SetActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "workspace", id, map[string]any{"active": active})
					return err
				},
				// Phase 2 TODO closeout: wire the workspace_user detail + add URLs
				// now that Phase 2 has registered those route constants.
				WorkspaceUserDetailURL: entydad.WorkspaceUserDetailURL,
				WorkspaceUserAddURL:    entydad.WorkspaceUserAddURL,
				UploadFile:             uploadFile,
				ListAttachments:        listAttachments,
				CreateAttachment:       createAttachment,
				DeleteAttachment:       deleteAttachment,
				NewID:                  newAttachmentID,
			}
			if uc.WorkspaceUser.GetListPageData != nil {
				wsMod.GetWorkspaceUserListPageData = uc.WorkspaceUser.GetListPageData
			}
			workspacemod.NewModule(wsMod).RegisterRoutes(ctx.Routes)

			// Switch workspace (raw POST — uses session cookie, issues HX-Redirect)
			// Registers when EITHER the legacy SwitchWorkspace use case is wired
			// OR the host app has provided a SecureSwitch override (A1 fix
			// WKR-P0-1: service-admin wires SecureSwitch so the sidebar
			// workspace-switcher rotates + audits via executePrincipalSwitch).
			//
			// Two wire-up paths supported:
			//   1. BlockOption: WithSecureSwitch(fn, resolveUserID, setCookie)
			//      — explicit, used when host can construct closures before
			//      Block() is applied.
			//   2. AppContext fields: ctx.SecureWorkspaceSwitch + the two
			//      sibling fields — used when the host has the appBuilder
			//      ready inside buildAppContext() but constructs the entydad
			//      AppOption from a different call site. This lets
			//      service-admin populate them inside buildAppContext()
			//      without restructuring entydadBlock().
			secureSwitch := cfg.secureSwitch
			secureSwitchResolveUser := cfg.secureSwitchResolveUser
			secureSwitchSetCookie := cfg.secureSwitchSetCookie
			if secureSwitch == nil {
				if v, ok := ctx.SecureWorkspaceSwitch.(workspaceaction.SecureSwitchFn); ok {
					secureSwitch = v
				}
				if v, ok := ctx.SecureWorkspaceSwitchResolveUserID.(func(r *http.Request) string); ok {
					secureSwitchResolveUser = v
				}
				if v, ok := ctx.SecureWorkspaceSwitchSetSessionCookie.(func(w http.ResponseWriter, token string)); ok {
					secureSwitchSetCookie = v
				}
			}
			if uc.Workspace.Switch != nil || secureSwitch != nil {
				handleFunc(ctx.Routes, "POST", routes.Workspace.SwitchURL, workspaceaction.NewSwitchWorkspaceHandler(&workspaceaction.SwitchWorkspaceDeps{
					SecureSwitch:          secureSwitch,
					ResolveUserID:         secureSwitchResolveUser,
					SetSessionCookie:      secureSwitchSetCookie,
					SwitchWorkspace:       uc.Workspace.Switch,
					HomeURLForWorkspaceID: cfg.homeURLForWorkspaceID,
					HomeURL:               cfg.homeURL,
				}))
			}
		}

		if cfg.enableAll || cfg.workspaceUser {
			if uc.WorkspaceUser.GetListPageData == nil {
				log.Println("entydad.Block: warning: workspace_user use cases not initialized — workspace_user detail routes will be unavailable")
			} else {
				wuRoutes := routes.WorkspaceUser
				wuMod := &workspaceusermod.ModuleDeps{
					Routes:                       wuRoutes,
					WorkspaceDetailURL:           entydad.WorkspaceDetailURL,
					CommonLabels:                 ctx.Common,
					Labels:                       labels.WorkspaceUser,
					TableLabels:                  ctx.Table,
					GetListPageData:              uc.WorkspaceUser.GetListPageData,
					GetWorkspaceUserItemPageData: uc.WorkspaceUser.GetItemPageData,
					CreateWorkspaceUser:          uc.WorkspaceUser.Create,
					DeleteWorkspaceUser:          uc.WorkspaceUser.Delete,
					SetWorkspaceUserActive: func(fctx context.Context, id string, active bool) error {
						_, err := db.Update(fctx, "workspace_user", id, map[string]any{"active": active})
						return err
					},
					// Phase 3 closeout: wire WorkspaceUserRole routes now that Phase 3 has registered them.
					WorkspaceUserRoleAddURL:    entydad.WorkspaceUserRoleAddURL,
					WorkspaceUserRoleDeleteURL: entydad.WorkspaceUserRoleDeleteURL,
					UploadFile:                 uploadFile,
					ListAttachments:            listAttachments,
					CreateAttachment:           createAttachment,
					DeleteAttachment:           deleteAttachment,
					NewID:                      newAttachmentID,
				}
				// ListUsers — needed for the user-search autocomplete on the add form.
				if uc.User.List != nil {
					wuMod.ListUsers = uc.User.List
				}
				// Phase 3 closeout: wire workspace_user_role list page data.
				if uc.WorkspaceUserRole.GetListPageData != nil {
					wuMod.GetWorkspaceUserRoleListPageData = uc.WorkspaceUserRole.GetListPageData
				}
				workspaceusermod.NewModule(wuMod).RegisterRoutes(ctx.Routes)
				log.Println("  ✓ WorkspaceUser module initialized (entydad.Block)")
			}
		}

		if cfg.enableAll || cfg.workspaceUserRole {
			if uc.WorkspaceUserRole.Create == nil {
				log.Println("entydad.Block: warning: workspace_user_role use cases not initialized — workspace_user_role drawer routes will be unavailable")
			} else {
				wurRoutes := routes.WorkspaceUserRole
				wurMod := &workspaceuserrolemod.ModuleDeps{
					Routes:                  wurRoutes,
					Labels:                  labels.WorkspaceUserRole,
					CommonLabels:            ctx.Common,
					CreateWorkspaceUserRole: uc.WorkspaceUserRole.Create,
					DeleteWorkspaceUserRole: uc.WorkspaceUserRole.Delete,
				}
				if uc.WorkspaceUser.GetItemPageData != nil {
					wurMod.GetWorkspaceUserItemPageData = uc.WorkspaceUser.GetItemPageData
				}
				if uc.Role.List != nil {
					wurMod.ListRoles = uc.Role.List
				}
				workspaceuserrolemod.NewModule(wurMod).RegisterRoutes(ctx.Routes)
				log.Println("  ✓ WorkspaceUserRole module initialized (entydad.Block)")
			}
		}

		if cfg.enableAll || cfg.supplier {
			supplierDeps := &suppliermod.ModuleDeps{
				Routes:            routes.Supplier,
				SupplierTagRoutes: routes.SupplierTag,
				// User module owns the timezone search endpoint; the supplier
				// representative form reuses the same JSON handler.
				SearchTimezonesURL: routes.User.SearchTimezonesURL,
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
				ListPaymentTerms: func(fctx context.Context) ([]*suppliermod.PaymentTermOption, error) {
					rows, err := db.ListSimple(fctx, "payment_term")
					if err != nil {
						return nil, err
					}
					opts := make([]*suppliermod.PaymentTermOption, 0, len(rows))
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
						opts = append(opts, &suppliermod.PaymentTermOption{Id: id, Name: name})
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
			suppliermod.NewModule(supplierDeps).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.clientTag {
			clienttagDeps := &clienttagmod.ModuleDeps{
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
			clienttagmod.NewModule(clienttagDeps).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.supplierTag {
			suppliertagDeps := &suppliertagmod.ModuleDeps{
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
			suppliertagmod.NewModule(suppliertagDeps).RegisterRoutes(ctx.Routes)
		}

		if cfg.enableAll || cfg.paymentTerm {
			if ctx.SqlDB == nil {
				log.Println("entydad.Block: warning: SqlDB is nil — skipping payment_term module")
			} else {
				repoAny, err := registry.CreateRepository("postgresql", entityid.PaymentTerm, ctx.SqlDB, entityid.PaymentTerm)
				if err != nil {
					return fmt.Errorf("entydad.Block: failed to create payment_term repository: %w", err)
				}
				ptRepo, ok := repoAny.(paymenttermpb.PaymentTermDomainServiceServer)
				if !ok {
					return fmt.Errorf("entydad.Block: payment_term repository does not implement PaymentTermDomainServiceServer")
				}
				setPaymentTermActive := func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "payment_term", id, map[string]any{"active": active})
					return err
				}
				// Client-context payment term list: shows terms with entity_scope IN ('client', 'both')
				paymenttermmod.NewModule(&paymenttermmod.ModuleDeps{
					Routes:               routes.PaymentTerm,
					CommonLabels:         ctx.Common,
					SharedLabels:         labels.Shared,
					Labels:               labels.PaymentTerm,
					TableLabels:          ctx.Table,
					GetListPageData:      ptRepo.GetPaymentTermListPageData,
					GetInUseIDs:          refChecker.GetPaymentTermInUseIDs,
					CreatePaymentTerm:    ptRepo.CreatePaymentTerm,
					ReadPaymentTerm:      ptRepo.ReadPaymentTerm,
					UpdatePaymentTerm:    ptRepo.UpdatePaymentTerm,
					DeletePaymentTerm:    ptRepo.DeletePaymentTerm,
					SetPaymentTermActive: setPaymentTermActive,
					Scope:                "client",
				}).RegisterRoutes(ctx.Routes)
				// Supplier-context payment term list: shows terms with entity_scope IN ('supplier', 'both')
				paymenttermmod.NewModule(&paymenttermmod.ModuleDeps{
					Routes:               routes.SupplierPaymentTerm.ToPaymentTermRoutes(),
					CommonLabels:         ctx.Common,
					SharedLabels:         labels.Shared,
					Labels:               labels.PaymentTerm,
					TableLabels:          ctx.Table,
					GetListPageData:      ptRepo.GetPaymentTermListPageData,
					GetInUseIDs:          refChecker.GetPaymentTermInUseIDs,
					CreatePaymentTerm:    ptRepo.CreatePaymentTerm,
					ReadPaymentTerm:      ptRepo.ReadPaymentTerm,
					UpdatePaymentTerm:    ptRepo.UpdatePaymentTerm,
					DeletePaymentTerm:    ptRepo.DeletePaymentTerm,
					SetPaymentTermActive: setPaymentTermActive,
					Scope:                "supplier",
				}).RegisterRoutes(ctx.Routes)
			}
		}

		if cfg.enableAll || cfg.admin {
			adminDeps := &adminmod.ModuleDeps{
				Routes:               routes.Admin,
				CommonLabels:         ctx.Common,
				DashboardLabels:      labels.Admin,
				DashboardTitleLabels: labels.Dashboard,
				// Cross-entity URLs for quick actions and "view all" links.
				DashboardRoutes: admindashboardroutes.Routes{
					DashboardURL:         routes.Admin.DashboardURL,
					NewUserURL:           routes.WorkspaceUser.AddURL,
					NewWorkspaceURL:      routes.Workspace.AddURL,
					AssignRoleURL:        routes.WorkspaceUserRole.AddURL,
					PermissionListURL:    routes.Permission.ListURL,
					RoleListURL:          routes.Role.ListURL,
					WorkspaceListURL:     routes.Workspace.ListURL,
					WorkspaceUserListURL: routes.WorkspaceUser.ListURL,
				},
			}
			if uc.GetAdminDashboardPageData != nil {
				adminDeps.GetDashboardData = uc.GetAdminDashboardPageData
			}
			adminmod.NewModule(adminDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Tax Registration module (entydad — polymorphic client + workspace)
		// =====================================================================

		if cfg.enableAll || cfg.taxRegistration {
			taxRegDeps := &taxregistrationmod.ModuleDeps{
				Routes:       routes.TaxRegistration,
				Labels:       labels.TaxRegistration,
				CommonLabels: ctx.Common,
				TableLabels:  ctx.Table,
			}
			if uc.TaxRegistration.List != nil {
				taxRegDeps.ListTaxRegistrations = uc.TaxRegistration.List
			}
			taxregistrationmod.NewModule(taxRegDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Conversation module (entydad — secure messaging / ticketing, Plan-4)
		// =====================================================================

		if cfg.enableAll || cfg.conversation {
			convDeps := &conversationmod.ModuleDeps{
				Routes:                routes.Conversation,
				CommonLabels:          ctx.Common,
				TableLabels:           ctx.Table,
				Labels:                labels.Conversation,
				PostLabels:            labels.ConversationPost,
				ListConversations:     uc.Conversation.List,
				ReadConversation:      uc.Conversation.Read,
				CreateConversation:    uc.Conversation.Create,
				AssignConversation:    uc.Conversation.Assign,
				SetConversationStatus: uc.Conversation.SetStatus,
				ListConversationPosts: uc.Conversation.Post.List,
				SendConversationPost:  uc.Conversation.Post.Send,
				MarkConversationRead:  uc.Conversation.Receipt.MarkRead,
				NewClientToken:        newAttachmentID,
				// Staff new-conversation drawer reuses existing autocomplete
				// endpoints: client search (client module) + workspace-user
				// search (role module) for the assignee picker.
				ClientSearchURL:   routes.Client.SearchURL,
				AssigneeSearchURL: routes.Role.UsersSearchURL,
				// Client-portal row-scope: read the session's acting-as-client
				// id from context. The host's view_adapter populates it per
				// request via consumer.WithActingAsClientID (sourced from the
				// session binding / derived from client_portal_grant for a
				// direct client). appcontext is the dependency-free leaf so the
				// block never imports consumer. Fail-closed: "" => portal denies.
				ActingAsClientID: appcontext.GetActingAsClientIDFromContext,
			}
			// ClientNameByID — best-effort display-name resolver for the inbox
			// Client column. Backed by a single ListSimple("client") scan.
			if crudDB, ok := db.(CRUDSource); ok {
				convDeps.ClientNameByID = func(fctx context.Context, ids []string) map[string]string {
					out := map[string]string{}
					if len(ids) == 0 {
						return out
					}
					want := make(map[string]struct{}, len(ids))
					for _, id := range ids {
						want[id] = struct{}{}
					}
					rows, err := crudDB.ListSimple(fctx, "client")
					if err != nil {
						return out
					}
					for _, row := range rows {
						id, _ := row["id"].(string)
						if _, ok := want[id]; !ok {
							continue
						}
						if name, _ := row["name"].(string); name != "" {
							out[id] = name
						}
					}
					return out
				}
			}
			convMod := conversationmod.NewModule(convDeps)
			convMod.RegisterRoutes(ctx.Routes)

			// Client-portal routes are GATED on AUTHZ_ENFORCE=true. The second
			// precondition — the inherited 20260601 Phase-4 acting_as_client_id
			// wiring for direct PRINCIPAL_TYPE_CLIENT principals — is now MET:
			// the host's view_adapter calls consumer.WithActingAsClientID per
			// request (sourced from the session binding and, for a direct
			// client, derived read-only from client_portal_grant in
			// composition.lookupSessionPrincipalFull); convDeps.ActingAsClientID
			// reads it back via appcontext. The portal handlers ALSO fail-closed
			// on an empty acting_as_client_id, so this stays defence-in-depth.
			const portalPhase4Ready = true
			authzEnforce := os.Getenv("AUTHZ_ENFORCE") == "true" ||
				os.Getenv("AUTHZ_ENFORCE") == "1" ||
				os.Getenv("AUTHZ_ENFORCE") == "yes"
			if authzEnforce && portalPhase4Ready && convMod.PortalReady() {
				convMod.RegisterPortalRoutes(ctx.Routes)
			} else if authzEnforce {
				// Register a 503 stub so the URL exists for discovery tests but
				// never serves client data until Phase-4 lands.
				handleFunc(ctx.Routes, "GET", routes.Conversation.PortalListURL,
					func(w http.ResponseWriter, r *http.Request) {
						http.Error(w, "Client portal messaging not yet available", http.StatusServiceUnavailable)
					})
			}
		}

		log.Println("  ✓ Entity domain initialized (entydad.Block)")
		return nil
	}
}
