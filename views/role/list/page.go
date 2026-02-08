package list

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

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *rolepb.GetRoleListPageDataRequest) (*rolepb.GetRoleListPageDataResponse, error)
	RefreshURL      string
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
				ActiveSubNav:   "roles-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
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

// buildTableConfig fetches role data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps, status string) (*types.TableConfig, error) {
	resp, err := deps.GetListPageData(ctx, &rolepb.GetRoleListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list roles: %v", err)
		return nil, fmt.Errorf("failed to load roles: %w", err)
	}

	l := deps.Labels
	columns := roleColumns(l)
	rows := buildTableRows(resp.GetRoleList(), status, l)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status)

	refreshURL := fmt.Sprintf("/action/roles/table/%s", status)

	tableConfig := &types.TableConfig{
		ID:                   "roles-table",
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
			Label:     l.Buttons.AddRole,
			ActionURL: "/action/roles/add",
			Icon:      "icon-plus",
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
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(roles []*rolepb.Role, status string, l entydad.RoleLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range roles {
		active := r.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := r.GetId()
		name := r.GetName()
		description := r.GetDescription()
		color := r.GetColor()

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: "/app/roles/" + id},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: "/action/roles/edit/" + id, DrawerTitle: l.Actions.Edit},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: "/action/roles/set-status?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to deactivate %s?", name),
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: "/action/roles/set-status?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to activate %s?", name),
			})
		}
		actions = append(actions, types.TableAction{
			Type: "delete", Label: l.Actions.Delete, Action: "delete",
			URL: "/action/roles/delete", ItemName: name,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "text", Value: color},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":        name,
				"description": description,
				"color":       color,
				"status":      recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l entydad.RoleLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l entydad.RoleLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.RoleLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.RoleLabels, status string) string {
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

func buildBulkActions(l entydad.RoleLabels, common pyeza.CommonLabels, status string) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-shield-off",
			Variant:         "warning",
			Endpoint:        "/action/roles/bulk-set-status",
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  "Are you sure you want to deactivate {{count}} role(s)?",
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-shield",
			Variant:         "primary",
			Endpoint:        "/action/roles/bulk-set-status",
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  "Are you sure you want to activate {{count}} role(s)?",
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:            "delete",
		Label:          common.Bulk.Delete,
		Icon:           "icon-trash-2",
		Variant:        "danger",
		Endpoint:       "/action/roles/bulk-delete",
		ConfirmTitle:   common.Bulk.Delete,
		ConfirmMessage: "Are you sure you want to delete {{count}} role(s)? This action cannot be undone.",
	})

	return actions
}
