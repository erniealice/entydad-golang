package permissions

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetRoleItemPageData func(ctx context.Context, req *rolepb.GetRoleItemPageDataRequest) (*rolepb.GetRoleItemPageDataResponse, error)
	Labels              entydad.RolePermissionLabels
	CommonLabels        pyeza.CommonLabels
	TableLabels         types.TableLabels
}

// PageData holds the data for the role permissions page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	RoleID          string
	RoleName        string
}

// NewView creates the role permissions view (full page).
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
				HeaderIcon:     "icon-key",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "role-permissions-content",
			Table:           tableConfig,
			RoleID:          roleID,
			RoleName:        roleName,
		}

		return view.OK("role-permissions", pageData)
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

// buildTableConfig fetches role data with permissions and builds the table.
func buildTableConfig(ctx context.Context, deps *Deps, roleID string) (*types.TableConfig, string, error) {
	resp, err := deps.GetRoleItemPageData(ctx, &rolepb.GetRoleItemPageDataRequest{
		RoleId: roleID,
	})
	if err != nil {
		log.Printf("Failed to get role item page data: %v", err)
		return nil, "", fmt.Errorf("failed to load role: %w", err)
	}

	role := resp.GetRole()
	roleName := role.GetName()

	l := deps.Labels
	columns := permissionColumns(l)
	rows := buildTableRows(role, l)
	types.ApplyColumnStyles(columns, rows)

	refreshURL := fmt.Sprintf("/action/manage/roles/%s/permissions/table", roleID)

	tableConfig := &types.TableConfig{
		ID:                   "role-permissions-table",
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
		DefaultSortColumn:    "permissionName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Buttons.AssignPermission,
			ActionURL: fmt.Sprintf("/action/manage/roles/%s/permissions/assign", roleID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, roleName, nil
}

func permissionColumns(l entydad.RolePermissionLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "permissionName", Label: l.Columns.PermissionName, Sortable: true},
		{Key: "code", Label: l.Columns.Code, Sortable: true},
		{Key: "type", Label: l.Columns.Type, Sortable: true, Width: "120px"},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, Sortable: true, Width: "180px"},
	}
}

func buildTableRows(role *rolepb.Role, l entydad.RolePermissionLabels) []types.TableRow {
	rows := []types.TableRow{}

	for _, rp := range role.GetRolePermissions() {
		perm := rp.GetPermission()
		if perm == nil {
			continue
		}

		rpID := rp.GetId()
		permName := perm.GetName()
		permCode := perm.GetPermissionCode()
		permType := "Allow"
		dateAssigned := rp.GetDateCreatedString()

		// Permission type badge — from the embedded Permission object
		pt := perm.GetPermissionType()
		if pt == permissionpb.PermissionType_PERMISSION_TYPE_DENY {
			permType = "Deny"
		}
		typeVariant := "success"
		if permType == "Deny" {
			typeVariant = "danger"
		}

		actions := []types.TableAction{
			{
				Type: "delete", Label: l.Actions.Remove, Action: "delete",
				URL:            fmt.Sprintf("/action/manage/roles/%s/permissions/remove", role.GetId()),
				ItemName:       permName,
				ConfirmTitle:   l.Actions.Remove,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to remove %s from this role?", permName),
			},
		}

		rows = append(rows, types.TableRow{
			ID: rpID,
			Cells: []types.TableCell{
				{Type: "text", Value: permName},
				{Type: "text", Value: permCode},
				{Type: "badge", Value: permType, Variant: typeVariant},
				{Type: "text", Value: dateAssigned},
			},
			DataAttrs: map[string]string{
				"permissionName": permName,
				"code":           permCode,
				"type":           permType,
			},
			Actions: actions,
		})
	}
	return rows
}
