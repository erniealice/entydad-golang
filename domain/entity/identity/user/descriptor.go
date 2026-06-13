package user

import "github.com/erniealice/pyeza-golang/compose"

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
	}
}
