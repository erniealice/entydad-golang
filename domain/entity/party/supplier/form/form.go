package form

import pyeza "github.com/erniealice/pyeza-golang/types"

// Labels holds i18n labels for the supplier drawer form template.
type Labels struct {
	Name               string
	SupplierType       string
	TaxID              string
	RegistrationNumber string
	StreetAddress      string
	City               string
	Province           string
	PostalCode         string
	Country            string
	BillingCurrency    string
	PaymentTerms       string
	LeadTimeDays       string
	CreditLimit        string
	Status             string
	Website            string
	Notes              string
	FirstName          string
	LastName           string
	Email              string
	Phone              string
	Active             string

	// Section titles
	SectionCompany        string
	SectionRepresentative string
	SectionAccounting     string
	SectionAddress        string
	SectionOthers         string

	// Timezone autocomplete
	Timezone                  string
	TimezonePlaceholder       string
	TimezoneSearchPlaceholder string
	TimezoneNoResults         string
	TimezoneInfo              string

	// Placeholders
	NamePlaceholder               string
	SupplierTypePlaceholder       string
	StatusPlaceholder             string
	FirstNamePlaceholder          string
	LastNamePlaceholder           string
	EmailPlaceholder              string
	PhonePlaceholder              string
	PaymentTermsPlaceholder       string
	CreditLimitPlaceholder        string
	BillingCurrencyPlaceholder    string
	LeadTimeDaysPlaceholder       string
	TaxIDPlaceholder              string
	RegistrationNumberPlaceholder string
	StreetAddressPlaceholder      string
	CityPlaceholder               string
	ProvincePlaceholder           string
	PostalCodePlaceholder         string
	CountryPlaceholder            string
	WebsitePlaceholder            string
	NotesPlaceholder              string

	// Select option labels
	TypeCompany    string
	TypeIndividual string

	StatusActive  string
	StatusBlocked string
	StatusOnHold  string

	SelectPaymentTerm string

	Tags                  string
	TagsPlaceholder       string
	TagsSearchPlaceholder string
	TagsNoResults         string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo               string
	SupplierTypeInfo       string
	StatusInfo             string
	EmailInfo              string
	PhoneInfo              string
	PaymentTermsInfo       string
	CreditLimitInfo        string
	BillingCurrencyInfo    string
	LeadTimeDaysInfo       string
	TaxIDInfo              string
	RegistrationNumberInfo string
	NotesInfo              string
	ActiveInfo             string
}

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id   string
	Name string
}

// TagOption represents a tag available for selection in the form.
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

// Data is the template data for the supplier drawer form.
type Data struct {
	FormAction               string
	WorkspaceID              string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit                   bool
	ID                       string
	Name                     string
	Timezone                 string
	SearchTimezonesURL       string
	SupplierType             string
	TaxID                    string
	RegistrationNumber       string
	StreetAddress            string
	City                     string
	Province                 string
	PostalCode               string
	Country                  string
	BillingCurrency          string
	PaymentTerms             []*PaymentTermOption
	SelectedPaymentTermID    string
	LeadTimeDays             string
	CreditLimit              string
	Status                   string
	Website                  string
	Notes                    string
	FirstName                string
	LastName                 string
	Email                    string
	Phone                    string
	Active                   bool
	Labels                   Labels
	CommonLabels             any
	TagOptions               []TagOption
	SelectedTags             []SelectedTag
	StatusOptions            []pyeza.SelectOption
	SupplierTypeOptions      []pyeza.SelectOption
	PaymentTermSelectOptions []pyeza.SelectOption
	BillingCurrencyOptions   []pyeza.SelectOption
}

