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
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *paymenttermpb.GetPaymentTermListPageDataRequest) (*paymenttermpb.GetPaymentTermListPageDataResponse, error)
	RefreshURL      string
	Routes          entydad.PaymentTermRoutes
	Labels          entydad.PaymentTermLabels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
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
		tableConfig, err := buildTableConfig(ctx, deps)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "settings",
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
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "payment-terms"); kb != nil {
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

	l := deps.Labels
	columns := paymentTermColumns(l)
	rows := buildTableRows(resp.GetPaymentTermList(), l, deps.SharedLabels, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:            "delete",
			Label:          deps.CommonLabels.Bulk.Delete,
			Icon:           "icon-trash-2",
			Variant:        "danger",
			Endpoint:       deps.Routes.BulkDeleteURL,
			ConfirmTitle:   deps.CommonLabels.Bulk.Delete,
			ConfirmMessage: deps.SharedLabels.Confirm.BulkDelete,
		},
	}

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
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func paymentTermColumns(l entydad.PaymentTermLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true, Width: "200px"},
		{Key: "code", Label: l.Columns.Code, Sortable: true, Width: "120px"},
		{Key: "type", Label: l.Columns.Type, Sortable: true, Width: "120px"},
		{Key: "net_days", Label: l.Columns.NetDays, Sortable: true, Width: "80px"},
		{Key: "entity_scope", Label: l.Columns.EntityScope, Sortable: true, Width: "120px"},
		{Key: "is_default", Label: l.Columns.IsDefault, Sortable: true, Width: "80px"},
		{Key: "active", Label: l.Columns.Active, Sortable: true, Width: "80px"},
	}
}

func buildTableRows(items []*paymenttermpb.PaymentTerm, l entydad.PaymentTermLabels, sl entydad.SharedLabels, routes entydad.PaymentTermRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, pt := range items {
		id := pt.GetId()
		name := pt.GetName()
		code := pt.GetCode()
		termType := pt.GetType()
		netDays := strconv.FormatInt(int64(pt.GetNetDays()), 10)
		entityScope := pt.GetEntityScope()
		isDefault := "no"
		isDefaultVariant := "default"
		if pt.GetIsDefault() {
			isDefault = "yes"
			isDefaultVariant = "success"
		}
		activeStatus := "inactive"
		activeVariant := "warning"
		if pt.GetActive() {
			activeStatus = "active"
			activeVariant = "success"
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: code},
				{Type: "text", Value: termType},
				{Type: "text", Value: netDays},
				{Type: "text", Value: entityScope},
				{Type: "badge", Value: isDefault, Variant: isDefaultVariant},
				{Type: "badge", Value: activeStatus, Variant: activeVariant},
			},
			DataAttrs: map[string]string{
				"name":         name,
				"code":         code,
				"type":         termType,
				"entity_scope": entityScope,
				"active":       activeStatus,
			},
			Actions: []types.TableAction{
				{
					Type:        "edit",
					Label:       l.Actions.Edit,
					Action:      "edit",
					URL:         route.ResolveURL(routes.EditURL, "id", id),
					DrawerTitle: l.Actions.Edit,
					Disabled:    !perms.Can("payment_term", "update"),
					DisabledTooltip: sl.Badges.NoPermission,
				},
				{
					Type:     "delete",
					Label:    l.Actions.Delete,
					Action:   "delete",
					URL:      routes.DeleteURL,
					ItemName: name,
					Disabled: !perms.Can("payment_term", "delete"),
					DisabledTooltip: sl.Badges.NoPermission,
				},
			},
		})
	}
	return rows
}
