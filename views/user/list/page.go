package list

import (
	"context"
	"fmt"
	"log"
	"math"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	userpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/user"

	"github.com/erniealice/entydad-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

var userSearchFields = []string{"first_name", "last_name", "email_address"}

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes               entydad.UserRoutes
	GetListPageData      func(ctx context.Context, req *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error)
	GetUserWorkspacesMap func(ctx context.Context) (map[string][]types.ChipData, error)
	RefreshURL           string
	Labels               entydad.UserLabels
	SharedLabels         entydad.SharedLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
}

// PageData holds the data for the user list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the user list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("user", "list") {
			return view.Forbidden("user:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := userColumns(deps.Labels)
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
				ActiveNav:      "user",
				ActiveSubNav:   "users-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "user-list-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "user"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("user-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("user", "list") {
			return view.Forbidden("user:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := userColumns(deps.Labels)
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

// buildTableConfig fetches user data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, userSearchFields)

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

	resp, err := deps.GetListPageData(ctx, &userpb.GetUserListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return nil, fmt.Errorf("failed to load users: %w", err)
	}

	// Fetch user-workspace mappings (best-effort; nil map means no workspace data)
	var userWorkspacesMap map[string][]types.ChipData
	if deps.GetUserWorkspacesMap != nil {
		userWorkspacesMap, err = deps.GetUserWorkspacesMap(ctx)
		if err != nil {
			log.Printf("Warning: Failed to load user workspaces map: %v", err)
		}
	}

	l := deps.Labels
	rows := buildTableRows(resp.GetUserList(), status, l, deps.SharedLabels, userWorkspacesMap, deps.Routes, perms)
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
		DefaultSortColumn:    "date_created",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddUser,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("user", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "user:create"),
		},
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func userColumns(l entydad.UserLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "first_name", Label: l.Columns.Name, MinWidth: "9.375rem"},
		{Key: "email_address", Label: l.Columns.Email, MinWidth: "11.25rem"},
		{Key: "workspaces", Label: l.Columns.Workspaces, NoSort: true, NoFilter: true, MinWidth: "7.5rem"},
		{Key: "date_created", Label: l.Columns.DateCreated, WidthClass: "col-6xl"},
		{Key: "status", Label: l.Columns.Status, NoFilter: true, WidthClass: "col-2xl"},
	}
}

func buildTableRows(users []*userpb.User, status string, l entydad.UserLabels, sl entydad.SharedLabels, userWorkspacesMap map[string][]types.ChipData, routes entydad.UserRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, u := range users {
		active := u.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := u.GetId()
		name := u.GetFirstName() + " " + u.GetLastName()
		email := u.GetEmailAddress()

		// Build workspace chips for this user
		workspaceChips := userWorkspacesMap[id] // nil-safe: returns nil slice for missing key
		workspacesCell := types.BuildChipCellFromChips(workspaceChips, 3)
		workspacesCell.TestID = "workspaces-chips"

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit,
				Disabled: !perms.Can("user", "update"), DisabledTooltip: sl.Badges.NoPermission},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
				Disabled:       !perms.Can("user", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
				Disabled:       !perms.Can("user", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		}
		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: email},
				workspacesCell,
				types.DateTimeCell(u.GetDateCreatedString(), types.DateReadable),
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"testid": "user-row-" + id,
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

func buildBulkActions(l entydad.UserLabels, sl entydad.SharedLabels, common pyeza.CommonLabels, status string, routes entydad.UserRoutes, perms *types.UserPermissions) []types.BulkAction {
	actions := []types.BulkAction{}

	canUpdate := perms.Can("user", "update")

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-user-minus",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
			Disabled:        !canUpdate,
			DisabledTooltip: fmt.Sprintf(common.Errors.MissingPermission, "user:update"),
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-user-check",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
			Disabled:        !canUpdate,
			DisabledTooltip: fmt.Sprintf(common.Errors.MissingPermission, "user:update"),
		})
	}

	return actions
}
