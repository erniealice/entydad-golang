package list

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
	RefreshURL      string
	Labels          entydad.ClientLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the client list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the client list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.GetListPageData(ctx, &clientpb.GetClientListPageDataRequest{})
		if err != nil {
			log.Printf("Failed to list clients: %v", err)
			return view.Error(fmt.Errorf("failed to load clients: %w", err))
		}

		l := deps.Labels
		columns := clientColumns(l)
		rows := buildTableRows(resp.GetClientList(), status, l)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = []types.BulkAction{
			{
				Key:            "set-inactive",
				Label:          l.BulkActions.SetAsInactive,
				Icon:           "icon-user-minus",
				Variant:        "warning",
				Endpoint:       "/action/clients/bulk-delete",
				ConfirmTitle:   l.BulkActions.SetAsInactive,
				ConfirmMessage: "Are you sure you want to set {{count}} customer(s) as inactive?",
			},
			{
				Key:            "delete",
				Label:          deps.CommonLabels.Bulk.Delete,
				Icon:           "icon-trash-2",
				Variant:        "danger",
				Endpoint:       "/action/clients/bulk-delete",
				ConfirmTitle:   deps.CommonLabels.Bulk.Delete,
				ConfirmMessage: "Are you sure you want to delete {{count}} customer(s)? This action cannot be undone.",
			},
		}

		tableConfig := &types.TableConfig{
			ID:                   "clients-table",
			RefreshURL:           deps.RefreshURL,
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
				Label:     l.Buttons.AddNew,
				ActionURL: "/action/clients/add",
				Icon:      "icon-plus",
			},
			BulkActions: &bulkCfg,
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "clients",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: statusPageCaption(l, status),
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "client-list-content",
			Table:           tableConfig,
		}

		return view.OK("client-list", pageData)
	})
}

func clientColumns(l entydad.ClientLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.ClientName, Sortable: true},
		{Key: "email", Label: l.Form.Email, Sortable: true},
		{Key: "phone", Label: l.Form.Phone, Sortable: false},
		{Key: "status", Label: l.Detail.CompanyDetails.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(clients []*clientpb.Client, status string, l entydad.ClientLabels) []types.TableRow {
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

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: email},
				{Type: "text", Value: phone},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"email":  email,
				"status": recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Detail.Actions.ViewClient, Action: "view", Href: "/app/clients/" + id},
				{Type: "edit", Label: l.Detail.Actions.EditClient, Action: "edit", URL: "/action/clients/edit/" + id, DrawerTitle: l.Detail.Actions.EditClient},
				{Type: "delete", Label: l.Detail.Actions.DeleteClient, Action: "delete", URL: "/action/clients/delete", ItemName: name},
			},
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
