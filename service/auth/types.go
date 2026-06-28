// Package auth provides the auth module for the entydad service layer.
// It owns the auth-flow orchestration (login, signup, reset-password,
// change-password, logout, multi-principal chooser) behind narrow
// interfaces so any consumer app can implement them.
//
// The types in this file mirror the principal/acting-as domain concepts
// that currently live in the app's HTTP layer. They are auth-domain
// concepts that belong here, next to the handlers and views that use them.
package auth

import (
	"strings"

	pyezarender "github.com/erniealice/pyeza-golang/render"
)

// PrincipalType is the presentation-layer principal kind.
// Re-exported from pyeza/render (the authoritative definition).
type PrincipalType = pyezarender.PrincipalType

const (
	PrincipalTypeUnspecified      = pyezarender.PrincipalTypeUnspecified
	PrincipalTypeOperatorOwner    = pyezarender.PrincipalTypeOperatorOwner
	PrincipalTypeOperatorStaff    = pyezarender.PrincipalTypeOperatorStaff
	PrincipalTypeClient           = pyezarender.PrincipalTypeClient
	PrincipalTypeClientDelegate   = pyezarender.PrincipalTypeClientDelegate
	PrincipalTypeSupplier         = pyezarender.PrincipalTypeSupplier
	PrincipalTypeSupplierDelegate = pyezarender.PrincipalTypeSupplierDelegate
	PrincipalTypeStaff            = pyezarender.PrincipalTypeStaff
)

// ActingAsTarget is one target a delegate may act for. A single
// CLIENT_DELEGATE principal may have N>1 acting-as options (one per
// DelegateClient row); the chooser picks one before entering the portal.
type ActingAsTarget struct {
	// ID is the underlying party id (client_id or supplier_id).
	ID string
	// WorkspaceID is the workspace the party lives in (resolved from the
	// client.workspace_id or supplier.workspace_id back-edge).
	WorkspaceID string
	// DisplayName is a human label (e.g. "Jane Smith").
	DisplayName string
}

// Principal is one resolved binding a User holds. The auth flow uses the
// list of Principals to decide post-login routing: 0 = no access, 1 = auto
// route, 2+ = present chooser.
//
// PrincipalID identifies the underlying grant row (WorkspaceUser.id /
// ClientPortalGrant.id / SupplierPortalGrant.id / Delegate.id) — this is
// the value stored in session.principal_id when the session row is
// established or rotated.
//
// For delegate principals (CLIENT_DELEGATE / SUPPLIER_DELEGATE), PrincipalID
// is the Delegate.id. ActingAsTargets holds the available
// (client_id, workspace_id) or (supplier_id, workspace_id) pairs the
// delegate may act for. Empty ActingAsTargets on a delegate row means
// the Delegate exists but has no active DelegateClient/DelegateSupplier
// rows — the principal is effectively non-actionable and must be filtered
// out by the loader before returning the result.
type Principal struct {
	Type            PrincipalType
	PrincipalID     string
	WorkspaceID     string
	DisplayName     string
	ActingAsTargets []ActingAsTarget // only populated for delegate principals
}

