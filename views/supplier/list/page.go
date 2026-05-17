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

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"

	"github.com/erniealice/entydad-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes              entydad.SupplierRoutes
	GetListPageData     func(ctx context.Context, req *supplierpb.GetSupplierListPageDataRequest) (*supplierpb.GetSupplierListPageDataResponse, error)
	GetInUseIDs         func(ctx context.Context, ids []string) (map[string]bool, error)
	GetSupplierBalances func(ctx context.Context) (map[string]int64, error)
	Labels              entydad.SupplierLabels
	SharedLabels        entydad.SharedLabels
	CommonLabels        pyeza.CommonLabels
	TableLabels         types.TableLabels
}

// PageData holds the data for the supplier list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var supplierSearchFields = []string{"name", "internal_id", "u.first_name", "u.last_name", "u.email_address"}

// NewView creates the supplier list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("supplier", "list") {
			return view.Forbidden("supplier:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := supplierColumns(deps.Labels)
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
				ActiveNav:      "supplier",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-truck",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-list-content",
			Table:           tableConfig,
		}

		// KB help content — status-specific slug with generic fallback.
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				slug := "suppliers-list-" + status
				kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, slug)
				if kb == nil {
					kb, _ = provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "suppliers-list")
				}
				if kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("supplier-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("supplier", "list") {
			return view.Forbidden("supplier:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := supplierColumns(deps.Labels)
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

// buildTableConfig fetches supplier data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, supplierSearchFields)

	// Inject status filter for server-side pagination
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "s.status",
		FilterType: &commonpb.TypedFilter_StringFilter{
			StringFilter: &commonpb.StringFilter{
				Value:    status,
				Operator: commonpb.StringOperator_STRING_EQUALS,
			},
		},
	})
	// Exclude soft-deleted suppliers from every status list.
	// Supplier uses both `status` (active/blocked/on_hold) and `active` (bool).
	// DeleteSupplier flips `active` to false but leaves `status` intact, so a
	// status-only filter still surfaces deleted rows.
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "s.active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: true},
		},
	})

	var resp *supplierpb.GetSupplierListPageDataResponse
	if deps.GetListPageData != nil {
		var err error
		resp, err = deps.GetListPageData(ctx, &supplierpb.GetSupplierListPageDataRequest{
			Search:     listParams.Search,
			Filters:    listParams.Filters,
			Sort:       listParams.Sort,
			Pagination: listParams.Pagination,
		})
		if err != nil {
			log.Printf("Failed to list suppliers: %v", err)
			return nil, fmt.Errorf("failed to load suppliers: %w", err)
		}
	}

	// Check which items are in use
	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range resp.GetSupplierList() {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	// Fetch outstanding balances for all suppliers
	var supplierBalances map[string]int64
	if deps.GetSupplierBalances != nil {
		supplierBalances, _ = deps.GetSupplierBalances(ctx)
	}

	l := deps.Labels
	rows := buildTableRows(resp.GetSupplierList(), status, l, deps.SharedLabels, deps.CommonLabels, deps.Routes, inUseIDs, supplierBalances, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
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

	var primaryAction *types.PrimaryAction
	if status == "active" {
		primaryAction = &types.PrimaryAction{
			Label:           l.Buttons.AddNew,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("supplier", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "supplier:create"),
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "suppliers-table",
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

func supplierColumns(l entydad.SupplierLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "supplier_type", Label: l.Columns.SupplierType, WidthClass: "col-3xl"},
		{Key: "internal_id", Label: l.Columns.InternalID, WidthClass: "col-3xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		{Key: "category", Label: l.Columns.Category, WidthClass: "col-7xl"},
		{Key: "payment_terms", Label: l.Columns.PaymentTerms, WidthClass: "col-3xl"},
		{Key: "contact_name", Label: l.Columns.ContactName},
		{Key: "outstanding_balance", Label: "Outstanding", Align: "right", WidthClass: "col-4xl"},
		{Key: "date_created", Label: l.Columns.DateCreated, WidthClass: "col-3xl"},
	}
}

func buildTableRows(suppliers []*supplierpb.Supplier, status string, l entydad.SupplierLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.SupplierRoutes, inUseIDs map[string]bool, balances map[string]int64, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, s := range suppliers {
		recordStatus := supplierStatus(s)

		id := s.GetId()
		name := s.GetName()
		supplierType := s.GetSupplierType()
		internalID := s.GetInternalId()
		paymentTerms := s.GetPaymentTerms()
		dateCreated := s.GetDateCreatedString()
		isInUse := inUseIDs[id]

		contactName := ""
		if u := s.GetUser(); u != nil {
			contactName = u.GetFirstName() + " " + u.GetLastName()
		}

		balanceCell := types.TableCell{Type: "text", Value: "—"}
		if balance, ok := balances[id]; ok && balance != 0 {
			balanceCell = types.MoneyCell(float64(balance), "", true)
		}

		var catLabels []string
		for _, sc := range s.GetCategories() {
			if cat := sc.GetCategory(); cat != nil {
				if n := cat.GetName(); n != "" {
					catLabels = append(catLabels, n)
				}
			}
		}
		categoryCell := types.BuildChipCellFromLabels(catLabels, 3)

		// Render payment_terms as a badge (consistent with client list).
		// Suppliers populate s.PaymentTerms via SQL JOIN with payment_term.name
		// in the postgres adapter, so the value is already the human-readable
		// term name when present.
		ptCell := types.TableCell{Type: "text", Value: ""}
		if paymentTerms != "" {
			ptCell = types.TableCell{Type: "badge", Value: paymentTerms, Variant: "default"}
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: supplierType},
				{Type: "text", Value: internalID},
				{Type: "badge", Value: statusLabel(cl, recordStatus), Variant: statusVariant(recordStatus)},
				categoryCell,
				ptCell,
				{Type: "text", Value: contactName},
				balanceCell,
				types.DateTimeCell(dateCreated, types.DateReadable),
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			// Row actions key off the page's list filter (not per-row
			// recordStatus) so transitions stay correct even when rows have
			// stale/unmigrated status values.
			Actions: buildRowActions(id, name, status, isInUse, l, sl, cl, routes, perms),
		})
	}
	return rows
}

// supplierStatus returns the effective status string from a supplier record.
// Uses the explicit Status field if set, otherwise falls back to Active bool.
func supplierStatus(s *supplierpb.Supplier) string {
	// Check active flag first — a soft-deleted supplier is always blocked
	// regardless of its status field
	if !s.GetActive() {
		return "blocked"
	}
	if st := s.GetStatus(); st != "" {
		return st
	}
	return "active"
}

func statusPageTitle(l entydad.SupplierLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "blocked":
		return l.Page.HeadingBlocked
	case "on_hold":
		return l.Page.HeadingOnHold
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l entydad.SupplierLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "blocked":
		return l.Page.CaptionBlocked
	case "on_hold":
		return l.Page.CaptionOnHold
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.SupplierLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "blocked":
		return l.Empty.BlockedTitle
	case "on_hold":
		return l.Empty.OnHoldTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.SupplierLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "blocked":
		return l.Empty.BlockedMessage
	case "on_hold":
		return l.Empty.OnHoldMessage
	default:
		return l.Empty.ActiveMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "blocked":
		return "danger"
	case "on_hold":
		return "warning"
	default:
		return "default"
	}
}

