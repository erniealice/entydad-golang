package supplier

// labels.go — Supplier label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// party sub-context). Pure structural move — field names, json tags, and
// string literals are byte-identical. Entity-local rename: SupplierLabels ->
// Labels, Supplier<Xxx>Labels -> <Xxx>Labels, SupplierDashboardLabels ->
// DashboardLabels.

// Labels holds all translatable strings for the supplier module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Detail  DetailLabels `json:"detail"`
	Actions ActionLabels `json:"actions"`
}

type PageLabels struct {
	Heading        string `json:"heading"`
	HeadingActive  string `json:"heading_active"`
	HeadingBlocked string `json:"heading_blocked"`
	HeadingOnHold  string `json:"heading_on_hold"`
	Caption        string `json:"caption"`
	CaptionActive  string `json:"caption_active"`
	CaptionBlocked string `json:"caption_blocked"`
	CaptionOnHold  string `json:"caption_on_hold"`
}

type ButtonLabels struct {
	AddNew string `json:"add_new"`
}

type ColumnLabels struct {
	Name         string `json:"name"`
	SupplierType string `json:"supplier_type"`
	InternalID   string `json:"internal_id"`
	Status       string `json:"status"`
	Category     string `json:"category"`
	PaymentTerms string `json:"payment_terms"`
	ContactName  string `json:"contact_name"`
	DateCreated  string `json:"date_created"`
}

type EmptyLabels struct {
	ActiveTitle    string `json:"active_title"`
	ActiveMessage  string `json:"active_message"`
	BlockedTitle   string `json:"blocked_title"`
	BlockedMessage string `json:"blocked_message"`
	OnHoldTitle    string `json:"on_hold_title"`
	OnHoldMessage  string `json:"on_hold_message"`
}

