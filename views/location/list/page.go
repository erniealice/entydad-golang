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

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	RefreshURL      string
	Routes          entydad.LocationRoutes
	Labels          entydad.LocationLabels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the location list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the location list view (full page).
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
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "locations",
				ActiveSubNav:   "locations-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-map-pin",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "location-list-content",
			Table:           tableConfig,
		}

		return view.OK("location-list", pageData)
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

// buildTableConfig fetches location data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps, status string) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	resp, err := deps.GetListPageData(ctx, &locationpb.GetLocationListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list locations: %v", err)
		return nil, fmt.Errorf("failed to load locations: %w", err)
	}

	// Check which items are in use
	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range resp.GetLocationList() {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	l := deps.Labels
	columns := locationColumns(l)
	rows := buildTableRows(resp.GetLocationList(), status, l, deps.SharedLabels, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	tableConfig := &types.TableConfig{
		ID:                   "locations-table",
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
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddLocation,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("location", "create"),
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func locationColumns(l entydad.LocationLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "address", Label: l.Columns.Address, Sortable: true},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(locations []*locationpb.Location, status string, l entydad.LocationLabels, sl entydad.SharedLabels, routes entydad.LocationRoutes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, loc := range locations {
		active := loc.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := loc.GetId()
		name := loc.GetName()
		address := loc.GetAddress()

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit,
				Disabled: !perms.Can("location", "update"), DisabledTooltip: sl.Badges.NoPermission},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
				Disabled:       !perms.Can("location", "update"), DisabledTooltip: sl.Badges.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
				Disabled:       !perms.Can("location", "update"), DisabledTooltip: sl.Badges.NoPermission,
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
		} else if !perms.Can("location", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = "No permission"
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: address},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"address":   address,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l entydad.LocationLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l entydad.LocationLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.LocationLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.LocationLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
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

func buildBulkActions(l entydad.LocationLabels, sl entydad.SharedLabels, common pyeza.CommonLabels, status string, routes entydad.LocationRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-map-pin-off",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-map-pin",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  sl.Confirm.BulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:              "delete",
		Label:            common.Bulk.Delete,
		Icon:             "icon-trash-2",
		Variant:          "danger",
		Endpoint:         routes.BulkDeleteURL,
		ConfirmTitle:     common.Bulk.Delete,
		ConfirmMessage:   sl.Confirm.BulkDelete,
		RequiresDataAttr: "deletable",
	})

	return actions
}
