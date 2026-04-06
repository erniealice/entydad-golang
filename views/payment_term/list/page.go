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
				ActiveNav:      "client",
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
	bulkCfg.Actions = buildBulkActions(deps.Labels, deps.SharedLabels, deps.CommonLabels, deps.Routes)

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
		{Key: "type", Label: l.Columns.Type, Sortable: true, Width: "160px"},
		{Key: "net_days", Label: l.Columns.NetDays, Sortable: true, Width: "80px"},
		{Key: "is_default", Label: l.Columns.IsDefault, Sortable: true, Width: "80px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "100px"},
	}
}

func paymentTermTypeLabel(code string) string {
	if label, ok := paymentTermTypeLabels[code]; ok {
		return label
	}
	return code
}

func buildTableRows(items []*paymenttermpb.PaymentTerm, l entydad.PaymentTermLabels, sl entydad.SharedLabels, routes entydad.PaymentTermRoutes, perms *types.UserPermissions) []types.TableRow {
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
		actions = append(actions, types.TableAction{
			Type:            "delete",
			Label:           l.Actions.Delete,
			Action:          "delete",
			URL:             routes.DeleteURL,
			ItemName:        name,
			Disabled:        !perms.Can("payment_term", "delete"),
			DisabledTooltip: sl.Badges.NoPermission,
		})

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
				"name":   name,
				"code":   code,
				"type":   termType,
				"status": recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func buildBulkActions(l entydad.PaymentTermLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.PaymentTermRoutes) []types.BulkAction {
	return []types.BulkAction{
		{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-clock-off",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		},
		{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-clock",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		},
		{
			Key:            "delete",
			Label:          cl.Bulk.Delete,
			Icon:           "icon-trash-2",
			Variant:        "danger",
			Endpoint:       routes.BulkDeleteURL,
			ConfirmTitle:   cl.Bulk.Delete,
			ConfirmMessage: sl.Confirm.BulkDelete,
		},
	}
}
