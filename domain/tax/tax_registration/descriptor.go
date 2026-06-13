package tax_registration

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "tax.tax_registration",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "tax_registration"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "tax_registration.json", Key: ""},
		LabelName: "TaxRegistrationLabels",
		Templates: TemplatesFS,
	}
}
