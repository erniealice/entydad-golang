// Package block — principal_switch_sidebar.go (relocated from app composition, Wave B D2a)
//
// Sidebar workspace-switcher wiring: routes POST /action/admin/switch-workspace
// through the secure principal-switch primitive (executePrincipalSwitch) instead
// of the legacy in-place SwitchWorkspace use case. Resolves the target binding
// via ResolveBindingInWorkspace, rotates the session token on cross-workspace
// boundary, and writes a required audit row. Wired into the WorkspaceUnit switch
// handler via Infra.SecureSwitch (entydad EngineBlock, D2a-4).

package block

import (
	"context"
	"errors"
	"net/http"

	workspaceaction "github.com/erniealice/entydad-golang/domain/entity/identity/workspace/action"
	consumer "github.com/erniealice/espyna-golang/consumer"

	adapthttp "github.com/erniealice/espyna-golang/consumer/http"
)

// secureSidebarSwitchWired reports whether the secure principal-switch
// primitive is available for this build/dialect. Delegates the build-tag
// check and nil-guard logic to consumer.SecureSidebarSwitchWired.
func (d *authChainDeps) secureSidebarSwitchWired() bool {
	return d != nil && consumer.SecureSidebarSwitchWired(
		d.sessionMw != nil,
		d.uc != nil && d.uc.Service != nil &&
			d.uc.Service.Auth != nil &&
			d.uc.Service.Auth.SwitchPrincipal != nil,
		d.uc != nil && d.uc.Service != nil &&
			d.uc.Service.Auth != nil &&
			d.uc.Service.Auth.ResolvePrincipals != nil,
	)
}

// secureSidebarSwitchFn builds the SecureSwitchFn closure that the entydad
// switch-workspace handler invokes. It resolves the target binding then
// delegates to executePrincipalSwitch.
func (d *authChainDeps) secureSidebarSwitchFn() workspaceaction.SecureSwitchFn {
	return func(ctx context.Context, userID, sessionToken, targetWorkspaceID string) (*workspaceaction.SecureSwitchResult, error) {
		if userID == "" || targetWorkspaceID == "" {
			return nil, errors.New("sidebar switch: user_id and workspace_id required")
		}
		loader := d.newPrincipalLoader()
		if loader == nil || !loader.IsEnabled() {
			return nil, errors.New("sidebar switch: principal loader unavailable")
		}

		// EXPLICIT-SWITCH: resolve fresh against target workspace with empty
		// hint. Single binding -> auto-pick; multiple/delegate -> chooser;
		// no binding -> 403.
		target, err := loader.ResolveBindingInWorkspace(
			ctx, userID, targetWorkspaceID,
			adapthttp.PrincipalTypeUnspecified, "",
		)
		if err != nil {
			if errors.Is(err, adapthttp.ErrAmbiguousBinding) {
				return &workspaceaction.SecureSwitchResult{
					RedirectURL: "/auth/select-workspace-role",
				}, nil
			}
			return nil, err
		}
		if target == nil {
			return nil, adapthttp.ErrNoBinding
		}

		result, switchErr := d.executePrincipalSwitch(ctx, authPrincipalSwitchInput{
			UserID: userID,
			Token:  sessionToken,
			TargetPrincipal: adapthttp.Principal{
				Type:            target.Type,
				PrincipalID:     target.PrincipalID,
				WorkspaceID:     target.WorkspaceID,
				DisplayName:     target.DisplayName,
				ActingAsTargets: target.ActingAsTargets,
			},
			UseCase:      "switch_explicit_rotate",
			URLDriven:    false,
			RequireAudit: true,
		})
		if switchErr != nil {
			return nil, switchErr
		}
		if result == nil {
			return nil, errors.New("sidebar switch: nil result from principal switch")
		}

		redirect := d.homeURLForWorkspaceID(ctx, target.WorkspaceID)
		return &workspaceaction.SecureSwitchResult{
			NewToken:    result.NewToken,
			RedirectURL: redirect,
		}, nil
	}
}

// secureSidebarResolveUserID delegates to consumer.SecureSidebarResolveUserID.
func secureSidebarResolveUserID(r *http.Request) string {
	return consumer.SecureSidebarResolveUserID(r)
}

// secureSidebarSetSessionCookie writes the rotated session cookie and a fresh
// workspace-claim CSRF cookie. The CSRF cookie is issued via the host-stamped
// 4-arg issuer (ctx.CSRFIssuer) so the cookie bytes stay identical to the
// app-side path.
func (d *authChainDeps) secureSidebarSetSessionCookie(w http.ResponseWriter, token string) {
	if d.sessionMw != nil {
		d.sessionMw.SetSessionCookie(w, token)
	}
	if len(d.csrfSecret) > 0 && d.csrfIssuer != nil {
		d.csrfIssuer(w, d.csrfSecret, token, "")
	}
}
