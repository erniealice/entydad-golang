package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
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
	Routes            entydad.ClientRoutes
	GetListPageData   func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
	GetInUseIDs       func(ctx context.Context, ids []string) (map[string]bool, error)
	GetClientBalances func(ctx context.Context) (map[string]int64, error)
	Labels            entydad.ClientLabels
	SharedLabels      entydad.SharedLabels
	CommonLabels      pyeza.CommonLabels
	TableLabels       types.TableLabels
}

// PageData holds the data for the client list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var clientSearchFields = []string{"name", "u.first_name", "u.last_name", "u.email_address"}

// NewView creates the client list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := clientColumns(deps.Labels)
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

		// KB help content — status-specific slug with generic fallback.
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				slug := "clients-list-" + status
				kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, slug)
				if kb == nil {
					kb, _ = provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "clients-list")
				}
				if kb != nil {
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

		columns := clientColumns(deps.Labels)
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

// buildTableConfig fetches client data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, clientSearchFields)

	// Inject status filter for server-side pagination. Client lifecycle now
	// has 5 states (prospect/active/on_hold/blocked/inactive). The legacy
	// `active` boolean is kept in sync by SetStatus closure but `status` is
	// the source of truth for filter equality.
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "status",
		FilterType: &commonpb.TypedFilter_StringFilter{
			StringFilter: &commonpb.StringFilter{
				Value:    status,
				Operator: commonpb.StringOperator_STRING_EQUALS,
			},
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
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
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
	// Status column omitted on purpose — the list page is already scoped
	// by /list/{status}, so a per-row badge would be redundant.
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.ClientName},
		{Key: "representative", Label: l.Columns.Representative},
		{Key: "category", Label: l.Columns.Category, NoFilter: true, WidthClass: "col-7xl"},
		{Key: "payment_term", Label: l.Columns.PaymentTerm, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "outstanding_balance", Label: "Outstanding", NoSort: true, NoFilter: true, Align: "right", WidthClass: "col-4xl"},
	}
}

// buildTableRows builds row data for a status-filtered list. listStatus is
// the lifecycle filter the page is currently rendering (e.g. "blocked"); row
// actions key off that, not the proto field, so transitions stay correct even
// when individual rows have stale/unmigrated status values. The badge cell
// still reflects each row's own recordStatus.
func buildTableRows(clients []*clientpb.Client, listStatus string, l entydad.ClientLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.ClientRoutes, inUseIDs map[string]bool, balances map[string]int64, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, c := range clients {
		recordStatus := clientStatus(c)

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
			Actions: buildRowActions(id, displayName, listStatus, isInUse, l, sl, cl, routes, perms),
		})
	}
	return rows
}

// clientStatus returns the effective lifecycle status for a client. Falls back
// to the legacy `active` boolean when the proto status field is unset (for
// rows created before the status column existed and not yet backfilled).
func clientStatus(c *clientpb.Client) string {
	if st := c.GetStatus(); st != "" {
		return st
	}
	if c.GetActive() {
		return "active"
	}
	return "inactive"
}

func statusPageTitle(l entydad.ClientLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "prospect":
		return l.Page.HeadingProspect
	case "on_hold":
		return l.Page.HeadingOnHold
	case "blocked":
		return l.Page.HeadingBlocked
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
	case "on_hold":
		return l.Page.CaptionOnHold
	case "blocked":
		return l.Page.CaptionBlocked
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
	case "on_hold":
		return l.Empty.OnHoldTitle
	case "blocked":
		return l.Empty.BlockedTitle
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
	case "on_hold":
		return l.Empty.OnHoldMessage
	case "blocked":
		return l.Empty.BlockedMessage
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
	case "prospect":
		return "info"
	case "on_hold":
		return "warning"
	case "blocked":
		return "danger"
	case "inactive":
		return "default"
	default:
		return "default"
	}
}

// statusLabel resolves the user-facing badge text for a client lifecycle
// status through the centralized pyeza.StatusLabels (CommonLabels.Status).
func statusLabel(cl pyeza.CommonLabels, status string) string {
	switch status {
	case "active":
		return cl.Status.Active
	case "prospect":
		return cl.Status.Prospect
	case "on_hold":
		return cl.Status.OnHold
	case "blocked":
		return cl.Status.Blocked
	case "inactive":
		return cl.Status.Inactive
	default:
		return status
	}
}

// statusTransition describes a possible move to a target lifecycle status.
// Order in clientStatusTransitions controls the order in which actions appear
// in row + bulk menus.
type statusTransition struct {
	target  string // proto status value: prospect/active/on_hold/blocked/inactive
	iconKey string // table.html action Type (drives the icon)
	action  string // table.html data-{action}-url switch (activate or deactivate)
	variant string // bulk action variant (primary/warning/danger/default)
}

