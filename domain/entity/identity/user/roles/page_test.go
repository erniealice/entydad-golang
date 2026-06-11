package roles

import (
	"context"
	"errors"
	"testing"

	"github.com/erniealice/entydad-golang"
	user "github.com/erniealice/entydad-golang/domain/entity/identity/user"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

func TestBuildTableConfig_FallbackScenariosKeepAssignActionEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		listData     []*workspaceuserpb.WorkspaceUser
		itemErr      error
		wantRows     int
		wantUserName string
	}{
		{
			name: "workspace user found returns role rows",
			listData: []*workspaceuserpb.WorkspaceUser{
				{Id: "wu-1", UserId: "user-1"},
			},
			wantRows:     1,
			wantUserName: "Ada Lovelace",
		},
		{
			name:         "missing workspace user falls back to empty table",
			listData:     []*workspaceuserpb.WorkspaceUser{},
			wantRows:     0,
			wantUserName: "Ada Lovelace",
		},
		{
			name: "workspace user item error falls back to empty table",
			listData: []*workspaceuserpb.WorkspaceUser{
				{Id: "wu-1", UserId: "user-1"},
			},
			itemErr:      errors.New("lookup failed"),
			wantRows:     0,
			wantUserName: "Ada Lovelace",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := &Deps{
				Routes: user.DefaultRoutes(),
				ReadUser: func(context.Context, *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error) {
					return &userpb.ReadUserResponse{
						Data: []*userpb.User{
							{Id: "user-1", FirstName: "Ada", LastName: "Lovelace"},
						},
					}, nil
				},
				ListWorkspaceUsers: func(context.Context, *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error) {
					return &workspaceuserpb.ListWorkspaceUsersResponse{
						Data: tc.listData,
					}, nil
				},
				GetWorkspaceUserItemPageData: func(context.Context, *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error) {
					if tc.itemErr != nil {
						return nil, tc.itemErr
					}
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
				Labels:       user.RoleLabels{},
				SharedLabels: entydad.SharedLabels{Confirm: entydad.SharedConfirmLabels{Remove: "remove %s"}},
				CommonLabels: pyeza.CommonLabels{},
				TableLabels:  types.TableLabels{},
			}

			// Context permissions are intentionally restrictive to document current behavior:
			// this view builder does not consume permissions and leaves assign/remove enabled.
			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(nil))
			table, userName, err := buildTableConfig(ctx, deps, "user-1")
			if err != nil {
				t.Fatalf("buildTableConfig() error = %v", err)
			}

			if got, want := userName, tc.wantUserName; got != want {
				t.Fatalf("userName = %q, want %q", got, want)
			}
			if got, want := len(table.Rows), tc.wantRows; got != want {
				t.Fatalf("row count = %d, want %d", got, want)
			}
			if table.PrimaryAction == nil {
				t.Fatalf("PrimaryAction is nil")
			}
			if table.PrimaryAction.Disabled {
				t.Fatalf("PrimaryAction.Disabled = true, want false")
			}
			if tc.wantRows > 0 {
				actions := table.Rows[0].Actions
				if got, want := len(actions), 1; got != want {
					t.Fatalf("row actions count = %d, want %d", got, want)
				}
				if got, want := actions[0].Type, "delete"; got != want {
					t.Fatalf("row action type = %q, want %q", got, want)
				}
				if actions[0].Disabled {
					t.Fatalf("row remove action Disabled = true, want false")
				}
			}
		})
	}
}

func strPtr(v string) *string {
	return &v
}
