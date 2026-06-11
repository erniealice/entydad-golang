package form

import (
	workspace_user "github.com/erniealice/entydad-golang/domain/entity/identity/workspace_user"
)

// Data is the template data for the "Add user to workspace" drawer form.
type Data struct {
	FormAction    string
	WorkspaceID   string
	Labels        workspace_user.Labels
	UserSearchURL string
	CommonLabels  any
}
