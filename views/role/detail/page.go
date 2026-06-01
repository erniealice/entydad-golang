package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	roleusers "github.com/erniealice/entydad-golang/views/role/users"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
	rolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/role"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	ReadRole             func(ctx context.Context, req *rolepb.ReadRoleRequest) (*rolepb.ReadRoleResponse, error)
	RoleGetItemPageData  func(ctx context.Context, req *rolepb.GetRoleItemPageDataRequest) (*rolepb.GetRoleItemPageDataResponse, error)
	GetUsersByRoleID     func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error)
	Routes               entydad.RoleRoutes
	Labels               entydad.RoleLabels
	SharedLabels         entydad.SharedLabels
	RolePermissionLabels entydad.RolePermissionLabels
	RoleUserLabels       entydad.RoleUserLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels

	// Attachment operations (embedded from hybra)
	attachment.AttachmentOps

	// Audit log operations (embedded from hybra)
	auditlog.AuditOps
}

// PageData holds the data for the role detail page.
type PageData struct {
	types.PageData
	ContentTemplate  string
	Labels           entydad.RoleLabels
	ActiveTab        string
	TabItems         []pyeza.TabItem
	ID               string
	RoleName         string
	RoleDescription  string
	RoleColor        string
	RoleStatus       string
	StatusVariant    string
	PermissionsTable *types.TableConfig
	UsersTable       *types.TableConfig
	AttachmentTable  *types.TableConfig
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView creates the role detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("role", "read") {
			return view.Forbidden("role:read")
		}

		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "role-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("role-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/roles/{id}/tab/{tab}
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("role", "read") {
			return view.Forbidden("role:read")
		}

		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		// Return only the tab partial template
		templateName := "role-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// buildPageData loads role data and builds the PageData for the given active tab.
func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	perms := view.GetUserPermissions(ctx)
	resp, err := deps.ReadRole(ctx, &rolepb.ReadRoleRequest{
		Data: &rolepb.Role{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read role %s: %v", id, err)
		return nil, fmt.Errorf("failed to load role: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("role not found")
	}
	role := data[0]

	roleName := role.GetName()
	roleDesc := role.GetDescription()
	roleColor := role.GetColor()

	roleStatus := "active"
	if !role.GetActive() {
		roleStatus = "inactive"
	}
	statusVariant := "success"
	if roleStatus == "inactive" {
		statusVariant = "warning"
	}

	// Get counts for tab badges
	permCount := 0
	userCount := 0

	if deps.RoleGetItemPageData != nil {
		itemResp, err := deps.RoleGetItemPageData(ctx, &rolepb.GetRoleItemPageDataRequest{
			RoleId: id,
		})
		if err != nil {
			log.Printf("Failed to get role item page data for %s: %v", id, err)
		} else {
			permCount = len(itemResp.GetRole().GetRolePermissions())
		}
	}

	if deps.GetUsersByRoleID != nil {
		users, err := deps.GetUsersByRoleID(ctx, id)
		if err != nil {
			log.Printf("Failed to get users for role %s: %v", id, err)
		} else {
			userCount = len(users)
		}
	}

	tabItems := buildTabItems(id, deps.Labels, deps.Routes, permCount, userCount)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          roleName,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "user",
			ActiveSubNav:   "role",
			HeaderTitle:    roleName,
			HeaderSubtitle: roleDesc,
			HeaderIcon:     "icon-shield",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "role-detail-content",
		Labels:          deps.Labels,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		RoleName:        roleName,
		RoleDescription: roleDesc,
		RoleColor:       roleColor,
		RoleStatus:      roleStatus,
		StatusVariant:   statusVariant,
	}

	// Load tab-specific data
	switch activeTab {
	case "permissions":
		tableConfig, err := buildPermissionsTable(ctx, deps, id, perms)
		if err != nil {
			log.Printf("Failed to build permissions table for role %s: %v", id, err)
		} else {
			pageData.PermissionsTable = tableConfig
		}
	case "users":
		tableConfig, err := buildUsersTable(ctx, deps, id, perms)
		if err != nil {
			log.Printf("Failed to build users table for role %s: %v", id, err)
		} else {
			pageData.UsersTable = tableConfig
		}
	case "attachments":
		if deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
			if err != nil {
				log.Printf("Failed to list attachments: %v", err)
			}
			var items []*attachmentpb.Attachment
			if resp != nil {
				items = resp.GetData()
			}
			pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
		}
	case "audit-history":
		if deps.ListAuditHistory != nil {
			cursor := viewCtx.Request.URL.Query().Get("cursor")
			auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
				EntityType:  "role",
				EntityID:    id,
				Limit:       20,
				CursorToken: cursor,
			})
			if err != nil {
				log.Printf("Failed to load audit history: %v", err)
			}
			if auditResp != nil {
				pageData.AuditEntries = auditResp.Entries
				pageData.AuditHasNext = auditResp.HasNext
				pageData.AuditNextCursor = auditResp.NextCursor
			}
		}
		pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
	}

	return pageData, nil
}

