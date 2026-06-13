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
		Nav: compose.NavContrib{
			Permission: "location_area:list",
			Items: []compose.NavItem{
				{Key: "location-areas-active", Route: "location_area.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-layers", Permission: "location_area:list"},
				{Key: "location-areas-inactive", Route: "location_area.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-layers-off", Permission: "location_area:list"},
			},
		},
	}
}
