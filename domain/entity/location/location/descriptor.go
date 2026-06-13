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
		Nav: compose.NavContrib{
			Permission: "location:list",
			AppEntry: &compose.AppEntry{
				Key:        "location",
				Route:      "location.list",
				Label:      "Locations",
				Icon:       "icon-map-pin",
				Permission: "location:list",
			},
			Items: []compose.NavItem{
				{Key: "locations-active", Route: "location.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-map-pin"},
				{Key: "locations-inactive", Route: "location.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-map-pin"},
			},
		},
	}
}
