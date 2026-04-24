package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"

	"github.com/erniealice/entydad-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes             entydad.ClientRoutes
	GetListPageData    func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
	GetInUseIDs        func(ctx context.Context, ids []string) (map[string]bool, error)
	GetClientBalances  func(ctx context.Context) (map[string]int64, error)
	Labels             entydad.ClientLabels
	SharedLabels       entydad.SharedLabels
	CommonLabels       pyeza.CommonLabels
	TableLabels        types.TableLabels
}

// PageData holds the data for the client list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var clientAllowedSortCols = []string{
	"date_created", "date_modified", "name",
	"u.first_name", "u.last_name", "u.email_address",
}

var clientSearchFields = []string{"name", "u.first_name", "u.last_name", "u.email_address"}

// NewView creates the client list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		p, err := espynahttp.ParseTableParams(viewCtx.Request, clientAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "client",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "client-list-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "client"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("client-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		p, err := espynahttp.ParseTableParams(viewCtx.Request, clientAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches client data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, p espynahttp.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, clientSearchFields)

	// Inject status filter for server-side pagination
	activeValue := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "c.active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	})

	resp, err := deps.GetListPageData(ctx, &clientpb.GetClientListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list clients: %v", err)
		return nil, fmt.Errorf("failed to load clients: %w", err)
	}

	// Check which items are in use
	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range resp.GetClientList() {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	var clientBalances map[string]int64
	if deps.GetClientBalances != nil {
		clientBalances, _ = deps.GetClientBalances(ctx)
	}

	l := deps.Labels
	columns := clientColumns(l)
	rows := buildTableRows(resp.GetClientList(), status, l, deps.SharedLabels, deps.CommonLabels, deps.Routes, inUseIDs, clientBalances, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, status, deps.Routes)

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
			Label:           l.Buttons.AddNew,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("client", "create"),
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "clients-table",
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
		PrimaryAction:    primaryAction,
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func clientColumns(l entydad.ClientLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.ClientName, Sortable: true, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "representative", Label: l.Columns.Representative, Sortable: true, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "category", Label: l.Columns.Category, WidthClass: "col-7xl"},
		{Key: "payment_term", Label: l.Columns.PaymentTerm, WidthClass: "col-3xl"},
		{Key: "outstanding_balance", Label: "Outstanding", Sortable: false, Align: "right", WidthClass: "col-4xl"},
	}
}

func buildTableRows(clients []*clientpb.Client, status string, l entydad.ClientLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.ClientRoutes, inUseIDs map[string]bool, balances map[string]int64, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, c := range clients {
		active := c.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := c.GetId()
		u := c.GetUser()
		name := c.GetName()
		repName := ""
		repEmail := ""
		if u != nil {
			first := u.GetFirstName()
			last := u.GetLastName()
			repName = strings.TrimSpace(first + " " + last)
			repEmail = u.GetEmailAddress()
		}
		displayName := name
		if displayName == "" {
			displayName = repName
		}
		isInUse := inUseIDs[id]

		// Build category chip cell
		var catLabels []string
		for _, cc := range c.GetCategories() {
			if n := cc.GetCategory().GetName(); n != "" {
				catLabels = append(catLabels, n)
			}
		}
		categoryCell := types.BuildChipCellFromLabels(catLabels, 3)

		// Build payment term badge cell
		ptCell := types.TableCell{Type: "text", Value: ""}
		if pt := c.GetPaymentTerm(); pt != nil && pt.GetName() != "" {
			ptCell = types.TableCell{Type: "badge", Value: pt.GetName(), Variant: "default"}
		}

		// Build outstanding balance cell
		balanceCell := types.TableCell{Type: "text", Value: "—"}
		if balance, ok := balances[id]; ok && balance != 0 {
			balanceCell = types.MoneyCell(float64(balance), "", true)
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: displayName},
				{Type: "single-person", Person: &types.PersonData{Name: repName, Email: repEmail}},
				categoryCell,
				ptCell,
				balanceCell,
			},
			DataAttrs: map[string]string{
				"name":      displayName,
				"email":     repEmail,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: buildRowActions(id, displayName, active, isInUse, l, sl, cl, routes, perms),
		})
	}
	return rows
}

func statusPageTitle(l entydad.ClientLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "prospect":
		return l.Page.HeadingProspect
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l entydad.ClientLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "prospect":
		return l.Page.CaptionProspect
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.ClientLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "prospect":
		return l.Empty.ProspectTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.ClientLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "prospect":
		return l.Empty.ProspectMessage
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

func buildRowActions(id, name string, active, isInUse bool, l entydad.ClientLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.ClientRoutes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Detail.Actions.ViewClient, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Detail.Actions.EditClient, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Detail.Actions.EditClient,
			Disabled: !perms.Can("client", "update"), DisabledTooltip: sl.Badges.NoPermission},
	}
	if active {
		actions = append(actions, types.TableAction{
			Type:            "clone",
			Label:           cl.Actions.Clone,
			Action:          "clone",
			URL:             route.ResolveURL(routes.EditURL, "id", id),
			DrawerTitle:     cl.Actions.Clone,
			Disabled:        !perms.Can("client", "create"),
			DisabledTooltip: sl.Badges.NoPermission,
		})
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Detail.Actions.DeactivateClient, Action: "deactivate",
			URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
			ConfirmTitle:   l.Detail.Actions.DeactivateClient,
			ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
			Disabled:       !perms.Can("client", "update"), DisabledTooltip: sl.Badges.NoPermission,
		})
	} else {
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Detail.Actions.ActivateClient, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: name,
			ConfirmTitle:   l.Detail.Actions.ActivateClient,
			ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
			Disabled:       !perms.Can("client", "update"), DisabledTooltip: sl.Badges.NoPermission,
		})
	}
	deleteAction := types.TableAction{
		Type:     "delete",
		Label:    l.Detail.Actions.DeleteClient,
		Action:   "delete",
		URL:      routes.DeleteURL,
		ItemName: name,
	}
	if isInUse {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = sl.Errors.CannotDeleteInUse
	} else if !perms.Can("client", "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = sl.Badges.NoPermission
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l entydad.ClientLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, status string, routes entydad.ClientRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.BulkActions.SetAsInactive,
			Icon:            "icon-user-minus",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.BulkActions.SetAsInactive,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           cl.Bulk.Activate,
			Icon:            "icon-user-check",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    cl.Bulk.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:              "delete",
		Label:            cl.Bulk.Delete,
		Icon:             "icon-trash-2",
		Variant:          "danger",
		Endpoint:         routes.BulkDeleteURL,
		ConfirmTitle:     cl.Bulk.Delete,
		ConfirmMessage:   sl.Confirm.BulkDelete,
		RequiresDataAttr: "deletable",
	})

	return actions
}
