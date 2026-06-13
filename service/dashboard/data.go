// Package dashboard provides workspace-scoped data-fetching builders for the
// admin and user dashboard views. Moved from service-admin's composition layer
// (2026-06-13) because entydad owns the dashboard views and types.
//
// The builders accept a *sql.DB and an optional WorkspaceIDFromCtx extractor
// so that the domain package stays free of consumer/middleware imports.
package dashboard

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"

	roleusers "github.com/erniealice/entydad-golang/domain/entity/identity/role/users"
	userdashboard "github.com/erniealice/entydad-golang/domain/entity/identity/user/dashboard"
)

// WorkspaceIDExtractor resolves the workspace ID from context. Callers
// typically pass identity.Must(ctx).WorkspaceID from espyna.
type WorkspaceIDExtractor func(ctx context.Context) string

// IconRenderer converts an icon template name to rendered HTML. Callers
// typically pass pyeza.HTMLRenderer.RenderIcon.
type IconRenderer func(name string) template.HTML

// BuildGetUsersByRoleID creates a closure that queries for users assigned to a
// specific role. Joins workspace_user_role -> workspace_user -> user to get
// user info. The result is workspace-scoped via the workspaceID from ctx.
func BuildGetUsersByRoleID(sqlDB *sql.DB, wsID WorkspaceIDExtractor) func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error) {
	if sqlDB == nil {
		return nil
	}

	return func(ctx context.Context, roleID string) ([]roleusers.UserByRole, error) {
		// Multi-tenancy: scope to the caller's workspace. workspace_user_role has no
		// workspace_id column (it is keyed via its parent workspace_user), so the
		// predicate is applied on the workspace_user join (wu.workspace_id), exactly
		// as the dashboard stats/activity/chart queries below do. workspaceID is
		// sourced from context (the workspace_path middleware sets it) — never from
		// caller input. Fail closed: an empty workspaceID yields no rows.
		workspaceID := wsID(ctx)
		query := `
			SELECT wur.id, wu.id, wu.user_id,
			       COALESCE(u.first_name || ' ' || u.last_name, u.email_address) as user_name,
			       COALESCE(u.email_address, '') as email,
			       COALESCE(TO_CHAR(wur.date_created, 'Mon DD, YYYY'), '') as date_assigned
			FROM workspace_user_role wur
			JOIN workspace_user wu ON wur.workspace_user_id = wu.id
			JOIN "user" u ON wu.user_id = u.id
			WHERE wur.role_id = $1 AND wur.active = true AND wu.workspace_id = $2
			ORDER BY u.first_name, u.last_name
		`
		rows, err := sqlDB.QueryContext(ctx, query, roleID, workspaceID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var users []roleusers.UserByRole
		for rows.Next() {
			var u roleusers.UserByRole
			if err := rows.Scan(&u.WorkspaceUserRoleID, &u.WorkspaceUserID, &u.UserID, &u.UserName, &u.Email, &u.DateAssigned); err != nil {
				continue
			}
			users = append(users, u)
		}
		return users, rows.Err()
	}
}

// BuildGetDashboardData creates a closure that queries the database for
// dashboard stats, recent activity, and chart data. The renderIcon callback
// is used to render SVG icons for activity items.
func BuildGetDashboardData(sqlDB *sql.DB, wsID WorkspaceIDExtractor, renderIcon IconRenderer) func(ctx context.Context) (*userdashboard.DashboardData, error) {
	if sqlDB == nil {
		return nil
	}

	return func(ctx context.Context) (*userdashboard.DashboardData, error) {
		workspaceID := wsID(ctx)

		data := &userdashboard.DashboardData{}

		// Stats: count users and roles scoped to the current workspace
		var totalUsers, activeUsers, inactiveUsers, totalRoles int
		// totalUsers = all workspace members; activeUsers = active subset;
		// inactiveUsers = inactive subset (total = active + inactive). Previously the
		// first subquery duplicated the active filter, making totalUsers == activeUsers.
		row := sqlDB.QueryRowContext(ctx, `SELECT
			COALESCE((SELECT COUNT(*) FROM workspace_user WHERE workspace_id = $1), 0),
			COALESCE((SELECT COUNT(*) FROM workspace_user WHERE active = true AND workspace_id = $1), 0),
			COALESCE((SELECT COUNT(*) FROM workspace_user WHERE active = false AND workspace_id = $1), 0),
			COALESCE((SELECT COUNT(*) FROM role WHERE active = true AND workspace_id = $1), 0)`,
			workspaceID)
		if err := row.Scan(&totalUsers, &activeUsers, &inactiveUsers, &totalRoles); err != nil {
			log.Printf("Dashboard stats query error: %v", err)
		}
		data.Stats = userdashboard.DashboardStats{
			TotalUsers:    totalUsers,
			ActiveUsers:   activeUsers,
			InactiveUsers: inactiveUsers,
			TotalRoles:    totalRoles,
		}

		// Recent activity: synthesized from user/role created/modified timestamps
		data.RecentActivity = buildRecentActivity(ctx, sqlDB, workspaceID, renderIcon)

		// Chart data: user creations per month for last 12 months
		data.Chart = buildChartData(ctx, sqlDB, workspaceID)

		return data, nil
	}
}

func buildRecentActivity(ctx context.Context, sqlDB *sql.DB, workspaceID string, renderIcon IconRenderer) []userdashboard.ActivityItem {
	query := `
		(SELECT 'user_created' as event_type, u.first_name || ' ' || u.last_name as name, wu.date_created as event_date
		 FROM workspace_user wu JOIN "user" u ON wu.user_id = u.id
		 WHERE wu.date_created IS NOT NULL AND wu.workspace_id = $1
		 ORDER BY wu.date_created DESC LIMIT 3)
		UNION ALL
		(SELECT 'role_modified' as event_type, r.name, r.date_modified as event_date
		 FROM role r
		 WHERE r.date_modified IS NOT NULL AND r.workspace_id = $1
		 ORDER BY r.date_modified DESC LIMIT 2)
		ORDER BY event_date DESC
		LIMIT 5
	`
	rows, err := sqlDB.QueryContext(ctx, query, workspaceID)
	if err != nil {
		log.Printf("Dashboard activity query error: %v", err)
		return nil
	}
	defer rows.Close()

	var items []userdashboard.ActivityItem
	for rows.Next() {
		var eventType, name string
		var eventDate time.Time
		if err := rows.Scan(&eventType, &name, &eventDate); err != nil {
			continue
		}
		item := MapActivityItem(eventType, name, eventDate, renderIcon)
		items = append(items, item)
	}
	return items
}

// MapActivityItem maps an event row to a dashboard ActivityItem.
func MapActivityItem(eventType, name string, eventDate time.Time, renderIcon IconRenderer) userdashboard.ActivityItem {
	item := userdashboard.ActivityItem{
		TimeAgo: FormatTimeAgo(eventDate),
	}
	switch eventType {
	case "user_created":
		item.IconHTML = renderIcon("icon-user-plus")
		item.Title = "New User Created"
		item.Description = name + " added"
	case "role_modified":
		item.IconHTML = renderIcon("icon-shield")
		item.Title = "Role Updated"
		item.Description = name + " role modified"
	default:
		item.IconHTML = renderIcon("icon-info")
		item.Title = "Activity"
		item.Description = name
	}
	return item
}

// FormatTimeAgo formats a timestamp as a human-readable relative time string.
func FormatTimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	}
}

func buildChartData(ctx context.Context, sqlDB *sql.DB, workspaceID string) userdashboard.ChartData {
	query := `
		SELECT TO_CHAR(date_trunc('month', wu.date_created), 'Mon') as month_label,
		       COUNT(*) as user_count
		FROM workspace_user wu
		WHERE wu.date_created >= NOW() - INTERVAL '12 months'
		  AND wu.workspace_id = $1
		GROUP BY date_trunc('month', wu.date_created)
		ORDER BY date_trunc('month', wu.date_created)
	`
	rows, err := sqlDB.QueryContext(ctx, query, workspaceID)
	if err != nil {
		log.Printf("Dashboard chart query error: %v", err)
		return userdashboard.ChartData{Period: "year"}
	}
	defer rows.Close()

	var labels []string
	var values []int
	for rows.Next() {
		var label string
		var count int
		if err := rows.Scan(&label, &count); err != nil {
			continue
		}
		labels = append(labels, label)
		values = append(values, count)
	}

	return userdashboard.ChartData{
		Labels: labels,
		Values: values,
		Period: "year",
	}
}
