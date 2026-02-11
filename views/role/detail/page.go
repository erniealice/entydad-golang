package detail

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	roleusers "github.com/erniealice/entydad-golang/views/role/users"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
)

// Deps holds view dependencies.
type Deps struct {
	ReadRole            func(ctx context.Context, req *rolepb.ReadRoleRequest) (*rolepb.ReadRoleResponse, error)
	RoleGetItemPageData func(ctx context.Context, req *rolepb.GetRoleItemPageDataRequest) (*rolepb.GetRoleItemPageDataResponse, error)
	GetUsersByRoleID    func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error)
	Labels              entydad.RoleLabels
	RolePermissionLabels entydad.RolePermissionLabels
	RoleUserLabels       entydad.RoleUserLabels
	CommonLabels        pyeza.CommonLabels
	TableLabels         types.TableLabels
}

// TabItem represents a tab in the detail view.
type TabItem struct {
	Key      string
	Label    string
	Href     string
	Icon     string
	Count    int
	Disabled bool
}

// PageData holds the data for the role detail page.
type PageData struct {
	types.PageData
	ContentTemplate  string
	Labels           entydad.RoleLabels
	ActiveTab        string
	TabItems         []TabItem
	ID               string
	RoleName         string
	RoleDescription  string
	RoleColor        string
	RoleStatus       string
	StatusVariant    string
	PermissionsTable *types.TableConfig
	UsersTable       *types.TableConfig
}

// NewView creates the role detail view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("role-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/roles/{id}/tab/{tab}
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		// Return only the tab partial template
		templateName := "role-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// buildPageData loads role data and builds the PageData for the given active tab.
func buildPageData(ctx context.Context, deps *Deps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadRole(ctx, &rolepb.ReadRoleRequest{
		Data: &rolepb.Role{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read role %s: %v", id, err)
		return nil, fmt.Errorf("failed to load role: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("role not found")
	}
	role := data[0]

	roleName := role.GetName()
	roleDesc := role.GetDescription()
	roleColor := role.GetColor()

	roleStatus := "active"
	if !role.GetActive() {
		roleStatus = "inactive"
	}
	statusVariant := "success"
	if roleStatus == "inactive" {
		statusVariant = "warning"
	}

	// Get counts for tab badges
	permCount := 0
	userCount := 0

	if deps.RoleGetItemPageData != nil {
		itemResp, err := deps.RoleGetItemPageData(ctx, &rolepb.GetRoleItemPageDataRequest{
			RoleId: id,
		})
		if err != nil {
			log.Printf("Failed to get role item page data for %s: %v", id, err)
		} else {
			permCount = len(itemResp.GetRole().GetRolePermissions())
		}
	}

	if deps.GetUsersByRoleID != nil {
		users, err := deps.GetUsersByRoleID(ctx, id)
		if err != nil {
			log.Printf("Failed to get users for role %s: %v", id, err)
		} else {
			userCount = len(users)
		}
	}

	tabItems := buildTabItems(id, deps.Labels, permCount, userCount)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          roleName,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "users",
			ActiveSubNav:   "roles-active",
			HeaderTitle:    roleName,
			HeaderSubtitle: roleDesc,
			HeaderIcon:     "icon-shield",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "role-detail-content",
		Labels:          deps.Labels,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		RoleName:        roleName,
		RoleDescription: roleDesc,
		RoleColor:       roleColor,
		RoleStatus:      roleStatus,
		StatusVariant:   statusVariant,
	}

	// Load tab-specific data
	switch activeTab {
	case "permissions":
		tableConfig, err := buildPermissionsTable(ctx, deps, id)
		if err != nil {
			log.Printf("Failed to build permissions table for role %s: %v", id, err)
		} else {
			pageData.PermissionsTable = tableConfig
		}
	case "users":
		tableConfig, err := buildUsersTable(ctx, deps, id)
		if err != nil {
			log.Printf("Failed to build users table for role %s: %v", id, err)
		} else {
			pageData.UsersTable = tableConfig
		}
	}

	return pageData, nil
}

func buildTabItems(id string, labels entydad.RoleLabels, permCount, userCount int) []TabItem {
	base := "/app/roles/detail/" + id
	return []TabItem{
		{Key: "info", Label: labels.Detail.Tabs.Info, Href: base + "?tab=info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "permissions", Label: labels.Detail.Tabs.Permissions, Href: base + "?tab=permissions", Icon: "icon-key", Count: permCount, Disabled: false},
		{Key: "users", Label: labels.Detail.Tabs.Users, Href: base + "?tab=users", Icon: "icon-user", Count: userCount, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Permissions tab table
// ---------------------------------------------------------------------------

func buildPermissionsTable(ctx context.Context, deps *Deps, roleID string) (*types.TableConfig, error) {
	if deps.RoleGetItemPageData == nil {
		return nil, fmt.Errorf("RoleGetItemPageData not available")
	}

	resp, err := deps.RoleGetItemPageData(ctx, &rolepb.GetRoleItemPageDataRequest{
		RoleId: roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load role permissions: %w", err)
	}

	role := resp.GetRole()
	l := deps.RolePermissionLabels

	columns := []types.TableColumn{
		{Key: "permissionName", Label: l.Columns.PermissionName, Sortable: true},
		{Key: "code", Label: l.Columns.Code, Sortable: true},
		{Key: "type", Label: l.Columns.Type, Sortable: true, Width: "120px"},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, Sortable: true, Width: "180px"},
	}

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
				URL:            fmt.Sprintf("/action/roles/detail/%s/permissions/remove", roleID),
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

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "role-permissions-table",
		RefreshURL:           fmt.Sprintf("/action/roles/detail/%s/permissions/table", roleID),
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
			ActionURL: fmt.Sprintf("/action/roles/detail/%s/permissions/assign", roleID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

// ---------------------------------------------------------------------------
// Users tab table
// ---------------------------------------------------------------------------

func buildUsersTable(ctx context.Context, deps *Deps, roleID string) (*types.TableConfig, error) {
	var users []roleusers.UserByRole
	if deps.GetUsersByRoleID != nil {
		var err error
		users, err = deps.GetUsersByRoleID(ctx, roleID)
		if err != nil {
			log.Printf("Failed to get users for role %s: %v", roleID, err)
			// Continue with empty table
		}
	}

	l := deps.RoleUserLabels

	columns := []types.TableColumn{
		{Key: "userName", Label: l.Columns.UserName, Sortable: true},
		{Key: "email", Label: l.Columns.Email, Sortable: true},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, Sortable: true, Width: "180px"},
	}

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

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "role-users-table",
		RefreshURL:           fmt.Sprintf("/action/roles/detail/%s/users/table", roleID),
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

	return tableConfig, nil
}
