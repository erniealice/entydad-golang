package auth

import (
	"testing"
)

// TestDelegateActingAsResolved_LoginAutoRoute_And_ChooserPOST_FailClosed locks in the
// codex RBC#1 (High-1) fix at the composition call sites. The login auto-route
// (domain_auth.go case-1 and the all-same-kind branch) and the chooser POST
// handler (POST /action/auth/switch-principal) both gate on
// DelegateActingAsResolved(target, actingAsClientID, actingAsSupplierID)
// BEFORE invoking executePrincipalSwitch. When that guard returns false the
// handler MUST fail closed (redirect to /auth/select-workspace-role) and MUST
// NOT call executePrincipalSwitch — which is the only thing that would persist
// the principal onto the session row. A persisted multi-target delegate with an
// empty acting_as_* fails closed to zero permissions downstream, so the only
// safe outcome is to route to the chooser.
//
// This test exercises the exact guard call the auth handlers make for the
// scenarios that reach those call sites, proving that:
//   - the login auto-route forwards NO acting-as (it calls the guard with two
//     empty strings); a multi-target delegate there is blocked → routed.
//   - the chooser POST forwards whatever acting-as the form posts (none today);
//     a multi-target delegate chosen there is blocked → routed.
//
// A full DB-backed handler test is intentionally avoided: the route-registration
// closure needs a live auth adapter, session middleware, and principal loader.
// The guard is the single narrow chokepoint both callers share (see
// domain_auth.go), so verifying its decision for the call-site inputs is the
// load-bearing proof that neither path can persist an empty acting-as.
//
// Migrated from apps/service-admin/internal/composition/auth_delegate_guard_test.go
// to entydad/service/auth/ because the guard function (DelegateActingAsResolved)
// and all supporting types (Principal, PrincipalType*, ActingAsTarget) live here.
func TestDelegateActingAsResolved_LoginAutoRoute_And_ChooserPOST_FailClosed(t *testing.T) {
	multiTargetClientDelegate := Principal{
		Type:        PrincipalTypeClientDelegate,
		PrincipalID: "del-1",
		WorkspaceID: "ws-X",
		ActingAsTargets: []ActingAsTarget{
			{ID: "client-A", WorkspaceID: "ws-X", DisplayName: "Client A"},
			{ID: "client-B", WorkspaceID: "ws-X", DisplayName: "Client B"},
		},
	}
	multiTargetSupplierDelegate := Principal{
		Type:        PrincipalTypeSupplierDelegate,
		PrincipalID: "del-2",
		WorkspaceID: "ws-X",
		ActingAsTargets: []ActingAsTarget{
			{ID: "supplier-A", WorkspaceID: "ws-X", DisplayName: "Supplier A"},
			{ID: "supplier-B", WorkspaceID: "ws-X", DisplayName: "Supplier B"},
		},
	}

	// --- Login auto-route call site: ALWAYS empty acting-as. ---
	// domain_auth.go case-1 and the all-same-kind branch both call:
	//   DelegateActingAsResolved(principals[0], "", "")
	// A multi-target delegate there MUST be blocked (→ chooser redirect), never
	// handed to executePrincipalSwitch.
	if DelegateActingAsResolved(multiTargetClientDelegate, "", "") {
		t.Fatal("login auto-route: multi-target CLIENT delegate with empty acting-as must be BLOCKED (would persist empty acting_as_*)")
	}
	if DelegateActingAsResolved(multiTargetSupplierDelegate, "", "") {
		t.Fatal("login auto-route: multi-target SUPPLIER delegate with empty acting-as must be BLOCKED")
	}

	// --- Chooser POST call site: forwards the form's acting-as fields. ---
	// The select-workspace-role.html form posts only principal_id +
	// principal_kind today, so actingAsClientID / actingAsSupplierID are empty.
	// POST /action/auth/switch-principal calls:
	//   DelegateActingAsResolved(target, actingAsClientID, actingAsSupplierID)
	// With empty acting-as a multi-target delegate MUST be blocked.
	if DelegateActingAsResolved(multiTargetClientDelegate, "" /*acting_as_client_id*/, "" /*acting_as_supplier_id*/) {
		t.Fatal("chooser POST: multi-target CLIENT delegate with empty acting-as must be BLOCKED")
	}
	// A foreign acting-as id (not one of the delegate's own targets) must not
	// unlock the switch either.
	if DelegateActingAsResolved(multiTargetClientDelegate, "client-NOT-MINE", "") {
		t.Fatal("chooser POST: multi-target CLIENT delegate with a non-owned acting-as id must be BLOCKED")
	}

	// --- Non-blocking cases must still proceed (no regression). ---
	singleTargetDelegate := Principal{
		Type:        PrincipalTypeClientDelegate,
		PrincipalID: "del-3",
		WorkspaceID: "ws-X",
		ActingAsTargets: []ActingAsTarget{
			{ID: "client-A", WorkspaceID: "ws-X", DisplayName: "Client A"},
		},
	}
	if !DelegateActingAsResolved(singleTargetDelegate, "", "") {
		t.Fatal("single-target delegate must resolve (auto-selected downstream) — no false block")
	}
	operatorStaff := Principal{
		Type:        PrincipalTypeOperatorStaff,
		PrincipalID: "wu-1",
		WorkspaceID: "ws-X",
	}
	if !DelegateActingAsResolved(operatorStaff, "", "") {
		t.Fatal("non-delegate principal (OPERATOR_STAFF) must always resolve — acting-as not applicable")
	}
	// An explicit, OWNED acting-as id resolves the multi-target ambiguity
	// (forward-compat for the FLAGGED per-target chooser UX).
	if !DelegateActingAsResolved(multiTargetClientDelegate, "client-B", "") {
		t.Fatal("multi-target delegate WITH a matching acting-as id must resolve")
	}
}
