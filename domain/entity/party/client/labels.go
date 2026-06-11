package client

// labels.go — Client label structs.
//
// Extracted verbatim from packages/entydad-golang/labels.go (entity domain,
// party sub-context). Pure structural move — field names, json tags, and
// string literals are byte-identical. Entity-local rename: ClientLabels ->
// Labels, Client<Xxx>Labels -> <Xxx>Labels, ClientDashboardLabels ->
// DashboardLabels.

import (
	"strings"
)

// Labels holds all translatable strings for the client module.
// JSON tags match the "client" wrapper key in retail/client.json.
type Labels struct {
	Page        PageLabels       `json:"page"`
	Buttons     ButtonLabels     `json:"buttons"`
	Columns     ColumnLabels     `json:"columns"`
	Empty       EmptyLabels      `json:"empty"`
	Form        FormLabels       `json:"form"`
	Detail      DetailLabels     `json:"detail"`
	BulkActions BulkActionLabels `json:"bulkActions"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingProspect string `json:"headingProspect"`
	HeadingOnHold   string `json:"headingOnHold"`
	HeadingBlocked  string `json:"headingBlocked"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionProspect string `json:"captionProspect"`
	CaptionOnHold   string `json:"captionOnHold"`
	CaptionBlocked  string `json:"captionBlocked"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddNew string `json:"addNew"`
}

type ColumnLabels struct {
	ClientName          string `json:"clientName"`
	Representative      string `json:"representative"`
	Status              string `json:"status"`
	Category            string `json:"category"`
	ActiveSubscriptions string `json:"activeSubscriptions"`
	PaymentTerm         string `json:"paymentTerm"`
	DateCreated         string `json:"dateCreated"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	ProspectTitle   string `json:"prospectTitle"`
	ProspectMessage string `json:"prospectMessage"`
	OnHoldTitle     string `json:"onHoldTitle"`
	OnHoldMessage   string `json:"onHoldMessage"`
	BlockedTitle    string `json:"blockedTitle"`
	BlockedMessage  string `json:"blockedMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
	Email                      string `json:"email"`
	Phone                      string `json:"phone"`
	Name                       string `json:"name"`
	NamePlaceholder            string `json:"namePlaceholder"`
	CompanyDetails             string `json:"companyDetails"`
	Representative             string `json:"representative"`
	StreetAddress              string `json:"streetAddress"`
	StreetAddressPlaceholder   string `json:"streetAddressPlaceholder"`
	City                       string `json:"city"`
	CityPlaceholder            string `json:"cityPlaceholder"`
	Province                   string `json:"province"`
	ProvincePlaceholder        string `json:"provincePlaceholder"`
	PostalCode                 string `json:"postalCode"`
	PostalCodePlaceholder      string `json:"postalCodePlaceholder"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notesPlaceholder"`
	Tags                       string `json:"tags"`
	TagsPlaceholder            string `json:"tagsPlaceholder"`
	TagsSearchPlaceholder      string `json:"tagsSearchPlaceholder"`
	TagsNoResults              string `json:"tagsNoResults"`
	Accounting                 string `json:"accounting"`
	BillingCurrency            string `json:"billingCurrency"`
	BillingCurrencyPlaceholder string `json:"billingCurrencyPlaceholder"`
	BillingCurrencyInfo        string `json:"billingCurrencyInfo"`
	// New fields
	Status             string `json:"status"`
	StatusPlaceholder  string `json:"statusPlaceholder"`
	StatusActive       string `json:"statusActive"`
	StatusBlocked      string `json:"statusBlocked"`
	StatusOnHold       string `json:"statusOnHold"`
	StatusInactive     string `json:"statusInactive"`
	StatusProspect     string `json:"statusProspect"`
	Country            string `json:"country"`
	CountryPlaceholder string `json:"countryPlaceholder"`
	// Phase 5 H2 — ISO 3166 alpha-2 country code separate from legacy Country
	CountryCode            string `json:"countryCode"`
	CountryCodePlaceholder string `json:"countryCodePlaceholder"`
	CountryCodeInfo        string `json:"countryCodeInfo"`
	Website                string `json:"website"`
	WebsitePlaceholder     string `json:"websitePlaceholder"`
	SectionCompany         string `json:"sectionCompany"`
	SectionAddress         string `json:"sectionAddress"`
	SectionRepresentative  string `json:"sectionRepresentative"`
	SectionAccounting      string `json:"sectionAccounting"`
	SectionOthers          string `json:"sectionOthers"`
}

