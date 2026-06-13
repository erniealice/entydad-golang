package workspace_user

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.workspace_user",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "workspace_user"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "workspace_user.json", Key: ""},
		LabelName: "WorkspaceUserLabels",
		Templates: TemplatesFS,
	}
}
