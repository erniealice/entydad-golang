package list

// Phase 5 (UI permission reflection) — page-controller permission-gate tests.
//
// Verifies buildBulkActions and the row-action permission gating on the
// workspace list. Workspace is the disabled-CTA pattern reference
// (PrimaryAction is rendered directly in the template; row + bulk actions
// use the standard Disabled+DisabledTooltip plumbing).

import (
	"fmt"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	"github.com/erniealice/entydad-golang"

	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
)

func workspaceTestCommonLabels() pyeza.CommonLabels {
	return pyeza.CommonLabels{
		Errors: pyeza.ErrorLabels{MissingPermission: "Missing permission: %s"},
		Bulk:   pyeza.BulkLabels{Delete: "Delete"},
	}
}

func workspaceTestSharedLabels() entydad.SharedLabels {
	return entydad.SharedLabels{
		Badges: entydad.SharedBadgeLabels{NoPermission: "No permission", Yes: "Yes", No: "No"},
		Confirm: entydad.SharedConfirmLabels{
			Activate:       "activate %s",
			Deactivate:     "deactivate %s",
			BulkActivate:   "bulk activate?",
			BulkDeactivate: "bulk deactivate?",
			BulkDelete:     "bulk delete?",
		},
	}
}

func workspaceTestLabels() entydad.WorkspaceLabels {
	return entydad.WorkspaceLabels{
		Actions: entydad.WorkspaceActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
	}
}

func findWorkspaceAction(actions []types.TableAction, typ string) *types.TableAction {
	for i := range actions {
		if actions[i].Type == typ {
			return &actions[i]
		}
	}
	return nil
}

// TestBuildTableRows_WorkspacePermissionMatrix exercises the
// {viewer, editor, admin} matrix against workspace row actions.
func TestBuildTableRows_WorkspacePermissionMatrix(t *testing.T) {
	t.Parallel()

	workspaces := []*workspacepb.Workspace{
		{Id: "ws-1", Name: "Acme Inc", Description: "Primary workspace", Private: false, Active: true},
	}

	sl := workspaceTestSharedLabels()
	l := workspaceTestLabels()
	routes := entydad.DefaultWorkspaceRoutes()

	cases := []struct {
		name             string
		perms            []string
		wantEditDisabled bool
		wantDeactDisabled bool
		wantDelDisabled  bool
	}{
		{
			name:              "viewer",
			perms:             []string{"workspace:list", "workspace:read"},
			wantEditDisabled:  true,
			wantDeactDisabled: true,
			wantDelDisabled:   true,
		},
		{
			name:              "editor (no delete)",
			perms:             []string{"workspace:list", "workspace:read", "workspace:create", "workspace:update"},
			wantEditDisabled:  false,
			wantDeactDisabled: false,
			wantDelDisabled:   true,
		},
		{
			name:              "admin",
			perms:             []string{"workspace:list", "workspace:read", "workspace:create", "workspace:update", "workspace:delete"},
			wantEditDisabled:  false,
			wantDeactDisabled: false,
			wantDelDisabled:   false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			perms := types.NewUserPermissions(tc.perms)
			rows := buildTableRows(workspaces, "active", l, sl, routes, perms)
			if len(rows) != 1 {
				t.Fatalf("rows = %d, want 1", len(rows))
			}
			actions := rows[0].Actions

			if edit := findWorkspaceAction(actions, "edit"); edit == nil {
				t.Fatalf("edit action not found")
			} else if edit.Disabled != tc.wantEditDisabled {
				t.Errorf("edit.Disabled = %v, want %v", edit.Disabled, tc.wantEditDisabled)
			}
			if deact := findWorkspaceAction(actions, "deactivate"); deact == nil {
				t.Fatalf("deactivate action not found")
			} else if deact.Disabled != tc.wantDeactDisabled {
				t.Errorf("deactivate.Disabled = %v, want %v", deact.Disabled, tc.wantDeactDisabled)
			}
			if del := findWorkspaceAction(actions, "delete"); del == nil {
				t.Fatalf("delete action not found")
			} else if del.Disabled != tc.wantDelDisabled {
				t.Errorf("delete.Disabled = %v, want %v", del.Disabled, tc.wantDelDisabled)
			}
		})
	}
}

// TestBuildBulkActions_WorkspacePermissionMatrix verifies bulk gating
// for the disabled-CTA pattern reference entity.
func TestBuildBulkActions_WorkspacePermissionMatrix(t *testing.T) {
	t.Parallel()

	sl := workspaceTestSharedLabels()
	cl := workspaceTestCommonLabels()
	l := workspaceTestLabels()
	routes := entydad.DefaultWorkspaceRoutes()

	cases := []struct {
		name                string
		status              string
		perms               []string
		wantStatusDisabled  bool
		wantDeleteDisabled  bool
		wantStatusTooltip   string
		wantDeleteTooltip   string
	}{
		{
			name:               "viewer on active list — every bulk action disabled",
			status:             "active",
			perms:              []string{"workspace:list", "workspace:read"},
			wantStatusDisabled: true,
			wantDeleteDisabled: true,
			wantStatusTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "workspace:update"),
			wantDeleteTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "workspace:delete"),
		},
		{
			name:               "admin on inactive list — all enabled",
			status:             "inactive",
			perms:              []string{"workspace:list", "workspace:read", "workspace:update", "workspace:delete"},
			wantStatusDisabled: false,
			wantDeleteDisabled: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			perms := types.NewUserPermissions(tc.perms)
			actions := buildBulkActions(l, sl, cl, tc.status, routes, perms)
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
