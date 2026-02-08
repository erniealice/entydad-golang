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
