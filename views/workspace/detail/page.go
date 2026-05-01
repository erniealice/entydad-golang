package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
)

// DetailViewDeps holds view dependencies for the workspace detail page.
type DetailViewDeps struct {
	Routes                        entydad.WorkspaceRoutes
	ReadWorkspace                 func(ctx context.Context, req *workspacepb.ReadWorkspaceRequest) (*workspacepb.ReadWorkspaceResponse, error)
	GetWorkspaceUserListPageData  func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserListPageDataRequest) (*workspaceuserpb.GetWorkspaceUserListPageDataResponse, error)
	Labels                        entydad.WorkspaceLabels
	CommonLabels                  pyeza.CommonLabels
	TableLabels                   types.TableLabels
	// WorkspaceUserDetailURL is the target route for the "View" row action on each workspace_user row.
	// Phase 2 will add this page; for now emit the URL so Phase 2 can register it.
	WorkspaceUserDetailURL string
	// WorkspaceUserAddURL is the drawer action for "Add user to workspace".
	// Phase 2 will add this handler; emit the URL so the button target is wired.
	WorkspaceUserAddURL string
}

// WorkspaceUserRow holds display data for a single workspace_user in the Users tab table.
type WorkspaceUserRow struct {
	ID         string
	UserName   string
	Email      string
	RoleCount  int
	Active     bool
	ViewURL    string
}

// PageData holds the data for the workspace detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	WorkspaceName   string
	WorkspaceID     string
	Description     string
	Currency        string
	Region          string
	Active          bool
	StatusVariant   string
	Labels          entydad.WorkspaceLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	EditURL         string
	// Users tab
	UsersTable          *types.TableConfig
	WorkspaceUserAddURL string
}

// tabLabels holds the resolved tab display strings, sourced from the lyngua
// workspace.json keys added in Phase 1.
type tabLabels struct {
	Info  string
	Users string
}

// resolveTabLabels returns display strings for the two tabs.
// The labels are loaded from WorkspaceLabels.Detail.Tabs (added in Phase 1).
// If the struct fields are empty (older lyngua build), fall back to English literals.
func resolveTabLabels(l entydad.WorkspaceLabels) tabLabels {
	info := l.Detail.Tabs.Info
	if info == "" {
		info = "Info"
	}
	users := l.Detail.Tabs.Users
	if users == "" {
		users = "Users"
	}
	return tabLabels{Info: info, Users: users}
}

// NewView creates the workspace detail view (full page load).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		ws, err := loadWorkspace(ctx, deps, id)
		if err != nil {
			return view.Error(err)
		}

		tl := resolveTabLabels(deps.Labels)
		tabItems := buildTabItems(id, deps, tl)

		pageData := buildPageData(viewCtx, id, ws, activeTab, tabItems, deps)

		switch activeTab {
		case "users":
			pageData.UsersTable = buildUsersTable(ctx, deps, id, tl)
		}

		return view.OK("workspace-detail", pageData)
	})
}

// NewTabAction creates the tab-action view (HTMX partial — returns only the tab content).
// Route: GET /action/workspace/{id}/tab/{tab}
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		ws, err := loadWorkspace(ctx, deps, id)
		if err != nil {
			return view.Error(err)
		}

		tl := resolveTabLabels(deps.Labels)
		tabItems := buildTabItems(id, deps, tl)
		pageData := buildPageData(viewCtx, id, ws, tab, tabItems, deps)

		switch tab {
		case "users":
			pageData.UsersTable = buildUsersTable(ctx, deps, id, tl)
			return view.OK("workspace-tab-users", pageData)
		default:
			return view.OK("workspace-tab-info", pageData)
		}
	})
}

// loadWorkspace fetches a single workspace by ID and returns it.
func loadWorkspace(ctx context.Context, deps *DetailViewDeps, id string) (*workspacepb.Workspace, error) {
	resp, err := deps.ReadWorkspace(ctx, &workspacepb.ReadWorkspaceRequest{
		Data: &workspacepb.Workspace{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read workspace %s: %v", id, err)
		return nil, fmt.Errorf("failed to load workspace: %w", err)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("workspace not found")
	}
	return data[0], nil
}

// buildPageData assembles the PageData struct from the loaded workspace and request context.
func buildPageData(viewCtx *view.ViewContext, id string, ws *workspacepb.Workspace, activeTab string, tabItems []pyeza.TabItem, deps *DetailViewDeps) *PageData {
	statusVariant := "success"
	if !ws.GetActive() {
		statusVariant = "warning"
	}
	return &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          ws.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "admin",
			ActiveSubNav:   "workspaces-active",
			HeaderTitle:    ws.GetName(),
			HeaderSubtitle: ws.GetDescription(),
			HeaderIcon:     "icon-briefcase",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:     "workspace-detail-content",
		WorkspaceName:       ws.GetName(),
		WorkspaceID:         id,
		Description:         ws.GetDescription(),
		Currency:            ws.GetFunctionalCurrency(),
		Region:              ws.GetComplianceRegion(),
		Active:              ws.GetActive(),
		StatusVariant:       statusVariant,
		Labels:              deps.Labels,
		ActiveTab:           activeTab,
		TabItems:            tabItems,
		EditURL:             route.ResolveURL(deps.Routes.EditURL, "id", id),
		WorkspaceUserAddURL: deps.WorkspaceUserAddURL + "?workspace_id=" + id,
	}
}

// buildTabItems constructs the two tab items (Info, Users) for the workspace detail page.
func buildTabItems(id string, deps *DetailViewDeps, tl tabLabels) []pyeza.TabItem {
	base := route.ResolveURL(deps.Routes.DetailURL, "id", id)
	action := route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: tl.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "users", Label: tl.Users, Href: base + "?tab=users", HxGet: action + "users", Icon: "icon-users"},
	}
}

