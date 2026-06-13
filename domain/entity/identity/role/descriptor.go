package role

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.role",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "role"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "role.json", Key: ""},
		LabelName: "RoleLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "user:list",
			Items: []compose.NavItem{
				{Key: "roles", Route: "role.list", Label: "Roles", Icon: "icon-shield"},
			},
		},
	}
}
