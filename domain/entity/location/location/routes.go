package location

// routes.go — Location route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, location sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: LocationRoutes ->
// Routes, DefaultLocationRoutes -> DefaultRoutes, Location<Xxx>URL -> <Xxx>URL.

// Default route constants for the location view.
const (
	DashboardURL        = "/locations/dashboard"
	DetailURL           = "/locations/detail/{id}"
	ListURL             = "/locations/list/{status}"
	TableURL            = "/action/location/table/{status}"
	AddURL              = "/action/location/add"
	EditURL             = "/action/location/edit/{id}"
	DeleteURL           = "/action/location/delete"
	BulkDeleteURL       = "/action/location/bulk-delete"
	SetStatusURL        = "/action/location/set-status"
	BulkSetStatusURL    = "/action/location/bulk-set-status"
	TabActionURL        = "/action/location/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/location/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/location/{id}/attachments/delete"
	EditDetailURL       = "/action/location/edit-detail/{id}"
)

// Routes holds all route paths for location management.
type Routes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`
	EditDetailURL    string `json:"edit_detail_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		DetailURL:        DetailURL,
		TableURL:         TableURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
		TabActionURL:     TabActionURL,
		EditDetailURL:    EditDetailURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"location.dashboard":       r.DashboardURL,
		"location.list":            r.ListURL,
		"location.detail":          r.DetailURL,
		"location.table":           r.TableURL,
		"location.add":             r.AddURL,
		"location.edit":            r.EditURL,
		"location.delete":          r.DeleteURL,
		"location.bulk_delete":     r.BulkDeleteURL,
		"location.set_status":      r.SetStatusURL,
		"location.bulk_set_status": r.BulkSetStatusURL,
		"location.tab_action":      r.TabActionURL,
		"location.edit_detail":     r.EditDetailURL,

		"location.attachment.upload": r.AttachmentUploadURL,
		"location.attachment.delete": r.AttachmentDeleteURL,
	}
}
