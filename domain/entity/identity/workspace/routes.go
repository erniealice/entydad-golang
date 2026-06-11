package workspace

// routes.go — Workspace route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, identity sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: WorkspaceRoutes ->
// Routes, DefaultWorkspaceRoutes -> DefaultRoutes, Workspace<Xxx>URL -> <Xxx>URL.

// Default route constants for the workspace view.
const (
	ListURL             = "/workspaces/list/{status}"
	TableURL            = "/action/workspace/table/{status}"
	AddURL              = "/action/workspace/add"
	EditURL             = "/action/workspace/edit/{id}"
	DeleteURL           = "/action/workspace/delete"
	BulkDeleteURL       = "/action/workspace/bulk-delete"
	SetStatusURL        = "/action/workspace/set-status"
	BulkSetStatusURL    = "/action/workspace/bulk-set-status"
	SwitchURL           = "/action/admin/switch-workspace"
	DetailURL           = "/workspaces/detail/{id}"
	TabActionURL        = "/action/workspace/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/workspace/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/workspace/{id}/attachments/delete"
)

// Routes holds all route paths for workspace management.
type Routes struct {
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	SwitchURL        string `json:"switch_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
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
		SwitchURL:        SwitchURL,
		DetailURL:        DetailURL,
		TabActionURL:     TabActionURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"workspace.list":            r.ListURL,
		"workspace.table":           r.TableURL,
		"workspace.add":             r.AddURL,
		"workspace.edit":            r.EditURL,
		"workspace.delete":          r.DeleteURL,
		"workspace.bulk_delete":     r.BulkDeleteURL,
		"workspace.set_status":      r.SetStatusURL,
		"workspace.bulk_set_status": r.BulkSetStatusURL,
		"workspace.switch_url":      r.SwitchURL,
		"workspace.detail":          r.DetailURL,
		"workspace.tab_action":      r.TabActionURL,

		"workspace.attachment.upload": r.AttachmentUploadURL,
		"workspace.attachment.delete": r.AttachmentDeleteURL,
	}
}