type DetailLabels struct {
	CompanyDetails  CompanyDetailLabels   `json:"companyDetails"`
	Actions         DetailActionLabels    `json:"actions"`
	Profile         DetailSectionLabels   `json:"profile"`
	Company         DetailSectionLabels   `json:"company"`
	Address         DetailSectionLabels   `json:"address"`
	Representative  string                `json:"representative"`
	NotesSection    DetailSectionLabels   `json:"notesSection"`
	Tags            DetailTagLabels       `json:"tags"`
	PurchaseHistory PurchaseHistoryLabels `json:"purchaseHistory"`
	Tabs            DetailTabLabels       `json:"tabs"`
	// Flat inline labels
	Name                 string `json:"name"`
	RecentOrders         string `json:"recentOrders"`
	ColDate              string `json:"colDate"`
	ColReference         string `json:"colReference"`
	ColAmount            string `json:"colAmount"`
	ColStatus            string `json:"colStatus"`
	PurchaseHistoryEmpty string `json:"purchaseHistoryEmpty"`
	// Subscriptions tab
	AddSubscription         string `json:"addSubscription"`
	EmptySubscriptionsTitle string `json:"emptySubscriptionsTitle"`
	EmptySubscriptions      string `json:"emptySubscriptions"`
	// Statement tab stat card labels
	OutstandingBalance string `json:"outstandingBalance"`
	TotalBilled        string `json:"totalBilled"`
	TotalReceived      string `json:"totalReceived"`
	Invoices           string `json:"invoices"`
	// Statement empty state
	EmptyStatementTitle   string `json:"emptyStatementTitle"`
	EmptyStatementMessage string `json:"emptyStatementMessage"`
	// PriceSchedules tab
	PriceSchedules PriceSchedulesLabels `json:"priceSchedules"`
	// Subscriptions tab column headers + confirm dialogs
	Subscriptions SubscriptionLabels `json:"subscriptions"`
	// Statement tab column headers + totals row
	Statement StatementLabels `json:"statement"`
	// OutstandingTable tab column headers + empty state
	OutstandingTable OutstandingTableLabels `json:"outstandingTable"`
	// RevenueRun drawer labels for the per-client Run Invoices flow
	RevenueRun RevenueRunLabels `json:"revenueRun"`
}

type CompanyDetailLabels struct {
	Status string `json:"status"`
}

// DetailSectionLabels holds a title for a detail page section.
type DetailSectionLabels struct {
	Title string `json:"title"`
}

// DetailTagLabels holds labels for the tags section on the detail page.
type DetailTagLabels struct {
	Title  string `json:"title"`
	NoTags string `json:"noTags"`
}

// PurchaseHistoryLabels holds labels for the purchase history section.
type PurchaseHistoryLabels struct {
	Title         string `json:"title"`
	LifetimeSpend string `json:"lifetimeSpend"`
	TotalOrders   string `json:"totalOrders"`
	AvgOrderValue string `json:"avgOrderValue"`
	LastPurchase  string `json:"lastPurchase"`
	Empty         string `json:"empty"`
}

