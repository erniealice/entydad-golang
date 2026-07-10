package form

import (
	pyeza "github.com/erniealice/pyeza-golang/types"
)

// Labels holds i18n labels for the client drawer form template.
type Labels struct {
	Name                       string
	NamePlaceholder            string
	CompanyDetails             string
	Representative             string
	FirstName                  string
	FirstNamePlaceholder       string
	LastName                   string
	LastNamePlaceholder        string
	Email                      string
	EmailPlaceholder           string
	Mobile                     string
	MobilePlaceholder          string
	Active                     string
	StreetAddress              string
	StreetAddressPlaceholder   string
	City                       string
	CityPlaceholder            string
	Province                   string
	ProvincePlaceholder        string
	PostalCode                 string
	PostalCodePlaceholder      string
	Notes                      string
	NotesPlaceholder           string
	PaymentTerms               string
	SelectPaymentTerm          string
	Tags                       string
	TagsPlaceholder            string
	TagsSearchPlaceholder      string
	TagsNoResults              string
	Accounting                 string
	BillingCurrency            string
	BillingCurrencyPlaceholder string
	BillingCurrencyInfo        string
	Timezone                   string
	TimezonePlaceholder        string
	TimezoneSearchPlaceholder  string
	TimezoneNoResults          string
	TimezoneInfo               string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo         string
	EmailInfo        string
	MobileInfo       string
	NotesInfo        string
	PaymentTermsInfo string
	TagsInfo         string
	ActiveInfo       string

	// Section + status + address field labels
	Status                string
	StatusPlaceholder     string
	StatusActive          string
	StatusBlocked         string
	StatusOnHold          string
	StatusInactive        string
	StatusProspect        string
	Country               string
	CountryPlaceholder    string
	Website               string
	WebsitePlaceholder    string
	SectionCompany        string
	SectionAddress        string
	SectionRepresentative string
	SectionAccounting     string
	SectionOthers         string

	// Accounting fields (mirrored from supplier)
	TaxID                         string
	TaxIDPlaceholder              string
	TaxIDInfo                     string
	RegistrationNumber            string
	RegistrationNumberPlaceholder string
	RegistrationNumberInfo        string
	CreditLimit                   string
	CreditLimitPlaceholder        string
	CreditLimitInfo               string
	LeadTimeDays                  string
	LeadTimeDaysPlaceholder       string
	LeadTimeDaysInfo              string

	// Tax identity fields (Phase 5)
	TIN            string
	TINPlaceholder string
	TINInfo        string

	// CountryCode — ISO 3166-1 alpha-2 (Phase 5 H2).
	// Separate from the legacy Country free-text field; drives jurisdiction lookup.
	CountryCode            string
	CountryCodePlaceholder string
	CountryCodeInfo        string
}

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id   string
	Name string
}

// TagOption represents a tag available for selection in the form.
// Fields named Value/Label to match the pyeza multi-select component template.
type TagOption struct {
	Value    string
	Label    string
	Selected bool
}

// SelectedTag represents a pre-selected tag for chip rendering in the multi-select.
type SelectedTag struct {
	Value string
	Label string
}

