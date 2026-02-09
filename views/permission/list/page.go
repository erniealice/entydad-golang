package list

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *permissionpb.GetPermissionListPageDataRequest) (*permissionpb.GetPermissionListPageDataResponse, error)
	RefreshURL      string
	Labels          entydad.PermissionLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the permission list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the permission list view (full page).
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
				ActiveNav:      "admin",
				ActiveSubNav:   "permissions-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-key",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "permission-list-content",
			Table:           tableConfig,
		}

		return view.OK("permission-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
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

// buildTableConfig fetches permission data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps, status string) (*types.TableConfig, error) {
	resp, err := deps.GetListPageData(ctx, &permissionpb.GetPermissionListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list permissions: %v", err)
		return nil, fmt.Errorf("failed to load permissions: %w", err)
	}

	l := deps.Labels
	columns := permissionColumns(l)
	rows := buildTableRows(resp.GetPermissionList(), status, l)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status)

	refreshURL := fmt.Sprintf("/action/permissions/table/%s", status)

	tableConfig := &types.TableConfig{
		ID:                   "permissions-table",
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
			Label:     l.Buttons.AddPermission,
			ActionURL: "/action/permissions/add",
			Icon:      "icon-plus",
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func permissionColumns(l entydad.PermissionLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "permission_code", Label: l.Columns.PermissionCode, Sortable: true},
		{Key: "permission_type", Label: l.Columns.Type, Sortable: true, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(permissions []*permissionpb.Permission, status string, l entydad.PermissionLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, p := range permissions {
		active := p.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := p.GetId()
		name := p.GetName()
		code := p.GetPermissionCode()
		permType := formatPermissionType(p.GetPermissionType())

		actions := []types.TableAction{
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: "/action/permissions/edit/" + id, DrawerTitle: l.Actions.Edit},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: "/action/permissions/set-status?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to deactivate %s?", name),
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: "/action/permissions/set-status?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to activate %s?", name),
			})
		}
		actions = append(actions, types.TableAction{
			Type: "delete", Label: l.Actions.Delete, Action: "delete",
			URL: "/action/permissions/delete", ItemName: name,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: code},
				{Type: "badge", Value: permType, Variant: permTypeVariant(permType)},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":            name,
				"permission_code": code,
				"permission_type": permType,
				"status":          recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func formatPermissionType(pt permissionpb.PermissionType) string {
	switch pt {
	case permissionpb.PermissionType_PERMISSION_TYPE_ALLOW:
		return "Allow"
	case permissionpb.PermissionType_PERMISSION_TYPE_DENY:
		return "Deny"
	default:
		return "Allow"
	}
}

func permTypeVariant(permType string) string {
	switch permType {
	case "Allow":
		return "success"
	case "Deny":
		return "danger"
	default:
		return "default"
	}
}

func statusTitle(l entydad.PermissionLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l entydad.PermissionLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.PermissionLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.PermissionLabels, status string) string {
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

func buildBulkActions(l entydad.PermissionLabels, common pyeza.CommonLabels, status string) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-key",
			Variant:         "warning",
			Endpoint:        "/action/permissions/bulk-set-status",
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  "Are you sure you want to deactivate {{count}} permission(s)?",
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-key",
			Variant:         "primary",
			Endpoint:        "/action/permissions/bulk-set-status",
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  "Are you sure you want to activate {{count}} permission(s)?",
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:            "delete",
		Label:          common.Bulk.Delete,
		Icon:           "icon-trash-2",
		Variant:        "danger",
		Endpoint:       "/action/permissions/bulk-delete",
		ConfirmTitle:   common.Bulk.Delete,
		ConfirmMessage: "Are you sure you want to delete {{count}} permission(s)? This action cannot be undone.",
	})

	return actions
}
