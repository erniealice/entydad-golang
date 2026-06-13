package client_tag

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.client_tag",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "client_tag"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "client_tag.json", Key: ""},
		LabelName: "ClientTagLabels",
		Templates: nil,
		Nav: compose.NavContrib{
			Permission: "client:list",
			Items: []compose.NavItem{
				{Key: "tags", Route: "client_tag.list", Label: "Tags", Icon: "icon-tag", Permission: "client:list"},
			},
		},
	}
}
