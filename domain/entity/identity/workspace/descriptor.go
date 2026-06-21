package workspace

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.workspace",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "workspace"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "workspace.json", Key: ""},
		LabelName: "WorkspaceLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "workspace:list",
			Items: []compose.NavItem{
				{Key: "workspaces-active", Route: "workspace.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-briefcase", Permission: "workspace:list"},
				{Key: "workspaces-inactive", Route: "workspace.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-briefcase", Permission: "workspace:list"},
			},
		},
	}
}
