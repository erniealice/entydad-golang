package location_area

// routes.go — LocationArea route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, location sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: LocationAreaRoutes ->
// Routes, DefaultLocationAreaRoutes -> DefaultRoutes, LocationArea<Xxx>URL ->
// <Xxx>URL.

// Default route constants for the location area view.
const (
	DashboardURL     = "/location-areas/dashboard"
	ListURL          = "/location-areas/list/{status}"
	TableURL         = "/action/location-area/table/{status}"
	DetailURL        = "/location-areas/detail/{id}"
	AddURL           = "/action/location-area/add"
	EditURL          = "/action/location-area/edit/{id}"
	DeleteURL        = "/action/location-area/delete"
	BulkDeleteURL    = "/action/location-area/bulk-delete"
	SetStatusURL     = "/action/location-area/set-status"
	BulkSetStatusURL = "/action/location-area/bulk-set-status"
)

// Routes holds all route paths for location area management.
type Routes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		TableURL:         TableURL,
		DetailURL:        DetailURL,
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
		"location_area.dashboard":       r.DashboardURL,
		"location_area.list":            r.ListURL,
		"location_area.table":           r.TableURL,
		"location_area.detail":          r.DetailURL,
		"location_area.add":             r.AddURL,
		"location_area.edit":            r.EditURL,
		"location_area.delete":          r.DeleteURL,
		"location_area.bulk_delete":     r.BulkDeleteURL,
		"location_area.set_status":      r.SetStatusURL,
		"location_area.bulk_set_status": r.BulkSetStatusURL,
	}
}
