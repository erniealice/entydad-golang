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
		SupplierType:       t("supplier.form.supplierType"),
		TaxID:              t("supplier.form.taxId"),
		RegistrationNumber: t("supplier.form.registrationNumber"),
		StreetAddress:      t("supplier.form.streetAddress"),
		City:               t("supplier.form.city"),
		Province:           t("supplier.form.province"),
		PostalCode:         t("supplier.form.postalCode"),
		Country:            t("supplier.form.country"),
		BillingCurrency:    t("supplier.form.billingCurrency"),
		PaymentTerms:       t("supplier.form.paymentTerms"),
		LeadTimeDays:       t("supplier.form.leadTimeDays"),
		CreditLimit:        t("supplier.form.creditLimit"),
		Status:             t("supplier.form.status"),
		Website:            t("supplier.form.website"),
		Notes:              t("supplier.form.notes"),
		FirstName:          t("supplier.form.firstName"),
		LastName:           t("supplier.form.lastName"),
		Email:              t("supplier.form.email"),
		Phone:              t("supplier.form.phone"),
		Active:             t("supplier.form.active"),

		// Section titles
		SectionCompany:        t("supplier.form.sectionCompany"),
		SectionRepresentative: t("supplier.form.sectionRepresentative"),
		SectionAccounting:     t("supplier.form.sectionAccounting"),
		SectionAddress:        t("supplier.form.sectionAddress"),
		SectionOthers:         t("supplier.form.sectionOthers"),

		// Timezone autocomplete
		Timezone:                  t("supplier.form.timezone"),
		TimezonePlaceholder:       t("supplier.form.timezonePlaceholder"),
		TimezoneSearchPlaceholder: t("supplier.form.timezoneSearchPlaceholder"),
		TimezoneNoResults:         t("supplier.form.timezoneNoResults"),
		TimezoneInfo:              t("supplier.form.timezoneInfo"),

		// Placeholders
		NamePlaceholder:               t("supplier.form.namePlaceholder"),
		SupplierTypePlaceholder:       t("supplier.form.supplierTypePlaceholder"),
		StatusPlaceholder:             t("supplier.form.statusPlaceholder"),
		FirstNamePlaceholder:          t("supplier.form.firstNamePlaceholder"),
		LastNamePlaceholder:           t("supplier.form.lastNamePlaceholder"),
		EmailPlaceholder:              t("supplier.form.emailPlaceholder"),
		PhonePlaceholder:              t("supplier.form.phonePlaceholder"),
		PaymentTermsPlaceholder:       t("supplier.form.paymentTermsPlaceholder"),
		CreditLimitPlaceholder:        t("supplier.form.creditLimitPlaceholder"),
		BillingCurrencyPlaceholder:    t("supplier.form.billingCurrencyPlaceholder"),
		LeadTimeDaysPlaceholder:       t("supplier.form.leadTimeDaysPlaceholder"),
		TaxIDPlaceholder:              t("supplier.form.taxIdPlaceholder"),
		RegistrationNumberPlaceholder: t("supplier.form.registrationNumberPlaceholder"),
		StreetAddressPlaceholder:      t("supplier.form.streetAddressPlaceholder"),
		CityPlaceholder:               t("supplier.form.cityPlaceholder"),
		ProvincePlaceholder:           t("supplier.form.provincePlaceholder"),
		PostalCodePlaceholder:         t("supplier.form.postalCodePlaceholder"),
		CountryPlaceholder:            t("supplier.form.countryPlaceholder"),
		WebsitePlaceholder:            t("supplier.form.websitePlaceholder"),
		NotesPlaceholder:              t("supplier.form.notesPlaceholder"),

		// Select option labels
		TypeCompany:    t("supplier.form.typeCompany"),
		TypeIndividual: t("supplier.form.typeIndividual"),

		StatusActive:  t("supplier.form.statusActive"),
		StatusBlocked: t("supplier.form.statusBlocked"),
		StatusOnHold:  t("supplier.form.statusOnHold"),

		SelectPaymentTerm: t("supplier.form.selectPaymentTerm"),

		Tags:                   t("supplier.form.tags"),
		TagsPlaceholder:        t("supplier.form.tagsPlaceholder"),
		TagsSearchPlaceholder:  t("supplier.form.tagsSearchPlaceholder"),
		TagsNoResults:          t("supplier.form.tagsNoResults"),
		NameInfo:               t("supplier.form.nameInfo"),
		SupplierTypeInfo:       t("supplier.form.supplierTypeInfo"),
		StatusInfo:             t("supplier.form.statusInfo"),
		EmailInfo:              t("supplier.form.emailInfo"),
		PhoneInfo:              t("supplier.form.phoneInfo"),
		PaymentTermsInfo:       t("supplier.form.paymentTermsInfo"),
		CreditLimitInfo:        t("supplier.form.creditLimitInfo"),
		BillingCurrencyInfo:    t("supplier.form.billingCurrencyInfo"),
		LeadTimeDaysInfo:       t("supplier.form.leadTimeDaysInfo"),
		TaxIDInfo:              t("supplier.form.taxIdInfo"),
		RegistrationNumberInfo: t("supplier.form.registrationNumberInfo"),
		NotesInfo:              t("supplier.form.notesInfo"),
		ActiveInfo:             t("supplier.form.activeInfo"),
	}
}
