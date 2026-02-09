package roles

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	ListWorkspaceUsers           func(ctx context.Context, req *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error)
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	ReadUser                     func(ctx context.Context, req *userpb.ReadUserRequest) (*userpb.ReadUserResponse, error)
	Labels                       entydad.UserRoleLabels
	CommonLabels                 pyeza.CommonLabels
	TableLabels                  types.TableLabels
}

// PageData holds the data for the user roles page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	UserID          string
	UserName        string
}

// NewView creates the user roles view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		userID := viewCtx.Request.PathValue("id")
		if userID == "" {
			return view.Error(fmt.Errorf("user ID is required"))
		}

		tableConfig, userName, err := buildTableConfig(ctx, deps, userID)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          fmt.Sprintf("%s - %s", deps.Labels.Page.Heading, userName),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "users",
				ActiveSubNav:   "users-active",
				HeaderTitle:    fmt.Sprintf("%s: %s", deps.Labels.Page.Heading, userName),
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-shield",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "user-roles-content",
			Table:           tableConfig,
			UserID:          userID,
			UserName:        userName,
		}

		return view.OK("user-roles", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		userID := viewCtx.Request.PathValue("id")
		if userID == "" {
			return view.Error(fmt.Errorf("user ID is required"))
		}

		tableConfig, _, err := buildTableConfig(ctx, deps, userID)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// findWorkspaceUserByUserID lists workspace_users and finds the one matching the given user ID.
func findWorkspaceUserByUserID(ctx context.Context, deps *Deps, userID string) (*workspaceuserpb.WorkspaceUser, error) {
	resp, err := deps.ListWorkspaceUsers(ctx, &workspaceuserpb.ListWorkspaceUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list workspace users: %w", err)
	}

	for _, wu := range resp.GetData() {
		if wu.GetUserId() == userID {
			return wu, nil
		}
	}

	return nil, fmt.Errorf("workspace user not found for user ID %s", userID)
}

// buildTableConfig fetches user data with roles and builds the table.
func buildTableConfig(ctx context.Context, deps *Deps, userID string) (*types.TableConfig, string, error) {
	// First, get the user's name
	userResp, err := deps.ReadUser(ctx, &userpb.ReadUserRequest{
		Data: &userpb.User{Id: userID},
	})
	if err != nil {
		log.Printf("Failed to read user %s: %v", userID, err)
		return nil, "", fmt.Errorf("failed to load user: %w", err)
	}
	user := userResp.GetData()[0]
	userName := user.GetFirstName() + " " + user.GetLastName()

	// Find the workspace_user record for this user
	wu, err := findWorkspaceUserByUserID(ctx, deps, userID)
	if err != nil {
		log.Printf("Failed to find workspace user for user %s: %v", userID, err)
		// If no workspace_user found, show empty table
		return buildEmptyTableConfig(deps, userID), userName, nil
	}

	// Get full workspace_user with roles via item page data
	wuResp, err := deps.GetWorkspaceUserItemPageData(ctx, &workspaceuserpb.GetWorkspaceUserItemPageDataRequest{
		WorkspaceUserId: wu.GetId(),
	})
	if err != nil {
		log.Printf("Failed to get workspace user item page data: %v", err)
		return buildEmptyTableConfig(deps, userID), userName, nil
	}

	workspaceUser := wuResp.GetWorkspaceUser()

	l := deps.Labels
	columns := roleColumns(l)
	rows := buildTableRows(workspaceUser, userID, l)
	types.ApplyColumnStyles(columns, rows)

	refreshURL := fmt.Sprintf("/action/manage/users/%s/roles/table", userID)

	tableConfig := &types.TableConfig{
		ID:                   "user-roles-table",
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
		DefaultSortColumn:    "roleName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Buttons.AssignRole,
			ActionURL: fmt.Sprintf("/action/manage/users/%s/roles/assign", userID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, userName, nil
}

func buildEmptyTableConfig(deps *Deps, userID string) *types.TableConfig {
	l := deps.Labels
	columns := roleColumns(l)

	refreshURL := fmt.Sprintf("/action/manage/users/%s/roles/table", userID)

	tableConfig := &types.TableConfig{
		ID:                   "user-roles-table",
		RefreshURL:           refreshURL,
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
			ActionURL: fmt.Sprintf("/action/manage/users/%s/roles/assign", userID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func roleColumns(l entydad.UserRoleLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "roleName", Label: l.Columns.RoleName, Sortable: true},
		{Key: "description", Label: l.Columns.Description, Sortable: true},
		{Key: "color", Label: l.Columns.Color, Sortable: true, Width: "120px"},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, Sortable: true, Width: "180px"},
	}
}

func buildTableRows(workspaceUser *workspaceuserpb.WorkspaceUser, userID string, l entydad.UserRoleLabels) []types.TableRow {
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

		// Color badge
		colorVariant := "default"
		if color != "" {
			colorVariant = "info"
		}

		actions := []types.TableAction{
			{
				Type: "delete", Label: l.Actions.Remove, Action: "delete",
				URL:            fmt.Sprintf("/action/manage/users/%s/roles/remove", userID),
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
	return rows
}
