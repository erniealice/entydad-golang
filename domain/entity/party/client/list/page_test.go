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
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	entityclient "github.com/erniealice/entydad-golang/domain/entity/party/client"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
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

func clientTestLabels() entityclient.Labels {
	return entityclient.Labels{
		Detail: entityclient.DetailLabels{
			Actions: entityclient.DetailActionLabels{
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
		BulkActions: entityclient.BulkActionLabels{
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

func TestBuildTableConfig_ActiveSubscriptionCountsUseOnlyReturnedPageIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		clients   []*clientpb.Client
		wantIDs   []string
		wantCalls int
	}{
		{
			name: "page IDs in returned order, skipping nil and empty IDs",
			clients: []*clientpb.Client{
				{Id: "client-2", Active: true}, nil, {Id: "", Active: true}, {Id: "client-1", Active: true},
			},
			wantIDs:   []string{"client-2", "client-1"},
			wantCalls: 1,
		},
		{
			name:      "empty page does not invoke count callback",
			clients:   nil,
			wantCalls: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotIDs []string
			calls := 0
			deps := &ListViewDeps{
				Routes:       entityclient.DefaultRoutes(),
				Labels:       clientTestLabels(),
				SharedLabels: clientTestSharedLabels(),
				CommonLabels: clientTestCommonLabels(),
				GetListPageData: func(context.Context, *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error) {
					return &clientpb.GetClientListPageDataResponse{
						ClientList: tc.clients,
						Pagination: &commonpb.PaginationResponse{TotalItems: int32(len(tc.clients))},
					}, nil
				},
				GetActiveSubscriptionCounts: func(_ context.Context, ids []string) (map[string]int32, error) {
					calls++
					gotIDs = append([]string(nil), ids...)
					return map[string]int32{}, nil
				},
			}

			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"client:list"}))
			_, err := buildTableConfig(ctx, deps, clientColumns(deps.Labels), "active", tableparams.TableQueryParams{Page: 1, PageSize: 20})
			if err != nil {
				t.Fatalf("buildTableConfig() error = %v", err)
			}
			if calls != tc.wantCalls {
				t.Fatalf("count callback calls = %d, want %d", calls, tc.wantCalls)
			}
			if !reflect.DeepEqual(gotIDs, tc.wantIDs) {
				t.Errorf("count callback IDs = %v, want %v", gotIDs, tc.wantIDs)
			}
		})
	}
}

func TestBuildTableConfig_ActiveSubscriptionCountErrorIsVisible(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("count query failed")
	deps := &ListViewDeps{
		Routes:       entityclient.DefaultRoutes(),
		Labels:       clientTestLabels(),
		SharedLabels: clientTestSharedLabels(),
		CommonLabels: clientTestCommonLabels(),
		GetListPageData: func(context.Context, *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error) {
			return &clientpb.GetClientListPageDataResponse{
				ClientList: []*clientpb.Client{{Id: "client-1", Active: true}},
				Pagination: &commonpb.PaginationResponse{TotalItems: 1},
			}, nil
		},
		GetActiveSubscriptionCounts: func(context.Context, []string) (map[string]int32, error) {
			return nil, wantErr
		},
	}

	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"client:list"}))
	_, err := buildTableConfig(ctx, deps, clientColumns(deps.Labels), "active", tableparams.TableQueryParams{Page: 1, PageSize: 20})
	if !errors.Is(err, wantErr) {
		t.Fatalf("buildTableConfig() error = %v, want wrapped %v", err, wantErr)
	}
}

// TestBuildRowActions_ClientPermissionMatrix exercises the
// {viewer, editor, admin} matrix against the client row actions.
func TestBuildRowActions_ClientPermissionMatrix(t *testing.T) {
	t.Parallel()

	sl := clientTestSharedLabels()
	cl := clientTestCommonLabels()
	l := clientTestLabels()
	routes := entityclient.DefaultRoutes()

	cases := []struct {
		name              string
		perms             []string
		wantEditDisabled  bool
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
	routes := entityclient.DefaultRoutes()

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
	routes := entityclient.DefaultRoutes()

	cases := []struct {
		name               string
		perms              []string
		wantStatusDisabled bool
		wantDeleteDisabled bool
		wantStatusTooltip  string
		wantDeleteTooltip  string
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
