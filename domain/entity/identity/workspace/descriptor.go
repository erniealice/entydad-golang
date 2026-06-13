package workspace

import "github.com/erniealice/pyeza-golang/compose"

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
	}
}
