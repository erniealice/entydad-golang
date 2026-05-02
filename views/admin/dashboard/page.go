// Package dashboard is the view layer for the entydad admin app dashboard
// (Phase 4b of the Pyeza dashboard plan).
//
// Admin is a composite app spanning permission/role/workspace_user/
// workspace_user_role. The view receives a single GetDashboardData
// callback that the container builds by composing the
// GetAdminDashboardPageDataUseCase.
package dashboard

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/entydad-golang"
	workspaceuserrolepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user_role"
)

// RolePermissionCount mirrors the use case's row type.
type RolePermissionCount struct {
	RoleID          string
	RoleName        string
	PermissionCount int64
}

// AdminDashboardData is the response shape consumed by the view.
type AdminDashboardData struct {
	WorkspaceUsers      int64
	Roles               int64
	Permissions         int64
	RecentRoleChanges7d int64
	UsersPerRole        map[string]int64
	TopRolesByPerms     []RolePermissionCount
	RecentAssignments   []*workspaceuserrolepb.WorkspaceUserRole
	// RoleNamesByID is an optional id→name map so the donut chart and
	// activity list can show role names instead of opaque IDs. The
	// container assembles it alongside the use case response.
	RoleNamesByID map[string]string
	// UserLabelsByID is an optional id→display map (e.g. email) for
	// recent-assignment row descriptions.
	UserLabelsByID map[string]string
}

// Routes holds the cross-entity route URLs the admin dashboard's quick
// actions and link buttons resolve to. Sourced from the orchestrator's
// composed entydad.UserRoutes / WorkspaceRoutes / etc.
type Routes struct {
	DashboardURL          string
	NewUserURL            string // workspace_user.AddURL
	NewWorkspaceURL       string // workspace.AddURL
	AssignRoleURL         string // workspace_user_role.AddURL
	AuditLogURL           string // optional — defaults to "#"
	PermissionListURL     string // permission.list (active)
	RoleListURL           string // role.list
	WorkspaceListURL      string // workspace.list (active)
	WorkspaceUserListURL  string // workspace_user.list (active)
}

// Deps holds view dependencies.
type Deps struct {
	DashboardLabels  entydad.DashboardLabels
	Dashboard        entydad.AdminDashboardLabels
	Routes           Routes
	CommonLabels     pyeza.CommonLabels
	GetDashboardData func(ctx context.Context) (*AdminDashboardData, error)
}

// PageData holds the data for the admin dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the admin dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Dashboard

		var data *AdminDashboardData
		if deps.GetDashboardData != nil {
			d, err := deps.GetDashboardData(ctx)
			if err != nil {
				log.Printf("admin dashboard: failed to load page data: %v", err)
			}
			data = d
		}
		if data == nil {
			data = &AdminDashboardData{}
		}

		// Stats — 4 cards.
		stats := []types.StatCardData{
			{
				Icon: "icon-users", Value: strconv.FormatInt(data.WorkspaceUsers, 10),
				Label: l.WorkspaceUsers, Color: "terracotta",
				TestID: "admin-stat-workspace-users",
			},
			{
				Icon: "icon-shield", Value: strconv.FormatInt(data.Roles, 10),
				Label: l.Roles, Color: "sage",
				TestID: "admin-stat-roles",
			},
			{
				Icon: "icon-key", Value: strconv.FormatInt(data.Permissions, 10),
				Label: l.Permissions, Color: "navy",
				TestID: "admin-stat-permissions",
			},
			{
				Icon: "icon-trending-up", Value: strconv.FormatInt(data.RecentRoleChanges7d, 10),
				Label: l.RecentRoleChanges, Color: "amber",
				TestID: "admin-stat-recent-changes",
			},
		}

		// Donut chart — users per role. Labels resolved through the
		// role-name map when provided, falling back to the role_id.
		donutLabels, donutValues := projectUsersPerRole(data.UsersPerRole, data.RoleNamesByID)
		donutChart := &types.ChartData{
			Labels: donutLabels,
			Series: []types.ChartSeries{{
				Name:   l.UsersPerRole,
				Values: donutValues,
				Color:  "terracotta",
			}},
			Donut: true,
		}
		donutChart.AutoScale()

		// Table widget — top roles by permission count (Type=custom).
		topRolesHTML := buildTopRolesHTML(data.TopRolesByPerms, deps.Routes.RoleListURL, l)

		// Activity list — recent role assignments.
		recentItems := buildRecentAssignmentsList(data.RecentAssignments, data.RoleNamesByID, data.UserLabelsByID, l)

		auditURL := deps.Routes.AuditLogURL
		if auditURL == "" {
			auditURL = "#"
		}
		newUserURL := deps.Routes.NewUserURL
		if newUserURL == "" {
			newUserURL = "#"
		}
		newWorkspaceURL := deps.Routes.NewWorkspaceURL
		if newWorkspaceURL == "" {
			newWorkspaceURL = "#"
		}
		assignRoleURL := deps.Routes.AssignRoleURL
		if assignRoleURL == "" {
			assignRoleURL = "#"
		}

		dash := types.DashboardData{
			Title:    deps.DashboardLabels.AdminTitle,
			Icon:     "icon-shield",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-user-plus", Label: l.QuickNewUser, Href: newUserURL, Variant: "primary", TestID: "admin-action-new-user"},
				{Icon: "icon-briefcase", Label: l.QuickNewWorkspace, Href: newWorkspaceURL, TestID: "admin-action-new-workspace"},
				{Icon: "icon-shield", Label: l.QuickAssignRole, Href: assignRoleURL, TestID: "admin-action-assign-role"},
				{Icon: "icon-clipboard", Label: l.QuickAuditLog, Href: auditURL, TestID: "admin-action-audit-log"},
			},
			Stats: stats,
			Widgets: []types.DashboardWidget{
				{
					ID: "users-per-role", Title: l.UsersPerRole,
					Type: "chart", ChartKind: "donut",
					ChartData: donutChart, Span: 1,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-users",
						Title: l.UsersPerRole,
					},
				},
				{
					ID: "top-roles", Title: l.RolesByPermissionCount,
					Type: "custom", Span: 2,
					Custom: topRolesHTML,
				},
				{
					ID: "recent-changes", Title: l.RecentRoleChangesList,
					Type: "list", Span: 3,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.WorkspaceUserListURL},
					},
					ListItems: recentItems,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-clock",
						Title: l.RecentRoleChangesList,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.DashboardLabels.AdminTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "admin",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.DashboardLabels.AdminTitle,
				HeaderIcon:   "icon-shield",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "admin-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("admin-dashboard", pageData)
	})
}

