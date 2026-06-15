package block

// auth_bridge.go — type bridges between espyna/adapthttp and entydad's auth
// module interfaces, plus the auth.Deps builder that registers the auth module
// DIRECTLY inside the entydad EngineBlock closure (Wave B D2a / D2-β).
//
// Relocated from app composition/auth_bridge.go. The former *appBuilder
// coupling is replaced by authChainDeps, which carries only the typed values
// the chain needs — sourced from the pyeza.AppContext slots the host stamps
// (UseCases, SessionManager, AuthAdapter, CSRFSecret/CSRFIssuer/CookieSecure,
// UI.Renderer, UI.AuthLabels). The type bridges are security-critical — they
// preserve every invariant from the original authAdapterBridge /
// principalResolverAdapter wiring.

import (
	"context"

	"github.com/erniealice/entydad-golang/service/auth"
	consumer "github.com/erniealice/espyna-golang/consumer"

	adapthttp "github.com/erniealice/espyna-golang/consumer/http"
)

// authChainDeps carries the typed values the relocated auth chain needs,
// replacing the former *appBuilder coupling. Built once inside the entydad
// EngineBlock closure from the pyeza.AppContext slots.
type authChainDeps struct {
	uc           *consumer.UseCases
	sessionMw    *consumer.SessionMiddleware
	authAdapter  *consumer.AuthAdapter
	renderer     auth.Renderer
	authLabels   auth.AuthLabels
	csrfSecret   []byte
	csrfIssuer   auth.CSRFIssuer
	cookieSecure bool
}

// authAdapterBridge adapts *consumer.AuthAdapter to auth.AuthAdapter.
// The only mismatch is Login's return type: consumer returns *authpb.Identity
// while auth.AuthAdapter returns auth.AuthIdentity. Since *authpb.Identity
// has GetId(), it already satisfies auth.AuthIdentity — the bridge just
// widens the return type.
type authAdapterBridge struct {
	inner *consumer.AuthAdapter
}

func (a *authAdapterBridge) Login(ctx context.Context, email, password string) (string, auth.AuthIdentity, error) {
	token, identity, err := a.inner.Login(ctx, email, password)
	if err != nil {
		return "", nil, err
	}
	// *authpb.Identity satisfies auth.AuthIdentity via its GetId() method.
	return token, identity, nil
}

func (a *authAdapterBridge) Register(ctx context.Context, email, password, firstName, lastName, mobileNumber string) (string, error) {
	return a.inner.Register(ctx, email, password, firstName, lastName, mobileNumber)
}

func (a *authAdapterBridge) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	return a.inner.RequestPasswordReset(ctx, email)
}

func (a *authAdapterBridge) ExecutePasswordReset(ctx context.Context, token, newPassword string) error {
	return a.inner.ExecutePasswordReset(ctx, token, newPassword)
}

func (a *authAdapterBridge) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	return a.inner.ChangePassword(ctx, userID, oldPassword, newPassword)
}

func (a *authAdapterBridge) ValidateSession(ctx context.Context, token string) (string, error) {
	return a.inner.ValidateSession(ctx, token)
}

func (a *authAdapterBridge) InvalidateSession(ctx context.Context, token string) error {
	return a.inner.InvalidateSession(ctx, token)
}

// principalResolverAdapter adapts adapthttp.DBPrincipalLoader (which returns
// []adapthttp.Principal) to auth.PrincipalResolver (which returns
// []auth.Principal). The two Principal types have identical fields but live
// in different packages — this adapter bridges them.
type principalResolverAdapter struct {
	loader adapthttp.PrincipalLoader
}

func (a *principalResolverAdapter) Resolve(ctx context.Context, userID string) ([]auth.Principal, error) {
	adaptPrincipals, err := a.loader.Resolve(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]auth.Principal, len(adaptPrincipals))
	for i, p := range adaptPrincipals {
		out[i] = toAuthPrincipal(p)
	}
	return out, nil
}

func (a *principalResolverAdapter) IsEnabled() bool {
	return a.loader.IsEnabled()
}

// toAuthPrincipal converts an adapthttp.Principal to an auth.Principal.
func toAuthPrincipal(p adapthttp.Principal) auth.Principal {
	targets := make([]auth.ActingAsTarget, len(p.ActingAsTargets))
	for i, t := range p.ActingAsTargets {
		targets[i] = auth.ActingAsTarget{
			ID:          t.ID,
			WorkspaceID: t.WorkspaceID,
			DisplayName: t.DisplayName,
		}
	}
	return auth.Principal{
		Type:            p.Type,
		PrincipalID:     p.PrincipalID,
		WorkspaceID:     p.WorkspaceID,
		DisplayName:     p.DisplayName,
		ActingAsTargets: targets,
	}
}

// toAdaptPrincipal converts an auth.Principal to an adapthttp.Principal.
func toAdaptPrincipal(p auth.Principal) adapthttp.Principal {
	targets := make([]adapthttp.ActingAsTarget, len(p.ActingAsTargets))
	for i, t := range p.ActingAsTargets {
		targets[i] = adapthttp.ActingAsTarget{
			ID:          t.ID,
			WorkspaceID: t.WorkspaceID,
			DisplayName: t.DisplayName,
		}
	}
	return adapthttp.Principal{
		Type:            p.Type,
		PrincipalID:     p.PrincipalID,
		WorkspaceID:     p.WorkspaceID,
		DisplayName:     p.DisplayName,
		ActingAsTargets: targets,
	}
}

