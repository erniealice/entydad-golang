package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *rolepb.GetRoleListPageDataRequest) (*rolepb.GetRoleListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	RefreshURL      string
	Routes          entydad.RoleRoutes
	Labels          entydad.RoleLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the role list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the role list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		tableConfig, err := buildTableConfig(ctx, deps)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "users",
				ActiveSubNav:   "roles",
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-shield",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "role-list-content",
			Table:           tableConfig,
		}

		return view.OK("role-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		tableConfig, err := buildTableConfig(ctx, deps)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches role data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	resp, err := deps.GetListPageData(ctx, &rolepb.GetRoleListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list roles: %v", err)
		return nil, fmt.Errorf("failed to load roles: %w", err)
	}

	// Check which items are in use
	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range resp.GetRoleList() {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	l := deps.Labels
	columns := roleColumns(l)
	rows := buildTableRows(resp.GetRoleList(), l, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, deps.Routes)

	tableConfig := &types.TableConfig{
		ID:                   "roles-table",
		RefreshURL:           deps.Routes.TableURL,
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
			Title:   l.Empty.ActiveTitle,
			Message: l.Empty.ActiveMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddRole,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("role", "create"),
			DisabledTooltip: "No permission",
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func roleColumns(l entydad.RoleLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "description", Label: l.Columns.Description, Sortable: true},
		{Key: "color", Label: l.Columns.Color, Sortable: true, Width: "120px"},
		{Key: "permissions", Label: l.Columns.Permissions, Sortable: false, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(roles []*rolepb.Role, l entydad.RoleLabels, routes entydad.RoleRoutes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range roles {
		active := r.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := r.GetId()
		name := r.GetName()
		description := r.GetDescription()
		color := r.GetColor()

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to deactivate %s?", name),
				Disabled: !perms.Can("role", "update"), DisabledTooltip: "No permission",
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to activate %s?", name),
				Disabled: !perms.Can("role", "update"), DisabledTooltip: "No permission",
			})
		}
		isInUse := inUseIDs[id]
		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if isInUse {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = "Cannot delete: role is assigned to users"
		} else if !perms.Can("role", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = "No permission"
		}
		actions = append(actions, deleteAction)

		permCount := len(r.GetRolePermissions())
		permCountStr := fmt.Sprintf("%d", permCount)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "text", Value: color},
				{Type: "badge", Value: permCountStr, Variant: "default", BadgeType: "count"},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":        name,
				"description": description,
				"color":       color,
				"permissions": permCountStr,
				"status":      recordStatus,
				"deletable":   strconv.FormatBool(!isInUse),
			},
			Actions: actions,
		})
	}
	return rows
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

func buildBulkActions(l entydad.RoleLabels, common pyeza.CommonLabels, routes entydad.RoleRoutes) []types.BulkAction {
	return []types.BulkAction{
		{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-shield",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  "Are you sure you want to activate {{count}} role(s)?",
			ExtraParamsJSON: `{"target_status":"active"}`,
		},
		{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-shield-off",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  "Are you sure you want to deactivate {{count}} role(s)?",
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		},
		{
			Key:              "delete",
			Label:            common.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         routes.BulkDeleteURL,
			ConfirmTitle:     common.Bulk.Delete,
			ConfirmMessage:   "Are you sure you want to delete {{count}} role(s)? This action cannot be undone.",
			RequiresDataAttr: "deletable",
		},
	}
}
