package tag

import (
	"context"
	"fmt"
	"log"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	suppliercategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier_category"
)

// Deps holds view dependencies.
type Deps struct {
	Routes                  entydad.SupplierTagRoutes
	GetCategoryListPageData func(ctx context.Context) ([]*categorypb.Category, error)
	ListSupplierCategories  func(ctx context.Context, req *suppliercategorypb.ListSupplierCategoriesRequest) (*suppliercategorypb.ListSupplierCategoriesResponse, error)
	GetInUseIDs             func(ctx context.Context, ids []string) (map[string]bool, error)
	RefreshURL              string
	Labels                  entydad.SupplierTagLabels
	SharedLabels            entydad.SharedLabels
	CommonLabels            pyeza.CommonLabels
	TableLabels             types.TableLabels
}

// PageData holds the data for the supplier tags list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the supplier tags list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		cats, err := deps.GetCategoryListPageData(ctx)
		if err != nil {
			log.Printf("Failed to list supplier tags: %v", err)
			return view.Error(fmt.Errorf("failed to load tags: %w", err))
		}

		// Build supplier count per category from supplier_category junction records
		supplierCounts := make(map[string]int)
		if deps.ListSupplierCategories != nil {
			scResp, err := deps.ListSupplierCategories(ctx, &suppliercategorypb.ListSupplierCategoriesRequest{})
			if err != nil {
				log.Printf("Failed to list supplier categories for counts: %v", err)
			} else {
				for _, sc := range scResp.GetData() {
					supplierCounts[sc.GetCategoryId()]++
				}
			}
		}

		// Check which items are in use (only for supplier-module categories)
		var inUseIDs map[string]bool
		if deps.GetInUseIDs != nil {
			var itemIDs []string
			for _, cat := range cats {
				if cat.GetModule() == "supplier" {
					itemIDs = append(itemIDs, cat.GetId())
				}
			}
			inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
		}

		l := deps.Labels
		columns := []types.TableColumn{
			{Key: "name", Label: l.Columns.TagName},
			{Key: "suppliers", Label: l.Columns.Suppliers, NoSort: true, WidthClass: "col-2xl"},
			{Key: "description", Label: l.Columns.Description},
			{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		}

		rows := buildTableRows(cats, supplierCounts, deps.Routes, inUseIDs, l, deps.SharedLabels)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, deps.Routes)

		tableConfig := &types.TableConfig{
			ID:                   "supplier-tags-table",
			RefreshURL:           deps.RefreshURL,
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           true,
			ShowActions:          true,
			DefaultSortColumn:    "name",
			DefaultSortDirection: "asc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     l.Buttons.AddTag,
				ActionURL: deps.Routes.AddURL,
				Icon:      "icon-plus",
			},
			BulkActions: &bulkCfg,
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "supplier",
				ActiveSubNav:   "tags",
				HeaderTitle:    l.Page.Heading,
				HeaderSubtitle: l.Page.Subtitle,
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-tag-list-content",
			Table:           tableConfig,
		}

		return view.OK("supplier-tag-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML (used for HTMX refresh).
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		cats, err := deps.GetCategoryListPageData(ctx)
		if err != nil {
			log.Printf("Failed to list supplier tags: %v", err)
			return view.Error(fmt.Errorf("failed to load tags: %w", err))
		}

		supplierCounts := make(map[string]int)
		if deps.ListSupplierCategories != nil {
			scResp, err := deps.ListSupplierCategories(ctx, &suppliercategorypb.ListSupplierCategoriesRequest{})
			if err != nil {
				log.Printf("Failed to list supplier categories for counts: %v", err)
			} else {
				for _, sc := range scResp.GetData() {
					supplierCounts[sc.GetCategoryId()]++
				}
			}
		}

		var inUseIDs map[string]bool
		if deps.GetInUseIDs != nil {
			var itemIDs []string
			for _, cat := range cats {
				if cat.GetModule() == "supplier" {
					itemIDs = append(itemIDs, cat.GetId())
				}
			}
			inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
		}

		l := deps.Labels
		columns := []types.TableColumn{
			{Key: "name", Label: l.Columns.TagName},
			{Key: "suppliers", Label: l.Columns.Suppliers, NoSort: true, WidthClass: "col-2xl"},
			{Key: "description", Label: l.Columns.Description},
			{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		}

		rows := buildTableRows(cats, supplierCounts, deps.Routes, inUseIDs, l, deps.SharedLabels)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, deps.Routes)

		tableConfig := &types.TableConfig{
			ID:                   "supplier-tags-table",
			RefreshURL:           deps.RefreshURL,
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           true,
			ShowActions:          true,
			DefaultSortColumn:    "name",
			DefaultSortDirection: "asc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     l.Buttons.AddTag,
				ActionURL: deps.Routes.AddURL,
				Icon:      "icon-plus",
			},
			BulkActions: &bulkCfg,
		}
		types.ApplyTableSettings(tableConfig)

		return view.OK("table-card", tableConfig)
	})
}

func buildTableRows(categories []*categorypb.Category, supplierCounts map[string]int, routes entydad.SupplierTagRoutes, inUseIDs map[string]bool, l entydad.SupplierTagLabels, sl entydad.SharedLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, cat := range categories {
		// Only show supplier-module categories
		if cat.GetModule() != "supplier" {
			continue
		}

		id := cat.GetId()
		name := cat.GetName()
		desc := cat.GetDescription()
		count := supplierCounts[id]
		active := cat.GetActive()
		statusLabel := "Active"
		variant := "success"
		recordStatus := "active"
		if !active {
			statusLabel = "Inactive"
			variant = "warning"
			recordStatus = "inactive"
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
			deleteAction.DisabledTooltip = l.Confirm.CannotDelete
		}

		actions := []types.TableAction{
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type:           "deactivate",
				Label:          l.Actions.Deactivate,
				Action:         "deactivate",
				URL:            routes.SetStatusURL + "?status=inactive",
				ItemName:       name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
			})
		} else {
			actions = append(actions, types.TableAction{
				Type:           "activate",
				Label:          l.Actions.Activate,
				Action:         "activate",
				URL:            routes.SetStatusURL + "?status=active",
				ItemName:       name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
			})
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: strconv.Itoa(count)},
				{Type: "text", Value: desc},
				{Type: "badge", Value: statusLabel, Variant: variant},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: actions,
		})
	}
	return rows
}

func buildBulkActions(l entydad.SupplierTagLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, routes entydad.SupplierTagRoutes) []types.BulkAction {
	return []types.BulkAction{
		{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-tag-off",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		},
		{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-tag",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		},
		{
			Key:              "delete",
			Label:            l.Actions.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         routes.BulkDeleteURL,
			ConfirmTitle:     l.Confirm.DeleteTitle,
			ConfirmMessage:   l.Confirm.DeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}
}
