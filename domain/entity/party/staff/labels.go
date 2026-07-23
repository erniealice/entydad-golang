package staff

// labels.go — Staff label structs.
// Mirrors party/delegate/labels.go trimmed to page/list/actions only.
// JSON tags match the "staff" wrapper key in lyngua staff.json.
// The drawer's Labels (user picker, employment_type, status options) live in
// the form/ subpackage per ui-drawer-form-subpackage-convention.md — this
// struct is only the list-page + row-action vocabulary, same split as
// party/delegate.

// Labels holds all translatable strings for the staff module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	List    ListLabels   `json:"list"`
	Actions ActionLabels `json:"actions"`
}

// PageLabels holds heading and caption for the staff list page.
type PageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

// ListLabels holds column header labels.
type ListLabels struct {
	Columns ColumnLabels `json:"columns"`
}

// ColumnLabels holds individual column header strings.
type ColumnLabels struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	EmploymentType string `json:"employment_type"`
	Status         string `json:"status"`
}

// ActionLabels holds add/edit/delete action labels.
type ActionLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}
