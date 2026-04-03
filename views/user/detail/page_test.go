package detail

import (
	"context"
	"testing"

	"github.com/erniealice/entydad-golang"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

func TestBuildRolesTable_DisablesAssignAndRemoveByUpdatePermission(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		permCodes    []string
		wantDisabled bool
	}{
		{
			name:         "missing user:update disables assign and remove",
			permCodes:    nil,
			wantDisabled: true,
		},
		{
			name:         "user:update enables assign and remove",
			permCodes:    []string{"user:update"},
			wantDisabled: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := &DetailViewDeps{
				Routes: entydad.DefaultUserRoutes(),
				ListWorkspaceUsers: func(context.Context, *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error) {
					return &workspaceuserpb.ListWorkspaceUsersResponse{
						Data: []*workspaceuserpb.WorkspaceUser{
							{Id: "wu-1", UserId: "user-1"},
						},
					}, nil
				},
				GetWorkspaceUserItemPageData: func(context.Context, *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error) {
					return &workspaceuserpb.GetWorkspaceUserItemPageDataResponse{
						WorkspaceUser: &workspaceuserpb.WorkspaceUser{
							Id: "wu-1",
							WorkspaceUserRoles: []*workspaceuserrolepb.WorkspaceUserRole{
								{
									Id:                "wur-1",
									DateCreatedString: strPtr("2026-01-01"),
									Role: &rolepb.Role{
										Id:          "role-1",
										Name:        "Admin",
										Description: "Administrator",
										Color:       "blue",
									},
								},
							},
						},
					}, nil
				},
				UserRoleLabels: entydad.UserRoleLabels{},
				SharedLabels: entydad.SharedLabels{
					Confirm: entydad.SharedConfirmLabels{Remove: "remove %s"},
					Badges:  entydad.SharedBadgeLabels{NoPermission: "No permission"},
				},
				CommonLabels: pyeza.CommonLabels{},
				TableLabels:  types.TableLabels{},
			}

			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(tc.permCodes))
			table, err := buildRolesTable(ctx, deps, "user-1", view.GetUserPermissions(ctx))
			if err != nil {
				t.Fatalf("buildRolesTable() error = %v", err)
			}

			if table.PrimaryAction == nil {
				t.Fatalf("PrimaryAction is nil")
			}
			if got, want := table.PrimaryAction.Disabled, tc.wantDisabled; got != want {
				t.Fatalf("PrimaryAction.Disabled = %v, want %v", got, want)
			}
			if got, want := len(table.Rows), 1; got != want {
				t.Fatalf("row count = %d, want %d", got, want)
			}
			actions := table.Rows[0].Actions
			if got, want := len(actions), 1; got != want {
				t.Fatalf("row actions count = %d, want %d", got, want)
			}
			if got, want := actions[0].Type, "delete"; got != want {
				t.Fatalf("row action type = %q, want %q", got, want)
			}
			if got, want := actions[0].Disabled, tc.wantDisabled; got != want {
				t.Fatalf("remove action disabled = %v, want %v", got, want)
			}
		})
	}
}

func TestBuildRolesTable_CurrentBehavior_EmptyFallbackDoesNotDisableAssign(t *testing.T) {
	t.Parallel()

	deps := &DetailViewDeps{
		Routes: entydad.DefaultUserRoutes(),
		ListWorkspaceUsers: func(context.Context, *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error) {
			return &workspaceuserpb.ListWorkspaceUsersResponse{
				Data: []*workspaceuserpb.WorkspaceUser{},
			}, nil
		},
		GetWorkspaceUserItemPageData: func(context.Context, *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error) {
			t.Fatalf("GetWorkspaceUserItemPageData should not be called when workspace user is missing")
			return nil, nil
		},
		UserRoleLabels: entydad.UserRoleLabels{},
		SharedLabels: entydad.SharedLabels{
			Badges: entydad.SharedBadgeLabels{NoPermission: "No permission"},
		},
		CommonLabels: pyeza.CommonLabels{},
		TableLabels:  types.TableLabels{},
	}

	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(nil))
	table, err := buildRolesTable(ctx, deps, "missing-user", view.GetUserPermissions(ctx))
	if err != nil {
		t.Fatalf("buildRolesTable() error = %v", err)
	}

	if got, want := len(table.Rows), 0; got != want {
		t.Fatalf("row count = %d, want %d", got, want)
	}
	if table.PrimaryAction == nil {
		t.Fatalf("PrimaryAction is nil")
	}
	if table.PrimaryAction.Disabled {
		t.Fatalf("PrimaryAction.Disabled = true, want false for current empty-table fallback behavior")
	}
}

func strPtr(v string) *string {
	return &v
}