// projectUsersPerRole sorts the role→count map descending by count and
// resolves labels through the optional id→name map.
func projectUsersPerRole(byRole map[string]int64, names map[string]string) ([]string, []float64) {
	if len(byRole) == 0 {
		return []string{"-"}, []float64{0}
	}
	type pair struct {
		label string
		count int64
	}
	pairs := make([]pair, 0, len(byRole))
	for roleID, n := range byRole {
		label := roleID
		if names != nil {
			if name, ok := names[roleID]; ok && name != "" {
				label = name
			}
		}
		pairs = append(pairs, pair{label: label, count: n})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].label < pairs[j].label
	})
	labels := make([]string, len(pairs))
	values := make([]float64, len(pairs))
	for i, p := range pairs {
		labels[i] = p.label
		values[i] = float64(p.count)
	}
	return labels, values
}

func buildTopRolesHTML(rows []RolePermissionCount, listURL string, l entydad.AdminDashboardLabels) template.HTML {
	if len(rows) == 0 {
		return template.HTML(`<div class="empty-state empty-state--inline" data-testid="admin-top-roles-empty">` +
			template.HTMLEscapeString(l.RolesByPermissionCount) +
			`</div>`)
	}
	var sb strings.Builder
	sb.WriteString(`<table class="dashboard-mini-table" data-testid="admin-top-roles-table"><thead><tr>`)
	sb.WriteString(`<th>` + template.HTMLEscapeString(l.ColumnRole) + `</th>`)
	sb.WriteString(`<th class="num">` + template.HTMLEscapeString(l.ColumnPermissionCount) + `</th>`)
	sb.WriteString(`</tr></thead><tbody>`)
	for _, r := range rows {
		rowID := template.HTMLEscapeString(r.RoleID)
		sb.WriteString(`<tr data-testid="admin-top-roles-row-` + rowID + `">`)
		sb.WriteString(`<td>` + template.HTMLEscapeString(r.RoleName) + `</td>`)
		sb.WriteString(`<td class="num">` + strconv.FormatInt(r.PermissionCount, 10) + `</td>`)
		sb.WriteString(`</tr>`)
	}
	sb.WriteString(`</tbody></table>`)
	return template.HTML(sb.String())
}

func buildRecentAssignmentsList(
	assignments []*workspaceuserrolepb.WorkspaceUserRole,
	roleNames, userLabels map[string]string,
	l entydad.AdminDashboardLabels,
) []types.ActivityItem {
	if len(assignments) == 0 {
		return nil
	}
	items := make([]types.ActivityItem, 0, len(assignments))
	for i, a := range assignments {
		roleName := a.GetRoleId()
		if roleNames != nil {
			if n, ok := roleNames[a.GetRoleId()]; ok && n != "" {
				roleName = n
			}
		}
		userLabel := a.GetWorkspaceUserId()
		if userLabels != nil {
			if u, ok := userLabels[a.GetWorkspaceUserId()]; ok && u != "" {
				userLabel = u
			}
		}
		t := ""
		if a.GetDateCreated() != 0 {
			t = formatRelative(time.UnixMilli(a.GetDateCreated()))
		}
		items = append(items, types.ActivityItem{
			IconName:    "icon-shield",
			IconVariant: "client",
			Title:       fmt.Sprintf("%s — %s", l.RoleAssigned, roleName),
			Description: userLabel,
			Time:        t,
			TestID:      fmt.Sprintf("admin-activity-%d", i+1),
		})
	}
	return items
}

func formatRelative(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	delta := time.Since(t)
	switch {
	case delta < time.Minute:
		return "just now"
	case delta < time.Hour:
		return fmt.Sprintf("%dm ago", int(delta.Minutes()))
	case delta < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(delta.Hours()))
	case delta < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(delta.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
}
