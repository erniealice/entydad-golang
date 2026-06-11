package supplier_tag

// routes.go — Supplier Tag route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, party sub-context). Pure structural move — route URL string
// values are byte-identical. Entity-local rename: SupplierTagRoutes -> Routes,
// DefaultSupplierTagRoutes -> DefaultRoutes, SupplierTag<Xxx>URL -> <Xxx>URL.

// Default route constants for the supplier tag view.
const (
	ListURL          = "/suppliers/settings/tags/list"
	TableURL         = "/action/supplier/tags/table"
	AddURL           = "/action/supplier/tags/add"
	EditURL          = "/action/supplier/tags/edit/{id}"
	DeleteURL        = "/action/supplier/tags/delete"
	BulkDeleteURL    = "/action/supplier/tags/bulk-delete"
	SetStatusURL     = "/action/supplier/tags/set-status"
	BulkSetStatusURL = "/action/supplier/tags/bulk-set-status"
)

// Routes holds all route paths for supplier tag (category) management.
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
		"supplier_tag.list":            r.ListURL,
		"supplier_tag.table":           r.TableURL,
		"supplier_tag.add":             r.AddURL,
		"supplier_tag.edit":            r.EditURL,
		"supplier_tag.delete":          r.DeleteURL,
		"supplier_tag.bulk_delete":     r.BulkDeleteURL,
		"supplier_tag.set_status":      r.SetStatusURL,
		"supplier_tag.bulk_set_status": r.BulkSetStatusURL,
	}
}
