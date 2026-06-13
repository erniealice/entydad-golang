package permission

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.permission",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "permission"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "permission.json", Key: ""},
		LabelName: "PermissionLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "permission:list",
			AppEntry: &compose.AppEntry{
				Key:        "admin",
				Route:      "permission.list",
				Label:      "Settings",
				Icon:       "icon-settings",
				Permission: "permission:list",
			},
			Items: []compose.NavItem{
				{Key: "permissions-active", Route: "permission.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-key", Permission: "permission:list"},
				{Key: "permissions-inactive", Route: "permission.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-key", Permission: "permission:list"},
			},
		},
	}
}
