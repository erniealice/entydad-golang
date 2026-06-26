package block

// engineblock.go — the entydad EngineBlock(opts...) pyeza.AppOption (Wave B D2a).
//
// Relocated + consolidated from app composition/adapters_entydad.go
// (entydadEngineBlock) + auth_bridge.go + principal_switch*.go +
// principal_loader_bridge.go + workspace_loader.go. It:
//
//  1. maps espyna's *consumer.UseCases -> entydad block *UseCases (the mapper in
//     usecases_from_consumer.go),
//  2. builds the auth chain deps from the pyeza.AppContext slots the host stamps
//     (UseCases, SessionManager, AuthAdapter, CSRFSecret/CSRFIssuer/CookieSecure,
//     UI.Renderer, UI.AuthLabels),
//  3. registers the auth module DIRECTLY (D2-β) — decoupled from the entity
//     overlay's all-or-nothing Assemble — asserting ctx.Routes to
//     auth.RouteRegistrar or boot-FATAL,
//  4. wires the secure sidebar workspace-switch (D2a-4) + WorkspaceLoader into
//     Infra,
//  5. runs AllUnits -> AssembleEngineBlock (preserving the ComposeResult
//     route-map merge; #18 fail-loud lives in the shared helper).
//
// infra.AuthDeps STAYS nil so AuthUnit.Mount is the inert no-op (no
// double-registration). The entity overlay can fail without locking out
// /auth/login (the auth path registered above, before AllUnits).

import (
	"context"
	"log"
	"net/http"
	"os"

	roleusers "github.com/erniealice/entydad-golang/domain/entity/identity/role/users"
	userdashboard "github.com/erniealice/entydad-golang/domain/entity/identity/user/dashboard"
	"github.com/erniealice/entydad-golang/service/auth"
	consumer "github.com/erniealice/espyna-golang/consumer"
	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/espyna-golang/ports"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	pytypes "github.com/erniealice/pyeza-golang/types"
)

// EngineBlockOption configures the entydad EngineBlock.
type EngineBlockOption func(*engineBlockConfig)

type engineBlockConfig struct {
	homeURL string
}

// WithEngineHomeURL sets the static fallback post-switch redirect URL (the home
// dashboard URL). The dynamic /w/{slug}/home resolver is derived internally from
// the workspace use case. (Distinct from the legacy BlockOption WithHomeURL.)
func WithEngineHomeURL(url string) EngineBlockOption {
	return func(c *engineBlockConfig) { c.homeURL = url }
}

// EngineBlock returns a pyeza.AppOption that registers all entydad entity domain
// modules via the compose engine AND registers the auth module directly (D2-β).
// Replaces the app-side entydadEngineBlock wrapper.
func EngineBlock(opts ...EngineBlockOption) consumerapp.AppOption {
	cfg := &engineBlockConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return func(ctx *consumerapp.AppContext) error {
		uc, err := consumerapp.RequireUseCases(ctx, "entydad.EngineBlock")
		if err != nil {
			return err
		}

		// ── Auth chain deps (from the pyeza.AppContext slots) ───────────────
		acd := buildAuthChainDeps(ctx, uc)

		// ── D2-β: register auth DIRECTLY, decoupled from the entity overlay ──
		// PRE-FLIGHT ASSERT (P2 / D2-β): ctx.Routes MUST satisfy auth.RouteRegistrar
		// (view.RouteRegistrar + HandleFunc) or boot fails LOUD — never serve with
		// /auth/* unregistered.
		ar, ok := ctx.Routes.(auth.RouteRegistrar)
		if !ok {
			log.Fatalf("FATAL entydad.EngineBlock: ctx.Routes (%T) does not satisfy "+
				"auth.RouteRegistrar (HandleFunc) — refusing to boot with /auth/* unregistered.", ctx.Routes)
		}
		auth.NewAuthModule(acd.buildAuthDeps()).RegisterRoutes(ar)

		// ── Infra ───────────────────────────────────────────────────────────
		infra := &Infra{}
		infra.UploadFile, _ = ctx.UploadFile.(func(context.Context, string, string, []byte, string) error)
		infra.ListAttachments, _ = ctx.ListAttachments.(func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error))
		infra.CreateAttachment, _ = ctx.CreateAttachment.(func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		infra.DeleteAttachment, _ = ctx.DeleteAttachment.(func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error))
		infra.NewAttachmentID, _ = ctx.NewAttachmentID.(func() string)
		if ctx.RefChecker != nil {
			if rc, ok := ctx.RefChecker.(ports.Checker); ok {
				infra.RefChecker = rc
			}
		}
		infra.GetUsersByRoleID, _ = ctx.GetUsersByRoleID.(func(context.Context, string) ([]roleusers.UserByRole, error))
		infra.GetDashboardData, _ = ctx.GetDashboardData.(func(context.Context) (*userdashboard.DashboardData, error))
		infra.HashPassword, _ = ctx.HashPassword.(func(string) (string, error))
		infra.GetUserWorkspacesMap, _ = ctx.GetUserWorkspacesMap.(func(context.Context) (map[string][]pytypes.ChipData, error))

		// HomeURLForWorkspaceID (post-switch /w/{slug}/home resolver).
		infra.HomeURLForWorkspaceID = acd.homeURLForWorkspaceIDFn()

		// SecureSwitch (D2a-4) — wire the secure sidebar workspace-switch
		// closures into Infra so WorkspaceUnit's switch handler routes through
		// the rotation+audit primitive. All-three-or-none (block.go semantics):
		// only set when the secure primitive is available for this build/dialect.
		if acd.secureSidebarSwitchWired() {
			infra.SecureSwitch = acd.secureSidebarSwitchFn()
			infra.SecureSwitchResolveUser = secureSidebarResolveUserID
			infra.SecureSwitchSetCookie = acd.secureSidebarSetSessionCookie
		}

		// AuthDeps STAYS nil — AuthUnit.Mount is the inert no-op; auth was
		// already registered directly above (no double-registration).
		infra.AuthDeps = nil

		// ── WorkspaceLoader (proto-backed; ctx slot for the Server finalize) ──
		// DBWorkspaceLoader satisfies consumerhttp.WorkspaceLoader structurally.
		if uc.Entity != nil && uc.Entity.Workspace != nil && uc.Entity.Workspace.ListUserWorkspaces != nil {
			ctx.WorkspaceLoader = NewDBWorkspaceLoader(uc.Entity.Workspace.ListUserWorkspaces)
		}

		// ── Map use cases + assemble (preserves the ComposeResult merge) ──────
		adapted := buildEntydadUseCases(uc, ctx.DB)
		units := AllUnits(adapted, infra)
		return consumerapp.AssembleEngineBlock("entydad", units, ctx)
	}
}