func statusLabel(cl pyeza.CommonLabels, status string) string {
	switch status {
	case "active":
		return cl.Status.Active
	case "inactive":
		return cl.Status.Inactive
	case "blocked":
		return cl.Status.Blocked
	case "on_hold":
		return cl.Status.OnHold
	default:
		return status
	}
}

// supplierStatusTransition describes a possible move to a target supplier
// lifecycle status (active/blocked/on_hold).
type supplierStatusTransition struct {
	target  string
	iconKey string
	action  string
	variant string
}

var supplierStatusTransitions = []supplierStatusTransition{
	{target: "active", iconKey: "activate", action: "activate", variant: "primary"},
	{target: "on_hold", iconKey: "hold", action: "deactivate", variant: "warning"},
	{target: "blocked", iconKey: "block", action: "deactivate", variant: "danger"},
}

func supplierTransitionLabels(target string, l entydad.SupplierLabels, sl entydad.SharedLabels) (string, string, string) {
	switch target {
	case "active":
		return l.Actions.Activate, sl.Confirm.Activate, sl.Confirm.BulkActivate
	case "on_hold":
		return l.Actions.SetOnHold, sl.Confirm.Hold, sl.Confirm.BulkHold
	case "blocked":
		return l.Actions.Block, sl.Confirm.Block, sl.Confirm.BulkBlock
	}
	return "", "", ""
}

