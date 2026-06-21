package payment_term

import "github.com/erniealice/espyna-golang/consumer/compose"

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
		Nav: compose.NavContrib{
			Permission: "client:list",
			Items: []compose.NavItem{
				// Client-context payment terms (under "client" app Settings section)
				{Key: "payment-terms", Route: "payment_term.list", Label: "Payment Terms", Icon: "icon-clock", Permission: "client:list"},
			},
		},
	}
}
