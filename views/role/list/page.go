package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"

	"github.com/erniealice/entydad-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	GetListPageData func(ctx context.Context, req *rolepb.GetRoleListPageDataRequest) (*rolepb.GetRoleListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	Routes          entydad.RoleRoutes
	Labels          entydad.RoleLabels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the role list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var roleAllowedSortCols = []string{
	"date_created", "name",
}

var roleSearchFields = []string{"name", "description"}

// NewView creates the role list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		p, err := espynahttp.ParseTableParams(viewCtx.Request, roleAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, p)
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

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "roles"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("role-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		p, err := espynahttp.ParseTableParams(viewCtx.Request, roleAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches role data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, p espynahttp.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, roleSearchFields)
	resp, err := deps.GetListPageData(ctx, &rolepb.GetRoleListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
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
	rows := buildTableRows(resp.GetRoleList(), l, deps.SharedLabels, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, deps.Routes)

	refreshURL := deps.Routes.TableURL

	// Build ServerPagination
	totalRows := int(resp.GetPagination().GetTotalItems())
	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "offset",
		CurrentPage:   p.Page,
		PageSize:      p.PageSize,
		TotalRows:     totalRows,
		TotalPages:    int(math.Ceil(float64(totalRows) / float64(p.PageSize))),
		SearchQuery:   p.Search,
		SortColumn:    p.SortColumn,
		SortDirection: p.SortDir,
		FiltersJSON:   p.FiltersRaw,
		PaginationURL: refreshURL,
	}
	sp.BuildDisplay()

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
			Title:   l.Empty.ActiveTitle,
			Message: l.Empty.ActiveMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddRole,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("role", "create"),
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		},
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func roleColumns(l entydad.RoleLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "description", Label: l.Columns.Description, Sortable: false, Filterable: false},
		{Key: "color", Label: l.Columns.Color, Sortable: false, Width: "120px"},
		{Key: "permissions", Label: l.Columns.Permissions, Sortable: false, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: false, Width: "120px"},
		{Key: "date_created", Label: "Date Created", Sortable: true, Filterable: true, FilterType: types.FilterTypeDate},
	}
}

func buildTableRows(roles []*rolepb.Role, l entydad.RoleLabels, sl entydad.SharedLabels, routes entydad.RoleRoutes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
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
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
				Disabled:       !perms.Can("role", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
				Disabled:       !perms.Can("role", "update"), DisabledTooltip: sl.Badges.NoPermission,
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
			deleteAction.DisabledTooltip = sl.Errors.CannotDeleteInUse
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

func buildBulkActions(l entydad.RoleLabels, sl entydad.SharedLabels, common pyeza.CommonLabels, routes entydad.RoleRoutes) []types.BulkAction {
	return []types.BulkAction{
		{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-shield",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		},
		{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-shield-off",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		},
		{
			Key:              "delete",
			Label:            common.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         routes.BulkDeleteURL,
			ConfirmTitle:     common.Bulk.Delete,
			ConfirmMessage:   sl.Confirm.BulkDelete,
			RequiresDataAttr: "deletable",
		},
	}
}