type FormLabels struct {
	Name               string `json:"name"`
	SupplierType       string `json:"supplier_type"`
	TaxID              string `json:"tax_id"`
	RegistrationNumber string `json:"registration_number"`
	StreetAddress      string `json:"street_address"`
	City               string `json:"city"`
	Province           string `json:"province"`
	PostalCode         string `json:"postal_code"`
	Country            string `json:"country"`
	BillingCurrency    string `json:"billing_currency"`
	PaymentTerms       string `json:"payment_terms"`
	LeadTimeDays       string `json:"lead_time_days"`
	CreditLimit        string `json:"credit_limit"`
	Status             string `json:"status"`
	Website            string `json:"website"`
	Notes              string `json:"notes"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Active             string `json:"active"`

	// Section titles
	SectionCompany        string `json:"section_company"`
	SectionRepresentative string `json:"section_representative"`
	SectionAccounting     string `json:"section_accounting"`
	SectionAddress        string `json:"section_address"`
	SectionOthers         string `json:"section_others"`

	// Timezone autocomplete
	Timezone                  string `json:"timezone"`
	TimezonePlaceholder       string `json:"timezone_placeholder"`
	TimezoneSearchPlaceholder string `json:"timezone_search_placeholder"`
	TimezoneNoResults         string `json:"timezone_no_results"`
	TimezoneInfo              string `json:"timezone_info"`

	// Placeholders
	NamePlaceholder               string `json:"name_placeholder"`
	SupplierTypePlaceholder       string `json:"supplier_type_placeholder"`
	StatusPlaceholder             string `json:"status_placeholder"`
	FirstNamePlaceholder          string `json:"first_name_placeholder"`
	LastNamePlaceholder           string `json:"last_name_placeholder"`
	EmailPlaceholder              string `json:"email_placeholder"`
	PhonePlaceholder              string `json:"phone_placeholder"`
	PaymentTermsPlaceholder       string `json:"payment_terms_placeholder"`
	CreditLimitPlaceholder        string `json:"credit_limit_placeholder"`
	BillingCurrencyPlaceholder    string `json:"billing_currency_placeholder"`
	LeadTimeDaysPlaceholder       string `json:"lead_time_days_placeholder"`
	TaxIDPlaceholder              string `json:"tax_id_placeholder"`
	RegistrationNumberPlaceholder string `json:"registration_number_placeholder"`
	StreetAddressPlaceholder      string `json:"street_address_placeholder"`
	CityPlaceholder               string `json:"city_placeholder"`
	ProvincePlaceholder           string `json:"province_placeholder"`
	PostalCodePlaceholder         string `json:"postal_code_placeholder"`
	CountryPlaceholder            string `json:"country_placeholder"`
	WebsitePlaceholder            string `json:"website_placeholder"`
	NotesPlaceholder              string `json:"notes_placeholder"`

	// Select option labels
	TypeCompany    string `json:"type_company"`
	TypeIndividual string `json:"type_individual"`

	StatusActive  string `json:"status_active"`
	StatusBlocked string `json:"status_blocked"`
	StatusOnHold  string `json:"status_on_hold"`

	TermsImmediate string `json:"terms_immediate"`
	TermsNet30     string `json:"terms_net30"`
	TermsNet60     string `json:"terms_net60"`
	Terms2_10Net30 string `json:"terms2_10_net30"`
}

type DetailLabels struct {
	InfoTab       string              `json:"info_tab"`
	CompanyInfo   DetailSectionLabels `json:"company_info"`
	ContactInfo   DetailSectionLabels `json:"contact_info"`
	FinancialInfo DetailSectionLabels `json:"financial_info"`
	AddressInfo   DetailSectionLabels `json:"address_info"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachments_tab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"audit_history_tab"`
	// Tab label for statement
	StatementTab string `json:"statement_tab"`
	// Inline labels
	DaysSuffix string `json:"days_suffix"`
	Website    string `json:"website"`
	// Purchase Orders tab labels
	PurchaseOrders PurchaseOrdersLabels `json:"purchase_orders"`
	// Statement tab stat card labels
	OutstandingBalance string `json:"outstanding_balance"`
	TotalBilled        string `json:"total_billed"`
	TotalPaid          string `json:"total_paid"`
	Bills              string `json:"bills"`
	// Statement empty state
	EmptyStatementTitle   string `json:"empty_statement_title"`
	EmptyStatementMessage string `json:"empty_statement_message"`
}

// DetailSectionLabels holds a title for a detail page section.
type DetailSectionLabels struct {
	Title string `json:"title"`
}

// PurchaseOrdersLabels holds labels for the purchase orders tab on the supplier detail page.
type PurchaseOrdersLabels struct {
	Title        string `json:"title"`
	ColPONumber  string `json:"col_ponumber"`
	ColOrderDate string `json:"col_order_date"`
	ColAmount    string `json:"col_amount"`
	ColCurrency  string `json:"col_currency"`
	ColStatus    string `json:"col_status"`
	EmptyPO      string `json:"empty_po"`
}

type ActionLabels struct {
	View      string `json:"view"`
	Edit      string `json:"edit"`
	Delete    string `json:"delete"`
	Activate  string `json:"activate"`
	Block     string `json:"block"`
	SetOnHold string `json:"set_on_hold"`
}

// DashboardLabels holds translatable strings for the supplier dashboard.
type DashboardLabels struct {
	TotalSuppliers   string `json:"total_suppliers"`
	Active           string `json:"active"`
	Blocked          string `json:"blocked"`
	OnHold           string `json:"on_hold"`
	SupplierActivity string `json:"supplier_activity"`
	TopSuppliers     string `json:"top_suppliers"`
	FilterWeek       string `json:"filter_week"`
	FilterMonth      string `json:"filter_month"`
	FilterYear       string `json:"filter_year"`
	RecentActivity   string `json:"recent_activity"`
	ViewAll          string `json:"view_all"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew        string `json:"quick_new"`
	QuickViewAll    string `json:"quick_view_all"`
	QuickTags       string `json:"quick_tags"`
	QuickCategories string `json:"quick_categories"`

	// Activity feed titles
	SupplierAdded     string `json:"supplier_added"`
	SupplierActivated string `json:"supplier_activated"`
	DetailsUpdated    string `json:"details_updated"`
	TagAssigned       string `json:"tag_assigned"`
}
