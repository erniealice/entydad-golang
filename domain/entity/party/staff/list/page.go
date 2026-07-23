// Package list renders the staff list page + table-only HTMX refresh.
// Mirrors party/delegate/list/page.go: the {status} URL segment drives the
// `active` BooleanFilter (active/inactive), independent of the entity's own
// AVAILABILITY `status` field (available|assigned|bench|offboarded) which is
// a form field/column only — not a list filter — in this minimal surface.
package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	staffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/staff"

	"github.com/erniealice/entydad-golang"
	entitystaff "github.com/erniealice/entydad-golang/domain/entity/party/staff"
)

// ListViewDeps holds view dependencies for the staff list.
// No GetInUseIDs (no ref-checker for Staff), no balance/subscription counts.
type ListViewDeps struct {
	Routes          entitystaff.Routes
	GetListPageData func(ctx context.Context, req *staffpb.GetStaffListPageDataRequest) (*staffpb.GetStaffListPageDataResponse, error)
	Labels          entitystaff.Labels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the staff list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var staffSearchFields = []string{"u.first_name", "u.last_name", "u.email_address"}

// NewView creates the staff list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("staff", "list") {
			return view.Forbidden("staff:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := staffColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "user_first_name", "asc")
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
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "staff",
				ActiveSubNav:   status,
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "staff-list-content",
			Table:           tableConfig,
		}

		return view.OK("staff-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("staff", "list") {
			return view.Forbidden("staff:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := staffColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "user_first_name", "asc")
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

// buildTableConfig fetches staff data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, staffSearchFields)

	// Inject active boolean filter from the status URL segment.
	// "active" → active=true, "inactive" → active=false.
	activeVal := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeVal},
		},
	})

	resp, err := deps.GetListPageData(ctx, &staffpb.GetStaffListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list staff: %v", err)
		return nil, fmt.Errorf("failed to load staff: %w", err)
	}

	l := deps.Labels
	rows := buildTableRows(resp.GetStaffList(), l, deps.CommonLabels, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(deps.CommonLabels, deps.Routes, perms)

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

	var primaryAction *types.PrimaryAction
	if status == "active" {
		primaryAction = &types.PrimaryAction{
			Label:           l.Actions.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("staff", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "staff:create"),
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "staffs-table",
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
			Title:   l.Page.Heading,
			Message: l.Page.Caption,
		},
		PrimaryAction:    primaryAction,
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func staffColumns(l entitystaff.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.List.Columns.Name, SortKey: "user_first_name"},
		{Key: "email", Label: l.List.Columns.Email, NoSort: true},
		{Key: "employment_type", Label: l.List.Columns.EmploymentType, NoSort: true},
		{Key: "status", Label: l.List.Columns.Status, NoSort: true},
	}
}

// buildTableRows builds row data for the staff list.
// Name/email come from the embedded User (populated by the GetStaffListPageData
// CTE join).
func buildTableRows(staffs []*staffpb.Staff, l entitystaff.Labels, cl pyeza.CommonLabels, routes entitystaff.Routes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, s := range staffs {
		id := s.GetId()
		u := s.GetUser()
		name := ""
		email := ""
		if u != nil {
			name = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			email = u.GetEmailAddress()
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: email},
				{Type: "text", Value: s.GetEmploymentType()},
				{Type: "text", Value: s.GetStatus()},
			},
			DataAttrs: map[string]string{
				"name":  name,
				"email": email,
			},
			Actions: buildRowActions(id, name, l, cl, routes, perms),
		})
	}
	return rows
}

func buildRowActions(id, name string, l entitystaff.Labels, cl pyeza.CommonLabels, routes entitystaff.Routes, perms *types.UserPermissions) []types.TableAction {
	canUpdate := perms.Can("staff", "update")
	canDelete := perms.Can("staff", "delete")

	return []types.TableAction{
		{
			Type:            "edit",
			Label:           l.Actions.Edit,
			Action:          "edit",
			URL:             route.ResolveURL(routes.EditURL, "id", id),
			DrawerTitle:     l.Actions.Edit,
			Disabled:        !canUpdate,
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "staff:update"),
		},
		{
			Type:            "delete",
			Label:           l.Actions.Delete,
			Action:          "delete",
			URL:             routes.DeleteURL,
			ItemName:        name,
			Disabled:        !canDelete,
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "staff:delete"),
		},
	}
}

func buildBulkActions(cl pyeza.CommonLabels, routes entitystaff.Routes, perms *types.UserPermissions) []types.BulkAction {
	canDelete := perms.Can("staff", "delete")
	return []types.BulkAction{
		{
			Key:             "delete",
			Label:           cl.Bulk.Delete,
			Icon:            "icon-trash-2",
			Variant:         "danger",
			Endpoint:        routes.BulkDeleteURL,
			ConfirmTitle:    cl.Bulk.Delete,
			Disabled:        !canDelete,
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "staff:delete"),
		},
	}
}
