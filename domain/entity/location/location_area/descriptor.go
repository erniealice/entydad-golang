package location_area

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "entity.location_area",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "location_area"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "location_area.json", Key: ""},
		LabelName: "LocationAreaLabels",
		Templates: TemplatesFS,
	}
}
