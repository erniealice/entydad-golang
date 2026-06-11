package form

import (
	locationarea "github.com/erniealice/entydad-golang/domain/entity/location/location_area"
)

// Data is the template data for the location area drawer form.
type Data struct {
	FormAction   string
	WorkspaceID  string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Active       bool
	Labels       locationarea.FormLabels
	CommonLabels any
}
