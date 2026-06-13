package payment_term

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.payment_term",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "payment_term"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "payment_term.json", Key: "paymentTerm"},
		LabelName: "PaymentTermLabels",
		Templates: TemplatesFS,
	}
}
