package form

import (
	"github.com/erniealice/entydad-golang"
)

// Data is the template data for the "Add user to workspace" drawer form.
type Data struct {
	FormAction    string
	WorkspaceID   string
	Labels        entydad.WorkspaceUserLabels
	UserSearchURL string
	CommonLabels  any
}
