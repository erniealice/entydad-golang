package staff

// routes.go — Staff route struct, URL consts, and constructors.
// Mirrors party/delegate/routes.go trimmed to the list+action surface
// (no dashboard, detail, tabs, attachments, statement, or revenue-run) plus
// a SearchURL for the user autocomplete (mirrors identity/workspace_user's
// SearchURL — Staff links to an EXISTING user, it does not embed a brand-new
// one the way Delegate/Client do).

// Default route constants for the staff view.
const (
	ListURL       = "/staffs/list/{status}"
	TableURL      = "/action/staff/table/{status}"
	AddURL        = "/action/staff/add"
	EditURL       = "/action/staff/edit/{id}"
	DeleteURL     = "/action/staff/delete"
	BulkDeleteURL = "/action/staff/bulk-delete"
	SearchURL     = "/action/staff/search"
)

// Routes holds the resolved URL strings for the staff module.
type Routes struct {
	ListURL       string `json:"list_url"`
	TableURL      string `json:"table_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
	SearchURL     string `json:"search_url"`
}

// DefaultRoutes returns a Routes populated from the package-level constants.
func DefaultRoutes() Routes {
	return Routes{
		ListURL:       ListURL,
		TableURL:      TableURL,
		AddURL:        AddURL,
		EditURL:       EditURL,
		DeleteURL:     DeleteURL,
		BulkDeleteURL: BulkDeleteURL,
		SearchURL:     SearchURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"staff.list":        r.ListURL,
		"staff.table":       r.TableURL,
		"staff.add":         r.AddURL,
		"staff.edit":        r.EditURL,
		"staff.delete":      r.DeleteURL,
		"staff.bulk_delete": r.BulkDeleteURL,
		"staff.search":      r.SearchURL,
	}
}
