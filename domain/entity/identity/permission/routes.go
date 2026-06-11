package permission

// routes.go — Permission route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, identity sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: PermissionRoutes ->
// Routes, DefaultPermissionRoutes -> DefaultRoutes, Permission<Xxx>URL ->
// <Xxx>URL.

// Default route constants for the permission view.
const (
	ListURL          = "/permissions/list/{status}"
	TableURL         = "/action/permission/table/{status}"
	AddURL           = "/action/permission/add"
	EditURL          = "/action/permission/edit/{id}"
	DeleteURL        = "/action/permission/delete"
	BulkDeleteURL    = "/action/permission/bulk-delete"
	SetStatusURL     = "/action/permission/set-status"
	BulkSetStatusURL = "/action/permission/bulk-set-status"
)

// Routes holds all route paths for permission management.
type Routes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants.
func DefaultRoutes() Routes {
	return Routes{
		ListURL:          ListURL,
		TableURL:         TableURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"permission.list":            r.ListURL,
		"permission.table":           r.TableURL,
		"permission.add":             r.AddURL,
		"permission.edit":            r.EditURL,
		"permission.delete":          r.DeleteURL,
		"permission.bulk_delete":     r.BulkDeleteURL,
		"permission.set_status":      r.SetStatusURL,
		"permission.bulk_set_status": r.BulkSetStatusURL,
	}
}
