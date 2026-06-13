package supplier_tag

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.supplier_tag",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "supplier_tag"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "supplier_tag.json", Key: ""},
		LabelName: "SupplierTagLabels",
		Templates: nil,
		Nav: compose.NavContrib{
			Permission: "supplier:list",
			Items: []compose.NavItem{
				{Key: "tags", Route: "supplier_tag.list", Label: "Tags", Icon: "icon-tag", Permission: "supplier:list"},
			},
		},
	}
}
