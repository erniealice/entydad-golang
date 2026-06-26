package list

import (
	"context"
	"fmt"
	"log"
	"math"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"

	"github.com/erniealice/entydad-golang"
	workspace "github.com/erniealice/entydad-golang/domain/entity/identity/workspace"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	GetListPageData func(ctx context.Context, req *workspacepb.GetWorkspaceListPageDataRequest) (*workspacepb.GetWorkspaceListPageDataResponse, error)
	RefreshURL      string
	Routes          workspace.Routes
	Labels          workspace.Labels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the workspace list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	Routes          workspace.Routes
	Labels          workspace.Labels
	Permissions     struct {
		HasWorkspaceCreate bool
	}
}

var workspaceSearchFields = []string{"name", "description"}

// NewView creates the workspace list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("workspace", "list") {
			return view.Forbidden("workspace:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := workspaceColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "name", "asc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "admin",
				ActiveSubNav:   "workspaces-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-briefcase",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "workspace-list-content",
			Table:           tableConfig,
			Routes:          deps.Routes,
			Labels:          deps.Labels,
		}

		// Populate permissions for the disabled-CTA pattern
		pageData.Permissions.HasWorkspaceCreate = view.GetUserPermissions(ctx).Can("workspace", "create")

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "workspace"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("workspace-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("workspace", "list") {
			return view.Forbidden("workspace:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := workspaceColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "name", "asc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches workspace data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, workspaceSearchFields)

	// Inject status filter for server-side pagination
	activeValue := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	})

	resp, err := deps.GetListPageData(ctx, &workspacepb.GetWorkspaceListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list workspaces: %v", err)
		return nil, fmt.Errorf("failed to load workspaces: %w", err)
	}

	l := deps.Labels
	rows := buildTableRows(resp.GetWorkspaceList(), status, l, deps.SharedLabels, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, status, deps.Routes, perms)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	// Build ServerPagination
	totalRows := int(resp.GetPagination().GetTotalItems())
	sp := &types.ServerPagination{
		Enabled:           true,
		Mode:              "offset",
		CurrentPage:       p.Page,
		PageSize:          p.PageSize,
		TotalRows:         totalRows,
		TotalPages:        int(math.Ceil(float64(totalRows) / float64(p.PageSize))),
		SearchQuery:       p.Search,
		SortColumn:        p.SortColumn,
		SortDirection:     p.SortDir,
		FiltersJSON:       p.FiltersRaw,
		PaginationURL:     refreshURL,
		PaginationBodyURL: refreshURL,
	}
	sp.BuildDisplay()

	tableConfig := &types.TableConfig{
		ID:                   "workspaces-table",
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
		// PrimaryAction is handled directly in workspace-list-content template
		// via the disabled-CTA pattern with conditional if/else rendering
		PrimaryAction:    nil,
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func workspaceColumns(l workspace.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "description", Label: l.Columns.Description},
		{Key: "private", Label: l.Columns.Private, WidthClass: "col-2xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
	}
}

func buildTableRows(workspaces []*workspacepb.Workspace, status string, l workspace.Labels, sl entydad.SharedLabels, routes workspace.Routes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, w := range workspaces {
		active := w.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := w.GetId()
		name := w.GetName()
		description := w.GetDescription()
		private := w.GetPrivate()

		privateLabel := sl.Badges.No
		privateVariant := "default"
		if private {
			privateLabel = sl.Badges.Yes
			privateVariant = "info"
		}

		actions := []types.TableAction{}
		if routes.DetailURL != "" {
			actions = append(actions, types.TableAction{
				Type: "view", Label: l.Actions.View, Action: "view",
				Href: route.ResolveURL(routes.DetailURL, "id", id),
			})
		}
		actions = append(actions, types.TableAction{
			Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit,
			Disabled: !perms.Can("workspace", "update"), DisabledTooltip: sl.Badges.NoPermission,
		})
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
				Disabled:       !perms.Can("workspace", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
				Disabled:       !perms.Can("workspace", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		}
		deleteAction := types.TableAction{
			Type: "delete", Label: l.Actions.Delete, Action: "delete",
			URL: routes.DeleteURL, ItemName: name,
		}
		if !perms.Can("workspace", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = sl.Badges.NoPermission
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "badge", Value: privateLabel, Variant: privateVariant},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":        name,
				"description": description,
				"private":     fmt.Sprintf("%v", private),
				"status":      recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l workspace.Labels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l workspace.Labels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l workspace.Labels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l workspace.Labels, status string) string {
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

func buildBulkActions(l workspace.Labels, sl entydad.SharedLabels, common pyeza.CommonLabels, status string, routes workspace.Routes, perms *types.UserPermissions) []types.BulkAction {
	actions := []types.BulkAction{}

	canUpdate := perms.Can("workspace", "update")
	canDelete := perms.Can("workspace", "delete")

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-briefcase",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
			Disabled:        !canUpdate,
			DisabledTooltip: fmt.Sprintf(common.Errors.MissingPermission, "workspace:update"),
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-briefcase",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
			Disabled:        !canUpdate,
			DisabledTooltip: fmt.Sprintf(common.Errors.MissingPermission, "workspace:update"),
		})
	}

	actions = append(actions, types.BulkAction{
		Key:             "delete",
		Label:           common.Bulk.Delete,
		Icon:            "icon-trash-2",
		Variant:         "danger",
		Endpoint:        routes.BulkDeleteURL,
		ConfirmTitle:    common.Bulk.Delete,
		ConfirmMessage:  sl.Confirm.BulkDelete,
		Disabled:        !canDelete,
		DisabledTooltip: fmt.Sprintf(common.Errors.MissingPermission, "workspace:delete"),
	})

	return actions
}
