package list

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error)
	GetUserRolesMap func(ctx context.Context) (map[string][]entydad.RoleBadge, error)
	RefreshURL      string
	Labels          entydad.UserLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the user list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the user list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		tableConfig, err := buildTableConfig(ctx, deps, status)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "users",
				ActiveSubNav:   "users-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "user-list-content",
			Table:           tableConfig,
		}

		return view.OK("user-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		tableConfig, err := buildTableConfig(ctx, deps, status)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches user data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps, status string) (*types.TableConfig, error) {
	resp, err := deps.GetListPageData(ctx, &userpb.GetUserListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return nil, fmt.Errorf("failed to load users: %w", err)
	}

	// Fetch user-role mappings (best-effort; nil map means no role data)
	var userRolesMap map[string][]entydad.RoleBadge
	if deps.GetUserRolesMap != nil {
		userRolesMap, err = deps.GetUserRolesMap(ctx)
		if err != nil {
			log.Printf("Warning: Failed to load user roles map: %v", err)
		}
	}

	l := deps.Labels
	columns := userColumns(l)
	rows := buildTableRows(resp.GetUserList(), status, l, userRolesMap)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status)

	refreshURL := fmt.Sprintf("/action/users/table/%s", status)

	tableConfig := &types.TableConfig{
		ID:                   "users-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Buttons.AddUser,
			ActionURL: "/action/users/add",
			Icon:      "icon-plus",
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func userColumns(l entydad.UserLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "email", Label: l.Columns.Email, Sortable: true},
		{Key: "roles", Label: l.Columns.Roles, Sortable: false},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(users []*userpb.User, status string, l entydad.UserLabels, userRolesMap map[string][]entydad.RoleBadge) []types.TableRow {
	rows := []types.TableRow{}
	for _, u := range users {
		active := u.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := u.GetId()
		name := u.GetFirstName() + " " + u.GetLastName()
		email := u.GetEmailAddress()

		// Build role chips for this user
		var roleNames []string
		if userRolesMap != nil {
			for _, rb := range userRolesMap[id] {
				roleNames = append(roleNames, rb.Name)
			}
		}
		rolesCell := types.BuildChipCellFromLabels(roleNames, 2)

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: "/app/users/" + id},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: "/action/users/edit/" + id, DrawerTitle: l.Actions.Edit},
			{Type: "view", Label: l.Actions.ManageRoles, Action: "view", Href: "/app/manage/users/" + id + "/roles"},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: "/action/users/set-status?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to deactivate %s?", name),
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: "/action/users/set-status?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to activate %s?", name),
			})
		}
		actions = append(actions, types.TableAction{
			Type: "delete", Label: l.Actions.Delete, Action: "delete",
			URL: "/action/users/delete", ItemName: name,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: email},
				rolesCell,
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"email":  email,
				"status": recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l entydad.UserLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l entydad.UserLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.UserLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.UserLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "inactive":
		return l.Empty.InactiveMessage
	default:
		return l.Empty.ActiveMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}

func buildBulkActions(l entydad.UserLabels, common pyeza.CommonLabels, status string) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-user-minus",
			Variant:         "warning",
			Endpoint:        "/action/users/bulk-set-status",
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  "Are you sure you want to deactivate {{count}} user(s)?",
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-user-check",
			Variant:         "primary",
			Endpoint:        "/action/users/bulk-set-status",
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  "Are you sure you want to activate {{count}} user(s)?",
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:            "delete",
		Label:          common.Bulk.Delete,
		Icon:           "icon-trash-2",
		Variant:        "danger",
		Endpoint:       "/action/users/bulk-delete",
		ConfirmTitle:   common.Bulk.Delete,
		ConfirmMessage: "Are you sure you want to delete {{count}} user(s)? This action cannot be undone.",
	})

	return actions
}
