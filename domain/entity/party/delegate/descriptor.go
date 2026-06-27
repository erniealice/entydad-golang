package delegate

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose.Unit descriptor for the Delegate entity.
// Mirrors party/client/descriptor.go with delegate-specific values.
//
// AppEntry.Route is "delegate.list" (not a dashboard) so Params must include
// {"status": "active"} for ResolveAppEntryURL to resolve the {status} segment.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := Labels{}
	return compose.Unit{
		Key:       "entity.delegate",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "delegate"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "delegate.json", Key: "delegate"},
		LabelName: "DelegateLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "delegate:list",
			AppEntry: &compose.AppEntry{
				Key:        "delegate",
				Route:      "delegate.list",
				Params:     map[string]string{"status": "active"},
				Label:      "Delegates",
				Icon:       "icon-users",
				Permission: "delegate:list",
			},
			Items: []compose.NavItem{
				{Key: "active", Route: "delegate.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-user-check"},
				{Key: "inactive", Route: "delegate.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-user-minus"},
			},
		},
	}
}
