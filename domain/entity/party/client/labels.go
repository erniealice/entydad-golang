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
	BulkActions BulkActionLabels `json:"bulk_actions"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"heading_active"`
	HeadingProspect string `json:"heading_prospect"`
	HeadingOnHold   string `json:"heading_on_hold"`
	HeadingBlocked  string `json:"heading_blocked"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionProspect string `json:"caption_prospect"`
	CaptionOnHold   string `json:"caption_on_hold"`
	CaptionBlocked  string `json:"caption_blocked"`
	CaptionInactive string `json:"caption_inactive"`
}

type ButtonLabels struct {
	AddNew string `json:"add_new"`
}

type ColumnLabels struct {
	ClientName          string `json:"client_name"`
	Representative      string `json:"representative"`
	Status              string `json:"status"`
	Category            string `json:"category"`
	ActiveSubscriptions string `json:"active_subscriptions"`
	PaymentTerm         string `json:"payment_term"`
	DateCreated         string `json:"date_created"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"active_title"`
	ActiveMessage   string `json:"active_message"`
	ProspectTitle   string `json:"prospect_title"`
	ProspectMessage string `json:"prospect_message"`
	OnHoldTitle     string `json:"on_hold_title"`
	OnHoldMessage   string `json:"on_hold_message"`
	BlockedTitle    string `json:"blocked_title"`
	BlockedMessage  string `json:"blocked_message"`
	InactiveTitle   string `json:"inactive_title"`
	InactiveMessage string `json:"inactive_message"`
}

type FormLabels struct {
	Email                      string `json:"email"`
	Phone                      string `json:"phone"`
	Name                       string `json:"name"`
	NamePlaceholder            string `json:"name_placeholder"`
	CompanyDetails             string `json:"company_details"`
	Representative             string `json:"representative"`
	StreetAddress              string `json:"street_address"`
	StreetAddressPlaceholder   string `json:"street_address_placeholder"`
	City                       string `json:"city"`
	CityPlaceholder            string `json:"city_placeholder"`
	Province                   string `json:"province"`
	ProvincePlaceholder        string `json:"province_placeholder"`
	PostalCode                 string `json:"postal_code"`
	PostalCodePlaceholder      string `json:"postal_code_placeholder"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notes_placeholder"`
	Tags                       string `json:"tags"`
	TagsPlaceholder            string `json:"tags_placeholder"`
	TagsSearchPlaceholder      string `json:"tags_search_placeholder"`
	TagsNoResults              string `json:"tags_no_results"`
	Accounting                 string `json:"accounting"`
	BillingCurrency            string `json:"billing_currency"`
	BillingCurrencyPlaceholder string `json:"billing_currency_placeholder"`
	BillingCurrencyInfo        string `json:"billing_currency_info"`
	// New fields
	Status             string `json:"status"`
	StatusPlaceholder  string `json:"status_placeholder"`
	StatusActive       string `json:"status_active"`
	StatusBlocked      string `json:"status_blocked"`
	StatusOnHold       string `json:"status_on_hold"`
	StatusInactive     string `json:"status_inactive"`
	StatusProspect     string `json:"status_prospect"`
	Country            string `json:"country"`
	CountryPlaceholder string `json:"country_placeholder"`
	// Phase 5 H2 — ISO 3166 alpha-2 country code separate from legacy Country
	CountryCode            string `json:"country_code"`
	CountryCodePlaceholder string `json:"country_code_placeholder"`
	CountryCodeInfo        string `json:"country_code_info"`
	Website                string `json:"website"`
	WebsitePlaceholder     string `json:"website_placeholder"`
	SectionCompany         string `json:"section_company"`
	SectionAddress         string `json:"section_address"`
	SectionRepresentative  string `json:"section_representative"`
	SectionAccounting      string `json:"section_accounting"`
	SectionOthers          string `json:"section_others"`
}

