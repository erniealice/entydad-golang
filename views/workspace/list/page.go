package list

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"

	"github.com/erniealice/entydad-golang"
)

// Deps holds view dependencies.
type Deps struct {
	GetListPageData func(ctx context.Context, req *workspacepb.GetWorkspaceListPageDataRequest) (*workspacepb.GetWorkspaceListPageDataResponse, error)
	RefreshURL      string
	Labels          entydad.WorkspaceLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the workspace list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the workspace list view (full page).
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
				ActiveNav:      "admin",
				ActiveSubNav:   "workspaces-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-briefcase",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "workspace-list-content",
			Table:           tableConfig,
		}

		return view.OK("workspace-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
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

// buildTableConfig fetches workspace data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *Deps, status string) (*types.TableConfig, error) {
	resp, err := deps.GetListPageData(ctx, &workspacepb.GetWorkspaceListPageDataRequest{})
	if err != nil {
		log.Printf("Failed to list workspaces: %v", err)
		return nil, fmt.Errorf("failed to load workspaces: %w", err)
	}

	l := deps.Labels
	columns := workspaceColumns(l)
	rows := buildTableRows(resp.GetWorkspaceList(), status, l)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := entydad.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status)

	refreshURL := fmt.Sprintf("/action/workspaces/table/%s", status)

	tableConfig := &types.TableConfig{
		ID:                   "workspaces-table",
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
			Label:     l.Buttons.AddWorkspace,
			ActionURL: "/action/workspaces/add",
			Icon:      "icon-plus",
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func workspaceColumns(l entydad.WorkspaceLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "description", Label: l.Columns.Description, Sortable: true},
		{Key: "private", Label: l.Columns.Private, Sortable: true, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(workspaces []*workspacepb.Workspace, status string, l entydad.WorkspaceLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, w := range workspaces {
		active := w.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := w.GetId()
		name := w.GetName()
		description := w.GetDescription()
		private := w.GetPrivate()

		privateLabel := "No"
		privateVariant := "default"
		if private {
			privateLabel = "Yes"
			privateVariant = "info"
		}

		actions := []types.TableAction{
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: "/action/workspaces/edit/" + id, DrawerTitle: l.Actions.Edit},
		}
		if active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: "/action/workspaces/set-status?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to deactivate %s?", name),
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: "/action/workspaces/set-status?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf("Are you sure you want to activate %s?", name),
			})
		}
		actions = append(actions, types.TableAction{
			Type: "delete", Label: l.Actions.Delete, Action: "delete",
			URL: "/action/workspaces/delete", ItemName: name,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "badge", Value: privateLabel, Variant: privateVariant},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":        name,
				"description": description,
				"private":     fmt.Sprintf("%v", private),
				"status":      recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l entydad.WorkspaceLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l entydad.WorkspaceLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l entydad.WorkspaceLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l entydad.WorkspaceLabels, status string) string {
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

func buildBulkActions(l entydad.WorkspaceLabels, common pyeza.CommonLabels, status string) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-briefcase",
			Variant:         "warning",
			Endpoint:        "/action/workspaces/bulk-set-status",
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  "Are you sure you want to deactivate {{count}} workspace(s)?",
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-briefcase",
			Variant:         "primary",
			Endpoint:        "/action/workspaces/bulk-set-status",
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  "Are you sure you want to activate {{count}} workspace(s)?",
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:            "delete",
		Label:          common.Bulk.Delete,
		Icon:           "icon-trash-2",
		Variant:        "danger",
		Endpoint:       "/action/workspaces/bulk-delete",
		ConfirmTitle:   common.Bulk.Delete,
		ConfirmMessage: "Are you sure you want to delete {{count}} workspace(s)? This action cannot be undone.",
	})

	return actions
}
