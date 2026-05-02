package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	GetListPageData func(ctx context.Context, status string, search string, page, pageSize int) (*LocationAreaListResult, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	RefreshURL      string
	Routes          entydad.LocationAreaRoutes
	Labels          entydad.LocationAreaLabels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// LocationAreaItem represents a single location area record for list display.
type LocationAreaItem struct {
	ID          string
	Name        string
	Description string
	Active      bool
	DateCreated string
}

// LocationAreaListResult holds the paginated list response.
type LocationAreaListResult struct {
	Items      []*LocationAreaItem
	TotalItems int
}

// PageData holds the data for the location area list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var locationAreaSearchFields = []string{"name", "description"}

// NewView creates the location area list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := locationAreaColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "date_created", "desc")
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
				ActiveNav:      "location",
				ActiveSubNav:   "location-areas-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-map",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "location-area-list-content",
			Table:           tableConfig,
		}

		return view.OK("location-area-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := locationAreaColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "date_created", "desc")
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

// buildTableConfig fetches location area data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	resp, err := deps.GetListPageData(ctx, status, p.Search, p.Page, p.PageSize)
	if err != nil {
		log.Printf("Failed to list location areas: %v", err)
		return nil, fmt.Errorf("failed to load location areas: %w", err)
	}

	// Check which items are in use
	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range resp.Items {
			itemIDs = append(itemIDs, item.ID)
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	l := deps.Labels
	rows := buildTableRows(resp.Items, status, l, deps.SharedLabels, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	totalRows := resp.TotalItems
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
		ID:                   "location-areas-table",
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
		DefaultSortColumn:    "date_created",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddLocationArea,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("location_area", "create"),
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		},
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func locationAreaColumns(l entydad.LocationAreaLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "description", Label: l.Columns.Description, NoSort: true, NoFilter: true},
		{Key: "date_created", Label: l.Columns.DateCreated},
	}
}

func buildTableRows(items []*LocationAreaItem, status string, l entydad.LocationAreaLabels, sl entydad.SharedLabels, routes entydad.LocationAreaRoutes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		active := item.Active
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := item.ID
		name := item.Name
		isInUse := inUseIDs[id]

		actions := []types.TableAction{
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit,
				Disabled: !perms.Can("location_area", "update"), DisabledTooltip: sl.Badges.NoPermission},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
				Disabled:       !perms.Can("location_area", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
				Disabled:       !perms.Can("location_area", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		}
		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if isInUse {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = sl.Errors.CannotDeleteInUse
		} else if !perms.Can("location_area", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = sl.Badges.NoPermission
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: item.Description},
				types.DateTimeCell(item.DateCreated, types.DateReadable),
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l entydad.LocationAreaLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l entydad.LocationAreaLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.LocationAreaLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.LocationAreaLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "inactive":
		return l.Empty.InactiveMessage
	default:
		return l.Empty.ActiveMessage
	}
}

func buildBulkActions(l entydad.LocationAreaLabels, sl entydad.SharedLabels, common pyeza.CommonLabels, status string, routes entydad.LocationAreaRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-map-off",
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
			Icon:            "icon-map",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:              "delete",
		Label:            common.Bulk.Delete,
		Icon:             "icon-trash-2",
		Variant:          "danger",
		Endpoint:         routes.BulkDeleteURL,
		ConfirmTitle:     common.Bulk.Delete,
		ConfirmMessage:   sl.Confirm.BulkDelete,
		RequiresDataAttr: "deletable",
	})

	return actions
}
