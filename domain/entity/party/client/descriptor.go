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
				{Key: "dashboard", Route: "client.dashboard", Label: "Dashboard", Icon: "icon-dashboard", LabelKey: "dashboard_label", IconKey: "dashboard_icon"},
				{Key: "active", Route: "client.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-user-check", LabelKey: "active_label", IconKey: "clients_active_icon"},
				{Key: "prospect", Route: "client.list", Params: map[string]string{"status": "prospect"}, Label: "Prospect", Icon: "icon-user-plus", LabelKey: "prospect_label", IconKey: "clients_prospect_icon"},
				{Key: "on_hold", Route: "client.list", Params: map[string]string{"status": "on_hold"}, Label: "On Hold", Icon: "icon-pause-circle", LabelKey: "on_hold_label", IconKey: "clients_on_hold_icon"},
				{Key: "blocked", Route: "client.list", Params: map[string]string{"status": "blocked"}, Label: "Blocked", Icon: "icon-x-circle", LabelKey: "blocked_label", IconKey: "clients_blocked_icon"},
				{Key: "inactive", Route: "client.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-user-minus", LabelKey: "inactive_label", IconKey: "clients_inactive_icon"},
				{Key: "payment-terms", Route: "client.payment_terms", Label: "Payment Terms", Icon: "icon-clock", Permission: "client:list", LabelKey: "payment_terms_label", IconKey: "payment_terms_icon"},
				{Key: "receivables-aging", Route: "client.receivables_aging", Label: "Receivables Aging", Icon: "icon-file-text", LabelKey: "receivables_aging_label", IconKey: "receivables_aging_icon"},
			},
		},
	}
}
