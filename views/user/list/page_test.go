package list

import (
	"context"
	"errors"
	"testing"

	"github.com/erniealice/entydad-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

func TestBuildTableConfig_UserPermissionGatesPrimaryAndRowActions(t *testing.T) {
	t.Parallel()

	users := []*userpb.User{
		{
			Id:                "user-active",
			FirstName:         "Ada",
			LastName:          "Lovelace",
			EmailAddress:      "ada@example.com",
			Active:            true,
			DateCreatedString: strPtr("2026-01-01"),
		},
		{
			Id:                "user-inactive",
			FirstName:         "Grace",
			LastName:          "Hopper",
			EmailAddress:      "grace@example.com",
			Active:            false,
			DateCreatedString: strPtr("2026-01-02"),
		},
	}

	tests := []struct {
		name                string
		status              string
		permCodes           []string
		wantPrimaryDisabled bool
		wantStatusAction    string
		wantRowDisabled     bool
	}{
		{
			name:                "active tab with no user permissions disables create and update actions",
			status:              "active",
			permCodes:           nil,
			wantPrimaryDisabled: true,
			wantStatusAction:    "deactivate",
			wantRowDisabled:     true,
		},
		{
			name:                "active tab with create and update permissions enables actions",
			status:              "active",
			permCodes:           []string{"user:create", "user:update"},
			wantPrimaryDisabled: false,
			wantStatusAction:    "deactivate",
			wantRowDisabled:     false,
		},
		{
			name:                "inactive tab uses activate action and still enforces update permission",
			status:              "inactive",
			permCodes:           nil,
			wantPrimaryDisabled: true,
			wantStatusAction:    "activate",
			wantRowDisabled:     true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var capturedReq *userpb.GetUserListPageDataRequest
			deps := newListViewDeps(users, func(ctx context.Context) (map[string][]types.ChipData, error) {
				return map[string][]types.ChipData{
					"user-active":   {{Label: "Acme Corp"}},
					"user-inactive": {{Label: "Globex"}},
				}, nil
			}, &capturedReq)

			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions(tc.permCodes))
			cols := userColumns(deps.Labels)
			table, err := buildTableConfig(ctx, deps, cols, tc.status, tableparams.TableQueryParams{
				Page:       2,
				PageSize:   10,
				Search:     "ada",
				SortColumn: "first_name",
				SortDir:    "asc",
				Timezone:   "UTC",
			})
			if err != nil {
				t.Fatalf("buildTableConfig() error = %v", err)
			}

			if capturedReq == nil {
				t.Fatalf("GetListPageData was not called")
			}
			if got, want := capturedReq.GetPagination().GetLimit(), int32(10); got != want {
				t.Fatalf("pagination limit = %d, want %d", got, want)
			}
			if got, want := capturedReq.GetPagination().GetOffset().GetPage(), int32(2); got != want {
				t.Fatalf("pagination page = %d, want %d", got, want)
			}
			if got, want := capturedReq.GetSort().GetFields()[0].GetField(), "first_name"; got != want {
				t.Fatalf("sort field = %q, want %q", got, want)
			}

			if table.PrimaryAction == nil {
				t.Fatalf("PrimaryAction is nil")
			}
			if got, want := table.PrimaryAction.Disabled, tc.wantPrimaryDisabled; got != want {
				t.Fatalf("PrimaryAction.Disabled = %v, want %v", got, want)
			}

			if got, want := len(table.Rows), 1; got != want {
				t.Fatalf("row count = %d, want %d", got, want)
			}
			actions := table.Rows[0].Actions
			if got, want := len(actions), 3; got != want {
				t.Fatalf("actions count = %d, want %d", got, want)
			}
			if got, want := actions[1].Type, "edit"; got != want {
				t.Fatalf("edit action type = %q, want %q", got, want)
			}
			if got, want := actions[1].Disabled, tc.wantRowDisabled; got != want {
				t.Fatalf("edit action disabled = %v, want %v", got, want)
			}
			if got, want := actions[2].Type, tc.wantStatusAction; got != want {
				t.Fatalf("status action type = %q, want %q", got, want)
			}
			if got, want := actions[2].Disabled, tc.wantRowDisabled; got != want {
				t.Fatalf("status action disabled = %v, want %v", got, want)
			}
		})
	}
}

func TestBuildTableConfig_GetListPageDataError(t *testing.T) {
	t.Parallel()

	deps := &ListViewDeps{
		Routes: entydad.DefaultUserRoutes(),
		GetListPageData: func(context.Context, *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error) {
			return nil, errors.New("db unavailable")
		},
		Labels:       entydad.UserLabels{},
		SharedLabels: entydad.SharedLabels{},
		CommonLabels: pyeza.CommonLabels{},
		TableLabels:  types.TableLabels{},
	}

	cols := userColumns(deps.Labels)
	_, err := buildTableConfig(context.Background(), deps, cols, "active", tableparams.TableQueryParams{
		Page:       1,
		PageSize:   25,
		SortColumn: "date_created",
		SortDir:    "desc",
		Timezone:   "UTC",
	})
	if err == nil {
		t.Fatalf("buildTableConfig() error = nil, want non-nil")
	}
}

func newListViewDeps(
	users []*userpb.User,
	getWorkspaces func(ctx context.Context) (map[string][]types.ChipData, error),
	capturedReq **userpb.GetUserListPageDataRequest,
) *ListViewDeps {
	return &ListViewDeps{
		Routes: entydad.DefaultUserRoutes(),
		GetListPageData: func(ctx context.Context, req *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error) {
			*capturedReq = req
			// Simulate server-side active filter so the test reflects real behaviour.
			wantActive := true
			for _, f := range req.GetFilters().GetFilters() {
				if f.GetField() == "active" {
					wantActive = f.GetBooleanFilter().GetValue()
				}
			}
			var filtered []*userpb.User
			for _, u := range users {
				if u.GetActive() == wantActive {
					filtered = append(filtered, u)
				}
			}
			return &userpb.GetUserListPageDataResponse{
				UserList: filtered,
				Pagination: &commonpb.PaginationResponse{
					TotalItems: int32(len(filtered)),
				},
				Success: true,
			}, nil
		},
		GetUserWorkspacesMap: getWorkspaces,
		Labels:          entydad.UserLabels{},
		SharedLabels: entydad.SharedLabels{
			Confirm: entydad.SharedConfirmLabels{
				Activate:   "activate %s",
				Deactivate: "deactivate %s",
			},
			Badges: entydad.SharedBadgeLabels{
				NoPermission: "No permission",
			},
		},
		CommonLabels: pyeza.CommonLabels{},
		TableLabels:  types.TableLabels{},
	}
}

func strPtr(v string) *string {
	return &v
}
