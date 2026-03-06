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

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	Routes          entydad.ClientRoutes
	GetListPageData func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
	GetInUseIDs     func(ctx context.Context, ids []string) (map[string]bool, error)
	Labels          entydad.ClientLabels
	SharedLabels    entydad.SharedLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the client list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the client list view (full page).
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
				ActiveNav:      "clients",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "client-list-content",
			Table:           tableConfig,
		}

		return view.OK("client-list", pageData)
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

// buildTableConfig fetches client data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps, status string) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)
	resp, err := deps.GetListPageData(ctx, &clientpb.GetClientListPageDataRequest{})
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

	l := deps.Labels
	columns := clientColumns(l)
	rows := buildTableRows(resp.GetClientList(), status, l, deps.SharedLabels, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.SharedLabels, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

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
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddNew,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("client", "create"),
			DisabledTooltip: deps.SharedLabels.Badges.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func clientColumns(l entydad.ClientLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.ClientName, Sortable: true},
		{Key: "email", Label: l.Form.Email, Sortable: true},
		{Key: "phone", Label: l.Form.Phone, Sortable: false},
		{Key: "status", Label: l.Detail.CompanyDetails.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(clients []*clientpb.Client, status string, l entydad.ClientLabels, sl entydad.SharedLabels, routes entydad.ClientRoutes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, c := range clients {
		active := c.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := c.GetId()
		u := c.GetUser()
		name := u.GetFirstName() + " " + u.GetLastName()
		email := u.GetEmailAddress()
		phone := u.GetMobileNumber()
		isInUse := inUseIDs[id]

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: email},
				{Type: "text", Value: phone},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"email":     email,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: buildRowActions(id, name, active, isInUse, l, sl, routes, perms),
		})
	}
	return rows
}

func statusPageTitle(l entydad.ClientLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "prospect":
		return l.Page.HeadingProspect
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

func buildRowActions(id, name string, active, isInUse bool, l entydad.ClientLabels, sl entydad.SharedLabels, routes entydad.ClientRoutes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Detail.Actions.ViewClient, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Detail.Actions.EditClient, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Detail.Actions.EditClient,
			Disabled: !perms.Can("client", "update"), DisabledTooltip: "No permission"},
	}
	if active {
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Detail.Actions.DeactivateClient, Action: "deactivate",
			URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
			ConfirmTitle:   l.Detail.Actions.DeactivateClient,
			ConfirmMessage: fmt.Sprintf(sl.Confirm.Deactivate, name),
			Disabled: !perms.Can("client", "update"), DisabledTooltip: sl.Badges.NoPermission,
		})
	} else {
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Detail.Actions.ActivateClient, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: name,
			ConfirmTitle:   l.Detail.Actions.ActivateClient,
			ConfirmMessage: fmt.Sprintf(sl.Confirm.Activate, name),
			Disabled: !perms.Can("client", "update"), DisabledTooltip: sl.Badges.NoPermission,
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

func buildBulkActions(l entydad.ClientLabels, sl entydad.SharedLabels, cl pyeza.CommonLabels, status string, routes entydad.ClientRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.BulkActions.SetAsInactive,
			Icon:            "icon-user-minus",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.BulkActions.SetAsInactive,
			ConfirmMessage:  sl.Confirm.BulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           cl.Bulk.Activate,
			Icon:            "icon-user-check",
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