// DetailTabLabels holds labels for the client detail page tabs.
type DetailTabLabels struct {
	Info               string `json:"info"`
	Representative     string `json:"representative"`
	Subscriptions      string `json:"subscriptions"`
	SubscriptionsSlug  string `json:"subscriptionsSlug"`
	Accounting         string `json:"accounting"`
	History            string `json:"history"`
	Statement          string `json:"statement"`
	PriceSchedules     string `json:"priceSchedules"`
	PriceSchedulesSlug string `json:"priceSchedulesSlug"`
	Attachments        string `json:"attachments"`
	AuditHistory       string `json:"auditHistory"`
	// Phase 2 — polymorphic tax registrations tab
	TaxRegistrations string `json:"taxRegistrations"`
}

// PriceSchedulesLabels holds labels for the PriceSchedules tab on the
// client detail page. All copy uses proto-generic vocabulary ("price schedule",
// "plan", "client"); tier-specific words ("rate card", "package") live only in
// lyngua JSON overrides.
type PriceSchedulesLabels struct {
	Empty           string `json:"empty"`
	AddAction       string `json:"addAction"`
	ColumnName      string `json:"columnName"`
	ColumnDateStart string `json:"columnDateStart"`
	ColumnDateEnd   string `json:"columnDateEnd"`
	ColumnPlanCount string `json:"columnPlanCount"`
}

// SubscriptionLabels holds column headers, actions, and confirm-dialog labels
// for the Subscriptions tab table on the client detail page.
type SubscriptionLabels struct {
	ColumnName           string `json:"columnName"`
	ColumnPlan           string `json:"columnPlan"`
	ColumnStartDate      string `json:"columnStartDate"`
	ColumnEndDate        string `json:"columnEndDate"`
	ConfirmDeleteTitle   string `json:"confirmDeleteTitle"`
	ConfirmDeleteMessage string `json:"confirmDeleteMessage"`
}

// StatementLabels holds column headers and totals-row label for the
// Statement tab table on the client detail page.
type StatementLabels struct {
	ColumnDate        string `json:"columnDate"`
	ColumnType        string `json:"columnType"`
	ColumnReference   string `json:"columnReference"`
	ColumnDescription string `json:"columnDescription"`
	ColumnBilled      string `json:"columnBilled"`
	ColumnReceived    string `json:"columnReceived"`
	ColumnBalance     string `json:"columnBalance"`
	TotalsRowLabel    string `json:"totalsRowLabel"`
}

// OutstandingTableLabels holds column headers and empty-state labels for
// the outstanding-revenue table on the Statement tab.
type OutstandingTableLabels struct {
	Columns          OutstandingTableColumnLabels `json:"columns"`
	Empty            OutstandingTableEmptyLabels  `json:"empty"`
	RunInvoicesLabel string                       `json:"runInvoicesLabel"`
}

