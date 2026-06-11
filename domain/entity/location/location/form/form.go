package form

import (
	pyeza "github.com/erniealice/pyeza-golang/types"

	location "github.com/erniealice/entydad-golang/domain/entity/location/location"
)

// Data is the template data for the location drawer form.
type Data struct {
	FormAction                string
	WorkspaceID               string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit                    bool
	ID                        string
	Name                      string
	Address                   string
	Description               string
	Timezone                  string
	Active                    bool
	SelectedLocationAreaID    string
	LocationAreaSelectOptions []pyeza.SelectOption
	Labels                    location.FormLabels
	CommonLabels              any
}