// homeURLForWorkspaceID looks up the workspace slug for a workspace_id and
// returns the post-login home URL: /w/{slug}/home. Falls back to /me/inbox
// (workspace-less personal landing per Q-WS-7 → B) if the lookup fails or
// returns empty — that path always works regardless of slug state.
//
// Relocated from app composition/auth_bridge.go; principal_switch_sidebar.go
// (block) also calls it.
func (d *authChainDeps) homeURLForWorkspaceID(ctx context.Context, workspaceID string) string {
	const fallback = "/me/inbox"
	if d == nil || d.uc == nil || workspaceID == "" {
		return fallback
	}
	if d.uc.Entity == nil || d.uc.Entity.Workspace == nil {
		return fallback
	}
	uc := d.uc.Entity.Workspace.ResolveSlugByWorkspaceID
	if uc == nil {
		return fallback
	}
	slug, err := uc.Execute(ctx, workspaceID)
	if err != nil || slug == "" {
		return fallback
	}
	return "/w/" + slug + "/home"
}

// homeURLForWorkspaceIDFn returns the entydad-block resolver closure used by
// WithHomeURLForWorkspaceID (the post-switch redirect resolver). It mirrors the
// app's homeURLForWorkspaceID, sourced from the same workspace use case.
func (d *authChainDeps) homeURLForWorkspaceIDFn() func(ctx context.Context, workspaceID string) string {
	return func(ctx context.Context, workspaceID string) string {
		return d.homeURLForWorkspaceID(ctx, workspaceID)
	}
}

// buildAuthDeps reconstructs the full 19-field auth.Deps from the authChainDeps.
// Every field is sourced per wave-b-deep-plan.md §D2.3a. Called inside the
// entydad EngineBlock closure (D2-β: auth registers directly, decoupled from
// the entity overlay).
//
// PRE-FLIGHT ASSERTS (D2a precondition P3): a non-mock boot with a nil renderer,
// empty auth-label set, nil session manager, nil CSRF issuer, or empty CSRF
// secret boot-FATALS. espyna's finalize must-asserts do NOT cover
// ctx.CSRFIssuer — the entydad block reads it — so a regressed slot must fail
// LOUD here, never ship a CSRF-less / blank-renderer login.
func (d *authChainDeps) buildAuthDeps() *auth.Deps {
	// Build the principal resolver adapter (bridges adapthttp → auth types).
	var resolver auth.PrincipalResolver
	if loader := d.newPrincipalLoader(); loader != nil {
		resolver = &principalResolverAdapter{loader: loader}
	}

	// Build the auth adapter bridge (widens *authpb.Identity → AuthIdentity).
	var authAdapter auth.AuthAdapter
	if d.authAdapter != nil {
		authAdapter = &authAdapterBridge{inner: d.authAdapter}
	}

	cookieSecure := d.cookieSecure
	return &auth.Deps{
		AuthAdapter:    authAdapter,
		SessionManager: d.sessionMw,

		PrincipalResolver: resolver,
		PrincipalSwitcher: func(ctx context.Context, input auth.PrincipalSwitchInput) (*auth.PrincipalSwitchResult, error) {
			result, err := d.executePrincipalSwitch(ctx, authPrincipalSwitchInput{
				UserID:             input.UserID,
				Token:              input.Token,
				TargetPrincipal:    toAdaptPrincipal(input.TargetPrincipal),
				ActingAsClientID:   input.ActingAsClientID,
				ActingAsSupplierID: input.ActingAsSupplierID,
				UseCase:            input.UseCase,
				RequestURL:         input.RequestURL,
				Referer:            input.Referer,
				SecFetchSite:       input.SecFetchSite,
				UserAgent:          input.UserAgent,
				RequireAudit:       input.RequireAudit,
			})
			if err != nil {
				return nil, err
			}
			return &auth.PrincipalSwitchResult{
				NewToken:    result.NewToken,
				RedirectURL: result.RedirectURL,
			}, nil
		},

		CSRFSecret: d.csrfSecret,
		CSRFIssuer: d.csrfIssuer,

		Renderer: d.renderer,

		UserIDByEmail: func(ctx context.Context, email string) string {
			if d.uc.Entity == nil || d.uc.Entity.User == nil {
				return ""
			}
			uc := d.uc.Entity.User.ResolveUserByEmail
			if uc == nil {
				return ""
			}
			id, _ := uc.Execute(ctx, email)
			return id
		},
		WorkspaceSlugResolver: func(ctx context.Context, wsID string) string {
			if d.uc.Entity == nil || d.uc.Entity.Workspace == nil {
				return ""
			}
			uc := d.uc.Entity.Workspace.ResolveSlugByWorkspaceID
			if uc == nil || wsID == "" {
				return ""
			}
			slug, _ := uc.Execute(ctx, wsID)
			return slug
		},

		Labels: d.authLabels,

		LogoText:     getEnv("ICHIZEN_LOGO_TEXT", "Ichizen"),
		AuthProvider: getEnv("CONFIG_AUTH_PROVIDER", ""),
		TestMode:     equalFoldTrue(getEnv("PASSWORD_AUTH_TEST_MODE", "")),
		// SecureCookies — derived from the host-resolved cookie-secure policy
		// (ctx.CookieSecure). Closure-captured so the auth module reads the SAME
		// policy the 4-arg CSRF issuer self-derives. WorkspaceCSRFCookieName left
		// "" → NewAuthModule defaults "ws_csrf" (the app default).
		SecureCookies: func() bool { return cookieSecure },
	}
}

// equalFoldTrue reports whether s case-insensitively equals "true".
func equalFoldTrue(s string) bool {
	return len(s) == 4 &&
		(s[0]|0x20) == 't' && (s[1]|0x20) == 'r' && (s[2]|0x20) == 'u' && (s[3]|0x20) == 'e'
}
