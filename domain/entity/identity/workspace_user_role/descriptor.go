package workspace_user_role

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.workspace_user_role",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "workspace_user_role"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "workspace_user_role.json", Key: "workspace_user_role"},
		LabelName: "WorkspaceUserRoleLabels",
		Templates: TemplatesFS,
	}
}
