package tag

import (
	"context"
	"fmt"
	"log"
	"strconv"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"

	categorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_category"
)

// Deps holds view dependencies.
type Deps struct {
	ListCategories       func(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error)
	ListClientCategories func(ctx context.Context, req *clientcategorypb.ListClientCategoriesRequest) (*clientcategorypb.ListClientCategoriesResponse, error)
	RefreshURL           string
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

		columns := []types.TableColumn{
			{Key: "name", Label: "Tag Name", Sortable: true},
			{Key: "customers", Label: "Customers", Sortable: false, Width: "120px"},
			{Key: "description", Label: "Description", Sortable: true},
			{Key: "status", Label: "Status", Sortable: true, Width: "120px"},
		}

		rows := buildTableRows(resp.GetData(), customerCounts)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = []types.BulkAction{
			{
				Key:            "delete",
				Label:          "Delete",
				Icon:           "icon-trash-2",
				Variant:        "danger",
				Endpoint:       entydad.ClientTagBulkDeleteURL,
				ConfirmTitle:   "Delete Tags",
				ConfirmMessage: "Are you sure you want to delete {{count}} tag(s)? This action cannot be undone.",
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
				Title:   "No tags found",
				Message: "Create your first tag to start organizing clients.",
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     "Add Tag",
				ActionURL: entydad.ClientTagAddURL,
				Icon:      "icon-plus",
			},
			BulkActions: &bulkCfg,
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Client Tags",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "clients",
				ActiveSubNav:   "tags",
				HeaderTitle:    "Client Tags",
				HeaderSubtitle: "Manage tags for organizing clients",
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "client-tag-list-content",
			Table:           tableConfig,
		}

		return view.OK("client-tag-list", pageData)
	})
}

func buildTableRows(categories []*categorypb.Category, customerCounts map[string]int) []types.TableRow {
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

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: strconv.Itoa(count)},
				{Type: "text", Value: desc},
				{Type: "badge", Value: status, Variant: variant},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"status": status,
			},
			Actions: []types.TableAction{
				{Type: "edit", Label: "Edit", Action: "edit", URL: "/action/clients/tags/edit/" + id, DrawerTitle: "Edit Tag"},
				{Type: "delete", Label: "Delete", Action: "delete", URL: "/action/clients/tags/delete", ItemName: name},
			},
		})
	}
	return rows
}
