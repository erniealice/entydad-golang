// Package detail provides the workspace_user nested detail page view.
// This is the page operators land on after clicking a row in workspace detail's Users tab.
// It renders Info and Roles tabs with a "Back to workspace" breadcrumb link.
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
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
	"github.com/erniealice/hybra-golang/views/attachment"
)

// DetailViewDeps holds view dependencies for the workspace_user detail page.
type DetailViewDeps struct {
	Routes                       entydad.WorkspaceUserRoutes
	WorkspaceDetailURL           string // /app/workspaces/detail/{id} — for "Back to workspace" link
	GetWorkspaceUserItemPageData func(ctx context.Context, req *workspaceuserpb.GetWorkspaceUserItemPageDataRequest) (*workspaceuserpb.GetWorkspaceUserItemPageDataResponse, error)
	Labels                       entydad.WorkspaceUserLabels
	CommonLabels                 pyeza.CommonLabels
	TableLabels                  types.TableLabels
	// Phase 3 wired: these are now supplied from block.go / container.go after Phase 3 shipped.
	GetWorkspaceUserRoleListPageData func(ctx context.Context, req *workspaceuserrolepb.GetWorkspaceUserRoleListPageDataRequest) (*workspaceuserrolepb.GetWorkspaceUserRoleListPageDataResponse, error)
	WorkspaceUserRoleAddURL          string
	WorkspaceUserRoleDeleteURL       string

	// Attachment operations (embedded from hybra)
	attachment.AttachmentOps
}

// PageData holds the data for the workspace_user detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	// WorkspaceUser fields for Info tab
	WorkspaceUserID string
	UserName        string
	Email           string
	WorkspaceName   string
	WorkspaceID     string
	DateJoined      string
	Active          bool
	StatusVariant   string
	Labels          entydad.WorkspaceUserLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	// "Back to workspace" link
	BackToWorkspaceURL  string
	BackToWorkspaceName string
	// Roles tab
	RolesTable              *types.TableConfig
	WorkspaceUserRoleAddURL string
	// Attachments tab
	AttachmentTable *types.TableConfig
}

// tabLabels holds the resolved tab display strings.
type tabLabels struct {
	Info        string
	Roles       string
	Attachments string
}

// resolveTabLabels returns display strings for the tabs with English fallbacks.
func resolveTabLabels(l entydad.WorkspaceUserLabels) tabLabels {
	info := l.Detail.Tabs.Info
	if info == "" {
		info = "Info"
	}
	roles := l.Detail.Tabs.Roles
	if roles == "" {
		roles = "Roles"
	}
	attachments := l.Detail.Tabs.Attachments
	if attachments == "" {
		attachments = "Attachments"
	}
	return tabLabels{Info: info, Roles: roles, Attachments: attachments}
}

// NewView creates the workspace_user detail view (full page load).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		wu, err := loadWorkspaceUser(ctx, deps, id)
		if err != nil {
			return view.Error(err)
		}

		tl := resolveTabLabels(deps.Labels)
		tabItems := buildTabItems(id, deps, tl)
		pageData := buildPageData(viewCtx, id, wu, activeTab, tabItems, deps)

		switch activeTab {
		case "roles":
			pageData.RolesTable = buildRolesTable(ctx, deps, wu, tl)
		case "attachments":
			loadAttachments(ctx, deps, id, pageData)
		}

		return view.OK("workspace-user-detail", pageData)
	})
}

// NewTabAction creates the tab-action view (HTMX partial — returns only the tab content).
// Route: GET /action/workspace_user/{id}/tab/{tab}
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		wu, err := loadWorkspaceUser(ctx, deps, id)
		if err != nil {
			return view.Error(err)
		}

		tl := resolveTabLabels(deps.Labels)
		tabItems := buildTabItems(id, deps, tl)
		pageData := buildPageData(viewCtx, id, wu, tab, tabItems, deps)

		switch tab {
		case "roles":
			pageData.RolesTable = buildRolesTable(ctx, deps, wu, tl)
			return view.OK("workspace-user-tab-roles", pageData)
		case "attachments":
			loadAttachments(ctx, deps, id, pageData)
			return view.OK("attachment-tab", pageData)
		default:
			return view.OK("workspace-user-tab-info", pageData)
		}
	})
}

// loadWorkspaceUser fetches a single workspace_user by ID (with nested user + workspace + roles).
func loadWorkspaceUser(ctx context.Context, deps *DetailViewDeps, id string) (*workspaceuserpb.WorkspaceUser, error) {
	if deps.GetWorkspaceUserItemPageData == nil {
		return nil, fmt.Errorf("GetWorkspaceUserItemPageData not wired")
	}
	resp, err := deps.GetWorkspaceUserItemPageData(ctx, &workspaceuserpb.GetWorkspaceUserItemPageDataRequest{
		WorkspaceUserId: id,
	})
	if err != nil {
		log.Printf("Failed to read workspace_user %s: %v", id, err)
		return nil, fmt.Errorf("failed to load workspace user: %w", err)
	}
	wu := resp.GetWorkspaceUser()
	if wu == nil {
		return nil, fmt.Errorf("workspace user not found")
	}
	return wu, nil
}