// DelegateActingAsResolved reports whether a principal about to be handed to
// executePrincipalSwitch carries an UNAMBIGUOUS acting-as identity.
//
// codex RBC#1 — High-1 (2026-06-02). The sidebar/URL path closes the
// multi-target-delegate hole inside pickBindingForSession (ResolveBindingInWorkspace),
// but the SAME delegate bug class lives on the login auto-route and the
// chooser POST: both forward a resolved delegate Principal to
// executePrincipalSwitch with NO explicit acting-as target. For a delegate
// holding N>1 acting-as targets that persists a delegate principal with an
// EMPTY acting_as_*, and the permission query then fails closed to zero
// permissions — an unresolved-principal state after switch.
//
// This is the SINGLE narrow chokepoint both callers (login auto-route in
// domain_auth.go and the /action/auth/switch-principal POST) MUST gate on
// before invoking executePrincipalSwitch, so the guard can't be bypassed.
//
// The rule (fail-closed):
//
//   - Non-delegate principal              → always resolved (true). The
//     acting-as concept does not apply; no acting-as id is needed.
//   - Delegate with 0 acting-as targets   → resolved (true). Such a delegate
//     is a dead-end the loader already drops; if one ever reaches here it is
//     unambiguous (nothing to choose) and the existing fail-closed permission
//     scoping handles it. We do NOT introduce a new reject for the zero case.
//   - Delegate with exactly 1 target      → resolved (true). The lone target
//     is unambiguous; SwitchPrincipal auto-selects it.
//   - Delegate with N>1 targets:
//   - explicit acting-as id supplied
//     that matches one of the targets  → resolved (true).
//   - no explicit acting-as id, or an
//     id that is NOT one of the
//     delegate's own targets           → NOT resolved (false) → caller
//     MUST fail closed (re-render the chooser / surface an explicit
//     "select a target" outcome), never persist an empty-acting-as
//     principal.
//
// Single-target delegates and non-delegate principals are unaffected.
func DelegateActingAsResolved(p Principal, actingAsClientID, actingAsSupplierID string) bool {
	switch p.Type {
	case PrincipalTypeClientDelegate, PrincipalTypeSupplierDelegate:
		// fall through to the multi-target check below.
	default:
		return true // non-delegate: acting-as not applicable.
	}

	if len(p.ActingAsTargets) <= 1 {
		// 0 or 1 target → unambiguous. (The loader drops zero-target
		// delegates; a lone target is auto-selected downstream.)
		return true
	}

	// N>1 targets. An explicit acting-as id resolves the ambiguity ONLY when
	// it names one of THIS delegate's own targets (never trust a free-form
	// id — it must be one the delegate actually holds).
	want := strings.TrimSpace(actingAsClientID)
	if p.Type == PrincipalTypeSupplierDelegate {
		want = strings.TrimSpace(actingAsSupplierID)
	}
	if want == "" {
		return false
	}
	for _, t := range p.ActingAsTargets {
		if t.ID == want {
			return true
		}
	}
	return false
}

// PrincipalTypeString returns the canonical lowercase token for a principal type. The
// token is used in route URLs (`/portal/{kind}/`) and in `data-testid`
// attributes (`select-workspace-role-{kind}`). Keep these stable.
//
// Defined as a standalone function (not a method) because PrincipalType is a
// type alias to pyezarender.PrincipalType, and Go does not allow methods on
// non-local type aliases. Once the go.mod bumps to a Go version that supports
// methods on aliases, this can become a method.
func PrincipalTypeString(t PrincipalType) string {
	switch t {
	case PrincipalTypeOperatorOwner:
		return "operator_owner"
	case PrincipalTypeOperatorStaff:
		return "operator_staff"
	case PrincipalTypeClient:
		return "client"
	case PrincipalTypeClientDelegate:
		return "client_delegate"
	case PrincipalTypeSupplier:
		return "supplier"
	case PrincipalTypeSupplierDelegate:
		return "supplier_delegate"
	case PrincipalTypeStaff:
		return "staff"
	default:
		return "unspecified"
	}
}

// PrincipalSwitchInput mirrors the composition-internal principalSwitchInput.
// Exported so the auth module can construct switch requests.
type PrincipalSwitchInput struct {
	UserID             string
	Token              string
	TargetPrincipal    Principal
	ActingAsClientID   string
	ActingAsSupplierID string
	UseCase            string
	RequestURL         string
	Referer            string
	SecFetchSite       string
	UserAgent          string
	RequireAudit       bool
}

// PrincipalSwitchResult mirrors the composition-internal principalSwitchResult.
type PrincipalSwitchResult struct {
	NewToken    string
	RedirectURL string
}
