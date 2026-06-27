package delegate

// labels.go — Delegate label structs.
// Mirrors party/client/labels.go trimmed to page/list/actions only.
// JSON tags match the "delegate" wrapper key in lyngua delegate.json.
// No Detail/Form/Empty/BulkActions/Dashboard/Subscriptions label types
// (Delegate has only active bool — no multi-status lifecycle).

// Labels holds all translatable strings for the delegate module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	List    ListLabels   `json:"list"`
	Actions ActionLabels `json:"actions"`
}

// PageLabels holds heading and caption for the delegate list page.
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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Students string `json:"students"`
	Active   string `json:"active"`
}

// ActionLabels holds add/edit/delete action labels.
type ActionLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}
