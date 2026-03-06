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
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
)

// Deps holds view dependencies.
type Deps struct {
	Routes               entydad.ClientTagRoutes
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	GetInUseIDs          func(ctx context.Context, ids []string) (map[string]bool, error)
	RefreshURL           string
	Labels               entydad.ClientTagLabels
	SharedLabels         entydad.SharedLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
}

// PageData holds the data for the client tags list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the client tags list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		resp, err := deps.ListCategories(ctx, &categorypb.ListCategoriesRequest{})
		if err != nil {
			log.Printf("Failed to list client tags: %v", err)
			return view.Error(fmt.Errorf("failed to load tags: %w", err))
		}

		// Build customer count per category from client_category junction records
		customerCounts := make(map[string]int)
		if deps.ListClientCategories != nil {
			ccResp, err := deps.ListClientCategories(ctx, &clientcategorypb.ListClientCategoriesRequest{})
			if err != nil {
				log.Printf("Failed to list client categories for counts: %v", err)
			} else {
				for _, cc := range ccResp.GetData() {
					customerCounts[cc.GetCategoryId()]++
				}
			}
		}

		// Check which items are in use (only for client-module categories)
		var inUseIDs map[string]bool
		if deps.GetInUseIDs != nil {
			var itemIDs []string
			for _, cat := range resp.GetData() {
				if cat.GetModule() == "client" {
					itemIDs = append(itemIDs, cat.GetId())
				}
			}
			inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
		}

		l := deps.Labels
		columns := []types.TableColumn{
			{Key: "name", Label: l.Columns.TagName, Sortable: true},
			{Key: "customers", Label: l.Columns.Customers, Sortable: false, Width: "120px"},
			{Key: "description", Label: l.Columns.Description, Sortable: true},
			{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
		}

		rows := buildTableRows(resp.GetData(), customerCounts, deps.Routes, inUseIDs, l)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = []types.BulkAction{
			{
				Key:              "delete",
				Label:            l.Actions.Delete,
				Icon:             "icon-trash-2",
				Variant:          "danger",
				Endpoint:         deps.Routes.BulkDeleteURL,
				ConfirmTitle:     l.Confirm.DeleteTitle,
				ConfirmMessage:   l.Confirm.DeleteMessage,
				RequiresDataAttr: "deletable",
			},
		}

		tableConfig := &types.TableConfig{
			ID:                   "client-tags-table",
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
				ActiveNav:      "clients",
				ActiveSubNav:   "tags",
				HeaderTitle:    l.Page.Heading,
				HeaderSubtitle: l.Page.Subtitle,
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "client-tag-list-content",
			Table:           tableConfig,
		}

		return view.OK("client-tag-list", pageData)
	})
}

func buildTableRows(categories []*categorypb.Category, customerCounts map[string]int, routes entydad.ClientTagRoutes, inUseIDs map[string]bool, l entydad.ClientTagLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, cat := range categories {
		// Only show client-module categories
		if cat.GetModule() != "client" {
			continue
		}

		id := cat.GetId()
		name := cat.GetName()
		desc := cat.GetDescription()
		count := customerCounts[id]
		status := "active"
		variant := "success"
		if !cat.GetActive() {
			status = "inactive"
			variant = "warning"
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

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: strconv.Itoa(count)},
				{Type: "text", Value: desc},
				{Type: "badge", Value: status, Variant: variant},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    status,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: []types.TableAction{
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit},
				deleteAction,
			},
		})
	}
	return rows
}
