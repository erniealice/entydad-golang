package workspace_user

// routes.go — WorkspaceUser route struct, URL consts, and constructors.
//
// Extracted verbatim from packages/entydad-golang/{routes.go,routes_config.go}
// (entity domain, identity sub-context). Pure structural move — route URL
// string values are byte-identical. Entity-local rename: WorkspaceUserRoutes ->
// Routes, DefaultWorkspaceUserRoutes -> DefaultRoutes, WorkspaceUser<Xxx>URL ->
// <Xxx>URL.

// Default route constants for the workspace_user view.
const (
	ListURL      = "/workspace-users/list/{status}"
	DetailURL    = "/workspace-users/detail/{id}"
	TabActionURL = "/action/workspace_user/{id}/tab/{tab}"
	AddURL       = "/action/workspace_user/add"
	DeleteURL    = "/action/workspace_user/delete/{id}"
	SetStatusURL = "/action/workspace_user/set-status/{id}"
	SearchURL    = "/action/workspace_user/search"

	AttachmentUploadURL = "/action/workspace_user/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/workspace_user/{id}/attachments/delete"
)

// Routes holds all route paths for workspace-user nested detail management.
type Routes struct {
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	AddURL       string `json:"add_url"`
	DeleteURL    string `json:"delete_url"`
	SetStatusURL string `json:"set_status_url"`
	SearchURL    string `json:"search_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultRoutes returns a Routes populated from the
// package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ListURL:      ListURL,
		DetailURL:    DetailURL,
		TabActionURL: TabActionURL,
		AddURL:       AddURL,
		DeleteURL:    DeleteURL,
		SetStatusURL: SetStatusURL,
		SearchURL:    SearchURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route path values.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"workspace_user.list":       r.ListURL,
		"workspace_user.detail":     r.DetailURL,
		"workspace_user.tab_action": r.TabActionURL,
		"workspace_user.add":        r.AddURL,
		"workspace_user.delete":     r.DeleteURL,
		"workspace_user.set_status": r.SetStatusURL,
		"workspace_user.search":     r.SearchURL,

		"workspace_user.attachment.upload": r.AttachmentUploadURL,
		"workspace_user.attachment.delete": r.AttachmentDeleteURL,
	}
}
