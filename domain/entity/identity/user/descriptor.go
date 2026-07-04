package user

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.user",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "user"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "user.json", Key: ""},
		LabelName: "UserLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "user:list",
			AppEntry: &compose.AppEntry{
				Key:        "user",
				Route:      "user.dashboard",
				Label:      "Users",
				Icon:       "icon-shield",
				Permission: "user:list",
			},
			Items: []compose.NavItem{
				{Key: "dashboard", Route: "user.dashboard", Label: "Dashboard", Icon: "icon-dashboard", LabelKey: "dashboard_label", IconKey: "dashboard_icon"},
				{Key: "active", Route: "user.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-shield", LabelKey: "active_label", IconKey: "users_active_icon"},
				{Key: "inactive", Route: "user.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-user-minus", LabelKey: "inactive_label", IconKey: "users_inactive_icon"},
			},
		},
	}
}
