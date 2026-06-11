package client_tag

// routes.go — Client Tag route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, party sub-context). Pure structural move — route URL string
// values are byte-identical. Entity-local rename: ClientTagRoutes -> Routes,
// DefaultClientTagRoutes -> DefaultRoutes, ClientTag<Xxx>URL -> <Xxx>URL.

// Default route constants for the client tag view.
const (
	ListURL          = "/clients/settings/tags/list"
	TableURL         = "/action/client/tags/table"
	AddURL           = "/action/client/tags/add"
	EditURL          = "/action/client/tags/edit/{id}"
	DeleteURL        = "/action/client/tags/delete"
	BulkDeleteURL    = "/action/client/tags/bulk-delete"
	SetStatusURL     = "/action/client/tags/set-status"
	BulkSetStatusURL = "/action/client/tags/bulk-set-status"
)

// Routes holds all route paths for client tag (category) management.
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

// DefaultRoutes returns a Routes populated from the
// package-level route constants.
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
		"client_tag.list":            r.ListURL,
		"client_tag.table":           r.TableURL,
		"client_tag.add":             r.AddURL,
		"client_tag.edit":            r.EditURL,
		"client_tag.delete":          r.DeleteURL,
		"client_tag.bulk_delete":     r.BulkDeleteURL,
		"client_tag.set_status":      r.SetStatusURL,
		"client_tag.bulk_set_status": r.BulkSetStatusURL,
	}
}