// BuildLabels constructs a Labels struct from a translation function.
// t is typically viewCtx.T — a narrow func(string) string with no Deps or storage access.
func BuildLabels(t func(string) string) Labels {
	return Labels{
		Name:               t("supplier.form.name"),
		SupplierType:       t("supplier.form.supplier_type"),
		TaxID:              t("supplier.form.tax_id"),
		RegistrationNumber: t("supplier.form.registration_number"),
		StreetAddress:      t("supplier.form.street_address"),
		City:               t("supplier.form.city"),
		Province:           t("supplier.form.province"),
		PostalCode:         t("supplier.form.postal_code"),
		Country:            t("supplier.form.country"),
		BillingCurrency:    t("supplier.form.billing_currency"),
		PaymentTerms:       t("supplier.form.payment_terms"),
		LeadTimeDays:       t("supplier.form.lead_time_days"),
		CreditLimit:        t("supplier.form.credit_limit"),
		Status:             t("supplier.form.status"),
		Website:            t("supplier.form.website"),
		Notes:              t("supplier.form.notes"),
		FirstName:          t("supplier.form.first_name"),
		LastName:           t("supplier.form.last_name"),
		Email:              t("supplier.form.email"),
		Phone:              t("supplier.form.phone"),
		Active:             t("supplier.form.active"),

		// Section titles
		SectionCompany:        t("supplier.form.section_company"),
		SectionRepresentative: t("supplier.form.section_representative"),
		SectionAccounting:     t("supplier.form.section_accounting"),
		SectionAddress:        t("supplier.form.section_address"),
		SectionOthers:         t("supplier.form.section_others"),

		// Timezone autocomplete
		Timezone:                  t("supplier.form.timezone"),
		TimezonePlaceholder:       t("supplier.form.timezone_placeholder"),
		TimezoneSearchPlaceholder: t("supplier.form.timezone_search_placeholder"),
		TimezoneNoResults:         t("supplier.form.timezone_no_results"),
		TimezoneInfo:              t("supplier.form.timezone_info"),

		// Placeholders
		NamePlaceholder:               t("supplier.form.name_placeholder"),
		SupplierTypePlaceholder:       t("supplier.form.supplier_type_placeholder"),
		StatusPlaceholder:             t("supplier.form.status_placeholder"),
		FirstNamePlaceholder:          t("supplier.form.first_name_placeholder"),
		LastNamePlaceholder:           t("supplier.form.last_name_placeholder"),
		EmailPlaceholder:              t("supplier.form.email_placeholder"),
		PhonePlaceholder:              t("supplier.form.phone_placeholder"),
		PaymentTermsPlaceholder:       t("supplier.form.payment_terms_placeholder"),
		CreditLimitPlaceholder:        t("supplier.form.credit_limit_placeholder"),
		BillingCurrencyPlaceholder:    t("supplier.form.billing_currency_placeholder"),
		LeadTimeDaysPlaceholder:       t("supplier.form.lead_time_days_placeholder"),
		TaxIDPlaceholder:              t("supplier.form.tax_id_placeholder"),
		RegistrationNumberPlaceholder: t("supplier.form.registration_number_placeholder"),
		StreetAddressPlaceholder:      t("supplier.form.street_address_placeholder"),
		CityPlaceholder:               t("supplier.form.city_placeholder"),
		ProvincePlaceholder:           t("supplier.form.province_placeholder"),
		PostalCodePlaceholder:         t("supplier.form.postal_code_placeholder"),
		CountryPlaceholder:            t("supplier.form.country_placeholder"),
		WebsitePlaceholder:            t("supplier.form.website_placeholder"),
		NotesPlaceholder:              t("supplier.form.notes_placeholder"),

		// Select option labels
		TypeCompany:    t("supplier.form.type_company"),
		TypeIndividual: t("supplier.form.type_individual"),

		StatusActive:  t("supplier.form.status_active"),
		StatusBlocked: t("supplier.form.status_blocked"),
		StatusOnHold:  t("supplier.form.status_on_hold"),

		SelectPaymentTerm: t("supplier.form.select_payment_term"),

		Tags:                   t("supplier.form.tags"),
		TagsPlaceholder:        t("supplier.form.tags_placeholder"),
		TagsSearchPlaceholder:  t("supplier.form.tags_search_placeholder"),
		TagsNoResults:          t("supplier.form.tags_no_results"),
		NameInfo:               t("supplier.form.name_info"),
		SupplierTypeInfo:       t("supplier.form.supplier_type_info"),
		StatusInfo:             t("supplier.form.status_info"),
		EmailInfo:              t("supplier.form.email_info"),
		PhoneInfo:              t("supplier.form.phone_info"),
		PaymentTermsInfo:       t("supplier.form.payment_terms_info"),
		CreditLimitInfo:        t("supplier.form.credit_limit_info"),
		BillingCurrencyInfo:    t("supplier.form.billing_currency_info"),
		LeadTimeDaysInfo:       t("supplier.form.lead_time_days_info"),
		TaxIDInfo:              t("supplier.form.tax_id_info"),
		RegistrationNumberInfo: t("supplier.form.registration_number_info"),
		NotesInfo:              t("supplier.form.notes_info"),
		ActiveInfo:             t("supplier.form.active_info"),
	}
}
