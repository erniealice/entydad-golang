package staff

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose.Unit descriptor for the Staff entity.
// Mirrors party/delegate/descriptor.go with staff-specific values.
//
// AppEntry.Route is "staff.list" (not a dashboard) so Params must include
// {"status": "active"} for ResolveAppEntryURL to resolve the {status} segment.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.staff",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "staff"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "staff.json", Key: "staff"},
		LabelName: "StaffLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "staff:list",
			AppEntry: &compose.AppEntry{
				Key:        "staff",
				Route:      "staff.list",
				Params:     map[string]string{"status": "active"},
				Label:      "Staff",
				Icon:       "icon-users",
				Permission: "staff:list",
			},
			Items: []compose.NavItem{
				{Key: "active", Route: "staff.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-user-check", LabelKey: "active_label", IconKey: "clients_active_icon"},
				{Key: "inactive", Route: "staff.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-user-minus", LabelKey: "inactive_label", IconKey: "clients_inactive_icon"},
			},
		},
	}
}
