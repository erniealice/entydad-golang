package client

import "github.com/erniealice/espyna-golang/consumer/compose"

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
		Nav: compose.NavContrib{
			Permission: "client:list",
			AppEntry: &compose.AppEntry{
				Key:        "client",
				Route:      "client.dashboard",
				Label:      "Clients",
				Icon:       "icon-users",
				Permission: "client:list",
			},
			Items: []compose.NavItem{
				{Key: "dashboard", Route: "client.dashboard", Label: "Dashboard", Icon: "icon-dashboard"},
				{Key: "active", Route: "client.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-user-check"},
				{Key: "prospect", Route: "client.list", Params: map[string]string{"status": "prospect"}, Label: "Prospect", Icon: "icon-user-plus"},
				{Key: "on_hold", Route: "client.list", Params: map[string]string{"status": "on_hold"}, Label: "On Hold", Icon: "icon-pause-circle"},
				{Key: "blocked", Route: "client.list", Params: map[string]string{"status": "blocked"}, Label: "Blocked", Icon: "icon-x-circle"},
				{Key: "inactive", Route: "client.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-user-minus"},
				{Key: "payment-terms", Route: "client.payment_terms", Label: "Payment Terms", Icon: "icon-clock", Permission: "client:list"},
				{Key: "receivables-aging", Route: "client.receivables_aging", Label: "Receivables Aging", Icon: "icon-file-text"},
			},
		},
	}
}