// buildUsersTable loads workspace_user rows filtered by workspace_id and returns a TableConfig.
// If GetWorkspaceUserListPageData is nil (Phase 2 not yet wired) it returns an empty table
// skeleton so the page renders without error.
//
// TODO(Phase 2): once espyna exposes a workspace-scoped list page data, replace the manual
// in-process filter below with a proper server-side filter request.
func buildUsersTable(ctx context.Context, deps *DetailViewDeps, workspaceID string, tl tabLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "user_name", Label: "Name"},
		{Key: "email", Label: "Email"},
		{Key: "role_count", Label: "Roles", WidthClass: "col-2xl"},
		{Key: "status", Label: "Status", WidthClass: "col-2xl"},
	}

	var rows []types.TableRow

	if deps.GetWorkspaceUserListPageData != nil {
		// Filter client-side by workspace_id using a string filter.
		resp, err := deps.GetWorkspaceUserListPageData(ctx, &workspaceuserpb.GetWorkspaceUserListPageDataRequest{
			Filters: &commonpb.FilterRequest{
				Filters: []*commonpb.TypedFilter{
					{
						Field: "workspace_id",
						FilterType: &commonpb.TypedFilter_StringFilter{
							StringFilter: &commonpb.StringFilter{
								Value:         workspaceID,
								Operator:      commonpb.StringOperator_STRING_EQUALS,
								CaseSensitive: true,
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Printf("Failed to list workspace users for workspace %s: %v", workspaceID, err)
		} else {
			for _, wu := range resp.GetWorkspaceUserList() {
				if wu.GetWorkspaceId() != workspaceID {
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

				viewURL := ""
				if deps.WorkspaceUserDetailURL != "" {
					viewURL = route.ResolveURL(deps.WorkspaceUserDetailURL, "id", wu.GetId())
				}

				cells := []types.TableCell{
					{Type: "text", Value: userName},
					{Type: "text", Value: email},
					{Type: "text", Value: fmt.Sprintf("%d", roleCount)},
					{Type: "badge", Value: statusValue, Variant: statusVariant},
				}
				actions := []types.TableAction{}
				if viewURL != "" {
					actions = append(actions, types.TableAction{
						Type:   "view",
						Label:  deps.CommonLabels.Actions.View,
						Action: "view",
						Href:   viewURL,
					})
				}
				rows = append(rows, types.TableRow{
					ID:    wu.GetId(),
					Cells: cells,
					DataAttrs: map[string]string{
						"user_name": userName,
						"email":     email,
						"testid":    "workspace-user-row-" + wu.GetId(),
					},
					Actions: actions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	addURL := ""
	addLabel := tl.Users + " — Add"
	if deps.Labels.Detail.Users.AddButton != "" {
		addLabel = deps.Labels.Detail.Users.AddButton
	}

	tc := &types.TableConfig{
		ID:                   "workspace-users-table",
		Columns:              columns,
		Rows:                 rows,
		Labels:               deps.TableLabels,
		ShowSearch:           false,
		ShowActions:          true,
		ShowSort:             false,
		ShowColumns:          false,
		ShowDensity:          false,
		ShowEntries:          false,
		DefaultSortColumn:    "user_name",
		DefaultSortDirection: "asc",
		EmptyState: types.TableEmptyState{
			Title:   "No users",
			Message: "No users have been added to this workspace yet.",
		},
	}
	if addURL != "" || deps.WorkspaceUserAddURL != "" {
		resolvedAddURL := deps.WorkspaceUserAddURL
		if resolvedAddURL != "" {
			resolvedAddURL = resolvedAddURL + "?workspace_id=" + workspaceID
		}
		tc.PrimaryAction = &types.PrimaryAction{
			Label:     addLabel,
			ActionURL: resolvedAddURL,
			Icon:      "icon-plus",
		}
	}
	types.ApplyTableSettings(tc)
	return tc
}
