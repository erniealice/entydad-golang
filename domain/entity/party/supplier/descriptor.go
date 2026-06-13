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
		Nav: compose.NavContrib{
			Permission: "supplier:list",
			AppEntry: &compose.AppEntry{
				Key:        "supplier",
				Route:      "supplier.dashboard",
				Label:      "Suppliers",
				Icon:       "icon-truck",
				Permission: "supplier:list",
			},
			Items: []compose.NavItem{
				{Key: "dashboard", Route: "supplier.dashboard", Label: "Dashboard", Icon: "icon-dashboard"},
				{Key: "suppliers-active", Route: "supplier.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-truck"},
				{Key: "suppliers-blocked", Route: "supplier.list", Params: map[string]string{"status": "blocked"}, Label: "Blocked", Icon: "icon-x-circle"},
				{Key: "suppliers-on-hold", Route: "supplier.list", Params: map[string]string{"status": "on_hold"}, Label: "On Hold", Icon: "icon-pause-circle"},
				{Key: "payment-terms", Route: "supplier.payment_terms", Label: "Payment Terms", Icon: "icon-clock", Permission: "supplier:list"},
				{Key: "payables-aging", Route: "supplier.payables_aging", Label: "Payables Aging", Icon: "icon-file-text"},
			},
		},
	}
}
