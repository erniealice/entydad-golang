package permission

// labels.go — Permission label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// identity sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: PermissionLabels -> Labels, Permission<Xxx>Labels -> <Xxx>Labels.

// Labels holds all translatable strings for the permission module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
}

type ButtonLabels struct {
	AddPermission string `json:"add_permission"`
}

type ColumnLabels struct {
	Name           string `json:"name"`
	Entity         string `json:"entity"`
	PermissionCode string `json:"permission_code"`
	Type           string `json:"type"`
	Status         string `json:"status"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"active_title"`
	ActiveMessage   string `json:"active_message"`
	InactiveTitle   string `json:"inactive_title"`
	InactiveMessage string `json:"inactive_message"`
}

type FormLabels struct {
	Name                      string `json:"name"`
	NamePlaceholder           string `json:"name_placeholder"`
	PermissionCode            string `json:"permission_code"`
	PermissionCodePlaceholder string `json:"permission_code_placeholder"`
	PermissionCodeHint        string `json:"permission_code_hint"`
	PermissionType            string `json:"permission_type"`
	Description               string `json:"description"`
	DescriptionPlaceholder    string `json:"description_placeholder"`
	Active                    string `json:"active"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}
