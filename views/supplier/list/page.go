package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	Routes          entydad.SupplierRoutes
	GetListPageData func(ctx context.Context, req *supplierpb.GetSupplierListPageDataRequest) (*supplierpb.GetSupplierListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	Labels          entydad.SupplierLabels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the supplier list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the supplier list view (full page).
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
				Title:          statusPageTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "suppliers",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-truck",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-list-content",
			Table:           tableConfig,
		}

		return view.OK("supplier-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
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

// buildTableConfig fetches supplier data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps, status string) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	resp, err := deps.GetListPageData(ctx, &supplierpb.GetSupplierListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list suppliers: %v", err)
		return nil, fmt.Errorf("failed to load suppliers: %w", err)
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

	l := deps.Labels
	columns := supplierColumns(l)
	rows := buildTableRows(resp.GetSupplierList(), status, l, deps.SharedLabels, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

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
		DefaultSortColumn:    "company_name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddNew,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("supplier", "create"),
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func supplierColumns(l entydad.SupplierLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "company_name", Label: l.Columns.CompanyName, Sortable: true},
		{Key: "supplier_type", Label: l.Columns.SupplierType, Sortable: true, Width: "130px"},
		{Key: "internal_id", Label: l.Columns.InternalID, Sortable: true, Width: "130px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
		{Key: "payment_terms", Label: l.Columns.PaymentTerms, Sortable: true, Width: "140px"},
		{Key: "contact_name", Label: l.Columns.ContactName, Sortable: true},
		{Key: "date_created", Label: l.Columns.DateCreated, Sortable: true, Width: "140px"},
	}
}

func buildTableRows(suppliers []*supplierpb.Supplier, status string, l entydad.SupplierLabels, sl entydad.SharedLabels, routes entydad.SupplierRoutes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, s := range suppliers {
		recordStatus := supplierStatus(s)
		if recordStatus != status {
			continue
		}

		id := s.GetId()
		companyName := s.GetCompanyName()
		supplierType := s.GetSupplierType()
		internalID := s.GetInternalId()
		paymentTerms := s.GetPaymentTerms()
		dateCreated := s.GetDateCreatedString()
		isInUse := inUseIDs[id]

		contactName := ""
		if u := s.GetUser(); u != nil {
			contactName = u.GetFirstName() + " " + u.GetLastName()
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: companyName},
				{Type: "text", Value: supplierType},
				{Type: "text", Value: internalID},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
				{Type: "text", Value: paymentTerms},
				{Type: "text", Value: contactName},
				{Type: "text", Value: dateCreated},
			},
			DataAttrs: map[string]string{
				"company_name": companyName,
				"status":       recordStatus,
				"deletable":    strconv.FormatBool(!isInUse),
			},
			Actions: buildRowActions(id, companyName, recordStatus, isInUse, l, sl, routes, perms),
		})
	}
	return rows
}

// supplierStatus returns the effective status string from a supplier record.
// Uses the explicit Status field if set, otherwise falls back to Active bool.
func supplierStatus(s *supplierpb.Supplier) string {
	if st := s.GetStatus(); st != "" {
		return st
	}
	if s.GetActive() {
		return "active"
	}
	return "blocked"
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

func buildRowActions(id, companyName, status string, isInUse bool, l entydad.SupplierLabels, sl entydad.SharedLabels, routes entydad.SupplierRoutes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit,
			Disabled: !perms.Can("supplier", "update"), DisabledTooltip: sl.Badges.NoPermission},
	}

	switch status {
	case "active":
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Actions.Block, Action: "block",
			URL: routes.SetStatusURL + "?status=blocked", ItemName: companyName,
			ConfirmTitle:   l.Actions.Block,
			ConfirmMessage: fmt.Sprintf(sl.Confirm.Block, companyName),
			Disabled: !perms.Can("supplier", "update"), DisabledTooltip: sl.Badges.NoPermission,
		})
	case "blocked":
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Actions.Activate, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: companyName,
			ConfirmTitle:   l.Actions.Activate,
			ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, companyName),
			Disabled: !perms.Can("supplier", "update"), DisabledTooltip: sl.Badges.NoPermission,
		})
	case "on_hold":
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Actions.Activate, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: companyName,
			ConfirmTitle:   l.Actions.Activate,
			ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, companyName),
			Disabled: !perms.Can("supplier", "update"), DisabledTooltip: sl.Badges.NoPermission,
		})
	}

	deleteAction := types.TableAction{
		Type:     "delete",
		Label:    l.Actions.Delete,
		Action:   "delete",
		URL:      routes.DeleteURL,
		ItemName: companyName,
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

func buildBulkActions(l entydad.SupplierLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, status string, routes entydad.SupplierRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "block",
			Label:           l.Actions.Block,
			Icon:            "icon-slash",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Block,
			ConfirmMessage:  sl.Confirm.BulkBlock,
			ExtraParamsJSON: `{"target_status":"blocked"}`,
		})
	case "blocked", "on_hold":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           cl.Bulk.Activate,
			Icon:            "icon-check-circle",
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
