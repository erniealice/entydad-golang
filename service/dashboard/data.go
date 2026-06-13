// Package dashboard provides utility functions for the dashboard views.
//
// The raw-SQL builder functions (BuildGetUsersByRoleID, BuildGetDashboardData)
// that previously lived here have been migrated to espyna use cases at
// packages/espyna-golang/internal/application/usecases/service/dashboard/home/.
// The SQL now lives in the postgres adapter at contrib/postgres/internal/
// adapter/entity/workspace_user_home_dashboard.go and
// workspace_user_role_home_dashboard.go.
//
// What remains: MapActivityItem and FormatTimeAgo — utility functions used by
// the composition layer to map use-case result rows to entydad view types.
package dashboard

import (
	"fmt"
	"html/template"
	"time"

	userdashboard "github.com/erniealice/entydad-golang/domain/entity/identity/user/dashboard"
)

// IconRenderer converts an icon template name to rendered HTML. Callers
// typically pass pyeza.HTMLRenderer.RenderIcon.
type IconRenderer func(name string) template.HTML

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