type DetailLabels struct {
	CompanyDetails  CompanyDetailLabels   `json:"company_details"`
	Actions         DetailActionLabels    `json:"actions"`
	Profile         DetailSectionLabels   `json:"profile"`
	Company         DetailSectionLabels   `json:"company"`
	Address         DetailSectionLabels   `json:"address"`
	Representative  string                `json:"representative"`
	NotesSection    DetailSectionLabels   `json:"notes_section"`
	Tags            DetailTagLabels       `json:"tags"`
	PurchaseHistory PurchaseHistoryLabels `json:"purchase_history"`
	Tabs            DetailTabLabels       `json:"tabs"`
	// Flat inline labels
	Name                 string `json:"name"`
	RecentOrders         string `json:"recent_orders"`
	ColDate              string `json:"col_date"`
	ColReference         string `json:"col_reference"`
	ColAmount            string `json:"col_amount"`
	ColStatus            string `json:"col_status"`
	PurchaseHistoryEmpty string `json:"purchase_history_empty"`
	// Subscriptions tab
	AddSubscription         string `json:"add_subscription"`
	EmptySubscriptionsTitle string `json:"empty_subscriptions_title"`
	EmptySubscriptions      string `json:"empty_subscriptions"`
	// Statement tab stat card labels
	OutstandingBalance string `json:"outstanding_balance"`
	TotalBilled        string `json:"total_billed"`
	TotalReceived      string `json:"total_received"`
	Invoices           string `json:"invoices"`
	// Statement empty state
	EmptyStatementTitle   string `json:"empty_statement_title"`
	EmptyStatementMessage string `json:"empty_statement_message"`
	// PriceSchedules tab
	PriceSchedules PriceSchedulesLabels `json:"price_schedules"`
	// Subscriptions tab column headers + confirm dialogs
	Subscriptions SubscriptionLabels `json:"subscriptions"`
	// Statement tab column headers + totals row
	Statement StatementLabels `json:"statement"`
	// OutstandingTable tab column headers + empty state
	OutstandingTable OutstandingTableLabels `json:"outstanding_table"`
	// RevenueRun drawer labels for the per-client Run Invoices flow
	RevenueRun RevenueRunLabels `json:"revenue_run"`
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
	NoTags string `json:"no_tags"`
}

// PurchaseHistoryLabels holds labels for the purchase history section.
type PurchaseHistoryLabels struct {
	Title         string `json:"title"`
	LifetimeSpend string `json:"lifetime_spend"`
	TotalOrders   string `json:"total_orders"`
	AvgOrderValue string `json:"avg_order_value"`
	LastPurchase  string `json:"last_purchase"`
	Empty         string `json:"empty"`
}

// DetailTabLabels holds labels for the client detail page tabs.
type DetailTabLabels struct {
	Info               string `json:"info"`
	Representative     string `json:"representative"`
	Subscriptions      string `json:"subscriptions"`
	SubscriptionsSlug  string `json:"subscriptions_slug"`
	Accounting         string `json:"accounting"`
	History            string `json:"history"`
	Statement          string `json:"statement"`
	PriceSchedules     string `json:"price_schedules"`
	PriceSchedulesSlug string `json:"price_schedules_slug"`
	Attachments        string `json:"attachments"`
	AuditHistory       string `json:"audit_history"`
	// Phase 2 — polymorphic tax registrations tab
	TaxRegistrations string `json:"tax_registrations"`
}

// PriceSchedulesLabels holds labels for the PriceSchedules tab on the
// client detail page. All copy uses proto-generic vocabulary ("price schedule",
// "plan", "client"); tier-specific words ("rate card", "package") live only in
// lyngua JSON overrides.
type PriceSchedulesLabels struct {
	Empty           string `json:"empty"`
	AddAction       string `json:"add_action"`
	ColumnName      string `json:"column_name"`
	ColumnDateStart string `json:"column_date_start"`
	ColumnDateEnd   string `json:"column_date_end"`
	ColumnPlanCount string `json:"column_plan_count"`
}

// SubscriptionLabels holds column headers, actions, and confirm-dialog labels
// for the Subscriptions tab table on the client detail page.
type SubscriptionLabels struct {
	ColumnName           string `json:"column_name"`
	ColumnPlan           string `json:"column_plan"`
	ColumnStartDate      string `json:"column_start_date"`
	ColumnEndDate        string `json:"column_end_date"`
	ConfirmDeleteTitle   string `json:"confirm_delete_title"`
	ConfirmDeleteMessage string `json:"confirm_delete_message"`
}

// StatementLabels holds column headers and totals-row label for the
// Statement tab table on the client detail page.
type StatementLabels struct {
	ColumnDate        string `json:"column_date"`
	ColumnType        string `json:"column_type"`
	ColumnReference   string `json:"column_reference"`
	ColumnDescription string `json:"column_description"`
	ColumnBilled      string `json:"column_billed"`
	ColumnReceived    string `json:"column_received"`
	ColumnBalance     string `json:"column_balance"`
	TotalsRowLabel    string `json:"totals_row_label"`
}

