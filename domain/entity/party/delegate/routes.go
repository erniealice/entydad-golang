package delegate

// routes.go — Delegate route struct, URL consts, and constructors.
// Mirrors party/client/routes.go trimmed to the list+action surface
// (no dashboard, detail, tabs, attachments, statement, or revenue-run).

// Default route constants for the delegate view.
const (
	ListURL       = "/delegates/list/{status}"
	TableURL      = "/action/delegate/table/{status}"
	AddURL        = "/action/delegate/add"
	EditURL       = "/action/delegate/edit/{id}"
	DeleteURL     = "/action/delegate/delete"
	BulkDeleteURL = "/action/delegate/bulk-delete"
)

// Routes holds the resolved URL strings for the delegate module.
type Routes struct {
	ListURL       string `json:"list_url"`
	TableURL      string `json:"table_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
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
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"delegate.list":        r.ListURL,
		"delegate.table":       r.TableURL,
		"delegate.add":         r.AddURL,
		"delegate.edit":        r.EditURL,
		"delegate.delete":      r.DeleteURL,
		"delegate.bulk_delete": r.BulkDeleteURL,
	}
}
