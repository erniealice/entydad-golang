package location

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.location",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "location"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "location.json", Key: ""},
		LabelName: "LocationLabels",
		Templates: TemplatesFS,
	}
}
