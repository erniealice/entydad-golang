package payment_term

// labels.go — PaymentTerm label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// commerce sub-context). Pure structural move — no behaviour change; field
// names, json tags, and string literals are byte-identical. Entity-local
// rename: PaymentTermLabels -> Labels, PaymentTerm<Xxx>Labels -> <Xxx>Labels.
//
// Note: the root labels.go defines no DefaultPaymentTermLabels() constructor
// (these labels are loaded from lyngua JSON, not a Go default), so this file
// carries the struct definitions only.

// Labels holds all translatable strings for the payment term module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
}

type PageLabels struct {
	Heading  string `json:"heading"`
	Subtitle string `json:"subtitle"`
}

type ButtonLabels struct {
	AddPaymentTerm string `json:"addPaymentTerm"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	NetDays   string `json:"netDays"`
	IsDefault string `json:"isDefault"`
	Status    string `json:"status"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FormLabels struct {
	SectionInfo            string `json:"sectionInfo"`
	SectionTerms           string `json:"sectionTerms"`
	SectionSettings        string `json:"sectionSettings"`
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Code                   string `json:"code"`
	CodePlaceholder        string `json:"codePlaceholder"`
	Type                   string `json:"type"`
	NetDays                string `json:"netDays"`
	DiscountDays           string `json:"discountDays"`
	DiscountPercentBps     string `json:"discountPercentBps"`
	TypeHint               string `json:"typeHint"`
	NetDaysHint            string `json:"netDaysHint"`
	DiscountDaysHint       string `json:"discountDaysHint"`
	DiscountPercentBpsHint string `json:"discountPercentBpsHint"`
	PriorityHint           string `json:"priorityHint"`
	EntityScope            string `json:"entityScope"`
	IsDefault              string `json:"isDefault"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	DisplayOrder           string `json:"displayOrder"`
	Active                 string `json:"active"`

	// Type select options
	TypeDueOnReceipt        string `json:"typeDueOnReceipt"`
	TypeNet                 string `json:"typeNet"`
	TypeCOD                 string `json:"typeCOD"`
	TypeProximate           string `json:"typeProximate"`
	ProximateDay            string `json:"proximateDay"`
	ProximateDayPlaceholder string `json:"proximateDayPlaceholder"`

	// Entity scope select options
	ScopesBoth         string `json:"scopesBoth"`
	ScopesSupplierOnly string `json:"scopesSupplierOnly"`
	ScopesClientOnly   string `json:"scopesClientOnly"`

	// Field-level info text for the drawer form.
	NameInfo        string `json:"nameInfo"`
	CodeInfo        string `json:"codeInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
}

type ActionLabels struct {
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}