// Data is the template data for the client drawer form.
type Data struct {
	FormAction         string
	WorkspaceID        string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit             bool
	ID                 string
	Mode               string
	Name               string
	FirstName          string
	LastName           string
	Email              string
	Mobile             string
	Timezone           string
	Active             bool
	Status             string
	Country            string
	Website            string
	StreetAddress      string
	City               string
	Province           string
	PostalCode         string
	Notes              string
	BillingCurrency    string
	TaxID              string
	RegistrationNumber string
	CreditLimit        string // form-input string; converted to int64 centavos in POST
	LeadTimeDays       string
	SearchTimezonesURL string
	// Tax identity fields (Phase 5)
	TIN                      string
	CountryCode              string
	PaymentTerms             []*PaymentTermOption
	SelectedPaymentTermID    string
	PaymentTermSelectOptions []pyeza.SelectOption
	StatusOptions            []pyeza.SelectOption
	BillingCurrencyOptions   []pyeza.SelectOption
	TagOptions               []TagOption
	SelectedTags             []SelectedTag
	Labels                   Labels
	CommonLabels             any
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — a narrow func(string) string with no Deps or storage access.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		Name:                          t("client.form.name"),
		NamePlaceholder:               t("client.form.name_placeholder"),
		CompanyDetails:                t("client.form.company_details"),
		Representative:                t("client.form.representative"),
		FirstName:                     t("client.form.first_name"),
		FirstNamePlaceholder:          t("client.form.first_name_placeholder"),
		LastName:                      t("client.form.last_name"),
		LastNamePlaceholder:           t("client.form.last_name_placeholder"),
		Email:                         t("client.form.email"),
		EmailPlaceholder:              t("client.form.email_placeholder"),
		Mobile:                        t("client.form.phone"),
		MobilePlaceholder:             t("client.form.phone_placeholder"),
		Active:                        t("client.form.active"),
		StreetAddress:                 t("client.form.street_address"),
		StreetAddressPlaceholder:      t("client.form.street_address_placeholder"),
		City:                          t("client.form.city"),
		CityPlaceholder:               t("client.form.city_placeholder"),
		Province:                      t("client.form.province"),
		ProvincePlaceholder:           t("client.form.province_placeholder"),
		PostalCode:                    t("client.form.postal_code"),
		PostalCodePlaceholder:         t("client.form.postal_code_placeholder"),
		Notes:                         t("client.form.notes"),
		NotesPlaceholder:              t("client.form.notes_placeholder"),
		PaymentTerms:                  t("client.form.payment_terms"),
		SelectPaymentTerm:             t("client.form.select_payment_term"),
		Tags:                          t("client.form.tags"),
		TagsPlaceholder:               t("client.form.tags_placeholder"),
		TagsSearchPlaceholder:         t("client.form.tags_search_placeholder"),
		TagsNoResults:                 t("client.form.tags_no_results"),
		NameInfo:                      t("client.form.name_info"),
		EmailInfo:                     t("client.form.email_info"),
		MobileInfo:                    t("client.form.mobile_info"),
		NotesInfo:                     t("client.form.notes_info"),
		PaymentTermsInfo:              t("client.form.payment_terms_info"),
		TagsInfo:                      t("client.form.tags_info"),
		ActiveInfo:                    t("client.form.active_info"),
		Accounting:                    t("client.form.accounting"),
		BillingCurrency:               t("client.form.billing_currency"),
		BillingCurrencyPlaceholder:    t("client.form.billing_currency_placeholder"),
		BillingCurrencyInfo:           t("client.form.billing_currency_info"),
		Timezone:                      t("client.form.timezone"),
		TimezonePlaceholder:           t("client.form.timezone_placeholder"),
		TimezoneSearchPlaceholder:     t("client.form.timezone_search_placeholder"),
		TimezoneNoResults:             t("client.form.timezone_no_results"),
		TimezoneInfo:                  t("client.form.timezone_info"),
		Status:                        t("client.form.status"),
		StatusPlaceholder:             t("client.form.status_placeholder"),
		StatusActive:                  t("client.form.status_active"),
		StatusBlocked:                 t("client.form.status_blocked"),
		StatusOnHold:                  t("client.form.status_on_hold"),
		StatusInactive:                t("client.form.status_inactive"),
		StatusProspect:                t("client.form.status_prospect"),
		Country:                       t("client.form.country"),
		CountryPlaceholder:            t("client.form.country_placeholder"),
		Website:                       t("client.form.website"),
		WebsitePlaceholder:            t("client.form.website_placeholder"),
		SectionCompany:                t("client.form.section_company"),
		SectionAddress:                t("client.form.section_address"),
		SectionRepresentative:         t("client.form.section_representative"),
		SectionAccounting:             t("client.form.section_accounting"),
		SectionOthers:                 t("client.form.section_others"),
		TaxID:                         t("client.form.tax_id"),
		TaxIDPlaceholder:              t("client.form.tax_id_placeholder"),
		TaxIDInfo:                     t("client.form.tax_id_info"),
		RegistrationNumber:            t("client.form.registration_number"),
		RegistrationNumberPlaceholder: t("client.form.registration_number_placeholder"),
		RegistrationNumberInfo:        t("client.form.registration_number_info"),
		CreditLimit:                   t("client.form.credit_limit"),
		CreditLimitPlaceholder:        t("client.form.credit_limit_placeholder"),
		CreditLimitInfo:               t("client.form.credit_limit_info"),
		LeadTimeDays:                  t("client.form.lead_time_days"),
		LeadTimeDaysPlaceholder:       t("client.form.lead_time_days_placeholder"),
		LeadTimeDaysInfo:              t("client.form.lead_time_days_info"),
		// Tax identity (Phase 5)
		TIN:            t("client.form.tin"),
		TINPlaceholder: t("client.form.tin_placeholder"),
		TINInfo:        t("client.form.tin_info"),
		// CountryCode (Phase 5 H2)
		CountryCode:            t("client.form.country_code"),
		CountryCodePlaceholder: t("client.form.country_code_placeholder"),
		CountryCodeInfo:        t("client.form.country_code_info"),
	}
}
