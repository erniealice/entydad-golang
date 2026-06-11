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
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddPermission string `json:"addPermission"`
}

type ColumnLabels struct {
	Name           string `json:"name"`
	Entity         string `json:"entity"`
	PermissionCode string `json:"permissionCode"`
	Type           string `json:"type"`
	Status         string `json:"status"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
	Name                      string `json:"name"`
	NamePlaceholder           string `json:"namePlaceholder"`
	PermissionCode            string `json:"permissionCode"`
	PermissionCodePlaceholder string `json:"permissionCodePlaceholder"`
	PermissionCodeHint        string `json:"permissionCodeHint"`
	PermissionType            string `json:"permissionType"`
	Description               string `json:"description"`
	DescriptionPlaceholder    string `json:"descriptionPlaceholder"`
	Active                    string `json:"active"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}
