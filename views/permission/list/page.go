package list

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *permissionpb.GetPermissionListPageDataRequest) (*permissionpb.GetPermissionListPageDataResponse, error)
	RefreshURL      string
	Routes          entydad.PermissionRoutes
	Labels          entydad.PermissionLabels
	SharedLabels    entydad.SharedLabels
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
	perms := view.GetUserPermissions(ctx)

	resp, err := deps.GetListPageData(ctx, &permissionpb.GetPermissionListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list permissions: %v", err)
		return nil, fmt.Errorf("failed to load permissions: %w", err)
	}

	l := deps.Labels
	columns := permissionColumns(l)
	rows := buildTableRows(resp.GetPermissionList(), status, l, deps.SharedLabels, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

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
		DefaultSortColumn:    "permission_code",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddPermission,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("permission", "create"),
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func permissionColumns(l entydad.PermissionLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "entity", Label: l.Columns.Entity, Sortable: true, Width: "140px"},
		{Key: "permission_code", Label: l.Columns.PermissionCode, Sortable: true},
		{Key: "permission_type", Label: l.Columns.Type, Sortable: true, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

// extractEntity returns the entity prefix from a colon-notation permission code.
// For example, "client:read" returns "client". If no colon is found, returns the full code.
func extractEntity(code string) string {
	if idx := strings.Index(code, ":"); idx >= 0 {
		return code[:idx]
	}
	return code
}

func buildTableRows(permissions []*permissionpb.Permission, status string, l entydad.PermissionLabels, sl entydad.SharedLabels, routes entydad.PermissionRoutes, perms *types.UserPermissions) []types.TableRow {
	// Filter permissions by status first
	filtered := make([]*permissionpb.Permission, 0, len(permissions))
	for _, p := range permissions {
		active := p.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus == status {
			filtered = append(filtered, p)
		}
	}

	// Sort by entity prefix, then by full permission code for readability
	sort.Slice(filtered, func(i, j int) bool {
		codeI := filtered[i].GetPermissionCode()
		codeJ := filtered[j].GetPermissionCode()
		entityI := extractEntity(codeI)
		entityJ := extractEntity(codeJ)
		if entityI != entityJ {
			return entityI < entityJ
		}
		return codeI < codeJ
	})

	rows := make([]types.TableRow, 0, len(filtered))
	for _, p := range filtered {
		active := p.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := p.GetId()
		name := p.GetName()
		code := p.GetPermissionCode()
		entity := extractEntity(code)
		permType := formatPermissionType(p.GetPermissionType(), sl)

		actions := []types.TableAction{
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit,
				Disabled: !perms.Can("permission", "update"), DisabledTooltip: sl.Badges.NoPermission},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
				Disabled:       !perms.Can("permission", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
				Disabled:       !perms.Can("permission", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		}
		deleteAction := types.TableAction{
			Type: "delete", Label: l.Actions.Delete, Action: "delete",
			URL: routes.DeleteURL, ItemName: name,
		}
		if !perms.Can("permission", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = "No permission"
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "badge", Value: entity, Variant: "default", BadgeType: "type"},
				{Type: "text", Value: code},
				{Type: "badge", Value: permType, Variant: permTypeVariant(p.GetPermissionType())},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":            name,
				"entity":          entity,
				"permission_code": code,
				"permission_type": permType,
				"status":          recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func formatPermissionType(pt permissionpb.PermissionType, sl entydad.SharedLabels) string {
	switch pt {
	case permissionpb.PermissionType_PERMISSION_TYPE_ALLOW:
		return sl.Badges.Allow
	case permissionpb.PermissionType_PERMISSION_TYPE_DENY:
		return sl.Badges.Deny
	default:
		return sl.Badges.Allow
	}
}

func permTypeVariant(pt permissionpb.PermissionType) string {
	switch pt {
	case permissionpb.PermissionType_PERMISSION_TYPE_ALLOW:
		return "success"
	case permissionpb.PermissionType_PERMISSION_TYPE_DENY:
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

func buildBulkActions(l entydad.PermissionLabels, sl entydad.SharedLabels, common pyeza.CommonLabels, status string, routes entydad.PermissionRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-key",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-key",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:            "delete",
		Label:          common.Bulk.Delete,
		Icon:           "icon-trash-2",
		Variant:        "danger",
		Endpoint:       routes.BulkDeleteURL,
		ConfirmTitle:   common.Bulk.Delete,
		ConfirmMessage: sl.Confirm.BulkDelete,
	})

	return actions
}
