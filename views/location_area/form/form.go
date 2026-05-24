package form

import (
	"github.com/erniealice/entydad-golang"
)

// Data is the template data for the location area drawer form.
type Data struct {
	FormAction   string
	WorkspaceID   string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Active       bool
	Labels       entydad.LocationAreaFormLabels
	CommonLabels any
}
