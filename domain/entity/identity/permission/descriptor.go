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
	}
}
