package client

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.client",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "client"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "client.json", Key: "client"},
		LabelName: "ClientLabels",
		Templates: TemplatesFS,
	}
}
