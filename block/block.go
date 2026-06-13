// Package block provides the Block() function for registering all entydad
// entity modules into a pyeza app via AppOption composition.
//
// It lives in a sub-package to avoid the import cycle that would arise if it
// were placed in the root entydad package (the view sub-packages already import
// the root package for route/label types).
//
// Companion files in this directory:
//   - helpers.go       — the categoryListPageDataGetter local interface and
//     getDefaultWorkspaceID. (The former DataSource/UpdateableSource/CRUDSource
//     ducks were deleted 2026-06-12; active/status writes now go through the
//     narrow typed UseCases.SetActive/SetStatus primitives bound by
//     service-admin.)
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

	"github.com/erniealice/espyna-golang/shared/identity"

	roleusers "github.com/erniealice/entydad-golang/domain/entity/identity/role/users"
	userdashboard "github.com/erniealice/entydad-golang/domain/entity/identity/user/dashboard"
	workspaceaction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/action"
	tax "github.com/erniealice/entydad-golang/domain/tax"
	adminmod "github.com/erniealice/entydad-golang/service/dashboard/views/admin"
	admindashboardroutes "github.com/erniealice/entydad-golang/service/dashboard/views/admin/dashboard"
	"github.com/erniealice/espyna-golang/reference"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/client_statement"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
	conversationmod "github.com/erniealice/hybra-golang/views/conversation"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
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
	secureSwitch            workspaceaction.SecureSwitchFn
	secureSwitchResolveUser func(r *http.Request) string
	secureSwitchSetCookie   func(w http.ResponseWriter, token string)
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
		// FAIL-CLOSED completeness gate (mirrors AUTHZ_ENFORCE boot-guard):
		// dev/test PANIC, prod log-FATAL + return so boot halts.
		if err := cfg.useCases.MustValidate(cfg); err != nil {
			return err
		}
		uc := cfg.useCases // local alias for brevity

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

		// Party sub-context (client, supplier, client_tag, supplier_tag).
		// Lifted to party.go (block-go-anatomy: B split). Pure code-move.
		wirePartyModule(ctx, partyWiring{
			cfg:                  cfg,
			uc:                   uc,
			labels:               labels,
			routes:               routes,
			refChecker:           refChecker,
			getClientStatement:   getClientStatement,
			getClientBalances:    getClientBalances,
			getSupplierStatement: getSupplierStatement,
			getSupplierBalances:  getSupplierBalances,
			uploadFile:           uploadFile,
			listAttachments:      listAttachments,
			createAttachment:     createAttachment,
			deleteAttachment:     deleteAttachment,
			newAttachmentID:      newAttachmentID,
		})

		// Commerce/location sub-context (location, location_area, payment_term).
		// Lifted to commerce.go (block-go-anatomy: B split). Pure code-move.
		if err := wireCommerceModule(ctx, commerceWiring{
			cfg:              cfg,
			uc:               uc,
			labels:           labels,
			routes:           routes,
			refChecker:       refChecker,
			uploadFile:       uploadFile,
			listAttachments:  listAttachments,
			createAttachment: createAttachment,
			deleteAttachment: deleteAttachment,
			newAttachmentID:  newAttachmentID,
		}); err != nil {
			return err
		}

		// Identity sub-context (user, role, permission, workspace,
		// workspace_user, workspace_user_role). Lifted to identity.go
		// (block-go-anatomy: B split). Pure code-move.
		wireIdentityModule(ctx, identityWiring{
			cfg:                  cfg,
			uc:                   uc,
			labels:               labels,
			routes:               routes,
			refChecker:           refChecker,
			getUserWorkspacesMap: getUserWorkspacesMap,
			getDashboardData:     getDashboardData,
			hashPassword:         hashPassword,
			getUsersByRoleID:     getUsersByRoleID,
			uploadFile:           uploadFile,
			listAttachments:      listAttachments,
			createAttachment:     createAttachment,
			deleteAttachment:     deleteAttachment,
			newAttachmentID:      newAttachmentID,
		})

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
			taxRegDeps := &tax.TaxRegistrationModuleDeps{
				Routes:       routes.TaxRegistration,
				Labels:       labels.TaxRegistration,
				CommonLabels: ctx.Common,
				TableLabels:  ctx.Table,
			}
			if uc.TaxRegistration.List != nil {
				taxRegDeps.ListTaxRegistrations = uc.TaxRegistration.List
			}
			tax.NewTaxRegistrationModule(taxRegDeps).RegisterRoutes(ctx.Routes)
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
				// direct client). identity is the dependency-free leaf so the
				// block never imports consumer. Fail-closed: "" => portal denies.
				ActingAsClientID: func(ctx context.Context) string {
					return identity.Must(ctx).ActingAsClientID
				},
			}
			// ClientNameByID — best-effort display-name resolver for the inbox
			// Client column. Backed by a single typed ListClients scan
			// (uc.Client.List), replacing the deleted duck's
			// ListSimple("client"). Wired only when the typed list is bound;
			// nil-safe — the inbox Client column falls back to ids otherwise.
			if listClients := uc.Client.List; listClients != nil {
				convDeps.ClientNameByID = func(fctx context.Context, ids []string) map[string]string {
					out := map[string]string{}
					if len(ids) == 0 {
						return out
					}
					want := make(map[string]struct{}, len(ids))
					for _, id := range ids {
						want[id] = struct{}{}
					}
					resp, err := listClients(fctx, &clientpb.ListClientsRequest{})
					if err != nil || resp == nil {
						return out
					}
					for _, c := range resp.GetData() {
						if c == nil {
							continue
						}
						id := c.GetId()
						if _, ok := want[id]; !ok {
							continue
						}
						if name := c.GetName(); name != "" {
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
			// reads it back via identity. The portal handlers ALSO fail-closed
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
