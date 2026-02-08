package entydad

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// ---------------------------------------------------------------------------
// Client labels
// ---------------------------------------------------------------------------

// ClientLabels holds all translatable strings for the client module.
// JSON tags match the "client" wrapper key in retail/client.json.
type ClientLabels struct {
	Page        ClientPageLabels        `json:"page"`
	Buttons     ClientButtonLabels      `json:"buttons"`
	Columns     ClientColumnLabels      `json:"columns"`
	Empty       ClientEmptyLabels       `json:"empty"`
	Form        ClientFormLabels        `json:"form"`
	Detail      ClientDetailLabels      `json:"detail"`
	BulkActions ClientBulkActionLabels  `json:"bulkActions"`
}

type ClientPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingProspect string `json:"headingProspect"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionProspect string `json:"captionProspect"`
	CaptionInactive string `json:"captionInactive"`
}

type ClientButtonLabels struct {
	AddNew string `json:"addNew"`
}

type ClientColumnLabels struct {
	ClientName string `json:"clientName"`
}

type ClientEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	ProspectTitle   string `json:"prospectTitle"`
	ProspectMessage string `json:"prospectMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type ClientFormLabels struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ClientDetailLabels struct {
	CompanyDetails ClientCompanyDetailLabels `json:"companyDetails"`
	Actions        ClientDetailActionLabels  `json:"actions"`
}

type ClientCompanyDetailLabels struct {
	Status string `json:"status"`
}

type ClientDetailActionLabels struct {
	ViewClient   string `json:"viewClient"`
	EditClient   string `json:"editClient"`
	DeleteClient string `json:"deleteClient"`
}

type ClientBulkActionLabels struct {
	SetAsInactive string `json:"setAsInactive"`
}

// ---------------------------------------------------------------------------
// User labels
// ---------------------------------------------------------------------------

// UserLabels holds all translatable strings for the user module.
// JSON tags match retail/user.json (no wrapper key).
type UserLabels struct {
	Page    UserPageLabels    `json:"page"`
	Buttons UserButtonLabels  `json:"buttons"`
	Columns UserColumnLabels  `json:"columns"`
	Empty   UserEmptyLabels   `json:"empty"`
	Form    UserFormLabels    `json:"form"`
	Actions UserActionLabels  `json:"actions"`
}

type UserPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type UserButtonLabels struct {
	AddUser string `json:"addUser"`
}

type UserColumnLabels struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type UserEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type UserFormLabels struct {
	Mobile string `json:"mobile"`
}

type UserActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// Location labels
// ---------------------------------------------------------------------------

// LocationLabels holds all translatable strings for the location module.
// JSON tags match the "location" wrapper key in retail/location.json.
type LocationLabels struct {
	Page    LocationPageLabels    `json:"page"`
	Buttons LocationButtonLabels  `json:"buttons"`
	Columns LocationColumnLabels  `json:"columns"`
	Empty   LocationEmptyLabels   `json:"empty"`
	Form    LocationFormLabels    `json:"form"`
	Actions LocationActionLabels  `json:"actions"`
}

type LocationPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type LocationButtonLabels struct {
	AddLocation string `json:"addLocation"`
}

type LocationColumnLabels struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

type LocationEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type LocationFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Address                string `json:"address"`
	AddressPlaceholder     string `json:"addressPlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Active                 string `json:"active"`
}

type LocationActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// Role labels
// ---------------------------------------------------------------------------

// RoleLabels holds all translatable strings for the role module.
// JSON tags match the "role" wrapper key in retail/role.json.
type RoleLabels struct {
	Page    RolePageLabels    `json:"page"`
	Buttons RoleButtonLabels  `json:"buttons"`
	Columns RoleColumnLabels  `json:"columns"`
	Empty   RoleEmptyLabels   `json:"empty"`
	Form    RoleFormLabels    `json:"form"`
	Actions RoleActionLabels  `json:"actions"`
}

type RolePageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type RoleButtonLabels struct {
	AddRole string `json:"addRole"`
}

type RoleColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Status      string `json:"status"`
}

type RoleEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type RoleFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Color                  string `json:"color"`
	ColorPlaceholder       string `json:"colorPlaceholder"`
	Active                 string `json:"active"`
}

type RoleActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

// MapTableLabels maps common labels into the flat types.TableLabels structure.
func MapTableLabels(common pyeza.CommonLabels) types.TableLabels {
	return types.TableLabels{
		Search:             common.Table.Search,
		SearchPlaceholder:  common.Table.SearchPlaceholder,
		Filters:            common.Table.Filters,
		FilterConditions:   common.Table.FilterConditions,
		ClearAll:           common.Table.ClearAll,
		AddCondition:       common.Table.AddCondition,
		Clear:              common.Table.Clear,
		ApplyFilters:       common.Table.ApplyFilters,
		Sort:               common.Table.Sort,
		Columns:            common.Table.Columns,
		Export:              common.Table.Export,
		DensityDefault:     common.Table.Density.Default,
		DensityComfortable: common.Table.Density.Comfortable,
		DensityCompact:     common.Table.Density.Compact,
		Show:               common.Table.Show,
		Entries:             common.Table.Entries,
		Showing:            common.Table.Showing,
		To:                 common.Table.To,
		Of:                 common.Table.Of,
		EntriesLabel:       common.Table.EntriesLabel,
		SelectAll:          common.Table.SelectAll,
		Actions:            common.Table.Actions,
		Prev:               common.Pagination.Prev,
		Next:               common.Pagination.Next,
	}
}

// MapBulkConfig returns a BulkActionsConfig with labels from common bulk labels.
func MapBulkConfig(common pyeza.CommonLabels) types.BulkActionsConfig {
	return types.BulkActionsConfig{
		Enabled:        true,
		SelectAllLabel: common.Bulk.SelectAll,
		SelectedLabel:  common.Bulk.Selected,
		CancelLabel:    common.Bulk.ClearSelection,
	}
}
