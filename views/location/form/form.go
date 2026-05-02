package form

import (
	pyeza "github.com/erniealice/pyeza-golang/types"

	"github.com/erniealice/entydad-golang"
)

// Data is the template data for the location drawer form.
type Data struct {
	FormAction                string
	IsEdit                    bool
	ID                        string
	Name                      string
	Address                   string
	Description               string
	Timezone                  string
	Active                    bool
	SelectedLocationAreaID    string
	LocationAreaSelectOptions []pyeza.SelectOption
	Labels                    entydad.LocationFormLabels
	CommonLabels              any
}
