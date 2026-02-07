package list

import (
	"context"
	"fmt"
	"log"

	"leapfor.xyz/entydad"

	userpb "leapfor.xyz/esqyma/golang/v1/domain/entity/user"

	"github.com/erniealice/pyeza-golang/types"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *userpb.GetUserListPageDataRequest) (*userpb.GetUserListPageDataResponse, error)
}

// PageData holds the data for the user list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the user list view.
func NewView(deps *Deps) entydad.View {
	return entydad.ViewFunc(func(ctx context.Context, viewCtx *entydad.ViewContext) entydad.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.GetListPageData(ctx, &userpb.GetUserListPageDataRequest{})
		if err != nil {
			log.Printf("Failed to list users: %v", err)
			return entydad.Error(fmt.Errorf("failed to load users: %w", err))
		}

		columns := userColumns()
		rows := buildTableRows(resp.GetUserList(), status)
		types.ApplyColumnStyles(columns, rows)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "users",
				ActiveSubNav:   "users-" + status,
				HeaderTitle:    statusTitle(status),
				HeaderSubtitle: statusSubtitle(status),
				HeaderIcon:     "icon-users",
			},
			ContentTemplate: "user-list-content",
			Table: &types.TableConfig{
				ID:         "users-table",
				Columns:    columns,
				Rows:       rows,
				ShowSearch: true,
				ShowActions: true,
				EmptyState: types.TableEmptyState{
					Title:   "No users found",
					Message: "No " + status + " users to display.",
				},
				PrimaryAction: &types.PrimaryAction{
					Label:     "Add User",
					ActionURL: "/action/users/add",
					Icon:      "icon-plus",
				},
			},
		}

		return entydad.OK("user-list", pageData)
	})
}

func userColumns() []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: "Name", Sortable: true},
		{Key: "email", Label: "Email", Sortable: true},
		{Key: "mobile", Label: "Mobile", Sortable: true, Width: "150px"},
		{Key: "status", Label: "Status", Sortable: true, Width: "120px"},
	}
}

func buildTableRows(users []*userpb.User, status string) []types.TableRow {
	rows := []types.TableRow{}
	for _, u := range users {
		active := u.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := u.GetId()
		name := u.GetFirstName() + " " + u.GetLastName()
		email := u.GetEmailAddress()
		mobile := u.GetMobileNumber()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: email},
				{Type: "text", Value: mobile},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"email":  email,
				"mobile": mobile,
				"status": recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: "View User", Action: "view", Href: "/app/users/" + id},
				{Type: "edit", Label: "Edit User", Action: "edit", URL: "/action/users/edit/" + id, DrawerTitle: "Edit User"},
				{Type: "delete", Label: "Delete User", Action: "delete", URL: "/action/users/delete", ItemName: name},
			},
		})
	}
	return rows
}

func statusTitle(status string) string {
	switch status {
	case "active":
		return "Active Users"
	case "inactive":
		return "Inactive Users"
	default:
		return "Users"
	}
}

func statusSubtitle(status string) string {
	switch status {
	case "active":
		return "Manage your active user accounts"
	case "inactive":
		return "View inactive user accounts"
	default:
		return "User management"
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