// buildAuthChainDeps assembles the authChainDeps from the pyeza.AppContext
// slots the host stamps. Performs the D2a precondition-P3 pre-flight asserts:
// a non-mock boot with a missing/wrong-type renderer, auth-label set, session
// manager, CSRF issuer, or CSRF secret boot-FATALS (espyna finalize asserts do
// NOT cover ctx.CSRFIssuer — the entydad block reads it).
func buildAuthChainDeps(ctx *consumerapp.AppContext, uc *consumer.UseCases) *authChainDeps {
	provider := getEnv("CONFIG_AUTH_PROVIDER", "")
	nonMock := provider != "" && provider != "mock"

	// UI bundle (Renderer + AuthLabels).
	ui, _ := ctx.UI.(*consumerapp.AppUIBundle)
	if nonMock && ui == nil {
		log.Fatalf("FATAL entydad.EngineBlock: ctx.UI is %T, want non-nil *pyeza.AppUIBundle "+
			"(the host must stamp the auth-half UI bundle). Refusing to boot auth with a blank renderer/labels.", ctx.UI)
	}

	var renderer auth.Renderer
	var authLabels auth.AuthLabels
	if ui != nil {
		r, rok := ui.Renderer.(auth.Renderer)
		if nonMock && !rok {
			log.Fatalf("FATAL entydad.EngineBlock: ctx.UI.Renderer is %T, want auth.Renderer "+
				"(*pyeza.HTMLRenderer) — refusing to boot the auth-shell render path with a nil renderer.", ui.Renderer)
		}
		renderer = r
		al, alok := ui.AuthLabels.(auth.AuthLabels)
		if nonMock && !alok {
			log.Fatalf("FATAL entydad.EngineBlock: ctx.UI.AuthLabels is %T, want auth.AuthLabels "+
				"— refusing to boot the login screen with blank labels.", ui.AuthLabels)
		}
		authLabels = al
	}

	// CSRF (D1.5 slots). The host stamps ctx.CSRFIssuer as the app's CONCRETE
	// 4-arg issuer (middleware.IssueWorkspaceCSRFCookie), whose dynamic type is
	// the UNNAMED func signature — NOT the named auth.CSRFIssuer. Assert against
	// the underlying signature, then convert; a direct .(auth.CSRFIssuer) assert
	// would fail (Go requires the exact named type for that).
	csrfSecret, _ := ctx.CSRFSecret.([]byte)
	var csrfIssuer auth.CSRFIssuer
	issuerFn, issuerOK := ctx.CSRFIssuer.(func(w http.ResponseWriter, secret []byte, sessionToken, workspaceID string) string)
	if issuerOK {
		csrfIssuer = issuerFn
	}
	if nonMock && !issuerOK {
		log.Fatalf("FATAL entydad.EngineBlock: ctx.CSRFIssuer is %T, want "+
			"func(http.ResponseWriter, []byte, string, string) string (the 4-arg workspace-CSRF "+
			"issuer) — refusing to boot a CSRF-less /auth/* + switch surface.", ctx.CSRFIssuer)
	}
	if nonMock && len(csrfSecret) == 0 {
		log.Fatalf("FATAL entydad.EngineBlock: provider %q runs the workspace-CSRF chain but ctx.CSRFSecret "+
			"is empty — refusing to register auth + switch with ws_csrf silently absent (fail-open).", provider)
	}
	cookieSecure, _ := ctx.CookieSecure.(bool)

	// Session manager + auth adapter (P1 slots).
	sessionMw, _ := ctx.SessionManager.(*consumer.SessionMiddleware)
	if nonMock && sessionMw == nil {
		log.Fatalf("FATAL entydad.EngineBlock: ctx.SessionManager is %T, want non-nil *consumer.SessionMiddleware "+
			"— refusing to boot auth without a session manager.", ctx.SessionManager)
	}
	authAdapter, _ := ctx.AuthAdapter.(*consumer.AuthAdapter)

	return &authChainDeps{
		uc:           uc,
		sessionMw:    sessionMw,
		authAdapter:  authAdapter,
		renderer:     renderer,
		authLabels:   authLabels,
		csrfSecret:   csrfSecret,
		csrfIssuer:   csrfIssuer,
		cookieSecure: cookieSecure,
	}
}

// getEnv reads an environment variable with a fallback default.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