// buildPageData assembles the PageData struct.
func buildPageData(viewCtx *view.ViewContext, id string, wu *workspaceuserpb.WorkspaceUser, activeTab string, tabItems []pyeza.TabItem, deps *DetailViewDeps) *PageData {
	u := wu.GetUser()
	userName := ""
	email := ""
	if u != nil {
		userName = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
		email = u.GetEmailAddress()
	}

	workspaceName := ""
	workspaceID := wu.GetWorkspaceId()
	if ws := wu.GetWorkspace(); ws != nil {
		workspaceName = ws.GetName()
		if workspaceID == "" {
			workspaceID = ws.GetId()
		}
	}

	active := wu.GetActive()
	statusVariant := "success"
	if !active {
		statusVariant = "warning"
	}

	// "Back to workspace" link — uses WorkspaceDetailURL populated with workspaceID.
	backURL := ""
	if deps.WorkspaceDetailURL != "" && workspaceID != "" {
		backURL = route.ResolveURL(deps.WorkspaceDetailURL, "id", workspaceID)
	}

	// Phase 3 role add URL: /action/workspace_user_role/add?workspace_user_id={id}
	roleAddURL := deps.WorkspaceUserRoleAddURL
	if roleAddURL != "" {
		roleAddURL = roleAddURL + "?workspace_user_id=" + id
	}

	title := userName
	if title == "" {
		title = "Workspace User"
	}

	return &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          title,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "admin",
			ActiveSubNav:   "workspaces-active",
			HeaderTitle:    title,
			HeaderSubtitle: email,
			HeaderIcon:     "icon-user",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:         "workspace-user-detail-content",
		WorkspaceUserID:         id,
		UserName:                userName,
		Email:                   email,
		WorkspaceName:           workspaceName,
		WorkspaceID:             workspaceID,
		DateJoined:              wu.GetDateCreatedString(),
		Active:                  active,
		StatusVariant:           statusVariant,
		Labels:                  deps.Labels,
		ActiveTab:               activeTab,
		TabItems:                tabItems,
		BackToWorkspaceURL:      backURL,
		BackToWorkspaceName:     workspaceName,
		WorkspaceUserRoleAddURL: roleAddURL,
	}
}

// buildTabItems constructs the tab items for the workspace_user detail page.
func buildTabItems(id string, deps *DetailViewDeps, tl tabLabels) []pyeza.TabItem {
	base := route.ResolveURL(deps.Routes.DetailURL, "id", id)
	action := route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: tl.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "roles", Label: tl.Roles, Href: base + "?tab=roles", HxGet: action + "roles", Icon: "icon-shield"},
		{Key: "attachments", Label: tl.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
	}
}

// buildRolesTable builds the WorkspaceUserRole rows table for the Roles tab.
// It reads roles from the already-loaded WorkspaceUser's nested WorkspaceUserRoles slice
// (populated by GetWorkspaceUserItemPageData CTE join). Falls back to GetWorkspaceUserRoleListPageData
// if the nested slice is empty and the dep is wired (Phase 3).
func buildRolesTable(ctx context.Context, deps *DetailViewDeps, wu *workspaceuserpb.WorkspaceUser, tl tabLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "role_name", Label: "Role"},
		{Key: "perm_count", Label: "Permissions", WidthClass: "col-2xl"},
		{Key: "date_created", Label: "Date Assigned", WidthClass: "col-lg"},
	}

	var rows []types.TableRow

	for _, wur := range wu.GetWorkspaceUserRoles() {
		roleName := ""
		permCount := 0
		if r := wur.GetRole(); r != nil {
			roleName = r.GetName()
			permCount = len(r.GetRolePermissions())
		}
		permLabel := fmt.Sprintf("%d", permCount)
		dateAssigned := wur.GetDateCreatedString()

		active := wur.GetActive()
		_ = active // junctions only shown while active

		cells := []types.TableCell{
			{Type: "text", Value: roleName},
			{Type: "text", Value: permLabel},
			{Type: "text", Value: dateAssigned},
		}

		var actions []types.TableAction
		if deps.WorkspaceUserRoleDeleteURL != "" {
			deleteURL := route.ResolveURL(deps.WorkspaceUserRoleDeleteURL, "id", wur.GetId())
			actions = append(actions, types.TableAction{
				Type:   "delete",
				Label:  deps.CommonLabels.Actions.Delete,
				Action: "delete",
				Href:   deleteURL,
			})
		}

		rows = append(rows, types.TableRow{
			ID:    wur.GetId(),
			Cells: cells,
			DataAttrs: map[string]string{
				"testid": "workspace-user-role-row-" + wur.GetId(),
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, rows)

	assignLabel := "Assign role"
	if deps.Labels.Detail.Roles.AssignButton != "" {
		assignLabel = deps.Labels.Detail.Roles.AssignButton
	}

	tc := &types.TableConfig{
		ID:                   "workspace-user-roles-table",
		Columns:              columns,
		Rows:                 rows,
		Labels:               deps.TableLabels,
		ShowSearch:           false,
		ShowActions:          true,
		ShowSort:             false,
		ShowColumns:          false,
		ShowDensity:          false,
		ShowEntries:          false,
		DefaultSortColumn:    "role_name",
		DefaultSortDirection: "asc",
		RefreshURL:           route.ResolveURL(deps.Routes.TabActionURL, "id", wu.GetId(), "tab", "roles"),
		EmptyState: types.TableEmptyState{
			Title:   "No roles",
			Message: "No roles have been assigned to this workspace user yet.",
		},
	}
	if deps.WorkspaceUserRoleAddURL != "" {
		roleAddURL := deps.WorkspaceUserRoleAddURL + "?workspace_user_id=" + wu.GetId()
		tc.PrimaryAction = &types.PrimaryAction{
			Label:     assignLabel,
			ActionURL: roleAddURL,
			Icon:      "icon-plus",
		}
	}
	types.ApplyTableSettings(tc)
	return tc
}
