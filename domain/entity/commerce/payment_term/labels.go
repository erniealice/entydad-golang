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
	AddPaymentTerm string `json:"add_payment_term"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	NetDays   string `json:"net_days"`
	IsDefault string `json:"is_default"`
	Status    string `json:"status"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FormLabels struct {
	SectionInfo            string `json:"section_info"`
	SectionTerms           string `json:"section_terms"`
	SectionSettings        string `json:"section_settings"`
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"name_placeholder"`
	Code                   string `json:"code"`
	CodePlaceholder        string `json:"code_placeholder"`
	Type                   string `json:"type"`
	NetDays                string `json:"net_days"`
	DiscountDays           string `json:"discount_days"`
	DiscountPercentBps     string `json:"discount_percent_bps"`
	TypeHint               string `json:"type_hint"`
	NetDaysHint            string `json:"net_days_hint"`
	DiscountDaysHint       string `json:"discount_days_hint"`
	DiscountPercentBpsHint string `json:"discount_percent_bps_hint"`
	PriorityHint           string `json:"priority_hint"`
	EntityScope            string `json:"entity_scope"`
	IsDefault              string `json:"is_default"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	DisplayOrder           string `json:"display_order"`
	Active                 string `json:"active"`

	// Type select options
	TypeDueOnReceipt        string `json:"type_due_on_receipt"`
	TypeNet                 string `json:"type_net"`
	TypeCOD                 string `json:"type_cod"`
	TypeProximate           string `json:"type_proximate"`
	ProximateDay            string `json:"proximate_day"`
	ProximateDayPlaceholder string `json:"proximate_day_placeholder"`

	// Entity scope select options
	ScopesBoth         string `json:"scopes_both"`
	ScopesSupplierOnly string `json:"scopes_supplier_only"`
	ScopesClientOnly   string `json:"scopes_client_only"`

	// Field-level info text for the drawer form.
	NameInfo        string `json:"name_info"`
	CodeInfo        string `json:"code_info"`
	DescriptionInfo string `json:"description_info"`
}

type ActionLabels struct {
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}