// OutstandingTableLabels holds column headers and empty-state labels for
// the outstanding-revenue table on the Statement tab.
type OutstandingTableLabels struct {
	Columns          OutstandingTableColumnLabels `json:"columns"`
	Empty            OutstandingTableEmptyLabels  `json:"empty"`
	RunInvoicesLabel string                       `json:"run_invoices_label"`
}

// OutstandingTableColumnLabels holds column header labels for the
// outstanding-revenue table.
type OutstandingTableColumnLabels struct {
	Date        string `json:"date"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
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
	SubtitleTemplate        string                `json:"subtitle_template"`
	AsOfDateLabel           string                `json:"as_of_date_label"`
	AsOfDateHint            string                `json:"as_of_date_hint"`
	BillThroughTodayLabel   string                `json:"bill_through_today_label"`
	ColumnSubscription      string                `json:"column_subscription"`
	ColumnPeriod            string                `json:"column_period"`
	ColumnAmount            string                `json:"column_amount"`
	ColumnLines             string                `json:"column_lines"`
	GroupTotalLabel         string                `json:"group_total_label"`
	GroupNoPending          string                `json:"group_no_pending"`
	GroupCurrencyMismatch   string                `json:"group_currency_mismatch"`
	ColumnSelectAriaLabel   string                `json:"column_select_aria_label"`
	EmptyTitle              string                `json:"empty_title"`
	EmptyMessage            string                `json:"empty_message"`
	IntroMessage            string                `json:"intro_message"`
	GenerateButton          string                `json:"generate_button"`
	GenerateButtonCountOne  string                `json:"generate_button_count_one"`
	GenerateButtonCountMany string                `json:"generate_button_count_many"`
	CancelButton            string                `json:"cancel_button"`
	ToastSuccess            string                `json:"toast_success"`
	ViewRunLink             string                `json:"view_run_link"`
	Errors                  RevenueRunErrorLabels `json:"errors"`
}

// RevenueRunErrorLabels — error copy surfaced in the drawer.
type RevenueRunErrorLabels struct {
	PermissionDenied   string `json:"permission_denied"`
	IDRequired         string `json:"id_required"`
	InvalidFormData    string `json:"invalid_form_data"`
	UseCaseUnavailable string `json:"use_case_unavailable"`
	SelectOne          string `json:"select_one"`
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
	ViewClient       string `json:"view_client"`
	EditClient       string `json:"edit_client"`
	DeleteClient     string `json:"delete_client"`
	DeactivateClient string `json:"deactivate_client"`
	ActivateClient   string `json:"activate_client"`
	BlockClient      string `json:"block_client"`
	HoldClient       string `json:"hold_client"`
	SetProspect      string `json:"set_prospect"`
}

type BulkActionLabels struct {
	SetAsInactive string `json:"set_as_inactive"`
	SetAsActive   string `json:"set_as_active"`
	SetAsBlocked  string `json:"set_as_blocked"`
	SetAsOnHold   string `json:"set_as_on_hold"`
	SetAsProspect string `json:"set_as_prospect"`
}

// DashboardLabels holds translatable strings for the client dashboard.
type DashboardLabels struct {
	TotalClients   string `json:"total_clients"`
	Active         string `json:"active"`
	Inactive       string `json:"inactive"`
	NewThisMonth   string `json:"new_this_month"`
	ClientGrowth   string `json:"client_growth"`
	FilterWeek     string `json:"filter_week"`
	FilterMonth    string `json:"filter_month"`
	FilterYear     string `json:"filter_year"`
	RecentActivity string `json:"recent_activity"`
	ViewAll        string `json:"view_all"`

	// Quick action labels (Phase 1b — pyeza dashboard block refactor)
	QuickNew          string `json:"quick_new"`
	QuickViewAll      string `json:"quick_view_all"`
	QuickTags         string `json:"quick_tags"`
	QuickPaymentTerms string `json:"quick_payment_terms"`

	// Activity feed titles
	ClientAdded     string `json:"client_added"`
	ClientActivated string `json:"client_activated"`
	ProfileUpdated  string `json:"profile_updated"`
	TagAssigned     string `json:"tag_assigned"`
}