// clientStatusTransitions enumerates the canonical transition order. Row
// actions filter out the source status, so a row in the "active" filter shows
// transitions to the other 4 states.
var clientStatusTransitions = []statusTransition{
	{target: "active", iconKey: "activate", action: "activate", variant: "primary"},
	{target: "prospect", iconKey: "prospect", action: "deactivate", variant: "info"},
	{target: "on_hold", iconKey: "hold", action: "deactivate", variant: "warning"},
	{target: "blocked", iconKey: "block", action: "deactivate", variant: "danger"},
	{target: "inactive", iconKey: "deactivate", action: "deactivate", variant: "default"},
}

// transitionLabels returns (rowLabel, bulkLabel, confirmRow, confirmBulk) for
// a given target status, sourced from typed labels so every string is
// translatable.
func transitionLabels(target string, l entydad.ClientLabels, sl entydad.SharedLabels) (string, string, string, string) {
	switch target {
	case "active":
		return l.Detail.Actions.ActivateClient, l.BulkActions.SetAsActive, sl.Confirm.Activate, sl.Confirm.BulkActivate
	case "prospect":
		return l.Detail.Actions.SetProspect, l.BulkActions.SetAsProspect, sl.Confirm.Prospect, sl.Confirm.BulkProspect
	case "on_hold":
		return l.Detail.Actions.HoldClient, l.BulkActions.SetAsOnHold, sl.Confirm.Hold, sl.Confirm.BulkHold
	case "blocked":
		return l.Detail.Actions.BlockClient, l.BulkActions.SetAsBlocked, sl.Confirm.Block, sl.Confirm.BulkBlock
	case "inactive":
		return l.Detail.Actions.DeactivateClient, l.BulkActions.SetAsInactive, sl.Confirm.Deactivate, sl.Confirm.BulkDeactivate
	}
	return "", "", "", ""
}

func buildRowActions(id, name, status string, isInUse bool, l entydad.ClientLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.ClientRoutes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Detail.Actions.ViewClient, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Detail.Actions.EditClient, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Detail.Actions.EditClient,
			Disabled: !perms.Can("client", "update"), DisabledTooltip: sl.Badges.NoPermission},
	}

	// Clone is only meaningful for healthy (active) clients.
	if status == "active" {
		actions = append(actions, types.TableAction{
			Type:            "clone",
			Label:           cl.Actions.Clone,
			Action:          "clone",
			URL:             route.ResolveURL(routes.EditURL, "id", id),
			DrawerTitle:     cl.Actions.Clone,
			Disabled:        !perms.Can("client", "create"),
			DisabledTooltip: sl.Badges.NoPermission,
		})
	}

	canUpdate := perms.Can("client", "update")
	tooltip := sl.Badges.NoPermission

	// Cross-status transitions: every status filter exposes moves to all 4
	// other lifecycle states, so users can always reach any status from any
	// list without round-tripping through detail.
	for _, tr := range clientStatusTransitions {
		if tr.target == status {
			continue
		}
		rowLabel, _, confirmRow, _ := transitionLabels(tr.target, l, sl)
		actions = append(actions, types.TableAction{
			Type: tr.iconKey, Label: rowLabel, Action: tr.action,
			URL: routes.SetStatusURL + "?status=" + tr.target, ItemName: name,
			ConfirmTitle:   rowLabel,
			ConfirmMessage: fmt.Sprintf(confirmRow, name),
			Disabled:       !canUpdate, DisabledTooltip: tooltip,
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

// bulkActionIcon maps a transition iconKey to the bulk-action icon name.
// Bulk action icons live outside table.html (no shared mapping), so each
// status keeps an explicit icon-* string.
func bulkActionIcon(iconKey string) string {
	switch iconKey {
	case "activate":
		return "icon-user-check"
	case "prospect":
		return "icon-user-plus"
	case "hold":
		return "icon-pause-circle"
	case "block":
		return "icon-x-circle"
	case "deactivate":
		return "icon-user-minus"
	}
	return "icon-edit"
}

func buildBulkActions(l entydad.ClientLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, status string, routes entydad.ClientRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	for _, tr := range clientStatusTransitions {
		if tr.target == status {
			continue
		}
		_, bulkLabel, _, confirmBulk := transitionLabels(tr.target, l, sl)
		actions = append(actions, types.BulkAction{
			Key:             tr.target,
			Label:           bulkLabel,
			Icon:            bulkActionIcon(tr.iconKey),
			Variant:         tr.variant,
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    bulkLabel,
			ConfirmMessage:  confirmBulk,
			ExtraParamsJSON: `{"target_status":"` + tr.target + `"}`,
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
