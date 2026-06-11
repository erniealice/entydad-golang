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
	HeadingActive  string `json:"headingActive"`
	HeadingBlocked string `json:"headingBlocked"`
	HeadingOnHold  string `json:"headingOnHold"`
	Caption        string `json:"caption"`
	CaptionActive  string `json:"captionActive"`
	CaptionBlocked string `json:"captionBlocked"`
	CaptionOnHold  string `json:"captionOnHold"`
}

type ButtonLabels struct {
	AddNew string `json:"addNew"`
}

type ColumnLabels struct {
	Name         string `json:"name"`
	SupplierType string `json:"supplierType"`
	InternalID   string `json:"internalId"`
	Status       string `json:"status"`
	Category     string `json:"category"`
	PaymentTerms string `json:"paymentTerms"`
	ContactName  string `json:"contactName"`
	DateCreated  string `json:"dateCreated"`
}

type EmptyLabels struct {
	ActiveTitle    string `json:"activeTitle"`
	ActiveMessage  string `json:"activeMessage"`
	BlockedTitle   string `json:"blockedTitle"`
	BlockedMessage string `json:"blockedMessage"`
	OnHoldTitle    string `json:"onHoldTitle"`
	OnHoldMessage  string `json:"onHoldMessage"`
}

type FormLabels struct {
	Name               string `json:"name"`
	SupplierType       string `json:"supplierType"`
	TaxID              string `json:"taxId"`
	RegistrationNumber string `json:"registrationNumber"`
	StreetAddress      string `json:"streetAddress"`
	City               string `json:"city"`
	Province           string `json:"province"`
	PostalCode         string `json:"postalCode"`
	Country            string `json:"country"`
	BillingCurrency    string `json:"billingCurrency"`
	PaymentTerms       string `json:"paymentTerms"`
	LeadTimeDays       string `json:"leadTimeDays"`
	CreditLimit        string `json:"creditLimit"`
	Status             string `json:"status"`
	Website            string `json:"website"`
	Notes              string `json:"notes"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Active             string `json:"active"`

	// Section titles
	SectionCompany        string `json:"sectionCompany"`
	SectionRepresentative string `json:"sectionRepresentative"`
	SectionAccounting     string `json:"sectionAccounting"`
	SectionAddress        string `json:"sectionAddress"`
	SectionOthers         string `json:"sectionOthers"`

	// Timezone autocomplete
	Timezone                  string `json:"timezone"`
	TimezonePlaceholder       string `json:"timezonePlaceholder"`
	TimezoneSearchPlaceholder string `json:"timezoneSearchPlaceholder"`
	TimezoneNoResults         string `json:"timezoneNoResults"`
	TimezoneInfo              string `json:"timezoneInfo"`

	// Placeholders
	NamePlaceholder               string `json:"namePlaceholder"`
	SupplierTypePlaceholder       string `json:"supplierTypePlaceholder"`
	StatusPlaceholder             string `json:"statusPlaceholder"`
	FirstNamePlaceholder          string `json:"firstNamePlaceholder"`
	LastNamePlaceholder           string `json:"lastNamePlaceholder"`
	EmailPlaceholder              string `json:"emailPlaceholder"`
	PhonePlaceholder              string `json:"phonePlaceholder"`
	PaymentTermsPlaceholder       string `json:"paymentTermsPlaceholder"`
	CreditLimitPlaceholder        string `json:"creditLimitPlaceholder"`
	BillingCurrencyPlaceholder    string `json:"billingCurrencyPlaceholder"`
	LeadTimeDaysPlaceholder       string `json:"leadTimeDaysPlaceholder"`
	TaxIDPlaceholder              string `json:"taxIdPlaceholder"`
	RegistrationNumberPlaceholder string `json:"registrationNumberPlaceholder"`
	StreetAddressPlaceholder      string `json:"streetAddressPlaceholder"`
	CityPlaceholder               string `json:"cityPlaceholder"`
	ProvincePlaceholder           string `json:"provincePlaceholder"`
	PostalCodePlaceholder         string `json:"postalCodePlaceholder"`
	CountryPlaceholder            string `json:"countryPlaceholder"`
	WebsitePlaceholder            string `json:"websitePlaceholder"`
	NotesPlaceholder              string `json:"notesPlaceholder"`

	// Select option labels
	TypeCompany    string `json:"typeCompany"`
	TypeIndividual string `json:"typeIndividual"`

	StatusActive  string `json:"statusActive"`
	StatusBlocked string `json:"statusBlocked"`
	StatusOnHold  string `json:"statusOnHold"`

	TermsImmediate string `json:"termsImmediate"`
	TermsNet30     string `json:"termsNet30"`
	TermsNet60     string `json:"termsNet60"`
	Terms2_10Net30 string `json:"terms2_10Net30"`
}

type DetailLabels struct {
	InfoTab       string              `json:"infoTab"`
	CompanyInfo   DetailSectionLabels `json:"companyInfo"`
	ContactInfo   DetailSectionLabels `json:"contactInfo"`
	FinancialInfo DetailSectionLabels `json:"financialInfo"`
	AddressInfo   DetailSectionLabels `json:"addressInfo"`
	// Tab label for attachments
	AttachmentsTab string `json:"attachmentsTab"`
	// Tab label for audit history
	AuditHistoryTab string `json:"auditHistoryTab"`
	// Tab label for statement
	StatementTab string `json:"statementTab"`
	// Inline labels
	DaysSuffix string `json:"daysSuffix"`
	Website    string `json:"website"`
	// Purchase Orders tab labels
	PurchaseOrders PurchaseOrdersLabels `json:"purchaseOrders"`
	// Statement tab stat card labels
	OutstandingBalance string `json:"outstandingBalance"`
	TotalBilled        string `json:"totalBilled"`
	TotalPaid          string `json:"totalPaid"`
	Bills              string `json:"bills"`
	// Statement empty state
	EmptyStatementTitle   string `json:"emptyStatementTitle"`
	EmptyStatementMessage string `json:"emptyStatementMessage"`
}

// DetailSectionLabels holds a title for a detail page section.
type DetailSectionLabels struct {
	Title string `json:"title"`
}

// PurchaseOrdersLabels holds labels for the purchase orders tab on the supplier detail page.
type PurchaseOrdersLabels struct {
	Title        string `json:"title"`
	ColPONumber  string `json:"colPONumber"`
	ColOrderDate string `json:"colOrderDate"`
	ColAmount    string `json:"colAmount"`
	ColCurrency  string `json:"colCurrency"`
	ColStatus    string `json:"colStatus"`
	EmptyPO      string `json:"emptyPO"`
}

type ActionLabels struct {
	View      string `json:"view"`
	Edit      string `json:"edit"`
	Delete    string `json:"delete"`
	Activate  string `json:"activate"`
	Block     string `json:"block"`
	SetOnHold string `json:"setOnHold"`
}

// DashboardLabels holds translatable strings for the supplier dashboard.
type DashboardLabels struct {
	TotalSuppliers   string `json:"totalSuppliers"`
	Active           string `json:"active"`
	Blocked          string `json:"blocked"`
	OnHold           string `json:"onHold"`
	SupplierActivity string `json:"supplierActivity"`
	TopSuppliers     string `json:"topSuppliers"`
	FilterWeek       string `json:"filterWeek"`
	FilterMonth      string `json:"filterMonth"`
	FilterYear       string `json:"filterYear"`
	RecentActivity   string `json:"recentActivity"`
	ViewAll          string `json:"viewAll"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew        string `json:"quickNew"`
	QuickViewAll    string `json:"quickViewAll"`
	QuickTags       string `json:"quickTags"`
	QuickCategories string `json:"quickCategories"`

	// Activity feed titles
	SupplierAdded     string `json:"supplierAdded"`
	SupplierActivated string `json:"supplierActivated"`
	DetailsUpdated    string `json:"detailsUpdated"`
	TagAssigned       string `json:"tagAssigned"`
}
