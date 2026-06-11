// Package list provides the workspace_user list view.
// This is a minimal list view to make workspace_user reachable for debugging.
// The primary workspace_user surface is the nested detail page (Phase 2).
package list

import (
	"context"
	"fmt"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	workspace_user "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
)

// ListViewDeps holds dependencies for the workspace_user list view.
type ListViewDeps struct {
	Routes          workspace_user.Routes
	Labels          workspace_user.Labels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
	GetListPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserListPageDataRequest) (*workspaceuserpb.GetWorkspaceUserListPageDataResponse, error)
}

// PageData is the template data for the workspace_user list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	Labels          workspace_user.Labels
	Status          string
	DetailBaseURL   string
}

// NewView creates the workspace_user list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("workspace_user", "list") {
			return view.Forbidden("workspace_user:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := []types.TableColumn{
			{Key: "user_name", Label: deps.Labels.Columns.UserName},
			{Key: "email", Label: deps.Labels.Columns.Email},
			{Key: "roles", Label: deps.Labels.Columns.Roles, WidthClass: "col-2xl"},
			{Key: "status", Label: deps.Labels.Columns.Status, WidthClass: "col-2xl"},
		}

		var rows []types.TableRow

		if deps.GetListPageData != nil {
			resp, err := deps.GetListPageData(ctx, &workspaceuserpb.GetWorkspaceUserListPageDataRequest{})
			if err == nil {
				for _, wu := range resp.GetWorkspaceUserList() {
					// Client-side status filter
					if status == "active" && !wu.GetActive() {
						continue
					}
					if status == "inactive" && wu.GetActive() {
						continue
					}
					u := wu.GetUser()
					userName := ""
					email := ""
					if u != nil {
						userName = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
						email = u.GetEmailAddress()
					}
					roleCount := len(wu.GetWorkspaceUserRoles())
					active := wu.GetActive()
					statusValue := "active"
					statusVariant := "success"
					if !active {
						statusValue = "inactive"
						statusVariant = "warning"
					}

					roleLabel := fmt.Sprintf("%d roles", roleCount)
					if roleCount == 1 {
						roleLabel = "1 role"
					}
					cells := []types.TableCell{
						{Type: "text", Value: userName},
						{Type: "text", Value: email},
						{Type: "text", Value: roleLabel},
						{Type: "badge", Value: statusValue, Variant: statusVariant},
					}

					rows = append(rows, types.TableRow{
						ID:    wu.GetId(),
						Cells: cells,
						DataAttrs: map[string]string{
							"testid": "workspace-user-row-" + wu.GetId(),
						},
					})
				}
			}
		}

		types.ApplyColumnStyles(columns, rows)

		tc := &types.TableConfig{
			ID:                   "workspace-users-table",
			Columns:              columns,
			Rows:                 rows,
			Labels:               deps.TableLabels,
			ShowSearch:           true,
			ShowActions:          false,
			ShowSort:             true,
			DefaultSortColumn:    "user_name",
			DefaultSortDirection: "asc",
			EmptyState: types.TableEmptyState{
				Title:   "No workspace users",
				Message: "No workspace-user assignments exist.",
			},
		}
		types.ApplyTableSettings(tc)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        "Workspace Users",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "admin",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "workspace-user-list-content",
			Table:           tc,
			Labels:          deps.Labels,
			Status:          status,
		}

		return view.OK("workspace-user-list", pageData)
	})
}
