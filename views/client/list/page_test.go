package list

// Phase 5 (UI permission reflection) — page-controller permission-gate tests.
//
// Verifies buildRowActions and buildBulkActions on the client list apply
// the correct Disabled flag and AWS-style "Missing permission: <code>"
// tooltips across the {viewer, editor, admin} permission matrix.
//
// Client has 5 row actions (View, Edit, Clone, status transitions, Delete)
// + bulk Activate/Deactivate/Block/Hold/Prospect/Delete — the highest
// combinatorial value in entydad after workspace/user/role.

import (
	"fmt"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	"github.com/erniealice/entydad-golang"
)

func clientTestCommonLabels() pyeza.CommonLabels {
	return pyeza.CommonLabels{
		Errors: pyeza.ErrorLabels{
			MissingPermission: "Missing permission: %s",
		},
		Actions: pyeza.ActionLabels{Clone: "Clone"},
		Bulk:    pyeza.BulkLabels{Delete: "Delete"},
	}
}

func clientTestSharedLabels() entydad.SharedLabels {
	return entydad.SharedLabels{
		Badges:  entydad.SharedBadgeLabels{NoPermission: "No permission"},
		Errors:  entydad.SharedErrorLabels{CannotDeleteInUse: "Cannot delete: in use"},
		Confirm: entydad.SharedConfirmLabels{Activate: "activate %s", Deactivate: "deactivate %s", Prospect: "prospect %s", Hold: "hold %s", Block: "block %s", BulkDelete: "delete?"},
	}
}

func clientTestLabels() entydad.ClientLabels {
	return entydad.ClientLabels{
		Detail: entydad.ClientDetailLabels{
			Actions: entydad.ClientDetailActionLabels{
				ViewClient:       "View",
				EditClient:       "Edit",
				DeleteClient:     "Delete",
				ActivateClient:   "Activate",
				DeactivateClient: "Deactivate",
				HoldClient:       "Hold",
				BlockClient:      "Block",
				SetProspect:      "Set prospect",
			},
		},
		BulkActions: entydad.ClientBulkActionLabels{
			SetAsActive:   "Bulk activate",
			SetAsProspect: "Bulk prospect",
			SetAsOnHold:   "Bulk hold",
			SetAsBlocked:  "Bulk block",
			SetAsInactive: "Bulk deactivate",
		},
	}
}

func findClientAction(actions []types.TableAction, typ string) *types.TableAction {
	for i := range actions {
		if actions[i].Type == typ {
			return &actions[i]
		}
	}
	return nil
}

// TestBuildRowActions_ClientPermissionMatrix exercises the
// {viewer, editor, admin} matrix against the client row actions.
func TestBuildRowActions_ClientPermissionMatrix(t *testing.T) {
	t.Parallel()

	sl := clientTestSharedLabels()
	cl := clientTestCommonLabels()
	l := clientTestLabels()
	routes := entydad.DefaultClientRoutes()

	cases := []struct {
		name             string
		perms            []string
		wantEditDisabled bool
		wantCloneDisabled bool
		wantHoldDisabled  bool // representative transition (on the active list)
		wantDelDisabled   bool
	}{
		{
			name:              "viewer — every mutating action disabled",
			perms:             []string{"client:list", "client:read"},
			wantEditDisabled:  true,
			wantCloneDisabled: true,
			wantHoldDisabled:  true,
			wantDelDisabled:   true,
		},
		{
			name:              "editor (no delete) — edit/clone/transitions enabled, delete disabled",
			perms:             []string{"client:list", "client:read", "client:create", "client:update"},
			wantEditDisabled:  false,
			wantCloneDisabled: false,
			wantHoldDisabled:  false,
			wantDelDisabled:   true,
		},
		{
			name:              "admin — every action enabled",
			perms:             []string{"client:list", "client:read", "client:create", "client:update", "client:delete"},
			wantEditDisabled:  false,
			wantCloneDisabled: false,
			wantHoldDisabled:  false,
			wantDelDisabled:   false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			perms := types.NewUserPermissions(tc.perms)
			actions := buildRowActions("client-1", "Acme Corp", "active", false /*isInUse*/, l, sl, cl, routes, perms)

			if edit := findClientAction(actions, "edit"); edit == nil {
				t.Fatalf("edit action not found")
			} else if edit.Disabled != tc.wantEditDisabled {
				t.Errorf("edit.Disabled = %v, want %v", edit.Disabled, tc.wantEditDisabled)
			}

			if clone := findClientAction(actions, "clone"); clone == nil {
				t.Fatalf("clone action not found")
			} else if clone.Disabled != tc.wantCloneDisabled {
				t.Errorf("clone.Disabled = %v, want %v", clone.Disabled, tc.wantCloneDisabled)
			}

			// "hold" is a status transition — disabled iff !perms.Can("client","update")
			if hold := findClientAction(actions, "hold"); hold == nil {
				t.Fatalf("hold transition action not found")
			} else if hold.Disabled != tc.wantHoldDisabled {
				t.Errorf("hold.Disabled = %v, want %v", hold.Disabled, tc.wantHoldDisabled)
			}

			if del := findClientAction(actions, "delete"); del == nil {
				t.Fatalf("delete action not found")
			} else if del.Disabled != tc.wantDelDisabled {
				t.Errorf("delete.Disabled = %v, want %v", del.Disabled, tc.wantDelDisabled)
			}
		})
	}
}

