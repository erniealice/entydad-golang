package list

// Phase 5 (UI permission reflection) — page-controller permission-gate tests.
//
// Verifies buildRowActions and buildBulkActions on the supplier list apply
// the correct Disabled flag across the {viewer, editor, admin} matrix.

import (
	"fmt"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	"github.com/erniealice/entydad-golang"
	entitysupplier "github.com/erniealice/entydad-golang/domain/entity/party/supplier"
)

func supplierTestCommonLabels() pyeza.CommonLabels {
	return pyeza.CommonLabels{
		Errors:  pyeza.ErrorLabels{MissingPermission: "Missing permission: %s"},
		Actions: pyeza.ActionLabels{Clone: "Clone"},
		Bulk:    pyeza.BulkLabels{Delete: "Delete"},
	}
}

func supplierTestSharedLabels() entydad.SharedLabels {
	return entydad.SharedLabels{
		Badges:  entydad.SharedBadgeLabels{NoPermission: "No permission"},
		Errors:  entydad.SharedErrorLabels{CannotDeleteInUse: "Cannot delete: in use"},
		Confirm: entydad.SharedConfirmLabels{Activate: "activate %s", Hold: "hold %s", Block: "block %s", BulkActivate: "bulk activate?", BulkHold: "bulk hold?", BulkBlock: "bulk block?", BulkDelete: "bulk delete?"},
	}
}

func supplierTestLabels() entitysupplier.Labels {
	return entitysupplier.Labels{
		Actions: entitysupplier.ActionLabels{
			View:      "View",
			Edit:      "Edit",
			Delete:    "Delete",
			Activate:  "Activate",
			Block:     "Block",
			SetOnHold: "Set on hold",
		},
	}
}

func findSupplierAction(actions []types.TableAction, typ string) *types.TableAction {
	for i := range actions {
		if actions[i].Type == typ {
			return &actions[i]
		}
	}
	return nil
}

// TestBuildRowActions_SupplierPermissionMatrix exercises the
// {viewer, editor, admin} matrix against the supplier row actions.
func TestBuildRowActions_SupplierPermissionMatrix(t *testing.T) {
	t.Parallel()

	sl := supplierTestSharedLabels()
	cl := supplierTestCommonLabels()
	l := supplierTestLabels()
	routes := entitysupplier.DefaultRoutes()

	cases := []struct {
		name              string
		perms             []string
		wantEditDisabled  bool
		wantCloneDisabled bool
		wantHoldDisabled  bool
		wantDelDisabled   bool
	}{
		{
			name:              "viewer — every mutating action disabled",
			perms:             []string{"supplier:list", "supplier:read"},
			wantEditDisabled:  true,
			wantCloneDisabled: true,
			wantHoldDisabled:  true,
			wantDelDisabled:   true,
		},
		{
			name:              "editor (no delete)",
			perms:             []string{"supplier:list", "supplier:read", "supplier:create", "supplier:update"},
			wantEditDisabled:  false,
			wantCloneDisabled: false,
			wantHoldDisabled:  false,
			wantDelDisabled:   true,
		},
		{
			name:              "admin",
			perms:             []string{"supplier:list", "supplier:read", "supplier:create", "supplier:update", "supplier:delete"},
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
			actions := buildRowActions("supp-1", "Acme Supplies", "active", false, l, sl, cl, routes, perms)

			if edit := findSupplierAction(actions, "edit"); edit == nil {
				t.Fatalf("edit action not found")
			} else if edit.Disabled != tc.wantEditDisabled {
				t.Errorf("edit.Disabled = %v, want %v", edit.Disabled, tc.wantEditDisabled)
			}
			if clone := findSupplierAction(actions, "clone"); clone == nil {
				t.Fatalf("clone action not found")
			} else if clone.Disabled != tc.wantCloneDisabled {
				t.Errorf("clone.Disabled = %v, want %v", clone.Disabled, tc.wantCloneDisabled)
			}
			if hold := findSupplierAction(actions, "hold"); hold == nil {
				t.Fatalf("hold action not found")
			} else if hold.Disabled != tc.wantHoldDisabled {
				t.Errorf("hold.Disabled = %v, want %v", hold.Disabled, tc.wantHoldDisabled)
			}
			if del := findSupplierAction(actions, "delete"); del == nil {
				t.Fatalf("delete action not found")
			} else if del.Disabled != tc.wantDelDisabled {
				t.Errorf("delete.Disabled = %v, want %v", del.Disabled, tc.wantDelDisabled)
			}
		})
	}
}

// TestBuildBulkActions_SupplierPermissionMatrix exercises bulk perms.
func TestBuildBulkActions_SupplierPermissionMatrix(t *testing.T) {
	t.Parallel()

	sl := supplierTestSharedLabels()
	cl := supplierTestCommonLabels()
	l := supplierTestLabels()
	routes := entitysupplier.DefaultRoutes()

	cases := []struct {
		name               string
		perms              []string
		wantStatusDisabled bool
		wantDeleteDisabled bool
		wantStatusTooltip  string
		wantDeleteTooltip  string
	}{
		{
			name:               "viewer",
			perms:              []string{"supplier:list", "supplier:read"},
			wantStatusDisabled: true,
			wantDeleteDisabled: true,
			wantStatusTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "supplier:update"),
			wantDeleteTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "supplier:delete"),
		},
		{
			name:               "editor without delete",
			perms:              []string{"supplier:list", "supplier:read", "supplier:update"},
			wantStatusDisabled: false,
			wantDeleteDisabled: true,
		},
		{
			name:               "admin",
			perms:              []string{"supplier:list", "supplier:read", "supplier:update", "supplier:delete"},
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

			var statusAct, deleteAct *types.BulkAction
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
			if deleteAct.Disabled != tc.wantDeleteDisabled {
				t.Errorf("bulk delete.Disabled = %v, want %v", deleteAct.Disabled, tc.wantDeleteDisabled)
			}
			if tc.wantStatusTooltip != "" && statusAct.DisabledTooltip != tc.wantStatusTooltip {
				t.Errorf("bulk status.DisabledTooltip = %q, want %q", statusAct.DisabledTooltip, tc.wantStatusTooltip)
			}
			if tc.wantDeleteTooltip != "" && deleteAct.DisabledTooltip != tc.wantDeleteTooltip {
				t.Errorf("bulk delete.DisabledTooltip = %q, want %q", deleteAct.DisabledTooltip, tc.wantDeleteTooltip)
			}
		})
	}
}
