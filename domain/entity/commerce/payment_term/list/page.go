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

	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"

	"github.com/erniealice/entydad-golang"
	paymentterm "github.com/erniealice/entydad-golang/domain/entity/commerce/payment_term"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// paymentTermTypeLabels maps payment term type codes to human-readable labels.
var paymentTermTypeLabels = map[string]string{
	"DUE_ON_RECEIPT": "Due on Receipt",
	"NET":            "Net",
	"COD":            "Cash on Delivery",
	"PROXIMATE":      "Proximate",
}

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *paymenttermpb.GetPaymentTermListPageDataRequest) (*paymenttermpb.GetPaymentTermListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	RefreshURL      string
	Routes          paymentterm.Routes
	Labels          paymentterm.Labels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
	// Scope filters which payment terms are shown. Valid values: "client", "supplier".
	// When set, only terms with entity_scope == Scope or entity_scope == "both" are shown.
	// Leave empty to show all terms (no filtering).
	Scope string
}

// PageData holds the data for the payment term list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the payment term list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("payment_term", "list") {
			return view.Forbidden("payment_term:list")
		}

		tableConfig, err := buildTableConfig(ctx, deps)
		if err != nil {
			return view.Error(err)
		}

		activeNav := "client"
		if deps.Scope == "supplier" {
			activeNav = "supplier"
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      activeNav,
				ActiveSubNav:   "payment-terms",
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Subtitle,
				HeaderIcon:     "icon-clock",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "payment-term-list-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "payment-term"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("payment-term-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("payment_term", "list") {
			return view.Forbidden("payment_term:list")
		}

		tableConfig, err := buildTableConfig(ctx, deps)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches payment term data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	var resp *paymenttermpb.GetPaymentTermListPageDataResponse
	if deps.GetListPageData != nil {
		var err error
		resp, err = deps.GetListPageData(ctx, &paymenttermpb.GetPaymentTermListPageDataRequest{})
		if err != nil {
			log.Printf("Failed to list payment terms: %v", err)
			return nil, fmt.Errorf("failed to load payment terms: %w", err)
		}
	}

	// Filter by entity scope if a scope is specified.
	allTerms := resp.GetPaymentTermList()
	filteredTerms := allTerms
	if deps.Scope != "" {
		filteredTerms = make([]*paymenttermpb.PaymentTerm, 0, len(allTerms))
		for _, pt := range allTerms {
			scope := pt.GetEntityScope()
			if scope == deps.Scope || scope == "both" {
				filteredTerms = append(filteredTerms, pt)
			}
		}
	}

	// Check which items are in use
	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range filteredTerms {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	l := deps.Labels
	columns := paymentTermColumns(l)
	rows := buildTableRows(filteredTerms, l, deps.SharedLabels, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(deps.Labels, deps.SharedLabels, deps.CommonLabels, deps.Routes, perms)

	refreshURL := deps.RefreshURL
	if refreshURL == "" {
		refreshURL = route.ResolveURL(deps.Routes.TableURL)
	}

	tableConfig := &types.TableConfig{
		ID:                   "payment-terms-table",
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
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddPaymentTerm,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("payment_term", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "payment_term:create"),
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func paymentTermColumns(l paymentterm.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, WidthClass: "col-7xl"},
		{Key: "code", Label: l.Columns.Code, WidthClass: "col-2xl"},
		{Key: "type", Label: l.Columns.Type, WidthClass: "col-5xl"},
		{Key: "net_days", Label: l.Columns.NetDays, WidthClass: "col-md"},
		{Key: "is_default", Label: l.Columns.IsDefault, WidthClass: "col-md"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-lg"},
	}
}

func paymentTermTypeLabel(code string) string {
	if label, ok := paymentTermTypeLabels[code]; ok {
		return label
	}
	return code
}

func buildTableRows(items []*paymenttermpb.PaymentTerm, l paymentterm.Labels, sl entydad.SharedLabels, routes paymentterm.Routes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, pt := range items {
		id := pt.GetId()
		name := pt.GetName()
		code := pt.GetCode()
		termType := pt.GetType()
		netDays := strconv.FormatInt(int64(pt.GetNetDays()), 10)
		isDefault := "No"
		isDefaultVariant := "default"
		if pt.GetIsDefault() {
			isDefault = "Yes"
			isDefaultVariant = "success"
		}
		active := pt.GetActive()
		activeStatus := "Inactive"
		activeVariant := "warning"
		if active {
			activeStatus = "Active"
			activeVariant = "success"
		}
		recordStatus := "inactive"
		if active {
			recordStatus = "active"
		}

		actions := []types.TableAction{
			{
				Type:            "edit",
				Label:           l.Actions.Edit,
				Action:          "edit",
				URL:             route.ResolveURL(routes.EditURL, "id", id),
				DrawerTitle:     l.Actions.Edit,
				Disabled:        !perms.Can("payment_term", "update"),
				DisabledTooltip: sl.Badges.NoPermission,
			},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type:            "deactivate",
				Label:           l.Actions.Deactivate,
				Action:          "deactivate",
				URL:             routes.SetStatusURL + "?status=inactive",
				ItemName:        name,
				ConfirmTitle:    l.Actions.Deactivate,
				ConfirmMessage:  fmt.Sprintf(sl.Confirm.Deactivate, name),
				Disabled:        !perms.Can("payment_term", "update"),
				DisabledTooltip: sl.Badges.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type:            "activate",
				Label:           l.Actions.Activate,
				Action:          "activate",
				URL:             routes.SetStatusURL + "?status=active",
				ItemName:        name,
				ConfirmTitle:    l.Actions.Activate,
				ConfirmMessage:  fmt.Sprintf(sl.Confirm.Activate, name),
				Disabled:        !perms.Can("payment_term", "update"),
				DisabledTooltip: sl.Badges.NoPermission,
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
		} else if !perms.Can("payment_term", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = sl.Badges.NoPermission
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: code},
				{Type: "text", Value: paymentTermTypeLabel(termType)},
				{Type: "text", Value: netDays},
				{Type: "badge", Value: isDefault, Variant: isDefaultVariant},
				{Type: "badge", Value: activeStatus, Variant: activeVariant},
			},
			DataAttrs: map[string]string{
				"name":          name,
				"code":          code,
				"type":          termType,
				"status":        recordStatus,
				"deletable":     strconv.FormatBool(!isInUse),
				"deactivatable": strconv.FormatBool(active),
				"activatable":   strconv.FormatBool(!active),
			},
			Actions: actions,
		})
	}
	return rows
}

