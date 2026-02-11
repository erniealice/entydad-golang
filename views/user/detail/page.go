package detail

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
)

// Deps holds view dependencies.
type Deps struct {
	ReadUser                     func(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	ListWorkspaceUsers           func(ctx context.Context, req *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
	Labels                       entydad.UserLabels
	UserRoleLabels               entydad.UserRoleLabels
	CommonLabels                 pyeza.CommonLabels
	TableLabels                  types.TableLabels
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

// PageData holds the data for the user detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          entydad.UserLabels
	ActiveTab       string
	TabItems        []TabItem
	ID              string
	UserFirstName   string
	UserLastName    string
	UserEmail       string
	UserMobile      string
	UserStatus      string
	StatusVariant   string
	RoleNames       []string
	RolesTable      *types.TableConfig
}

// NewView creates the user detail view (full page).
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

		return view.OK("user-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/users/{id}/tab/{tab}
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
		templateName := "user-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// buildPageData loads user data and builds the PageData for the given active tab.
func buildPageData(ctx context.Context, deps *Deps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadUser(ctx, &userpb.ReadUserRequest{
		Data: &userpb.User{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read user %s: %v", id, err)
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	user := data[0]

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	email := user.GetEmailAddress()
	mobile := user.GetMobileNumber()
	displayName := firstName + " " + lastName
	if firstName == "" && lastName == "" {
		displayName = email
	}

	userStatus := "active"
	if !user.GetActive() {
		userStatus = "inactive"
	}
	statusVariant := "success"
	if userStatus == "inactive" {
		statusVariant = "warning"
	}

	// Get role count for the Roles tab badge
	roleCount, roleNames := getUserRoles(ctx, deps, id)

	tabItems := buildTabItems(id, deps.Labels, roleCount)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          displayName,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "users",
			ActiveSubNav:   "users-active",
			HeaderTitle:    displayName,
			HeaderSubtitle: email,
			HeaderIcon:     "icon-user",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "user-detail-content",
		Labels:          deps.Labels,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		UserFirstName:   firstName,
		UserLastName:    lastName,
		UserEmail:       email,
		UserMobile:      mobile,
		UserStatus:      userStatus,
		StatusVariant:   statusVariant,
		RoleNames:       roleNames,
	}

	// Load tab-specific data
	switch activeTab {
	case "roles":
		tableConfig, err := buildRolesTable(ctx, deps, id)
		if err != nil {
			log.Printf("Failed to build roles table for user %s: %v", id, err)
		} else {
			pageData.RolesTable = tableConfig
		}
	}

	return pageData, nil
}

func buildTabItems(id string, labels entydad.UserLabels, roleCount int) []TabItem {
	base := "/app/users/detail/" + id
	return []TabItem{
		{Key: "info", Label: labels.Detail.Tabs.Info, Href: base + "?tab=info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "roles", Label: labels.Detail.Tabs.Roles, Href: base + "?tab=roles", Icon: "icon-shield", Count: roleCount, Disabled: false},
	}
}

func getUserRoles(ctx context.Context, deps *Deps, userID string) (int, []string) {
	if deps.ListWorkspaceUsers == nil || deps.GetWorkspaceUserItemPageData == nil {
		return 0, nil
	}

	// Find workspace user for this user ID
	wuResp, err := deps.ListWorkspaceUsers(ctx, &workspaceuserpb.ListWorkspaceUsersRequest{})
	if err != nil {
		log.Printf("Failed to list workspace users: %v", err)
		return 0, nil
	}

	var wuID string
	for _, wu := range wuResp.GetData() {
		if wu.GetUserId() == userID {
			wuID = wu.GetId()
			break
		}
	}
	if wuID == "" {
		return 0, nil
	}

	itemResp, err := deps.GetWorkspaceUserItemPageData(ctx, &workspaceuserpb.GetWorkspaceUserItemPageDataRequest{
		WorkspaceUserId: wuID,
	})
	if err != nil {
		log.Printf("Failed to get workspace user item page data: %v", err)
		return 0, nil
	}

	roles := itemResp.GetWorkspaceUser().GetWorkspaceUserRoles()
	names := make([]string, 0, len(roles))
	for _, wur := range roles {
		if r := wur.GetRole(); r != nil {
			names = append(names, r.GetName())
		}
	}
	return len(names), names
}

// ---------------------------------------------------------------------------
// Roles tab table
// ---------------------------------------------------------------------------

func buildRolesTable(ctx context.Context, deps *Deps, userID string) (*types.TableConfig, error) {
	if deps.ListWorkspaceUsers == nil || deps.GetWorkspaceUserItemPageData == nil {
		return nil, fmt.Errorf("workspace user dependencies not available")
	}

	// Find workspace user for this user ID
	wuResp, err := deps.ListWorkspaceUsers(ctx, &workspaceuserpb.ListWorkspaceUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list workspace users: %w", err)
	}

	var wuID string
	for _, wu := range wuResp.GetData() {
		if wu.GetUserId() == userID {
			wuID = wu.GetId()
			break
		}
	}
	if wuID == "" {
		// No workspace user found, return empty table
		return buildEmptyRolesTable(deps, userID), nil
	}

	itemResp, err := deps.GetWorkspaceUserItemPageData(ctx, &workspaceuserpb.GetWorkspaceUserItemPageDataRequest{
		WorkspaceUserId: wuID,
	})
	if err != nil {
		log.Printf("Failed to get workspace user item page data: %v", err)
		return buildEmptyRolesTable(deps, userID), nil
	}

	workspaceUser := itemResp.GetWorkspaceUser()
	l := deps.UserRoleLabels

	columns := []types.TableColumn{
		{Key: "roleName", Label: l.Columns.RoleName, Sortable: true},
		{Key: "description", Label: l.Columns.Description, Sortable: true},
		{Key: "color", Label: l.Columns.Color, Sortable: true, Width: "120px"},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, Sortable: true, Width: "180px"},
	}

	rows := []types.TableRow{}
	for _, wur := range workspaceUser.GetWorkspaceUserRoles() {
		role := wur.GetRole()
		if role == nil {
			continue
		}

		wurID := wur.GetId()
		roleName := role.GetName()
		description := role.GetDescription()
		color := role.GetColor()
		dateAssigned := wur.GetDateCreatedString()

		colorVariant := "default"
		if color != "" {
			colorVariant = "info"
		}

		actions := []types.TableAction{
			{
				Type: "delete", Label: l.Actions.Remove, Action: "delete",
				URL:            fmt.Sprintf("/action/users/detail/%s/roles/remove", userID),
				ItemName:       roleName,
				ConfirmTitle:   l.Actions.Remove,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to remove %s from this user?", roleName),
			},
		}

		rows = append(rows, types.TableRow{
			ID: wurID,
			Cells: []types.TableCell{
				{Type: "text", Value: roleName},
				{Type: "text", Value: description},
				{Type: "badge", Value: color, Variant: colorVariant},
				{Type: "text", Value: dateAssigned},
			},
			DataAttrs: map[string]string{
				"roleName":    roleName,
				"description": description,
				"color":       color,
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "user-roles-table",
		RefreshURL:           fmt.Sprintf("/action/users/detail/%s/roles/table", userID),
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
		DefaultSortColumn:    "roleName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Buttons.AssignRole,
			ActionURL: fmt.Sprintf("/action/users/detail/%s/roles/assign", userID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func buildEmptyRolesTable(deps *Deps, userID string) *types.TableConfig {
	l := deps.UserRoleLabels
	columns := []types.TableColumn{
		{Key: "roleName", Label: l.Columns.RoleName, Sortable: true},
		{Key: "description", Label: l.Columns.Description, Sortable: true},
		{Key: "color", Label: l.Columns.Color, Sortable: true, Width: "120px"},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, Sortable: true, Width: "180px"},
	}

	tableConfig := &types.TableConfig{
		ID:                   "user-roles-table",
		RefreshURL:           fmt.Sprintf("/action/users/detail/%s/roles/table", userID),
		Columns:              columns,
		Rows:                 []types.TableRow{},
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "roleName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Buttons.AssignRole,
			ActionURL: fmt.Sprintf("/action/users/detail/%s/roles/assign", userID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}
