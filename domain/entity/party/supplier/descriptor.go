package supplier

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.supplier",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "supplier"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "supplier.json", Key: "supplier"},
		LabelName: "SupplierLabels",
		Templates: TemplatesFS,
	}
}