func buildTabItems(id string, labels entydad.RoleLabels, routes entydad.RoleRoutes, permCount, userCount int) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: labels.Detail.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "permissions", Label: labels.Detail.Tabs.Permissions, Href: base + "?tab=permissions", HxGet: action + "permissions", Icon: "icon-key", Count: permCount, Disabled: false},
		{Key: "users", Label: labels.Detail.Tabs.Users, Href: base + "?tab=users", HxGet: action + "users", Icon: "icon-user", Count: userCount, Disabled: false},
		{Key: "attachments", Label: labels.Detail.AttachmentsTab, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip", Count: 0, Disabled: false},
		{Key: "audit-history", Label: func() string {
			if labels.Detail.AuditHistoryTab != "" {
				return labels.Detail.AuditHistoryTab
			}
			return "History"
		}(), Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

// ---------------------------------------------------------------------------
// Permissions tab table
// ---------------------------------------------------------------------------

func buildPermissionsTable(ctx context.Context, deps *DetailViewDeps, roleID string, perms *types.UserPermissions) (*types.TableConfig, error) {
	if deps.RoleGetItemPageData == nil {
		return nil, fmt.Errorf("RoleGetItemPageData not available")
	}

	resp, err := deps.RoleGetItemPageData(ctx, &rolepb.GetRoleItemPageDataRequest{
		RoleId: roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load role permissions: %w", err)
	}

	role := resp.GetRole()
	l := deps.RolePermissionLabels

	columns := []types.TableColumn{
		{Key: "permissionName", Label: l.Columns.PermissionName},
		{Key: "code", Label: l.Columns.Code},
		{Key: "type", Label: l.Columns.Type, WidthClass: "col-2xl"},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, WidthClass: "col-6xl"},
	}

	rows := []types.TableRow{}
	for _, rp := range role.GetRolePermissions() {
		perm := rp.GetPermission()
		if perm == nil {
			continue
		}

		rpID := rp.GetId()
		permName := perm.GetName()
		permCode := perm.GetPermissionCode()
		permType := deps.SharedLabels.Badges.Allow
		dateAssigned := rp.GetDateCreatedString()

		pt := perm.GetPermissionType()
		if pt == permissionpb.PermissionType_PERMISSION_TYPE_DENY {
			permType = deps.SharedLabels.Badges.Deny
		}
		typeVariant := "success"
		if pt == permissionpb.PermissionType_PERMISSION_TYPE_DENY {
			typeVariant = "danger"
		}

		actions := []types.TableAction{
			{
				Type: "delete", Label: l.Actions.Remove, Action: "delete",
				URL:            route.ResolveURL(deps.Routes.DetailPermissionsRemoveURL, "id", roleID),
				ItemName:       permName,
				ConfirmTitle:   l.Actions.Remove,
				ConfirmMessage: fmt.Sprintf(deps.SharedLabels.Confirm.Remove, permName),
				Disabled:       !perms.Can("role_permission", "delete"), DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "role_permission:delete"),
			},
		}

		rows = append(rows, types.TableRow{
			ID: rpID,
			Cells: []types.TableCell{
				{Type: "text", Value: permName},
				{Type: "text", Value: permCode},
				{Type: "badge", Value: permType, Variant: typeVariant},
				{Type: "text", Value: dateAssigned},
			},
			DataAttrs: map[string]string{
				"permissionName": permName,
				"code":           permCode,
				"type":           permType,
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "role-permissions-table",
		RefreshURL:           route.ResolveURL(deps.Routes.DetailPermissionsTableURL, "id", roleID),
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "permissionName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AssignPermission,
			ActionURL:       route.ResolveURL(deps.Routes.DetailPermissionsAssignURL, "id", roleID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("role_permission", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "role_permission:create"),
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

// ---------------------------------------------------------------------------
// Users tab table
// ---------------------------------------------------------------------------

func buildUsersTable(ctx context.Context, deps *DetailViewDeps, roleID string, perms *types.UserPermissions) (*types.TableConfig, error) {
	var users []roleusers.UserByRole
	if deps.GetUsersByRoleID != nil {
		var err error
		users, err = deps.GetUsersByRoleID(ctx, roleID)
		if err != nil {
			log.Printf("Failed to get users for role %s: %v", roleID, err)
			// Continue with empty table
		}
	}

	l := deps.RoleUserLabels

	columns := []types.TableColumn{
		{Key: "userName", Label: l.Columns.UserName},
		{Key: "email", Label: l.Columns.Email},
		{Key: "dateAssigned", Label: l.Columns.DateAssigned, WidthClass: "col-6xl"},
	}

	rows := []types.TableRow{}
	for _, u := range users {
		actions := []types.TableAction{
			{
				Type: "delete", Label: l.Actions.Remove, Action: "delete",
				URL:            route.ResolveURL(deps.Routes.UsersRemoveURL, "id", roleID),
				ItemName:       u.UserName,
				ConfirmTitle:   l.Actions.Remove,
				ConfirmMessage: fmt.Sprintf(deps.SharedLabels.Confirm.Remove, u.UserName),
				Disabled:       !perms.Can("workspace_user_role", "delete"), DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "workspace_user_role:delete"),
			},
		}

		rows = append(rows, types.TableRow{
			ID: u.WorkspaceUserRoleID,
			Cells: []types.TableCell{
				{Type: "text", Value: u.UserName},
				{Type: "text", Value: u.Email},
				{Type: "text", Value: u.DateAssigned},
			},
			DataAttrs: map[string]string{
				"userName": u.UserName,
				"email":    u.Email,
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "role-users-table",
		RefreshURL:           route.ResolveURL(deps.Routes.UsersTableURL, "id", roleID),
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "userName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AssignUser,
			ActionURL:       route.ResolveURL(deps.Routes.UsersAssignURL, "id", roleID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("workspace_user_role", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "workspace_user_role:create"),
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}
