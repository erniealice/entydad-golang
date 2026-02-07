package list

import (
	"context"
	"fmt"
	"log"

	"leapfor.xyz/entydad"

	clientpb "leapfor.xyz/esqyma/golang/v1/domain/entity/client"

	"github.com/erniealice/pyeza-golang/types"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *clientpb.GetClientListPageDataRequest) (*clientpb.GetClientListPageDataResponse, error)
}

// PageData holds the data for the client list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the client list view.
func NewView(deps *Deps) entydad.View {
	return entydad.ViewFunc(func(ctx context.Context, viewCtx *entydad.ViewContext) entydad.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.GetListPageData(ctx, &clientpb.GetClientListPageDataRequest{})
		if err != nil {
			log.Printf("Failed to list clients: %v", err)
			return entydad.Error(fmt.Errorf("failed to load clients: %w", err))
		}

		t := viewCtx.T
		columns := clientColumns(t)
		rows := buildTableRows(resp.GetClientList(), status, t)
		types.ApplyColumnStyles(columns, rows)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(t, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "clients",
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(t, status),
				HeaderSubtitle: statusPageCaption(t, status),
				HeaderIcon:     "icon-users",
			},
			ContentTemplate: "client-list-content",
			Table: &types.TableConfig{
				ID:          "clients-table",
				Columns:     columns,
				Rows:        rows,
				ShowSearch:  true,
				ShowActions: true,
				EmptyState: types.TableEmptyState{
					Title:   statusEmptyTitle(t, status),
					Message: statusEmptyMessage(t, status),
				},
				PrimaryAction: &types.PrimaryAction{
					Label:     t("client.buttons.addNew"),
					ActionURL: "/action/clients/add",
					Icon:      "icon-plus",
				},
			},
		}

		return entydad.OK("client-list", pageData)
	})
}

// T is a translation function type (alias for readability).
type T = func(string) string

func clientColumns(t T) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: t("client.columns.clientName"), Sortable: true},
		{Key: "email", Label: t("client.form.email"), Sortable: true},
		{Key: "phone", Label: t("client.form.phone"), Sortable: false},
		{Key: "status", Label: t("client.detail.companyDetails.status"), Sortable: true, Width: "120px"},
	}
}

func buildTableRows(clients []*clientpb.Client, status string, t T) []types.TableRow {
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
				{Type: "view", Label: t("client.detail.actions.viewClient"), Action: "view", Href: "/app/clients/" + id},
				{Type: "edit", Label: t("client.detail.actions.editClient"), Action: "edit", URL: "/action/clients/edit/" + id, DrawerTitle: t("client.detail.actions.editClient")},
				{Type: "delete", Label: t("client.detail.actions.deleteClient"), Action: "delete", URL: "/action/clients/delete", ItemName: name},
			},
		})
	}
	return rows
}

func statusPageTitle(t T, status string) string {
	switch status {
	case "active":
		return t("client.page.headingActive")
	case "prospect":
		return t("client.page.headingProspect")
	case "inactive":
		return t("client.page.headingInactive")
	default:
		return t("client.page.heading")
	}
}

func statusPageCaption(t T, status string) string {
	switch status {
	case "active":
		return t("client.page.captionActive")
	case "prospect":
		return t("client.page.captionProspect")
	case "inactive":
		return t("client.page.captionInactive")
	default:
		return t("client.page.caption")
	}
}

func statusEmptyTitle(t T, status string) string {
	switch status {
	case "active":
		return t("client.empty.activeTitle")
	case "prospect":
		return t("client.empty.prospectTitle")
	case "inactive":
		return t("client.empty.inactiveTitle")
	default:
		return t("client.empty.activeTitle")
	}
}

func statusEmptyMessage(t T, status string) string {
	switch status {
	case "active":
		return t("client.empty.activeMessage")
	case "prospect":
		return t("client.empty.prospectMessage")
	case "inactive":
		return t("client.empty.inactiveMessage")
	default:
		return t("client.empty.activeMessage")
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
