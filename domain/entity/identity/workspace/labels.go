package workspace

// labels.go — Workspace label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// identity sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: WorkspaceLabels -> Labels, Workspace<Xxx>Labels -> <Xxx>Labels.

// Labels holds all translatable strings for the workspace module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
	Detail  DetailLabels `json:"detail"`
}

// DetailLabels holds i18n strings for the workspace detail page (Phase 1).
type DetailLabels struct {
	Tabs  DetailTabLabels  `json:"tabs"`
	Users DetailUserLabels `json:"users"`
	// Info holds the Info-tab section title + non-Form/Column field labels.
	// Name/Description/Status reuse Form.* and Columns.Status; Info supplies
	// the section heading + the Currency/Region fields not present elsewhere.
	Info DetailInfoLabels `json:"info"`
	// TaxReg holds the Tax Registrations tab panel copy (W4.5 label
	// remediation — previously hardcoded in detail.html).
	TaxReg DetailTaxRegLabels `json:"taxReg"`
}

// DetailTaxRegLabels holds the Tax Registrations tab panel copy
// (W4.5 label remediation).
type DetailTaxRegLabels struct {
	Loading       string `json:"loading"`
	Title         string `json:"title"`
	NotConfigured string `json:"notConfigured"`
}

// DetailInfoLabels holds the Info-tab labels on the workspace detail
// page (W4.5 label remediation). SectionTitle is the "Details" heading;
// Currency/Region are the optional workspace fields shown on the Info tab.
type DetailInfoLabels struct {
	SectionTitle string `json:"sectionTitle"`
	Currency     string `json:"currency"`
	Region       string `json:"region"`
}

// DetailTabLabels holds the tab display names for the workspace detail page.
type DetailTabLabels struct {
	Info        string `json:"info"`
	Users       string `json:"users"`
	Attachments string `json:"attachments"`
	// Phase 2 — polymorphic tax registrations tab
	TaxRegistrations string `json:"taxRegistrations"`
}

// DetailUserLabels holds i18n strings for the Users tab on the workspace detail page.
type DetailUserLabels struct {
	AddButton string `json:"addButton"`
	// Empty-state copy shown when the workspace has no user assignments yet
	// (W4.5 label remediation — previously hardcoded in users-tab.html).
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
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
	AddWorkspace string `json:"addWorkspace"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     string `json:"private"`
	Status      string `json:"status"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Private                string `json:"private"`
	Active                 string `json:"active"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}
