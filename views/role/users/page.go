package users

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"

	"github.com/erniealice/entydad-golang"
)

// UserByRole holds user info for display in the role-users table.
type UserByRole struct {
	WorkspaceUserRoleID string
	WorkspaceUserID     string
	UserID              string
	UserName            string
	Email               string
	DateAssigned        string
}

// Deps holds view dependencies.
type Deps struct {
	GetUsersByRoleID func(ctx context.Context, roleID string) ([]UserByRole, error)
	ReadRole         func(ctx context.Context, req *rolepb.ReadRoleRequest) (*rolepb.ReadRoleResponse, error)
	Labels           entydad.RoleUserLabels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels
}

// PageData holds the data for the role users page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	RoleID          string
	RoleName        string
}

// NewView creates the role users view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		roleID := viewCtx.Request.PathValue("id")
		if roleID == "" {
			return view.Error(fmt.Errorf("role ID is required"))
		}

		tableConfig, roleName, err := buildTableConfig(ctx, deps, roleID)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          fmt.Sprintf("%s - %s", deps.Labels.Page.Heading, roleName),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "users",
				ActiveSubNav:   "roles-active",
				HeaderTitle:    fmt.Sprintf("%s: %s", deps.Labels.Page.Heading, roleName),
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "role-users-content",
			Table:           tableConfig,
			RoleID:          roleID,
			RoleName:        roleName,
		}

		return view.OK("role-users", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		roleID := viewCtx.Request.PathValue("id")
		if roleID == "" {
			return view.Error(fmt.Errorf("role ID is required"))
		}

		tableConfig, _, err := buildTableConfig(ctx, deps, roleID)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches users assigned to a role and builds the table.
func buildTableConfig(ctx context.Context, deps *Deps, roleID string) (*types.TableConfig, string, error) {
	// Get role name for the header
	roleResp, err := deps.ReadRole(ctx, &rolepb.ReadRoleRequest{
		Data: &rolepb.Role{Id: roleID},
	})
	if err != nil {
		log.Printf("Failed to read role %s: %v", roleID, err)
		return nil, "", fmt.Errorf("failed to load role: %w", err)
	}
	data := roleResp.GetData()
	if len(data) == 0 {
		return nil, "", fmt.Errorf("role not found")
	}
	roleName := data[0].GetName()

	// Get users assigned to this role
	var users []UserByRole
	if deps.GetUsersByRoleID != nil {
		users, err = deps.GetUsersByRoleID(ctx, roleID)
		if err != nil {
			log.Printf("Failed to get users for role %s: %v", roleID, err)
			// Continue with empty table rather than erroring
		}
	}

	l := deps.Labels
	columns := userColumns(l)
	rows := buildTableRows(users, roleID, l)
	types.ApplyColumnStyles(columns, rows)

	refreshURL := fmt.Sprintf("/action/roles/detail/%s/users/table", roleID)

	tableConfig := &types.TableConfig{
		ID:                   "role-users-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "userName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Buttons.AssignUser,
			ActionURL: fmt.Sprintf("/action/roles/detail/%s/users/assign", roleID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, roleName, nil
}

func userColumns(l entydad.RoleUserLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "userName", Label: l.Columns.UserName, Sortable: true},
		{Key: "email", Label: l.Columns.Email, Sortable: true},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, Sortable: true, Width: "180px"},
	}
}

func buildTableRows(users []UserByRole, roleID string, l entydad.RoleUserLabels) []types.TableRow {
	rows := []types.TableRow{}

	for _, u := range users {
		actions := []types.TableAction{
			{
				Type: "delete", Label: l.Actions.Remove, Action: "delete",
				URL:            fmt.Sprintf("/action/roles/detail/%s/users/remove", roleID),
				ItemName:       u.UserName,
				ConfirmTitle:   l.Actions.Remove,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to remove %s from this role?", u.UserName),
			},
		}

		rows = append(rows, types.TableRow{
			ID: u.WorkspaceUserRoleID,
			Cells: []types.TableCell{
				{Type: "text", Value: u.UserName},
				{Type: "text", Value: u.Email},
				{Type: "text", Value: u.DateAssigned},
			},
			DataAttrs: map[string]string{
				"userName": u.UserName,
				"email":    u.Email,
			},
			Actions: actions,
		})
	}
	return rows
}