// OutstandingTableColumnLabels holds column header labels for the
// outstanding-revenue table.
type OutstandingTableColumnLabels struct {
	Date        string `json:"date"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	Billed      string `json:"billed"`
	Paid        string `json:"paid"`
	Outstanding string `json:"outstanding"`
	Status      string `json:"status"`
}

// OutstandingTableEmptyLabels holds empty-state labels for the
// outstanding-revenue table.
type OutstandingTableEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// RevenueRunLabels holds labels for the per-client Run Invoices
// drawer surfaced from the Statement-tab outstanding table.
type RevenueRunLabels struct {
	Title                   string                `json:"title"`
	SubtitleTemplate        string                `json:"subtitleTemplate"`
	AsOfDateLabel           string                `json:"asOfDateLabel"`
	AsOfDateHint            string                `json:"asOfDateHint"`
	BillThroughTodayLabel   string                `json:"billThroughTodayLabel"`
	ColumnSubscription      string                `json:"columnSubscription"`
	ColumnPeriod            string                `json:"columnPeriod"`
	ColumnAmount            string                `json:"columnAmount"`
	ColumnLines             string                `json:"columnLines"`
	GroupTotalLabel         string                `json:"groupTotalLabel"`
	GroupNoPending          string                `json:"groupNoPending"`
	GroupCurrencyMismatch   string                `json:"groupCurrencyMismatch"`
	ColumnSelectAriaLabel   string                `json:"columnSelectAriaLabel"`
	EmptyTitle              string                `json:"emptyTitle"`
	EmptyMessage            string                `json:"emptyMessage"`
	IntroMessage            string                `json:"introMessage"`
	GenerateButton          string                `json:"generateButton"`
	GenerateButtonCountOne  string                `json:"generateButtonCountOne"`
	GenerateButtonCountMany string                `json:"generateButtonCountMany"`
	CancelButton            string                `json:"cancelButton"`
	ToastSuccess            string                `json:"toastSuccess"`
	ViewRunLink             string                `json:"viewRunLink"`
	Errors                  RevenueRunErrorLabels `json:"errors"`
}

// RevenueRunErrorLabels — error copy surfaced in the drawer.
type RevenueRunErrorLabels struct {
	PermissionDenied   string `json:"permissionDenied"`
	IDRequired         string `json:"idRequired"`
	InvalidFormData    string `json:"invalidFormData"`
	UseCaseUnavailable string `json:"useCaseUnavailable"`
	SelectOne          string `json:"selectOne"`
}

// ResolveTabSlug returns the URL slug for a canonical tab key. Tier-specific
// slugs flow through here so URLs match the operator's vocabulary (e.g.
// professional ships "engagements" + "rate-cards"). Tabs without overrides
// round-trip through unchanged.
func (t DetailTabLabels) ResolveTabSlug(canonical string) string {
	switch canonical {
	case "subscriptions":
		if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" {
			return s
		}
	case "priceSchedules":
		if s := strings.TrimSpace(t.PriceSchedulesSlug); s != "" {
			return s
		}
	}
	return canonical
}

// CanonicalizeTab maps an incoming URL tab slug back to its canonical key so
// internal template lookups and equality checks stay tier-agnostic.
func (t DetailTabLabels) CanonicalizeTab(slug string) string {
	if slug == "" {
		return ""
	}
	if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" && slug == s {
		return "subscriptions"
	}
	if s := strings.TrimSpace(t.PriceSchedulesSlug); s != "" && slug == s {
		return "priceSchedules"
	}
	return slug
}

type DetailActionLabels struct {
	ViewClient       string `json:"viewClient"`
	EditClient       string `json:"editClient"`
	DeleteClient     string `json:"deleteClient"`
	DeactivateClient string `json:"deactivateClient"`
	ActivateClient   string `json:"activateClient"`
	BlockClient      string `json:"blockClient"`
	HoldClient       string `json:"holdClient"`
	SetProspect      string `json:"setProspect"`
}

type BulkActionLabels struct {
	SetAsInactive string `json:"setAsInactive"`
	SetAsActive   string `json:"setAsActive"`
	SetAsBlocked  string `json:"setAsBlocked"`
	SetAsOnHold   string `json:"setAsOnHold"`
	SetAsProspect string `json:"setAsProspect"`
}

// DashboardLabels holds translatable strings for the client dashboard.
type DashboardLabels struct {
	TotalClients   string `json:"totalClients"`
	Active         string `json:"active"`
	Inactive       string `json:"inactive"`
	NewThisMonth   string `json:"newThisMonth"`
	ClientGrowth   string `json:"clientGrowth"`
	FilterWeek     string `json:"filterWeek"`
	FilterMonth    string `json:"filterMonth"`
	FilterYear     string `json:"filterYear"`
	RecentActivity string `json:"recentActivity"`
	ViewAll        string `json:"viewAll"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew          string `json:"quickNew"`
	QuickViewAll      string `json:"quickViewAll"`
	QuickTags         string `json:"quickTags"`
	QuickPaymentTerms string `json:"quickPaymentTerms"`

	// Activity feed titles
	ClientAdded     string `json:"clientAdded"`
	ClientActivated string `json:"clientActivated"`
	ProfileUpdated  string `json:"profileUpdated"`
	TagAssigned     string `json:"tagAssigned"`
}