// TestBuildRowActions_Client_InUseBlocksDeleteWithStateTooltip verifies the
// in-use status tooltip wins over the permission tooltip when the user
// has delete permission but the client cannot be deleted.
func TestBuildRowActions_Client_InUseBlocksDeleteWithStateTooltip(t *testing.T) {
	t.Parallel()

	sl := clientTestSharedLabels()
	cl := clientTestCommonLabels()
	l := clientTestLabels()
	routes := entydad.DefaultClientRoutes()

	// Admin perms but client is in use.
	perms := types.NewUserPermissions([]string{"client:list", "client:read", "client:create", "client:update", "client:delete"})
	actions := buildRowActions("client-2", "InUseCorp", "active", true /*isInUse*/, l, sl, cl, routes, perms)

	del := findClientAction(actions, "delete")
	if del == nil {
		t.Fatalf("delete action not found")
	}
	if !del.Disabled {
		t.Error("delete should be disabled when client is in-use")
	}
	if del.DisabledTooltip != sl.Errors.CannotDeleteInUse {
		t.Errorf("delete.DisabledTooltip = %q, want CannotDeleteInUse %q", del.DisabledTooltip, sl.Errors.CannotDeleteInUse)
	}
}

// TestBuildBulkActions_ClientPermissionMatrix exercises the bulk-action
// matrix. Bulk activate/deactivate/etc. gate on client:update; bulk delete
// gates on client:delete. The AWS-style tooltip is interpolated.
func TestBuildBulkActions_ClientPermissionMatrix(t *testing.T) {
	t.Parallel()

	sl := clientTestSharedLabels()
	cl := clientTestCommonLabels()
	l := clientTestLabels()
	routes := entydad.DefaultClientRoutes()

	cases := []struct {
		name                string
		perms               []string
		wantStatusDisabled  bool
		wantDeleteDisabled  bool
		wantStatusTooltip   string
		wantDeleteTooltip   string
	}{
		{
			name:               "viewer — every bulk action disabled with interpolated tooltip",
			perms:              []string{"client:list", "client:read"},
			wantStatusDisabled: true,
			wantDeleteDisabled: true,
			wantStatusTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "client:update"),
			wantDeleteTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "client:delete"),
		},
		{
			name:               "editor without delete — status enabled, delete disabled",
			perms:              []string{"client:list", "client:read", "client:create", "client:update"},
			wantStatusDisabled: false,
			wantDeleteDisabled: true,
			wantDeleteTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "client:delete"),
		},
		{
			name:               "admin — all bulk actions enabled",
			perms:              []string{"client:list", "client:read", "client:create", "client:update", "client:delete"},
			wantStatusDisabled: false,
			wantDeleteDisabled: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			perms := types.NewUserPermissions(tc.perms)
			actions := buildBulkActions(l, sl, cl, "active", routes, perms)
			if len(actions) == 0 {
				t.Fatal("no bulk actions produced")
			}

			// Find the first status-transition action (non-delete key).
			var statusAct *types.BulkAction
			var deleteAct *types.BulkAction
			for i := range actions {
				a := &actions[i]
				if a.Key == "delete" {
					deleteAct = a
				} else if statusAct == nil {
					statusAct = a
				}
			}
			if statusAct == nil {
				t.Fatal("no status-transition bulk action found")
			}
			if deleteAct == nil {
				t.Fatal("no bulk delete action found")
			}

			if statusAct.Disabled != tc.wantStatusDisabled {
				t.Errorf("bulk status.Disabled = %v, want %v", statusAct.Disabled, tc.wantStatusDisabled)
			}
			if tc.wantStatusTooltip != "" && statusAct.DisabledTooltip != tc.wantStatusTooltip {
				t.Errorf("bulk status.DisabledTooltip = %q, want %q", statusAct.DisabledTooltip, tc.wantStatusTooltip)
			}

			if deleteAct.Disabled != tc.wantDeleteDisabled {
				t.Errorf("bulk delete.Disabled = %v, want %v", deleteAct.Disabled, tc.wantDeleteDisabled)
			}
			if tc.wantDeleteTooltip != "" && deleteAct.DisabledTooltip != tc.wantDeleteTooltip {
				t.Errorf("bulk delete.DisabledTooltip = %q, want %q", deleteAct.DisabledTooltip, tc.wantDeleteTooltip)
			}
		})
	}
}
