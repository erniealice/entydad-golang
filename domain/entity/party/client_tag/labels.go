package client_tag

// labels.go — Client Tag label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// party sub-context). Pure structural move — field names, json tags, and
// string literals are byte-identical. Entity-local rename: ClientTagLabels ->
// Labels, ClientTag<Xxx>Labels -> <Xxx>Labels.

// Labels holds all translatable strings for the client tag module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Empty   EmptyLabels   `json:"empty"`
	Actions ActionLabels  `json:"actions"`
	Confirm ConfirmLabels `json:"confirm"`
}

type PageLabels struct {
	Heading  string `json:"heading"`
	Subtitle string `json:"subtitle"`
}

type ButtonLabels struct {
	AddTag string `json:"addTag"`
}

type ColumnLabels struct {
	TagName     string `json:"tagName"`
	Customers   string `json:"customers"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ActionLabels struct {
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ConfirmLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
	CannotDelete  string `json:"cannotDelete"`
}