func supplierBulkActionIcon(iconKey string) string {
	switch iconKey {
	case "activate":
		return "icon-check-circle"
	case "hold":
		return "icon-pause-circle"
	case "block":
		return "icon-x-circle"
	}
	return "icon-edit"
}

func buildRowActions(id, name, status string, isInUse bool, l entydad.SupplierLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.SupplierRoutes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit,
			Disabled: !perms.Can("supplier", "update"), DisabledTooltip: sl.Badges.NoPermission},
	}

	if status == "active" {
		actions = append(actions, types.TableAction{
			Type:            "clone",
			Label:           cl.Actions.Clone,
			Action:          "clone",
			URL:             route.ResolveURL(routes.EditURL, "id", id),
			DrawerTitle:     cl.Actions.Clone,
			Disabled:        !perms.Can("supplier", "create"),
			DisabledTooltip: sl.Badges.NoPermission,
		})
	}

	canUpdate := perms.Can("supplier", "update")
	tooltip := sl.Badges.NoPermission

	// Cross-status transitions: every status filter exposes moves to all
	// other supplier lifecycle states (active/on_hold/blocked). Overflow:true
	// collapses these into the row's ⋮ menu so the inline action bar stays
	// compact (view / edit / clone / delete only).
	for _, tr := range supplierStatusTransitions {
		if tr.target == status {
			continue
		}
		rowLabel, confirmRow, _ := supplierTransitionLabels(tr.target, l, sl)
		actions = append(actions, types.TableAction{
			Type: tr.iconKey, Label: rowLabel, Action: tr.action,
			URL: routes.SetStatusURL + "?status=" + tr.target, ItemName: name,
			ConfirmTitle:   rowLabel,
			ConfirmMessage: fmt.Sprintf(confirmRow, name),
			Disabled:       !canUpdate, DisabledTooltip: tooltip,
			Overflow: true,
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
	} else if !perms.Can("supplier", "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = sl.Badges.NoPermission
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l entydad.SupplierLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, status string, routes entydad.SupplierRoutes, perms *types.UserPermissions) []types.BulkAction {
	actions := []types.BulkAction{}

	canUpdate := perms.Can("supplier", "update")
	canDelete := perms.Can("supplier", "delete")

	for _, tr := range supplierStatusTransitions {
		if tr.target == status {
			continue
		}
		bulkLabel, _, confirmBulk := supplierTransitionLabels(tr.target, l, sl)
		actions = append(actions, types.BulkAction{
			Key:             tr.target,
			Label:           bulkLabel,
			Icon:            supplierBulkActionIcon(tr.iconKey),
			Variant:         tr.variant,
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    bulkLabel,
			ConfirmMessage:  confirmBulk,
			ExtraParamsJSON: `{"target_status":"` + tr.target + `"}`,
			Disabled:        !canUpdate,
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "supplier:update"),
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
		Disabled:         !canDelete,
		DisabledTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "supplier:delete"),
	})

	return actions
}