func buildBulkActions(l paymentterm.Labels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes paymentterm.Routes, perms *types.UserPermissions) []types.BulkAction {
	canUpdate := perms.Can("payment_term", "update")
	canDelete := perms.Can("payment_term", "delete")
	return []types.BulkAction{
		{
			Key:              "deactivate",
			Label:            l.Actions.Deactivate,
			Icon:             "icon-clock-off",
			Variant:          "warning",
			Endpoint:         routes.BulkSetStatusURL,
			ConfirmTitle:     l.Actions.Deactivate,
			ConfirmMessage:   sl.Confirm.BulkDeactivate,
			ExtraParamsJSON:  `{"target_status":"inactive"}`,
			RequiresDataAttr: "deactivatable",
			Disabled:         !canUpdate,
			DisabledTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "payment_term:update"),
		},
		{
			Key:              "activate",
			Label:            l.Actions.Activate,
			Icon:             "icon-clock",
			Variant:          "primary",
			Endpoint:         routes.BulkSetStatusURL,
			ConfirmTitle:     l.Actions.Activate,
			ConfirmMessage:   sl.Confirm.BulkActivate,
			ExtraParamsJSON:  `{"target_status":"active"}`,
			RequiresDataAttr: "activatable",
			Disabled:         !canUpdate,
			DisabledTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "payment_term:update"),
		},
		{
			Key:              "delete",
			Label:            cl.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         routes.BulkDeleteURL,
			ConfirmTitle:     cl.Bulk.Delete,
			ConfirmMessage:   sl.Confirm.BulkDelete,
			RequiresDataAttr: "deletable",
			Disabled:         !canDelete,
			DisabledTooltip:  fmt.Sprintf(cl.Errors.MissingPermission, "payment_term:delete"),
		},
	}
}
